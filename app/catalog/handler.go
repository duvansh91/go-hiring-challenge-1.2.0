package catalog

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/app/dto"
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
	Products      []dto.Product `json:"products"`
	TotalProducts int64         `json:"total_products"`
	TotalPages    int           `json:"total_pages"`
}

// CatalogHandler handles HTTP requests for the product catalog.
type CatalogHandler struct {
	repo ProductsRepository
}

// ProductsRepository defines the contract for accessing product data.
type ProductsRepository interface {
	GetAllProducts(filterParams models.FilterParams, paginationParams models.PaginationParams) (models.PaginatedResult, error)
	GetProductByCode(code string) (models.Product, error)
}

// NewCatalogHandler creates a new CatalogHandler with the given repository.
func NewCatalogHandler(r ProductsRepository) *CatalogHandler {
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
		slog.Error("error getting products", "error", err)
		http.Error(w, "Error getting products", http.StatusInternalServerError)
		return
	}

	productsDTO := make([]dto.Product, len(result.Data))
	for i, p := range result.Data {
		productsDTO[i] = dto.MapProductToDTO(p, false)
	}

	err = json.NewEncoder(w).Encode(Response{
		Products:      productsDTO,
		TotalProducts: result.Total,
		TotalPages:    result.TotalPages,
	})
	if err != nil {
		slog.Error("error parsing products response", "error", err)
		http.Error(w, "error parsing products response", http.StatusInternalServerError)
		return
	}
}

// GetProductDetail retrieves a single product and returns it as JSON.
func (h *CatalogHandler) GetProductDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	productCode := r.PathValue("code")
	product, err := h.repo.GetProductByCode(productCode)
	if err != nil {
		if err == models.ProductNotFoundError {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		slog.Error("error getting product with code", "code", productCode, "error", err)
		http.Error(w, "Error getting product", http.StatusInternalServerError)
		return
	}

	productDTO := dto.MapProductToDTO(product, true)

	err = json.NewEncoder(w).Encode(productDTO)
	if err != nil {
		slog.Error("error parsing product response", "error", err)
		http.Error(w, "error parsing product response", http.StatusInternalServerError)
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
		slog.Error("error parsing offset:", "error", err)
		return models.PaginationParams{}, fmt.Errorf("Invalid offset param")
	}

	limitStr := query.Get("limit")
	if limitStr == "" {
		limitStr = limitDefault
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		slog.Error("error parsing limit:", "error", err)
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
