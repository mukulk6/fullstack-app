// handlers/user.go
package handlers

import (
	"context"
	"net/http"
	"os"
	"server/db"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "hell from backend 12345678910111"})
}

func GetUsers(c *gin.Context) {
	// You can fetch users from DB here
	c.JSON(http.StatusOK, gin.H{"users": []string{"Alice", "Bob", "Gitanjali"}})
}

func SignUpUser(c *gin.Context) {
	type SignupInput struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
		Role     string `json:"role" binding:"required,oneof=Admin User"`
	}

	var input SignupInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists
	var existingEmail string
	checkQuery := `SELECT email FROM users WHERE email = $1`
	err := db.Pool.QueryRow(context.Background(), checkQuery, input.Email).Scan(&existingEmail)

	if err == nil {
		// Email found
		c.JSON(http.StatusBadRequest, gin.H{"message": "User with this email already exists.", "error": true})
		return
	} else if err != pgx.ErrNoRows {
		// Some other DB error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing user", "detail": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Insert user
	insertQuery := `INSERT INTO users (email, role, password) VALUES ($1, $2, $3)`
	_, err = db.Pool.Exec(context.Background(), insertQuery, input.Email, input.Role, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User signed up successfully"})
}

func SignInUser(c *gin.Context) {
	type SignInInput struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	var input SignInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var storedHashedPassword string
	var role string
	var id int

	// Fetch all necessary fields in one query
	query := `SELECT id, password, role FROM users WHERE email = $1`
	err := db.Pool.QueryRow(context.Background(), query, input.Email).Scan(&id, &storedHashedPassword, &role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password", "error": err.Error()})
		return
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		return
	}

	// Generate token with user ID and role
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"role":    role,
		"email":   input.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
		"role":    role,
		"id":      id,
		"email":   input.Email,
	})
}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required"`
}

func CreateProduct(c *gin.Context) {
	var input Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	query := `
		INSERT INTO products (name, description, price, quantity, created_by)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := db.Pool.Exec(context.Background(), query,
		input.Name, input.Description, input.Price, input.Quantity, userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully"})
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var input Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `UPDATE products SET name = $1, description = $2, price = $3, quantity = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5`
	_, err := db.Pool.Exec(context.Background(), query, input.Name, input.Description, input.Price, input.Quantity, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	query := `DELETE FROM products WHERE id = $1`
	_, err := db.Pool.Exec(context.Background(), query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func GetAllProducts(c *gin.Context) {
	adminUsername := c.Query("admin")

	var query string
	var rows pgx.Rows
	var err error

	if adminUsername != "" {
		query = `SELECT p.id, p.name, p.description, p.price, p.quantity
FROM products p
JOIN users u ON p.created_by = u.id
WHERE u.email = $1
ORDER BY p.updated_at DESC`

		rows, err = db.Pool.Query(context.Background(), query, adminUsername)
	} else {
		query = `
			SELECT id, name, description, price, quantity
			FROM products ORDER BY updated_at DESC
		`
		rows, err = db.Pool.Query(context.Background(), query)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var products []map[string]interface{}

	for rows.Next() {
		var id int
		var name, description string
		var price float64
		var quantity int

		err := rows.Scan(&id, &name, &description, &price, &quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning product"})
			return
		}

		products = append(products, map[string]interface{}{
			"id":          id,
			"name":        name,
			"description": description,
			"price":       price,
			"quantity":    quantity,
		})
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func GetProductById(c *gin.Context) {
	id := c.Param("id")
	var product Product
	query := `SELECT name,description,price,quantity FROM products WHERE id = $1`
	err := db.Pool.QueryRow(context.Background(), query, id).Scan(&product.Name,
		&product.Description,
		&product.Price,
		&product.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch details of the user", "error": true})
	}
	c.JSON(http.StatusOK, gin.H{"product": product})
}

func GetAdminList(c *gin.Context) {
	query := `SELECT id, email FROM users WHERE role = 'Admin' ORDER BY email ASC`

	rows, err := db.Pool.Query(context.Background(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id int
		var email string
		if err := rows.Scan(&id, &email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading admin user"})
			return
		}
		result = append(result, gin.H{
			"label": email,
			"value": id,
		})
	}

	c.JSON(http.StatusOK, result)
}
