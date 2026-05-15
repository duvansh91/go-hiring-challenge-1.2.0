package dto

import (
	"github.com/mytheresa/go-hiring-challenge/models"
)

// Product represents a single product in the catalog response.
type Product struct {
	Code     string    `json:"code"`
	Price    float64   `json:"price"`
	Category Category  `json:"category"`
	Variants []Variant `json:"variants,omitempty"`
}

// Cataloge represents the catalog response containing products and pagination details.
type Cataloge struct {
	Products      []Product `json:"products"`
	TotalProducts int64     `json:"total_products"`
	TotalPages    int       `json:"total_pages"`
}

// MapProductToDTO converts a models.Product to a product dto.
func MapProductToDTO(product models.Product, addVariants bool) Product {
	productDto := Product{
		Code:  product.Code,
		Price: product.Price.InexactFloat64(),
		Category: Category{
			Code: product.Category.Code,
			Name: product.Category.Name,
		},
	}

	if !addVariants {
		return productDto
	}

	variantsDTO := make([]Variant, len(product.Variants))
	for i, v := range product.Variants {
		variantPrice := v.Price
		if v.Price.IsZero() {
			variantPrice = product.Price
		}
		variantsDTO[i] = Variant{
			ProductID: v.ProductID,
			Name:      v.Name,
			SKU:       v.SKU,
			Price:     variantPrice.InexactFloat64(),
		}
	}

	productDto.Variants = variantsDTO

	return productDto
}
