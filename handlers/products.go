package handlers

import (
	protos "github.com/redjoker011/online-cafe/currency/protos/currency"
	"log"
)

type Products struct {
	l  *log.Logger
	cc protos.CurrencyClient
}

func NewProducts(l *log.Logger, cc protos.CurrencyClient) *Products {
	return &Products{l, cc}
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
