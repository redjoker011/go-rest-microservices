package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/redjoker011/online-cafe/data"
)

// swagger:route PUT /products{id} products updateProduct
// responses:
// 201: noContent
// 404: errorResponse
// 422: validationResponse

// UpdateProduct update a product from the data store
func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	// Fetch parameters from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(rw, "Missing id parameter", http.StatusBadRequest)
	}

	p.l.Info("Handle PUT Product", id)

	// Fetch Product Object from request context
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	err = p.productDB.UpdateProduct(id, &prod)

	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}
