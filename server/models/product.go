package models

// Product represents a product item
type Product struct {
	Name        string  `json:"name" example:"Laptop"`
	Description string  `json:"description" example:"A high-end gaming laptop"`
	Price       float64 `json:"price" example:"1999.99"`
	Quantity    int     `json:"quantity" example:"10"`
	Id          int     `json:"id" example:"1"`
}
