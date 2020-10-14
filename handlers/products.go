package handlers

import (
	"github.com/hashicorp/go-hclog"
	"github.com/redjoker011/online-cafe/data"
)

type Products struct {
	l         hclog.Logger
	productDB *data.ProductsDB
}

func NewProducts(l hclog.Logger, pdb *data.ProductsDB) *Products {
	return &Products{l, pdb}
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
