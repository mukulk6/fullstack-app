// router/router.go
package router

import (
	"server/handlers"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/ping", handlers.PingHandler)
	r.GET("/users", handlers.GetUsers)
	r.POST("/signup", handlers.SignUpUser)
	r.POST("/login", handlers.SignInUser)

	// All product routes require authentication
	product := r.Group("/products")
	product.Use(middleware.AuthMiddleware()) // ✅ Add this
	{
		product.GET("", handlers.GetAllProducts)
		product.GET("/:id", handlers.GetProductById)

		// Admin-only routes
		product.POST("", middleware.AdminMiddleware(), handlers.CreateProduct)
		product.PUT("/:id", middleware.AdminMiddleware(), handlers.UpdateProduct)
		product.DELETE("/:id", middleware.AdminMiddleware(), handlers.DeleteProduct)
	}

	// Admin-only endpoints
	admin := r.Group("/admins")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		admin.GET("/list", handlers.GetAdminList)
	}
}
