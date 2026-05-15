package categories

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/dto"
	"github.com/mytheresa/go-hiring-challenge/models"
)

// CategoriesHandler handles HTTP requests for the product Categories.
type CategoriesHandler struct {
	repo CategoriesRepository
}

type CategoriesRepository interface {
	GetAllCategories() ([]models.Category, error)
	CreateCategory(category *models.Category) error
}

// NewCategoriesHandler creates a new CatalogHandler with the given repository.
func NewCategoriesHandler(r CategoriesRepository) *CategoriesHandler {
	return &CategoriesHandler{
		repo: r,
	}
}

func (h *CategoriesHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	categories, err := h.repo.GetAllCategories()
	if err != nil {
		slog.Error("error getting categories", "error", err)
		http.Error(w, "error getting categories", http.StatusInternalServerError)
		return
	}

	var categoriesDTO []dto.Category
	for _, category := range categories {
		categoriesDTO = append(categoriesDTO, dto.Category{
			Code: category.Code,
			Name: category.Name,
		})
	}

	json.NewEncoder(w).Encode(categoriesDTO)
}

func (h *CategoriesHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var categoryDTO dto.Category
	err := json.NewDecoder(r.Body).Decode(&categoryDTO)
	if err != nil {
		slog.Error("invalid request body", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	category := models.Category{
		Code: categoryDTO.Code,
		Name: categoryDTO.Name,
	}

	err = h.repo.CreateCategory(&category)
	if err != nil {
		slog.Error("error creating category", "error", err)

		if err == models.CategoryAlreadyExistsError {
			http.Error(w, "category already exists", http.StatusConflict)
			return
		}

		http.Error(w, "error creating category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.Category{
		Code: category.Code,
		Name: category.Name,
	})
}
