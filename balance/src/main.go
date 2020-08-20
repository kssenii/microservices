package main

import (
	"context"
	"handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	dbHandler := handlers.SetupDB()
	clientHandler := handlers.NewRequest(dbHandler)

	serveMux := mux.NewRouter()

	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", clientHandler.GetBalance)
	getRouter.Use(clientHandler.MiddlewareValidateClient)

	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/", clientHandler.UpdateBalance)
	putRouter.Use(clientHandler.MiddlewareValidateClient)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", clientHandler.AddClient)
	postRouter.Use(clientHandler.MiddlewareValidateClient)

	serverLogger := log.New(os.Stdout, "ServerLog ", log.LstdFlags)
	server := &http.Server{
		Addr:         ":9090",
		Handler:      serveMux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		serverLogger.Println("Starting server on port 9090")

		err := server.ListenAndServe()
		if err != nil {
			serverLogger.Printf("Error starting server: %s", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan // Block untill message is available to be consumed
	serverLogger.Println(" Received terminate, gracefull shutdown", sig)

	timeoutContext, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(timeoutContext)

}
