package models

// SignupInput defines the expected input for the user signup endpoint
type SignupInput struct {
	Email    string `json:"email" example:"user@example.com" binding:"required,email"`
	Password string `json:"password" example:"strongpassword123" binding:"required,min=8"`
	Role     string `json:"role" example:"User" binding:"required,oneof=Admin User"`
}
