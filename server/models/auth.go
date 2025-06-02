package models

// SignInInput represents login credentials
type SignInInput struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"securepassword"`
}
