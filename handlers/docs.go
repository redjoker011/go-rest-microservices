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

// swagger:parameters listProducts
type productQueryParam struct {
	// Currency used when returning the price of the product
	// when not specified currency is returned is USD
	// in: query
	// required: false
	Currency string `json:"currency"`
}

// No content is returned by this API endpoint
// swagger:response noContent
type productsNoContent struct {
}

// Generic error message returned as a string
// swagger:response errorResponse
type errorResponseWrapper struct {
	// Error Message
	// in: Body
	Body GenericError
}

// Validation error messages
// swagger:response validationResponse
type validationResponseWrappper struct {
	// Collection of error messages
	// in: Body
	Body ValidationError
}
