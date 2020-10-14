package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	protos "github.com/redjoker011/online-cafe/currency/protos/currency"
	"github.com/redjoker011/online-cafe/data"
	"github.com/redjoker011/online-cafe/handlers"
	"google.golang.org/grpc"
)

func main() {
	l := hclog.Default()

	// Create GRPC Client
	conn, err := grpc.Dial("localhost:9092", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	// Currency Client
	cc := protos.NewCurrencyClient(conn)

	// create db instance
	db := data.NewProductsDB(cc, l)

	// Initialize Handlers
	np := handlers.NewProducts(l, db)

	// Initialize new servemux and bind handlers using Gorilla
	sm := mux.NewRouter()
	// Initialize a subrouter and filter GET method requests
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", np.GetProducts)
	getRouter.HandleFunc("/product", np.GetProduct)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products/{id:[0-9]+}", np.UpdateProduct)
	putRouter.Use(np.MiddlewareProductValidation)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", np.AddProduct)
	postRouter.Use(np.MiddlewareProductValidation)

	// Go Open API Runtime to generate Redoc Client
	opts := middleware.RedocOpts{SpecURL: "/swagger.yml"}
	sh := middleware.Redoc(opts, nil)
	// Swagger Handler
	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yml", http.FileServer(http.Dir("./")))

	// CORS
	// Whitelist consumer URL's
	urls := []string{"http:localhost:8080"}
	ch := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins(urls))

	server := &http.Server{
		Addr:              ":9090",
		Handler:           ch(sm), // Use servemux as cors handler argument
		ErrorLog:          l.StandardLogger(&hclog.StandardLoggerOptions{}),
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
	}

	l.Info("Starting Server", "port", ":9090")

	go func() {
		err := server.ListenAndServe()

		if err != nil {
			l.Error("Error starting server", "error", err)
		}
	}()

	// Initialize Channel
	sigChan := make(chan os.Signal)
	// Add os listener to kill, interrupt command
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// Add Channel listener
	sig := <-sigChan

	l.Info("Received Terminate, gracefully shutting down", sig)

	// Graceful shutdown
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(tc)
}
