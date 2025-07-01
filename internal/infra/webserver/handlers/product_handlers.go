package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Mazzael/go-api/internal/dto"
	"github.com/Mazzael/go-api/internal/entity"
	"github.com/Mazzael/go-api/internal/infra/database"
)

type ProductHandler struct {
	GormProductRepository database.ProductRepository
}

func NewProductHandler(repo database.ProductRepository) *ProductHandler {
	return &ProductHandler{
		GormProductRepository: repo,
	}
}

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
