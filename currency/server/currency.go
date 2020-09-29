package server

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/redjoker011/online-cafe/currency/data"
	protos "github.com/redjoker011/online-cafe/currency/protos/currency"
)

type Currency struct {
	log   hclog.Logger
	rates *data.ExchangeRates
}

// Define Factory Function
func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {
	return &Currency{log: l, rates: r}
}

// Define service
func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())

	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())

	if err != nil {
		return nil, err
	}
	return &protos.RateResponse{Rate: rate}, nil
}
