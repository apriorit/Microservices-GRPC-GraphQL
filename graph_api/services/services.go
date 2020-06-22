package services

import (
	"io"
	"log"

	booksv1 "tutorial/gen/go/proto/books"
	holdersv1 "tutorial/gen/go/proto/holders"

	"google.golang.org/grpc"
)

type ServicesConfig struct {
	BooksSvc   string
	HoldersSvc string
}

type services struct {
	io.Closer
	booksClientConn   *grpc.ClientConn
	booksClient       booksv1.BooksAPIClient
	holdersClientConn *grpc.ClientConn
	holdersClient     holdersv1.HoldersAPIClient
}

type Services interface {
	Books() booksv1.BooksAPIClient
	Holders() holdersv1.HoldersAPIClient
}

func NewServicesKeeper(conf ServicesConfig) (Services, error) {
	log.Printf("Connection to Books Service: %s...", conf.BooksSvc)
	booksConnection, err := grpc.Dial(conf.BooksSvc, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	log.Printf("Connection to Holders Service: %s...", conf.HoldersSvc)
	holdersConnection, err := grpc.Dial(conf.HoldersSvc, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	ah := &services{
		booksClientConn:   booksConnection,
		booksClient:       booksv1.NewBooksAPIClient(booksConnection),
		holdersClientConn: holdersConnection,
		holdersClient:     holdersv1.NewHoldersAPIClient(holdersConnection),
	}
	return ah, nil
}

func (ah *services) Books() booksv1.BooksAPIClient {
	return ah.booksClient
}

func (ah *services) Holders() holdersv1.HoldersAPIClient {
	return ah.holdersClient
}

func (ah *services) Close() error {
	err := ah.booksClientConn.Close()
	if err != nil {
		log.Printf("An error occurred while closing connection on Books service: %s", err)
	}
	err = ah.holdersClientConn.Close()
	if err != nil {
		log.Printf("An error occurred while closing connection on Holders service: %s", err)
	}

	return nil
}
