package mocks

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/models"
)

// ProductsRepository is a mock implementation of the ProductsRepository interface for testing.
type ProductsRepository struct {
	GetAllFunc    func(ctx context.Context, filter models.FilterParams, pagination models.PaginationParams) (models.PaginatedResult, error)
	GetByCodeFunc func(ctx context.Context, code string) (models.Product, error)
}

// GetAll calls the mocked GetAllFunc.
func (m *ProductsRepository) GetAll(ctx context.Context, filter models.FilterParams, pagination models.PaginationParams) (models.PaginatedResult, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx, filter, pagination)
	}
	return models.PaginatedResult{}, nil
}

// GetByCode calls the mocked GetByCodeFunc.
func (m *ProductsRepository) GetByCode(ctx context.Context, code string) (models.Product, error) {
	if m.GetByCodeFunc != nil {
		return m.GetByCodeFunc(ctx, code)
	}
	return models.Product{}, nil
}
