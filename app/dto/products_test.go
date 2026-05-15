package dto

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/mytheresa/go-hiring-challenge/models"
)

func TestMapProductToDTO(t *testing.T) {
	t.Run("successfully maps all fields when addVariants is true", func(t *testing.T) {
		product := models.Product{
			ID:    1,
			Code:  "PROD001",
			Price: decimal.NewFromFloat(10.99),
			Category: &models.Category{
				ID:   1,
				Code: "CLOTHING",
				Name: "Clothing",
			},
			Variants: []models.Variant{
				{
					ID:        1,
					ProductID: 1,
					Name:      "Variant A",
					SKU:       "SKU001A",
					Price:     decimal.NewFromFloat(11.99),
				},
				{
					ID:        2,
					ProductID: 1,
					Name:      "Variant B",
					SKU:       "SKU001B",
					Price:     decimal.NewFromFloat(12.49),
				},
			},
		}

		result := MapProductToDTO(product, true)

		assert.Equal(t, "PROD001", result.Code, "Expected product code to match")
		assert.Equal(t, 10.99, result.Price, "Expected product price to match")
		assert.Equal(t, "CLOTHING", result.Category.Code, "Expected category code to match")
		assert.Equal(t, "Clothing", result.Category.Name, "Expected category name to match")

		assert.Len(t, result.Variants, 2, "Expected 2 variants")

		assert.Equal(t, uint(1), result.Variants[0].ProductID, "Expected first variant product ID to match")
		assert.Equal(t, "Variant A", result.Variants[0].Name, "Expected first variant name to match")
		assert.Equal(t, "SKU001A", result.Variants[0].SKU, "Expected first variant SKU to match")
		assert.Equal(t, 11.99, result.Variants[0].Price, "Expected first variant price to match")

		assert.Equal(t, uint(1), result.Variants[1].ProductID, "Expected second variant product ID to match")
		assert.Equal(t, "Variant B", result.Variants[1].Name, "Expected second variant name to match")
		assert.Equal(t, "SKU001B", result.Variants[1].SKU, "Expected second variant SKU to match")
		assert.Equal(t, 12.49, result.Variants[1].Price, "Expected second variant price to match")
	})

	t.Run("variants is nil when addVariants is false", func(t *testing.T) {
		product := models.Product{
			ID:    1,
			Code:  "PROD001",
			Price: decimal.NewFromFloat(10.99),
			Category: &models.Category{
				ID:   1,
				Code: "CLOTHING",
				Name: "Clothing",
			},
			Variants: []models.Variant{
				{
					ID:        1,
					ProductID: 1,
					Name:      "Variant A",
					SKU:       "SKU001A",
					Price:     decimal.NewFromFloat(11.99),
				},
			},
		}

		result := MapProductToDTO(product, false)

		assert.Equal(t, "PROD001", result.Code, "Expected product code to match")
		assert.Equal(t, 10.99, result.Price, "Expected product price to match")
		assert.Equal(t, "CLOTHING", result.Category.Code, "Expected category code to match")
		assert.Equal(t, "Clothing", result.Category.Name, "Expected category name to match")
		assert.Nil(t, result.Variants, "Expected variants to be nil when addVariants is false")
	})

	t.Run("variant with no price inherits from product price", func(t *testing.T) {
		product := models.Product{
			ID:    1,
			Code:  "PROD001",
			Price: decimal.NewFromFloat(10.99),
			Category: &models.Category{
				ID:   1,
				Code: "CLOTHING",
				Name: "Clothing",
			},
			Variants: []models.Variant{
				{
					ID:        1,
					ProductID: 1,
					Name:      "Variant A",
					SKU:       "SKU001A",
					Price:     decimal.Zero,
				},
				{
					ID:        2,
					ProductID: 1,
					Name:      "Variant B",
					SKU:       "SKU001B",
					Price:     decimal.NewFromFloat(15.00),
				},
			},
		}

		result := MapProductToDTO(product, true)

		assert.Len(t, result.Variants, 2, "Expected 2 variants")
		assert.Equal(t, 10.99, result.Variants[0].Price, "Expected zero-price variant to inherit from product price")
		assert.Equal(t, 15.00, result.Variants[1].Price, "Expected non-zero variant price to be preserved")
	})
}
