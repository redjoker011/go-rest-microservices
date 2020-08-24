package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// Convert JSON into Struct
func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

// Create embedded type for Products
type Products []*Product

// Create JSON encoder that writes into io writer
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Esspresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "abc324",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}

// Return Products Data
func GetProducts() Products {
	return productList
}

// Add new product into products array
func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

// Get Last Product ID
func getNextID() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

func UpdateProduct(id int, p *Product) error {
	prod, pos, err := findProduct(id)

	fmt.Printf("Old Product: %#v \n", prod)

	if err != nil {
		return err
	}

	// Set product struct ID
	p.ID = id
	fmt.Printf("Updated Product: %#v \n", p)
	// Get product by index and update value
	productList[pos] = p
	return nil
}

// Initialize Error object
var ErrProductNotFound = fmt.Errorf("Product not found")

// Find product by ID
func findProduct(id int) (*Product, int, error) {
	// Iterate over product list and match product by ID
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, 0, ErrProductNotFound
}
