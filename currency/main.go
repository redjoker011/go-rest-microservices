package main

import (
	"net"
	"os"

	"github.com/hashicorp/go-hclog"
	protos "github.com/redjoker011/online-cafe/currency/protos/currency"
	"github.com/redjoker011/online-cafe/currency/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()
	// Initialize new grpc server
	gs := grpc.NewServer()
	// Initialize new server instance
	cs := server.NewCurrency(log)

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
