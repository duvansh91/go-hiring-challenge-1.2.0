package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/mytheresa/go-hiring-challenge/app/dto"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/mytheresa/go-hiring-challenge/models/mocks"
)

func TestNewCatalogHandler(t *testing.T) {
	t.Run("creates handler with repository", func(t *testing.T) {
		mockRepo := &mocks.ProductsRepository{}
		handler := NewCatalogHandler(mockRepo)

		assert.NotNil(t, handler, "Handler should not be nil")
		assert.Equal(t, mockRepo, handler.repo, "Repository should be set correctly")
	})
}

func TestGet(t *testing.T) {
	t.Run("successfully get cataloge with default pagination", func(t *testing.T) {
		mockRepo := &mocks.ProductsRepository{
			GetAllFunc: func(ctx context.Context, filter models.FilterParams, pagination models.PaginationParams) (models.PaginatedResult, error) {
				return models.PaginatedResult{
					Products: []models.Product{
						{
							ID:         1,
							Code:       "PROD001",
							Price:      decimal.NewFromFloat(10.99),
							CategoryID: 1,
							Category: &models.Category{
								ID:   1,
								Code: "CLOTHING",
								Name: "Clothing",
							},
						},
						{
							ID:         2,
							Code:       "PROD002",
							Price:      decimal.NewFromFloat(12.49),
							CategoryID: 2,
							Category: &models.Category{
								ID:   2,
								Code: "SHOES",
								Name: "Shoes",
							},
						},
					},
					Total:      2,
					TotalPages: 1,
				}, nil
			},
		}

		handler := NewCatalogHandler(mockRepo)
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/catalog", nil)

		handler.Get(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200 OK")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Expected Content-Type to be application/json")

		var response dto.Cataloge
		err := json.NewDecoder(recorder.Body).Decode(&response)
		assert.NoError(t, err, "Response body should be valid JSON")

		assert.Equal(t, int64(2), response.TotalProducts, "Expected total_products to be 2")
		assert.Equal(t, 1, response.TotalPages, "Expected total_pages to be 1")
		assert.Len(t, response.Products, 2, "Expected 2 products in the response")

		assert.Equal(t, "PROD001", response.Products[0].Code, "Expected first product code to be PROD001")
		assert.Equal(t, 10.99, response.Products[0].Price, "Expected first product price to be 10.99")
		assert.Equal(t, "CLOTHING", response.Products[0].Category.Code, "Expected first product category code to be CLOTHING")
		assert.Equal(t, "Clothing", response.Products[0].Category.Name, "Expected first product category name to be Clothing")

		assert.Equal(t, "PROD002", response.Products[1].Code, "Expected second product code to be PROD002")
		assert.Equal(t, 12.49, response.Products[1].Price, "Expected second product price to be 12.49")
		assert.Equal(t, "SHOES", response.Products[1].Category.Code, "Expected second product category code to be SHOES")
		assert.Equal(t, "Shoes", response.Products[1].Category.Name, "Expected second product category name to be Shoes")
	})

	t.Run("returns error when pagination offset is invalid", func(t *testing.T) {
		mockRepo := &mocks.ProductsRepository{}
		handler := NewCatalogHandler(mockRepo)
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/catalog?offset=invalid", nil)

		handler.Get(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code, "Expected status code 400 Bad Request")
		assert.JSONEq(t, `{"error":"Invalid offset param"}`, recorder.Body.String(), "Expected error message in response body")
	})

	t.Run("returns error when filter price_less_than is invalid", func(t *testing.T) {
		mockRepo := &mocks.ProductsRepository{}
		handler := NewCatalogHandler(mockRepo)
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/catalog?price_less_than=invalid", nil)

		handler.Get(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code, "Expected status code 400 Bad Request")
		assert.JSONEq(t, `{"error":"Invalid price_less_than param"}`, recorder.Body.String(), "Expected error message in response body")
	})

	t.Run("returns error when repository fails on GetAll", func(t *testing.T) {
		mockRepo := &mocks.ProductsRepository{
			GetAllFunc: func(ctx context.Context, filter models.FilterParams, pagination models.PaginationParams) (models.PaginatedResult, error) {
				return models.PaginatedResult{}, errors.New("database error")
			},
		}

		handler := NewCatalogHandler(mockRepo)
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/catalog", nil)

		handler.Get(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code, "Expected status code 500 Internal Server Error")
	})

	t.Run("valid pagination parameters", func(t *testing.T) {
		mockRepo := &mocks.ProductsRepository{
			GetAllFunc: func(ctx context.Context, filter models.FilterParams, pagination models.PaginationParams) (models.PaginatedResult, error) {
				assert.Equal(t, 5, pagination.Offset, "Expected offset to be 5")
				assert.Equal(t, 20, pagination.Limit, "Expected limit to be 20")
				return models.PaginatedResult{
					Products:   []models.Product{},
					Total:      0,
					TotalPages: 0,
				}, nil
			},
		}

		handler := NewCatalogHandler(mockRepo)
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/catalog?offset=5&limit=20", nil)

		handler.Get(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200 OK")
	})

	t.Run("valid filter parameters", func(t *testing.T) {
		mockRepo := &mocks.ProductsRepository{
			GetAllFunc: func(ctx context.Context, filter models.FilterParams, pagination models.PaginationParams) (models.PaginatedResult, error) {
				assert.Equal(t, "CLOTHING", filter.CategoryCode, "Expected category code to be CLOTHING")
				assert.Equal(t, 50.0, filter.PriceLessThan, "Expected price_less_than to be 50.0")
				return models.PaginatedResult{
					Products:   []models.Product{},
					Total:      0,
					TotalPages: 0,
				}, nil
			},
		}

		handler := NewCatalogHandler(mockRepo)
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/catalog?category_code=CLOTHING&price_less_than=50.0", nil)

		handler.Get(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200 OK")
	})
}

func TestGetProductDetail(t *testing.T) {
	t.Run("successfully gets product detail", func(t *testing.T) {
		expectedProduct := models.Product{
			ID:         1,
			Code:       "PROD001",
			Price:      decimal.NewFromFloat(10.99),
			CategoryID: 1,
			Category: &models.Category{
				ID:   1,
				Code: "CLOTHING",
				Name: "Clothing",
			},
		}

		mockRepo := &mocks.ProductsRepository{
			GetByCodeFunc: func(ctx context.Context, code string) (models.Product, error) {
				assert.Equal(t, "PROD001", code, "Expected code to be PROD001")
				return expectedProduct, nil
			},
		}

		handler := NewCatalogHandler(mockRepo)
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/catalog/PROD001", nil)
		req.SetPathValue("code", "PROD001")

		handler.GetProductDetail(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200 OK")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Expected Content-Type to be application/json")

		var response dto.Product
		err := json.NewDecoder(recorder.Body).Decode(&response)
		assert.NoError(t, err, "Response body should be valid JSON")

		assert.Equal(t, "PROD001", response.Code, "Expected product code to be PROD001")
		assert.Equal(t, 10.99, response.Price, "Expected product price to be 10.99")
		assert.Equal(t, "CLOTHING", response.Category.Code, "Expected category code to be CLOTHING")
		assert.Equal(t, "Clothing", response.Category.Name, "Expected category name to be Clothing")
	})

	t.Run("returns not found when product does not exist", func(t *testing.T) {
		mockRepo := &mocks.ProductsRepository{
			GetByCodeFunc: func(ctx context.Context, code string) (models.Product, error) {
				return models.Product{}, models.ProductNotFoundError
			},
		}

		handler := NewCatalogHandler(mockRepo)
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/catalog/PRD100", nil)
		req.SetPathValue("code", "PRD100")

		handler.GetProductDetail(recorder, req)

		assert.Equal(t, http.StatusNotFound, recorder.Code, "Expected status code 404 Not Found")
	})

	t.Run("returns error when repository fails on GetByCode", func(t *testing.T) {
		mockRepo := &mocks.ProductsRepository{
			GetByCodeFunc: func(ctx context.Context, code string) (models.Product, error) {
				return models.Product{}, errors.New("database error")
			},
		}

		handler := NewCatalogHandler(mockRepo)
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/catalog/PROD001", nil)
		req.SetPathValue("code", "PROD001")

		handler.GetProductDetail(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code, "Expected status code 500 Internal Server Error")
	})
}

func TestGetPaginationParams(t *testing.T) {
	t.Run("uses default values when no parameters provided", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/catalog", nil)

		params, err := getPaginationParams(req)

		assert.NoError(t, err, "Should not return an error")
		assert.Equal(t, 0, params.Offset, "Expected default offset to be 0")
		assert.Equal(t, 10, params.Limit, "Expected default limit to be 10")
	})

	t.Run("sucessfully gets custom offset and limit values", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/catalog?offset=5&limit=25", nil)

		params, err := getPaginationParams(req)

		assert.NoError(t, err, "Should not return an error")
		assert.Equal(t, 5, params.Offset, "Expected offset to be 5")
		assert.Equal(t, 25, params.Limit, "Expected limit to be 25")
	})

	t.Run("returns error for invalid offset", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/catalog?offset=invalid", nil)

		_, err := getPaginationParams(req)

		assert.Error(t, err, "Should return an error")
		assert.Equal(t, "Invalid offset param", err.Error(), "Expected error message about invalid offset")
	})

	t.Run("returns error for invalid limit", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/catalog?limit=invalid", nil)

		_, err := getPaginationParams(req)

		assert.Error(t, err, "Should return an error")
		assert.Equal(t, "Invalid limit param", err.Error(), "Expected error message about invalid limit")
	})

	t.Run("returns error when limit exceeds maximum", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/catalog?limit=150", nil)

		_, err := getPaginationParams(req)

		assert.Error(t, err, "Should return an error")
		assert.Contains(t, err.Error(), "Limit must be an integer between", "Expected error about limit validation")
	})

	t.Run("returns error when limit is below minimum", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/catalog?limit=0", nil)

		_, err := getPaginationParams(req)

		assert.Error(t, err, "Should return an error")
		assert.Contains(t, err.Error(), "Limit must be an integer between", "Expected error about limit validation")
	})

	t.Run("returns error when offset is negative", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/catalog?offset=-1", nil)

		_, err := getPaginationParams(req)

		assert.Error(t, err, "Should return an error")
		assert.Equal(t, "Offset must be a non-negative integer", err.Error(), "Expected error about negative offset")
	})
}

