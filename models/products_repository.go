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

type FilterParams struct {
	CategoryCode  string  `json:"category_code"`
	PriceLessThan float64 `json:"price_less_than"`
}

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts(filterParams FilterParams, paginationParams PaginationParams) (PaginatedResult, error) {
	var total int64
	countQuery := r.db.Model(&Product{})
	countQuery = addFilter(filterParams, countQuery)
	if err := countQuery.Count(&total).Error; err != nil {
		return PaginatedResult{}, err
	}

	var products []Product
	getQuery := r.db.Preload("Variants").Preload("Category").
		Order("id ASC").
		Offset(paginationParams.Offset).Limit(paginationParams.Limit)
	getQuery = addFilter(filterParams, getQuery)
	err := getQuery.Find(&products).Error
	if err != nil {
		return PaginatedResult{}, err
	}

	totalPages := int(total) / paginationParams.Limit
	if int(total)%paginationParams.Limit != 0 {
		totalPages++
	}

	return PaginatedResult{
		Data:       products,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func addFilter(filterParams FilterParams, db *gorm.DB) *gorm.DB {
	if filterParams.CategoryCode != "" {
		db = db.Joins("JOIN categories ON products.category_id = categories.id").
			Where("categories.code = ?", filterParams.CategoryCode)
	}

	if filterParams.PriceLessThan > 0 {
		db = db.Where("price < ?", filterParams.PriceLessThan)
	}

	return db
}
