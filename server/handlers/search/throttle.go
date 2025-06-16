package search

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"server/config"
	"time"

	"github.com/gin-gonic/gin"
)

func Throttle(limit time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		query := c.Query("query")
		key := fmt.Sprintf("throttle:%s:%s", ip, query)

		if config.RedisClient == nil {
			log.Println("RedisClient is nil. Did you forget to initialize it?")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Redis connection is not initialized",
			})
			return
		}
		exists, _ := config.RedisClient.Exists(context.Background(), key).Result()
		if exists > 0 {
			c.AbortWithStatusJSON(429, gin.H{"error": "Too many requests"})
			return
		}

		config.RedisClient.Set(context.Background(), key, 1, limit)
		c.Next()
	}
}
