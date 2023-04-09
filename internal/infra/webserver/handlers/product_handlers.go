package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/felipedias-dev/fullcycle-go-expert-basic-api/internal/dto"
	"github.com/felipedias-dev/fullcycle-go-expert-basic-api/internal/entity"
	"github.com/felipedias-dev/fullcycle-go-expert-basic-api/internal/infra/database"
	pkgEntity "github.com/felipedias-dev/fullcycle-go-expert-basic-api/pkg/entity"
	"github.com/go-chi/chi"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

// Create product godoc
// @Summary 		Create product
// @Description Create product
// @Tags 				products
// @Accept  		json
// @Produce  		json
// @Param 			request		body			dto.CreateProductInput true "product request"
// @Success 		201
// @Failure 		400 			{object}	Error
// @Failure 		500 			{object}	Error
// @Router 			/products [post]
// @Security 		ApiKeyAuth
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	product, err := entity.NewProduct(input.Name, input.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	err = h.ProductDB.Create(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Get products godoc
// @Summary 		List products
// @Description List products
// @Tags 				products
// @Accept  		json
// @Produce  		json
// @Param 			page		query			string	false	"page number"
// @Param 			limit		query			string	false	"limit"
// @Param 			sort		query			string	false	"sort"
// @Success 		200 		{object}	[]entity.Product
// @Failure 		500 			{object}	Error
// @Router 			/products [get]
// @Security 		ApiKeyAuth
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	sort := r.URL.Query().Get("sort")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 0
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10
	}

	products, err := h.ProductDB.FindAll(pageInt, limitInt, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

// Get product godoc
// @Summary 		Get product
// @Description Get product
// @Tags 				products
// @Accept  		json
// @Produce  		json
// @Param 			id		path			string	true	"product ID" Format(uuid)
// @Success 		200 	{object}	entity.Product
// @Failure 		400 	{object}	Error
// @Failure 		404 	{object}	Error
// @Router 			/products/{id} [get]
// @Security 		ApiKeyAuth
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: "ID is required"}
		json.NewEncoder(w).Encode(error)
		return
	}
	product, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// Update product godoc
// @Summary 		Update product
// @Description Update product
// @Tags 				products
// @Accept  		json
// @Produce  		json
// @Param 			id		path			string	true	"product ID" Format(uuid)
// @Param 			request		body			dto.UpdateProductInput true "product request"
// @Success 		200
// @Failure 		400 	{object}	Error
// @Failure 		404 	{object}	Error
// @Failure 		500 	{object}	Error
// @Router 			/products/{id} [put]
// @Security 		ApiKeyAuth
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: "ID is required"}
		json.NewEncoder(w).Encode(error)
		return
	}
	var product entity.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	product.ID, err = pkgEntity.ParseID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	productFound, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	product.CreatedAt = productFound.CreatedAt
	err = h.ProductDB.Update(&product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
}

// Delete product godoc
// @Summary 		Delete product
// @Description Delete product
// @Tags 				products
// @Accept  		json
// @Produce  		json
// @Param 			id		path			string	true	"product ID" Format(uuid)
// @Success 		200
// @Failure 		400 	{object}	Error
// @Failure 		404 	{object}	Error
// @Failure 		500 	{object}	Error
// @Router 			/products/{id} [delete]
// @Security 		ApiKeyAuth
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: "ID is required"}
		json.NewEncoder(w).Encode(error)
		return
	}
	_, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	err = h.ProductDB.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
}
