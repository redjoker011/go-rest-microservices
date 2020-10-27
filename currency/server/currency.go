package server

import (
	"context"
	"io"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/redjoker011/online-cafe/currency/data"
	protos "github.com/redjoker011/online-cafe/currency/protos/currency"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Currency struct {
	log           hclog.Logger
	rates         *data.ExchangeRates
	subscriptions map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest
}

// Define Factory Function
func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {
	c := &Currency{l, r, make(map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest)}
	go c.handleUpdates()

	return c
}

// Define service
func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())

	// Return structured error for grpc
	if rr.Base == rr.Destination {
		err := status.Newf(
			codes.InvalidArgument,
			"Base currency %s cannot be the same as destination currency %s",
			rr.Base.String(),
			rr.Destination.String(),
		)

		err, wde := err.WithDetails(rr)

		if wde != nil {
			return nil, wde
		}

		return nil, err.Err()
	}

	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())

	if err != nil {
		c.log.Error("Unable to convert rate", "base", rr.GetBase(), "destination", rr.GetDestination(), "error", err)
		return &protos.RateResponse{Rate: 0.0}, err
	}
	return &protos.RateResponse{Rate: rate}, nil
}

func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {
	// Create a block for streaming
	for {
		rr, err := src.Recv()
		if err == io.EOF {
			c.log.Info("Client has closed the connection")
			break
		}

		if err != nil {
			c.log.Error("Unable to read from client", "error", err)
			break
		}

		c.log.Info("Handle client request", "request_base", rr.GetBase(), "request_dest", rr.GetDestination())

		rrs, ok := c.subscriptions[src]

		if !ok {
			rrs = []*protos.RateRequest{}
		}

		var validationError *status.Status
		// check that subscription does not exists
		for _, v := range rrs {
			if v.Base == rr.Base && v.Destination == rr.Destination {
				// subscription exists return errors
				validationError := status.Newf(
					codes.AlreadyExists,
					"Unable to subscribe for currency as subscription already exists",
				)

				// add original request as metadata
				validationError, err = validationError.WithDetails(rr)

				if err != nil {
					c.log.Error("Unable to add metadata to error", "error", err)
					break
				}
				break
			}
		}

		// if a validation error return error and continue
		if validationError != nil {
			src.Send(
				&protos.StreamingRateResponse{
					Message: &protos.StreamingRateResponse_Error{
						Error: validationError.Proto(),
					},
				},
			)
		}

		// all ok
		rrs = append(rrs, rr)
		c.subscriptions[src] = rrs
	}

	return nil
}

func (c *Currency) handleUpdates() {
	ru := c.rates.MonitorRates(5 * time.Second)
	for range ru {
		c.log.Info("Got Updated rates")

		// loop over subscribed clients
		for k, v := range c.subscriptions {

			// loop over subscribed rates
			for _, rr := range v {
				r, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
				if err != nil {
					c.log.Error("Unable to get update rate", "base", rr.GetBase().String(), "destination", rr.GetDestination().String())
				}

				// create the response and send to the client
				err = k.Send(&protos.StreamingRateResponse{
					Message: &protos.StreamingRateResponse_RateResponse{
						RateResponse: &protos.RateResponse{Rate: r},
					},
				})

				if err != nil {
					c.log.Error("Unable to send updated rate", "base", rr.GetBase().String(), "destination", rr.GetDestination().String())
				}
			}
		}

	}
}
