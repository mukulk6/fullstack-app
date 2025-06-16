package main

import (
	"log"
	"os"

	"server/config"
	"server/db"
	"server/docs"
	"server/router"
	"server/scheduler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"

	_ "server/docs" // Swagger docs

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initEnv() string {
	_ = godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

func initSwagger(port string) {
	docs.SwaggerInfo.Title = "Product Management APIs"
	docs.SwaggerInfo.Description = "List of APIs for Product Management"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:" + port
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}
}

func initDB() {
	db.Init()
}

func main() {
	port := initEnv()
	initSwagger(port)
	initDB()
	defer db.Pool.Close()
	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"), // your Kafka broker address
		Topic:    "product-notifications",
		Balancer: &kafka.LeastBytes{},
	}
	defer kafkaWriter.Close()

	// ✅ Start background job for weekly product check
	scheduler.StartProductCheckScheduler(db.Pool, "Asia/Kolkata")

	// ✅ Initialize Gin
	r := gin.Default()

	// ✅ Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ✅ CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// ✅ App routes and background indexing
	router.SetupRoutes(r)
	config.InitClients()
	config.BulkSyncProductsToES()
	log.Printf("🚀 Server running at http://localhost:%s", port)
	r.Run(":" + port)
}
