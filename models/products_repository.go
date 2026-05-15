package models

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

var (
	ProductNotFoundError = errors.New("Product not found")
)

// PaginationParams represents the pagination parameters for querying products.
type PaginationParams struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// PaginatedResult represents the paginated result of products.
type PaginatedResult struct {
	Products   []Product `json:"data"`
	Total      int64     `json:"total"`
	TotalPages int       `json:"total_pages"`
}

// FilterParams represents the filtering parameters for querying products.
type FilterParams struct {
	CategoryCode  string  `json:"category_code"`
	PriceLessThan float64 `json:"price_less_than"`
}

// ProductsRepository provides methods to interact with the products in the db.
type ProductsRepository struct {
	db *gorm.DB
}

// NewProductsRepository creates a new instance of ProductsRepository with the given connection.
func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

// GetAll retrieves all products from the db with optional filtering and pagination.
func (r *ProductsRepository) GetAll(ctx context.Context, filter FilterParams, pagination PaginationParams) (PaginatedResult, error) {
	g := new(errgroup.Group)
	var total int64
	var products []Product

	countQuery := r.db.WithContext(ctx).Model(&Product{})
	countQuery = addFilter(filter, countQuery)
	g.Go(func() error {
		err := countQuery.Count(&total).Error
		if err != nil {
			return err
		}
		return nil
	})

	getQuery := r.db.
		WithContext(ctx).
		Preload("Variants").
		Preload("Category").
		Order("id ASC").
		Offset(pagination.Offset).
		Limit(pagination.Limit)
	getQuery = addFilter(filter, getQuery)
	g.Go(func() error {
		err := getQuery.Find(&products).Error
		if err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return PaginatedResult{}, err
	}

	return PaginatedResult{
		Products:   products,
		Total:      total,
		TotalPages: getTotalPages(total, pagination.Limit),
	}, nil
}

// GetByCode retrieves a product by its code from the db.
func (r *ProductsRepository) GetByCode(ctx context.Context, code string) (Product, error) {
	var product Product

	err := r.db.
		WithContext(ctx).
		Preload("Variants").
		Preload("Category").
		Where("code = ?", code).
		First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Product{}, ProductNotFoundError
		}
		return Product{}, err
	}

	return product, err
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

func getTotalPages(total int64, limit int) int {
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return totalPages
}
