package handlers

import (
	"log"
	"net/http"
	"online-cafe/data"
	"regexp"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// Respond based on HTTP Method
func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	if r.Method == http.MethodPut {
		// expect product id in URI
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		// Prevent from having multiple matches
		if len(g) != 1 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		// Prevent from having multiple id parameter for nested resource
		// /1 - good
		// /1/2 - bad
		if len(g[0]) != 2 {
			p.l.Println("reg", g)
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		// Assign last match into a variable
		idString := g[0][1]
		// Convert string into integer
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		// Update Product based on ID parameter
		p.updateProduct(id, rw, r)
		return
	}

	// Raise error for unsupported HTTP Method
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	// Get products
	lp := data.GetProducts()
	// Convert struct into JSON and write into ResponseWriter
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")
	// Get parameters and convert into Product Object
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)

	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	data.AddProduct(prod)
}

func (p *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Product")
	// Get parameters and convert into Product Object

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)

	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}
	err = data.UpdateProduct(id, prod)

	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}
