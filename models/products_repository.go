package models

import (
	"gorm.io/gorm"
)

type PaginationParams struct {
	Offset int `json:"offset"` // number of records to skip
	Limit  int `json:"limit"`  // number of records per page
}

type PaginatedResult struct {
	Data       []Product `json:"data"`
	Total      int64     `json:"total"`
	TotalPages int       `json:"total_pages"`
}

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts(params PaginationParams) (PaginatedResult, error) {
	var products []Product
	err := r.db.Preload("Variants").Preload("Category").Offset(params.Offset).Limit(params.Limit).Find(&products).Error
	if err != nil {
		return PaginatedResult{}, err
	}

	var total int64
	if err := r.db.Model(&Product{}).Count(&total).Error; err != nil {
		return PaginatedResult{}, err
	}

	totalPages := int(total) / params.Limit
	if int(total)%params.Limit != 0 {
		totalPages++
	}

	return PaginatedResult{
		Data:       products,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}
