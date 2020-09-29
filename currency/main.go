package main

import (
	"net"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/redjoker011/online-cafe/currency/data"
	protos "github.com/redjoker011/online-cafe/currency/protos/currency"
	"github.com/redjoker011/online-cafe/currency/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()
	rates, err := data.NewRates(log)

	if err != nil {
		log.Error("Unable to generate rates", "error", err)
	}

	// Initialize new grpc server
	gs := grpc.NewServer()
	// Initialize new server instance
	cs := server.NewCurrency(rates, log)

	// Register GRPC Server
	protos.RegisterCurrencyServer(gs, cs)

	// Register grpc server under reflection
	reflection.Register(gs)

	l, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}
	// Accept Request
	gs.Serve(l)
}
