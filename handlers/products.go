package handlers

import (
	"log"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
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
