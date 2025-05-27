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
	product := r.Group("/products")
	{
		product.GET("", handlers.GetAllProducts)       // Paginated list of products
		product.POST("", handlers.CreateProduct)       // Admin only
		product.PUT("/:id", handlers.UpdateProduct)    // Admin only
		product.DELETE("/:id", handlers.DeleteProduct) // Admin only
		product.GET("/:id", handlers.GetProductById)   //Get Product By Id
	}
	admin := r.Group("/admins")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware()) // Restrict to Admins
	{
		admin.GET("/list", handlers.GetAdminList)
	}
}
