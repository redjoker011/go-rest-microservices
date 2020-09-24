package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	protos "github.com/redjoker011/online-cafe/currency/protos/currency"
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

func (p *Products) GetProduct(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	rr := &protos.RateRequest{
		Base:        protos.Currencies(protos.Currencies_value["EUR"]),
		Destination: protos.Currencies(protos.Currencies_value["USD"]),
	}
	curr, err := p.cc.GetRate(context.Background(), rr)

	if err != nil {
		p.l.Println("[Error] error getting new rate", err)
	}

	p.l.Println("Exchange Rate", curr.Rate)
	prod := data.Product{
		ID:          3,
		Name:        "Test",
		Description: "Taste my test",
		SKU:         "TEst",
		Price:       1 * curr.Rate,
	}

	e := json.NewEncoder(rw)
	err = e.Encode(prod)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}
