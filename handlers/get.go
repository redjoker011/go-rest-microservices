package handlers

import (
	"net/http"

	"github.com/redjoker011/online-cafe/data"
)

// swagger:route GET /products products listProducts
// Returns List of Products
// responses:
// 200:productsResponse

// Respond based on HTTP Method
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	// Get products
	lp := data.GetProducts()
	rw.Header().Add("Content-Type", "application/json")
	// Convert struct into JSON and write into ResponseWriter
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}
