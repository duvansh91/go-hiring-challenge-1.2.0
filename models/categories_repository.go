package models

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

const (
	duplicatedKeyErrorCode = "23505"
)

var (
	CategoryAlreadyExistsError = errors.New("category already exists")
)

type CategoriesRepository struct {
	db *gorm.DB
}

func NewCategoriesRepository(db *gorm.DB) *CategoriesRepository {
	return &CategoriesRepository{
		db: db,
	}
}

func (r *CategoriesRepository) GetAllCategories() ([]Category, error) {
	var categories []Category
	err := r.db.Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoriesRepository) CreateCategory(category *Category) error {
	err := r.db.Create(category).Error
	if err != nil {
		if strings.Contains(err.Error(), duplicatedKeyErrorCode) {
			return CategoryAlreadyExistsError
		}
		return err
	}
	return nil
}
