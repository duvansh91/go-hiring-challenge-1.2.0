package models

import (
	"context"

	"gorm.io/gorm"
)

// CategoriesRepository provides methods to interact with the categories table in the db.
type CategoriesRepository struct {
	db *gorm.DB
}

// NewCategoriesRepository creates a new instance of CategoriesRepository with the given connection.
func NewCategoriesRepository(db *gorm.DB) *CategoriesRepository {
	return &CategoriesRepository{
		db: db,
	}
}

// GetAll retrieves all categories from the db.
func (r *CategoriesRepository) GetAll(ctx context.Context) ([]Category, error) {
	var categories []Category
	err := r.db.WithContext(ctx).Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

// Create stores a new category in the db.
func (r *CategoriesRepository) Create(ctx context.Context, category *Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}
