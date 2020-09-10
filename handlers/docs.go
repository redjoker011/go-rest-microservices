// Package Classification of Products API
//
// Documentation for Product API
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package handlers

import "github.com/redjoker011/online-cafe/data"

// A List of Products
// swagger:response productsResponse
type productsResponseWrapper struct {
	// All Products in the system
	// in: body
	Body []data.Product
}

// swagger:parameters updateProduct
type productIDParameterWrapper struct {
	// The ID of the product to delete from the database
	// in: path
	// required: true
	ID int `json:"id"`
}

// No content is returned by this API endpoint
// swagger:response noContent
type productsNoContent struct {
}