func TestGetFilterParams(t *testing.T) {
	t.Run("returns empty filters when no parameters provided", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/catalog", nil)

		params, err := getFilterParams(req)

		assert.NoError(t, err, "Should not return an error")
		assert.Equal(t, "", params.CategoryCode, "Expected empty category code")
		assert.Equal(t, 0.0, params.PriceLessThan, "Expected default price_less_than to be 0")
	})

	t.Run("successfully gets custom filter values", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/catalog?category_code=CLOTHING&price_less_than=50.99", nil)

		params, err := getFilterParams(req)

		assert.NoError(t, err, "Should not return an error")
		assert.Equal(t, "CLOTHING", params.CategoryCode, "Expected category code to be CLOTHING")
		assert.Equal(t, 50.99, params.PriceLessThan, "Expected price_less_than to be 50.99")
	})

	t.Run("trims whitespace from category code", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/catalog?category_code=%20%20SHOES%20%20", nil)

		params, err := getFilterParams(req)

		assert.NoError(t, err, "Should not return an error")
		assert.Equal(t, "SHOES", params.CategoryCode, "Expected category code to be trimmed")
	})

	t.Run("returns error for invalid price_less_than", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/catalog?price_less_than=invalid", nil)

		_, err := getFilterParams(req)

		assert.Error(t, err, "Should return an error")
		assert.Equal(t, "Invalid price_less_than param", err.Error(), "Expected error message about invalid price")
	})
}

