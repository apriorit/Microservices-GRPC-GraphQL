package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apiHolder "tutorial/graph_api/api_holder"
	"tutorial/graph_api/gen"
	"tutorial/graph_api/resolvers"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
)

const defaultPort = "80"

func main() {
	// Load environment variables
	port := os.Getenv("GRAPH_API_PORT")
	if port == "" {
		port = defaultPort
	}
	booksSvc := os.Getenv("BOOKS_SERVICE")
	if booksSvc == "" {
		log.Fatalf("Failed to load environmet variable: %s", "BOOKS_SERVICE")
	}
	holdersSvc := os.Getenv("HOLDERS_SERVICE")
	if holdersSvc == "" {
		log.Fatalf("Failed to load environmet variable: %s", "HOLDERS_SERVICE")
	}

	// Connect to the services
	svcs, err := apiHolder.NewServicesHolder(apiHolder.ServicesConfig{
		BooksSvc:   booksSvc,
		HoldersSvc: holdersSvc,
	})
	if err != nil {
		log.Fatalf("Failed to create grpc api holder: %s", err)
	}

	// Create graphApi handlers
	router := mux.NewRouter()
	graphAPIHandler := handler.NewDefaultServer(gen.NewExecutableSchema(gen.Config{Resolvers: resolvers.NewResolver(svcs)}))
	router.Handle("/", playground.Handler("GraphQL playground", "/tutorial"))
	router.Handle("/tutorial", graphAPIHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
		Handler:      router,
	}

	// Start graph_api server
	log.Printf("Please connect to the http://localhost:%s/ for GraphQL playground", port)
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	// Block until cancel signal is received.
	<-sig
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	log.Print("Shutting down graph_api server")

	if err := srv.Shutdown(ctx); err != nil {
		log.Print(err)
	}
	<-ctx.Done()
	os.Exit(0)
}
