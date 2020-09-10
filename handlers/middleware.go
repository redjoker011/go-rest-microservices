package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/redjoker011/online-cafe/data"
)

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
