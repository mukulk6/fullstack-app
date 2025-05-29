package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role") // You must set this during JWT parsing
		if role != "Admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Admins only."})
			c.Abort()
			return
		}
		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	type Claims struct {
		UserID int    `json:"user_id"`
		Role   string `json:"role"`
		jwt.RegisteredClaims
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or malformed token"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("role", claims.Role)
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
