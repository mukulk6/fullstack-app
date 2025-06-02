package models

// AdminUserResponse represents a simplified admin user object for dropdown lists
type AdminUserResponse struct {
	Label string `json:"label" example:"admin@example.com"`
	Value int    `json:"value" example:"1"`
}
