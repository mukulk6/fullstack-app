package config

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"server/db"
	"server/models"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var (
	EsClient    *elasticsearch.Client
	RedisClient *redis.Client
)

func InitClients() {
	var err error

	EsClient, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"https://localhost:9200"},
		Username:  "elastic",
		Password:  "lOSV5DY*H5E7W3QAaX3O",
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // ✅ Accept self-signed certs
			},
		},
	})
	if err != nil {
		log.Fatalf("Elasticsearch error: %v", err)
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
}

func SearchHandler(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(400, gin.H{"error": "Query required"})
		return
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("search:%s", query)

	// Check Redis cache
	cached, err := RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedResult interface{}
		_ = json.Unmarshal([]byte(cached), &cachedResult)
		c.JSON(200, gin.H{"source": "cache", "results": cachedResult})
		return
	}

	// Build Elasticsearch query
	esQuery := fmt.Sprintf(`{
	"query": {
		"multi_match": {
			"query": "%s",
			"fields": ["name", "description"],
			"fuzziness":"AUTO"
		}
	}
}`, query)

	// Query Elasticsearch
	res, err := EsClient.Search(
		EsClient.Search.WithIndex("products"),
		EsClient.Search.WithBody(strings.NewReader(esQuery)),
		EsClient.Search.WithPretty(),
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "Elasticsearch error", "stack": err.Error()})
		return
	}
	defer res.Body.Close()
	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		c.JSON(int(res.StatusCode), gin.H{
			"error":   "Elasticsearch error",
			"details": string(bodyBytes),
		})
		return
	}

	var esResult map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&esResult); err != nil {
		c.JSON(500, gin.H{"error": "Failed to decode Elasticsearch response", "details": err.Error()})
		return
	}
	json.NewDecoder(res.Body).Decode(&esResult)

	// Cache result in Redis
	bytes, _ := json.Marshal(esResult)
	RedisClient.Set(ctx, cacheKey, bytes, 10*time.Second)

	hits := esResult["hits"].(map[string]interface{})["hits"].([]interface{})
	products := []interface{}{}

	for _, hit := range hits {
		product := hit.(map[string]interface{})["_source"]
		products = append(products, product)
	}

	c.JSON(200, gin.H{
		"source":  "elasticsearch",
		"results": products,
	})
}

func BulkSyncProductsToES() {
	ctx := context.Background()

	rows, err := db.Pool.Query(ctx, `SELECT id, name, description, price, quantity FROM products`)
	if err != nil {
		log.Fatalf("PostgreSQL query failed: %v", err)
	}
	defer rows.Close()

	var buffer bytes.Buffer

	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.Id, &p.Name, &p.Description, &p.Price, &p.Quantity)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		meta := []byte(fmt.Sprintf(`{ "index" : { "_index" : "products", "_id" : "%d" } }%s`, p.Id, "\n"))
		data, _ := json.Marshal(p)
		buffer.Write(meta)
		buffer.Write(data)
		buffer.WriteString("\n")
	}

	// Bulk upload to Elasticsearch
	res, err := EsClient.Bulk(strings.NewReader(buffer.String()), EsClient.Bulk.WithRefresh("true"))
	if err != nil {
		log.Fatalf("Elasticsearch bulk indexing failed: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Bulk indexing error: %s", res.String())
	} else {
		log.Println("Bulk indexing successful.")
	}
}
