package mocks

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/models"
)

// CategoriesRepository is a mock implementation of the CategoriesRepository interface for testing.
type CategoriesRepository struct {
	GetAllFunc func(ctx context.Context) ([]models.Category, error)
	CreateFunc func(ctx context.Context, category *models.Category) error
}

// GetAll calls the mocked GetAllFunc.
func (m *CategoriesRepository) GetAll(ctx context.Context) ([]models.Category, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return []models.Category{}, nil
}

// Create calls the mocked CreateFunc.
func (m *CategoriesRepository) Create(ctx context.Context, category *models.Category) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, category)
	}
	return nil
}
