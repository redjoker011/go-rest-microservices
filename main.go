package main

import (
	"context"
	"log"
	"net/http"
	"online-cafe/handlers"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	logger := log.New(os.Stdout, "products-api", log.LstdFlags)
	// Initialize Handlers
	np := handlers.NewProducts(logger)

	// Initialize new servemux and bind handlers using Gorilla
	sm := mux.NewRouter()
	// Initialize a subrouter and filter GET method requests
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", np.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products/{id:[0-9]+}", np.UpdateProduct)
	putRouter.Use(np.MiddlewareProductValidation)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", np.AddProduct)
	postRouter.Use(np.MiddlewareProductValidation)

	server := &http.Server{
		Addr:              ":9090",
		Handler:           sm,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	// Initialize Channel
	sigChan := make(chan os.Signal)
	// Add os listener to kill, interrupt command
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// Add Channel listener
	sig := <-sigChan

	logger.Println("Received Terminate, gracefully shutting down", sig)

	// Graceful shutdown
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(tc)
}
