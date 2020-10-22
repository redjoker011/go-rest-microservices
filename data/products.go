package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
	"github.com/hashicorp/go-hclog"
	protos "github.com/redjoker011/online-cafe/currency/protos/currency"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	Price       float64 `json:"price" validate:"gt=0"`
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

type ProductsDB struct {
	currency protos.CurrencyClient
	log      hclog.Logger
}

func NewProductsDB(c protos.CurrencyClient, l hclog.Logger) *ProductsDB {
	return &ProductsDB{c, l}
}

// Return Products Data
func (p *ProductsDB) GetProducts(currency string) (Products, error) {
	if currency == "" {
		return productList, nil
	}

	r, err := p.getRate(currency)

	if err != nil {
		p.log.Error("Unable to get rate", "currency", currency, "error", err)
		return nil, err
	}

	products := Products{}
	for _, p := range productList {
		np := *p
		np.Price = np.Price * r
		products = append(products, &np)
	}

	return products, nil
}

func (p *ProductsDB) getRate(dest string) (float64, error) {
	rr := &protos.RateRequest{
		Base:        protos.Currencies(protos.Currencies_value["USD"]),
		Destination: protos.Currencies(protos.Currencies_value[dest]),
	}
	curr, err := p.currency.GetRate(context.Background(), rr)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			md := s.Details()[0].(*protos.RateRequest)

			if s.Code() == codes.InvalidArgument {
				return -1, fmt.Errorf("Unable to get rate from currency server, destination and base cannot be the same, base: %s, dest: %s", md.Base.String(), md.Destination.String())
			}
			return -1, fmt.Errorf("Unable to get rate from currency server, base: %s, dest: %s", md.Base.String(), md.Destination.String())
		}

		return -1, err
	}

	return curr.Rate, err
}

// Add new product into products array
func (p *ProductsDB) AddProduct(prod *Product) {
	prod.ID = getNextID()
	productList = append(productList, prod)
}

// Get Last Product ID
func getNextID() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

func (p *ProductsDB) UpdateProduct(id int, product *Product) error {
	prod, pos, err := findProduct(id)

	fmt.Printf("Old Product: %#v \n", prod)

	if err != nil {
		return err
	}

	// Set product struct ID
	product.ID = id
	p.log.Debug("Updated Product: %#v \n", product)
	// Get product by index and update value
	productList[pos] = product
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
