package handlers

import (
	// "context"
	// protos "github.com/redjoker011/online-cafe/currency/protos/currency"

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
	lp, err := p.productDB.GetProducts("")

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{err.Error()}, rw)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	// Convert struct into JSON and write into ResponseWriter
	err = data.ToJSON(lp, rw)
	if err != nil {
		p.l.Error("Unable to serialize product", "error", err)
	}
}

func (p *Products) GetProduct(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	// rr := &protos.RateRequest{
	// 	Base:        protos.Currencies(protos.Currencies_value["EUR"]),
	// 	Destination: protos.Currencies(protos.Currencies_value["USD"]),
	// }
	// curr, err := p.cc.GetRate(context.Background(), rr)

	// if err != nil {
	// 	p.l.Error("[Error] error getting new rate", err)
	// }

	// p.l.Info("Exchange Rate", curr.Rate)

	prod := data.Product{
		ID:          3,
		Name:        "Test",
		Description: "Taste my test",
		SKU:         "TEst",
		Price:       1,
		// Price:       1 * curr.Rate,
	}

	err := data.ToJSON(prod, rw)
	if err != nil {
		p.l.Error("Unable to serialize product", err)
	}
}
