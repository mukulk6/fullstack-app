// router/router.go
package router

import (
	"server/handlers"
	"server/middleware"

	"github.com/gin-gonic/gin"
	// Swagger docs import (replace with your actual module path)
)

// SetupRoutes sets up all API endpoints
func SetupRoutes(r *gin.Engine) {
	r.POST("/signup", handlers.SignUpUser)
	r.POST("/login", handlers.SignInUser)

	product := r.Group("/products")
	product.Use(middleware.AuthMiddleware())
	{

		product.GET("", handlers.GetAllProducts)
		product.GET(":id", handlers.GetProductById)

		// Admin-only routes
		product.POST("", middleware.AdminMiddleware(), handlers.CreateProduct)
		product.PUT(":id", middleware.AdminMiddleware(), handlers.UpdateProduct)
		product.DELETE(":id", middleware.AdminMiddleware(), handlers.DeleteProduct)
	}

	admin := r.Group("/admins")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		admin.GET("/list", handlers.GetAdminList)
	}
}
