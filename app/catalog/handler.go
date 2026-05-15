package catalog

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/app/dto"
	"github.com/mytheresa/go-hiring-challenge/models"
)

const (
	offsetDefatult = "0"
	limitDefault   = "10"
	maxLimit       = 100
	minLimit       = 1
)

// CatalogHandler handles HTTP requests for the catalog.
type CatalogHandler struct {
	repo ProductsRepository
}

// ProductsRepository defines the contract for accessing product data.
type ProductsRepository interface {
	GetAll(ctx context.Context, filter models.FilterParams, pagination models.PaginationParams) (models.PaginatedResult, error)
	GetByCode(ctx context.Context, code string) (models.Product, error)
}

// NewCatalogHandler creates a new CatalogHandler with the given repository.
func NewCatalogHandler(r ProductsRepository) *CatalogHandler {
	return &CatalogHandler{
		repo: r,
	}
}

// Get retrieves all products and returns them as JSON.
func (h *CatalogHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	paginationParams, err := getPaginationParams(r)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())

		return
	}

	filterParams, err := getFilterParams(r)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())

		return
	}

	result, err := h.repo.GetAll(ctx, filterParams, paginationParams)
	if err != nil {
		slog.Error("error getting products", "error", err)
		api.ErrorResponse(w, http.StatusInternalServerError, "Error getting products")

		return
	}

	productsDTO := make([]dto.Product, len(result.Products))
	for i, p := range result.Products {
		productsDTO[i] = dto.MapProductToDTO(p, false)
	}

	api.OKResponse(w, dto.Cataloge{
		Products:      productsDTO,
		TotalProducts: result.Total,
		TotalPages:    result.TotalPages,
	})
}

// GetProductDetail retrieves a single product and returns it as JSON.
func (h *CatalogHandler) GetProductDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := strings.TrimSpace(r.PathValue("code"))
	product, err := h.repo.GetByCode(ctx, code)
	if err != nil {
		if err == models.ProductNotFoundError {
			api.ErrorResponse(w, http.StatusNotFound, "Product not found")

			return
		}

		slog.Error("error getting product with code", "code", code, "error", err)
		api.ErrorResponse(w, http.StatusInternalServerError, "Error getting product")

		return
	}

	productDTO := dto.MapProductToDTO(product, true)

	api.OKResponse(w, productDTO)
}

func getPaginationParams(r *http.Request) (models.PaginationParams, error) {
	query := r.URL.Query()

	offsetSTR := query.Get("offset")
	if offsetSTR == "" {
		offsetSTR = offsetDefatult
	}

	offset, err := strconv.Atoi(offsetSTR)
	if err != nil {
		slog.Error("error parsing offset", "error", err)

		return models.PaginationParams{}, fmt.Errorf("Invalid offset param")
	}

	limitSTR := query.Get("limit")
	if limitSTR == "" {
		limitSTR = limitDefault
	}

	limit, err := strconv.Atoi(limitSTR)
	if err != nil {
		slog.Error("error parsing limit", "error", err)

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

	priceLessThanSTR := query.Get("price_less_than")
	if priceLessThanSTR == "" {
		priceLessThanSTR = "0"
	}

	priceLessThan, err := strconv.ParseFloat(priceLessThanSTR, 64)
	if err != nil {
		slog.Error("error parsing price_less_than", "error", err)
		return models.FilterParams{}, fmt.Errorf("Invalid price_less_than param")
	}

	categoryCodeStr := query.Get("category_code")

	return models.FilterParams{
		CategoryCode:  strings.TrimSpace(categoryCodeStr),
		PriceLessThan: priceLessThan,
	}, nil
}
