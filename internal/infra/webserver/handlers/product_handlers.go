package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Mazzael/go-api/internal/dto"
	"github.com/Mazzael/go-api/internal/entity"
	"github.com/Mazzael/go-api/internal/infra/database"
	entityPkg "github.com/Mazzael/go-api/pkg/entity"
	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	GormProductRepository database.ProductRepository
}

func NewProductHandler(repo database.ProductRepository) *ProductHandler {
	return &ProductHandler{
		GormProductRepository: repo,
	}
}

// Create Product godoc
// @Summary      Create product
// @Description  Create products
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        request     body      dto.CreateProductInput  true  "product request"
// @Success      201
// @Failure      500         {object}  Error
// @Router       /products [post]
// @Security ApiKeyAuth
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.GormProductRepository.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetProduct godoc
// @Summary      Get a product
// @Description  Get a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "product ID" Format(uuid)
// @Success      200  {object}  entity.Product
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /products/{id} [get]
// @Security ApiKeyAuth
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := Error{Message: "Product ID is required"}
		json.NewEncoder(w).Encode(err)
		return
	}

	product, err := h.GormProductRepository.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// List Products godoc
// @Summary      List products
// @Description  get all products
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        page      query     string  false  "page number"
// @Param        limit     query     string  false  "limit"
// @Success      200       {array}   entity.Product
// @Failure      404       {object}  Error
// @Failure      500       {object}  Error
// @Router       /products [get]
// @Security ApiKeyAuth
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "0"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 0
	}

	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "10"
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10
	}

	sort := r.URL.Query().Get("sort")

	products, err := h.GormProductRepository.FindAll(pageInt, limitInt, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	println(pageInt, limitInt, sort)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

// UpdateProduct godoc
// @Summary      Update a product
// @Description  Update a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id        	path      string                  true  "product ID" Format(uuid)
// @Param        request     body      dto.CreateProductInput  true  "product request"
// @Success      200
// @Failure      400       {object}  Error
// @Failure      404	   {object}  Error
// @Failure      500       {object}  Error
// @Router       /products/{id} [put]
// @Security ApiKeyAuth
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := Error{Message: "Product ID is required"}
		json.NewEncoder(w).Encode(err)
		return
	}

	var product entity.Product

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	product.ID, err = entityPkg.ParseID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	_, err = h.GormProductRepository.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	err = h.GormProductRepository.Update(&product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteProduct godoc
// @Summary      Delete a product
// @Description  Delete a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id        path      string                  true  "product ID" Format(uuid)
// @Success      200
// @Failure      400	   {object}  Error
// @Failure      404       {object}  Error
// @Failure      500       {object}  Error
// @Router       /products/{id} [delete]
// @Security ApiKeyAuth
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := Error{Message: "Product ID is required"}
		json.NewEncoder(w).Encode(err)
		return
	}

	_, err := h.GormProductRepository.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	err = h.GormProductRepository.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
