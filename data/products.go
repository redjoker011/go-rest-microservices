package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

// Product defines the structure of an API product
// swagger:model
type Product struct {
	// The id for this product
	//
	// require: true
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// Product Struct Validator
func (p *Product) Validate() error {
	validate := validator.New()
	// Register custom validator
	validate.RegisterValidation("sku", validateSKU)

	// specify validator and validate struct
	return validate.Struct(p)
}

// Custom validator
func validateSKU(fl validator.FieldLevel) bool {
	// sku is of format abc-absd-ddsf
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]`)
	matches := re.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}

	return true
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
