package categories

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/app/dto"
	"github.com/mytheresa/go-hiring-challenge/models"
)

// CategoriesHandler handles HTTP requests for the Categories.
type CategoriesHandler struct {
	repo CategoriesRepository
}

// CategoriesRepository defines the methods of the categories repository.
type CategoriesRepository interface {
	GetAll(ctx context.Context) ([]models.Category, error)
	Create(ctx context.Context, category *models.Category) error
}

// NewCategoriesHandler creates a new CatalogHandler with the given repository.
func NewCategoriesHandler(r CategoriesRepository) *CategoriesHandler {
	return &CategoriesHandler{
		repo: r,
	}
}

// GetAll retrieves all categories and returns them as JSON.
func (h *CategoriesHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categories, err := h.repo.GetAll(ctx)
	if err != nil {
		log.Printf("error getting categories: %v", err)
		api.ErrorResponse(w, http.StatusInternalServerError, "error getting categories")

		return
	}

	var categoriesDTO []dto.Category
	for _, category := range categories {
		categoriesDTO = append(categoriesDTO, dto.Category{
			Code: category.Code,
			Name: category.Name,
		})
	}

	api.OKResponse(w, categoriesDTO)
}

// Create creates a new category and stores it in db.
func (h *CategoriesHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var categoryDTO dto.Category
	err := json.NewDecoder(r.Body).Decode(&categoryDTO)
	if err != nil {
		log.Printf("invalid request body: %v", err)
		api.ErrorResponse(w, http.StatusBadRequest, "invalid request body")

		return
	}

	category := models.Category{
		Code: categoryDTO.Code,
		Name: categoryDTO.Name,
	}

	err = validateCategoryFields(&category)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())

		return
	}

	err = h.repo.Create(ctx, &category)
	if err != nil {
		log.Printf("error creating category: %v", err)
		api.ErrorResponse(w, http.StatusInternalServerError, "error creating category")

		return
	}

	api.OKResponse(w, dto.Category{
		Code: category.Code,
		Name: category.Name,
	})
}

func validateCategoryFields(category *models.Category) error {
	if category.Code == "" {
		return errors.New("Category code is required")
	}

	if category.Name == "" {
		return errors.New("Category name is required")
	}

	return nil
}
