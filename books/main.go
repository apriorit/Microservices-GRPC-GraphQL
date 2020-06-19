package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tutorial/books/server"
	"tutorial/db"
)

func main() {
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFn()

	database, err := db.Connect(ctx, db.GetDbConfig())
	if err != nil {
		log.Fatalf("db.OpenDatabase failed with error: %s", err)
	}

	srv, err := server.NewServer(ctx, database)
	if err != nil {
		log.Fatalf("NewServer failed with error: %s", err)
	}

	srv.Run()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	signal := <-sigChan
	log.Printf("shutting down books server with signal: %s", signal)
}