func TestValidatePaginationParams(t *testing.T) {
	t.Run("accepts valid pagination params", func(t *testing.T) {
		params := models.PaginationParams{
			Offset: 0,
			Limit:  10,
		}

		err := validatePaginationParams(params)

		assert.NoError(t, err, "Should not return an error for valid params")
	})

	t.Run("rejects negative offset", func(t *testing.T) {
		params := models.PaginationParams{
			Offset: -1,
			Limit:  10,
		}

		err := validatePaginationParams(params)

		assert.Error(t, err, "Should return an error for negative offset")
		assert.Equal(t, "Offset must be a non-negative integer", err.Error(), "Expected error message about negative offset")
	})

	t.Run("rejects limit below minimum", func(t *testing.T) {
		params := models.PaginationParams{
			Offset: 0,
			Limit:  0,
		}

		err := validatePaginationParams(params)

		assert.Error(t, err, "Should return an error for limit below minimum")
		assert.Contains(t, err.Error(), "Limit must be an integer between", "Expected error about limit range")
	})

	t.Run("rejects limit above maximum", func(t *testing.T) {
		params := models.PaginationParams{
			Offset: 0,
			Limit:  maxLimit + 1,
		}

		err := validatePaginationParams(params)

		assert.Error(t, err, "Should return an error for limit above maximum")
		assert.Contains(t, err.Error(), "Limit must be an integer between", "Expected error about limit range")
	})
}
