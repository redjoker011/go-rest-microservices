package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/redjoker011/online-cafe/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// Create empty struct which act as request context key identifier
type KeyProduct struct{}

// Generic Error response from server
type GenericError struct {
	Message string `json:"message"`
}

// Validation Error
type ValidationError struct {
	Nessage []string `json:"messages"`
}

// swagger:route PUT /products{id} products updateProduct
// responses:
// 200: noContent

// UpdateProduct update a product from the data store
func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	// Fetch parameters from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(rw, "Missing id parameter", http.StatusBadRequest)
	}

	p.l.Println("Handle PUT Product", id)

	// Fetch Product Object from request context
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	err = data.UpdateProduct(id, &prod)

	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

// Request Middleware
func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// Initialize empty Product Object
		prod := data.Product{}
		// Convert request parameters into product object
		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to marshal json", http.StatusBadRequest)
			return
		}

		err = prod.Validate()

		if err != nil {
			http.Error(
				rw,
				fmt.Sprintf("Error validating product: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		// Store Product Object into request context using KeyProduct as identifier
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		// Invoke Request handler
		next.ServeHTTP(rw, req)
	})
}
