package dto

// Variant represents a product variant response.
type Variant struct {
	ProductID uint    `json:"product_id"`
	Name      string  `json:"name"`
	SKU       string  `json:"sku"`
	Price     float64 `json:"price"`
}
