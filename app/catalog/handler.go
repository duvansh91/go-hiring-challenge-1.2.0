package catalog

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/models"
)

const (
	offsetDefatult = "0"
	limitDefault   = "10"
	maxLimit       = 100
	minLimit       = 1
)

// Response represents the catalog API response containing products.
type Response struct {
	Products      []ProductDTO `json:"products"`
	TotalProducts int64        `json:"total_products"`
	TotalPages    int          `json:"total_pages"`
}

// ProductDTO represents a single product in the catalog response.
type ProductDTO struct {
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

// CatalogHandler handles HTTP requests for the product catalog.
type CatalogHandler struct {
	repo Repository
}

// Repository defines the contract for accessing product data.
type Repository interface {
	GetAllProducts(filterParams models.FilterParams, paginationParams models.PaginationParams) (models.PaginatedResult, error)
}

// NewCatalogHandler creates a new CatalogHandler with the given repository.
func NewCatalogHandler(r Repository) *CatalogHandler {
	return &CatalogHandler{
		repo: r,
	}
}

// GetCatalog retrieves all products and returns them as JSON.
func (h *CatalogHandler) GetCatalog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	paginationParams, err := getPaginationParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	filterParams, err := getFilterParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.repo.GetAllProducts(filterParams, paginationParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	productsDTO := make([]ProductDTO, len(result.Data))
	for i, p := range result.Data {
		productsDTO[i] = ProductDTO{
			Code:     p.Code,
			Price:    p.Price.InexactFloat64(),
			Category: p.Category.Name,
		}
	}

	response := Response{
		Products:      productsDTO,
		TotalProducts: result.Total,
		TotalPages:    result.TotalPages,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getPaginationParams(r *http.Request) (models.PaginationParams, error) {
	query := r.URL.Query()

	offsetStr := query.Get("offset")
	if offsetStr == "" {
		offsetStr = offsetDefatult
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		log.Printf("Error parsing offset: %v", err)
		return models.PaginationParams{}, fmt.Errorf("Invalid offset param")
	}

	limitStr := query.Get("limit")
	if limitStr == "" {
		limitStr = limitDefault
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Printf("Error parsing limit: %v", err)
		return models.PaginationParams{}, fmt.Errorf("Invalid limit param")
	}

	params := models.PaginationParams{
		Offset: offset,
		Limit:  limit,
	}

	err = validatePaginationParams(params)
	if err != nil {
		return models.PaginationParams{}, err
	}

	return params, nil
}

func validatePaginationParams(params models.PaginationParams) error {
	if params.Offset < 0 {
		return fmt.Errorf("Offset must be a non-negative integer")
	}
	if params.Limit < minLimit || params.Limit > maxLimit {
		return fmt.Errorf("Limit must be an integer between %d and %d", minLimit, maxLimit)
	}

	return nil
}

func getFilterParams(r *http.Request) (models.FilterParams, error) {
	query := r.URL.Query()

	priceLessThanStr := query.Get("price_less_than")
	if priceLessThanStr == "" {
		priceLessThanStr = "0"
	}
	priceLessThan, err := strconv.ParseFloat(priceLessThanStr, 64)
	if err != nil {
		log.Printf("Error parsing price_less_than: %v", err)
		return models.FilterParams{}, fmt.Errorf("Invalid price_less_than param")
	}

	categoryCodeStr := query.Get("category_code")

	return models.FilterParams{
		CategoryCode:  strings.TrimSpace(categoryCodeStr),
		PriceLessThan: priceLessThan,
	}, nil
}
