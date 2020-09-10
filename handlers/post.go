package handlers

import (
	"net/http"

	"github.com/redjoker011/online-cafe/data"
)

// swagger:route POST /products products createProduct
// Create a new Product
// response:
// 200: productsResponse
// 422: errorValidation
// 501: errorResponse

// Create handles POST request to add new products
func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")
	// Fetch Product Object from request context
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	data.AddProduct(&prod)
}
