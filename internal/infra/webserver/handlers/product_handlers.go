package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/felipedias-dev/fullcycle-go-expert-basic-api/internal/dto"
	"github.com/felipedias-dev/fullcycle-go-expert-basic-api/internal/entity"
	"github.com/felipedias-dev/fullcycle-go-expert-basic-api/internal/infra/database"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product, err := entity.NewProduct(input.Name, input.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.ProductDB.Create(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
