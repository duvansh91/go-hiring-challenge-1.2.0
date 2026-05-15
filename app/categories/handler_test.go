package categories

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mytheresa/go-hiring-challenge/app/dto"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/mytheresa/go-hiring-challenge/models/mocks"
)

func TestNewCategoriesHandler(t *testing.T) {
	t.Run("creates handler with repository", func(t *testing.T) {
		mockRepo := &mocks.CategoriesRepository{}
		handler := NewCategoriesHandler(mockRepo)

		assert.NotNil(t, handler, "Handler should not be nil")
		assert.Equal(t, mockRepo, handler.repo, "Repository should be set correctly")
	})
}

func TestGetAll(t *testing.T) {
	t.Run("successfully gets all categories", func(t *testing.T) {
		mockRepo := &mocks.CategoriesRepository{
			GetAllFunc: func(ctx context.Context) ([]models.Category, error) {
				return []models.Category{
					{
						ID:   1,
						Code: "CLOTHING",
						Name: "Clothing",
					},
					{
						ID:   2,
						Code: "SHOES",
						Name: "Shoes",
					},
					{
						ID:   3,
						Code: "ACCESSORIES",
						Name: "Accessories",
					},
				}, nil
			},
		}

		handler := NewCategoriesHandler(mockRepo)
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/categories", nil)

		handler.GetAll(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200 OK")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Expected Content-Type to be application/json")

		var response []dto.Category
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.NoError(t, err, "Should be valid JSON")
		assert.Len(t, response, 3, "Expected 3 categories in response")

		assert.Equal(t, "CLOTHING", response[0].Code, "Expected first category code to be CLOTHING")
		assert.Equal(t, "Clothing", response[0].Name, "Expected first category name to be Clothing")
		assert.Equal(t, "SHOES", response[1].Code, "Expected second category code to be SHOES")
		assert.Equal(t, "Shoes", response[1].Name, "Expected second category name to be Shoes")
		assert.Equal(t, "ACCESSORIES", response[2].Code, "Expected third category code to be ACCESSORIES")
		assert.Equal(t, "Accessories", response[2].Name, "Expected third category name to be Accessories")
	})

	t.Run("returns empty list when no categories exist", func(t *testing.T) {
		mockRepo := &mocks.CategoriesRepository{
			GetAllFunc: func(ctx context.Context) ([]models.Category, error) {
				return []models.Category{}, nil
			},
		}

		handler := NewCategoriesHandler(mockRepo)
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/categories", nil)

		handler.GetAll(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200 OK")

		var response []dto.Category
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.NoError(t, err, "Should be valid JSON")
		assert.Len(t, response, 0, "Expected empty list")
	})

	t.Run("returns error when repository fails on GetAll", func(t *testing.T) {
		mockRepo := &mocks.CategoriesRepository{
			GetAllFunc: func(ctx context.Context) ([]models.Category, error) {
				return nil, errors.New("database error")
			},
		}

		handler := NewCategoriesHandler(mockRepo)
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/categories", nil)

		handler.GetAll(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code, "Expected status code 500 Internal Server Error")
	})

	t.Run("successfully maps categories to DTOs", func(t *testing.T) {
		mockRepo := &mocks.CategoriesRepository{
			GetAllFunc: func(ctx context.Context) ([]models.Category, error) {
				return []models.Category{
					{
						ID:   1,
						Code: "CLOTHING",
						Name: "Clothing",
					},
				}, nil
			},
		}

		handler := NewCategoriesHandler(mockRepo)
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/categories", nil)

		handler.GetAll(recorder, req)

		var response []dto.Category
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(t, "CLOTHING", response[0].Code, "Expected category code to match")
		assert.Equal(t, "Clothing", response[0].Name, "Expected category name to match")
	})
}

func TestCreate(t *testing.T) {
	t.Run("successfully creates category", func(t *testing.T) {
		mockRepo := &mocks.CategoriesRepository{
			CreateFunc: func(ctx context.Context, category *models.Category) error {
				category.ID = 1
				return nil
			},
		}

		handler := NewCategoriesHandler(mockRepo)
		recorder := httptest.NewRecorder()

		categoryDTO := dto.Category{
			Code: "CLOTHING",
			Name: "Clothing",
		}
		body, err := json.Marshal(categoryDTO)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest("POST", "/categories", bytes.NewReader(body))

		handler.Create(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200 OK")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Expected Content-Type to be application/json")

		var response dto.Category
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, "CLOTHING", response.Code, "Expected category code in response")
		assert.Equal(t, "Clothing", response.Name, "Expected category name in response")
	})

	t.Run("returns error when request body is invalid", func(t *testing.T) {
		mockRepo := &mocks.CategoriesRepository{}
		handler := NewCategoriesHandler(mockRepo)
		recorder := httptest.NewRecorder()

		invalidBody := bytes.NewReader([]byte("invalid json"))
		req := httptest.NewRequest("POST", "/categories", invalidBody)

		handler.Create(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code, "Expected status code 400 Bad Request")
	})

	t.Run("returns error when repository fails on Create", func(t *testing.T) {
		mockRepo := &mocks.CategoriesRepository{
			CreateFunc: func(ctx context.Context, category *models.Category) error {
				return errors.New("database error")
			},
		}

		handler := NewCategoriesHandler(mockRepo)
		recorder := httptest.NewRecorder()

		categoryDTO := dto.Category{
			Code: "CLOTHING",
			Name: "Clothing",
		}
		body, err := json.Marshal(categoryDTO)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest("POST", "/categories", bytes.NewReader(body))

		handler.Create(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code, "Expected status code 500 Internal Server Error")
	})

	t.Run("successfully maps DTO to model", func(t *testing.T) {
		var expectedCategory *models.Category
		mockRepo := &mocks.CategoriesRepository{
			CreateFunc: func(ctx context.Context, category *models.Category) error {
				expectedCategory = category
				return nil
			},
		}

		handler := NewCategoriesHandler(mockRepo)
		recorder := httptest.NewRecorder()

		categoryDTO := dto.Category{
			Code: "SHOES",
			Name: "Shoes",
		}
		body, err := json.Marshal(categoryDTO)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest("POST", "/categories", bytes.NewReader(body))

		handler.Create(recorder, req)

		assert.NotNil(t, expectedCategory, "Category should be captured by repository")
		assert.Equal(t, "SHOES", expectedCategory.Code, "Expected category code to match")
		assert.Equal(t, "Shoes", expectedCategory.Name, "Expected category name to match")
	})

	t.Run("returns error for empty request body", func(t *testing.T) {
		mockRepo := &mocks.CategoriesRepository{}
		handler := NewCategoriesHandler(mockRepo)
		recorder := httptest.NewRecorder()

		req := httptest.NewRequest("POST", "/categories", bytes.NewReader([]byte("")))

		handler.Create(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code, "Expected status code 400 Bad Request")
	})

	t.Run("successfully creates category", func(t *testing.T) {
		var expectedCategory *models.Category
		mockRepo := &mocks.CategoriesRepository{
			CreateFunc: func(ctx context.Context, category *models.Category) error {
				expectedCategory = category
				return nil
			},
		}

		handler := NewCategoriesHandler(mockRepo)
		recorder := httptest.NewRecorder()

		categoryDTO := dto.Category{
			Code: "ACCESSORIES",
			Name: "Accessories",
		}
		body, err := json.Marshal(categoryDTO)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest("POST", "/categories", bytes.NewReader(body))

		handler.Create(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code, "Expected successful creation")
		assert.NotNil(t, expectedCategory, "Category should be created")
		assert.Equal(t, "ACCESSORIES", expectedCategory.Code, "Expected category code")
		assert.Equal(t, "Accessories", expectedCategory.Name, "Expected category name")
	})

	t.Run("returns error when missing code field in request", func(t *testing.T) {
		mockRepo := &mocks.CategoriesRepository{}
		handler := NewCategoriesHandler(mockRepo)
		recorder := httptest.NewRecorder()

		categoryDTO := map[string]string{
			"name": "Clothing",
		}
		body, err := json.Marshal(categoryDTO)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest("POST", "/categories", bytes.NewReader(body))

		handler.Create(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code, "Should be bad request when missing code field")

		assert.Equal(t, http.StatusBadRequest, recorder.Code, "Expected status code 400 Bad Request")
		assert.JSONEq(t, `{"error":"Category code is required"}`, recorder.Body.String(), "Expected error message in response body")
	})

	t.Run("returns error when missing name field in request", func(t *testing.T) {
		mockRepo := &mocks.CategoriesRepository{}
		handler := NewCategoriesHandler(mockRepo)
		recorder := httptest.NewRecorder()

		categoryDTO := map[string]string{
			"code": "CLOTHING",
		}
		body, err := json.Marshal(categoryDTO)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest("POST", "/categories", bytes.NewReader(body))

		handler.Create(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code, "Should be bad request when missing code field")

		assert.Equal(t, http.StatusBadRequest, recorder.Code, "Expected status code 400 Bad Request")
		assert.JSONEq(t, `{"error":"Category name is required"}`, recorder.Body.String(), "Expected error message in response body")
	})
}
