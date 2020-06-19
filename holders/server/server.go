package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"tutorial/db"
	holdersv1 "tutorial/gen/go/proto/holders"

	"github.com/arangodb/go-driver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	holdersCollectionName = "Holders"
	defaultPort           = "60002"
)

type Server struct {
	database          driver.Database
	holdersCollection driver.Collection
}

func NewServer(ctx context.Context, database driver.Database) (*Server, error) {
	collection, err := db.AttachCollection(ctx, database, holdersCollectionName)
	if err != nil {
		return nil, err
	}

	return &Server{
		database:          database,
		holdersCollection: collection,
	}, nil
}

func (s *Server) Run() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = defaultPort
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Print("net.Listen failed")
		return
	}
	grpcServer := grpc.NewServer()
	holdersv1.RegisterHoldersAPIServer(grpcServer, s) // use authogenerated code to register the server
	reflection.Register(grpcServer)

	log.Printf("Starting Holders server on port %s", port)
	go func() {
		grpcServer.Serve(listener)
	}()
}

func (s *Server) ListHolders(ctx context.Context, in *holdersv1.ListHoldersRequest) (*holdersv1.ListHoldersResponse, error) {
	if in == nil {
		return nil, fmt.Errorf("Request is empty")
	}

	cursor, err := s.database.Query(ctx, db.ListRecords(s.holdersCollection.Name()), nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to iterate over documents: %s", err)
	}
	defer cursor.Close()

	allHolders := []*holdersv1.Holder{}
	for {
		holder := new(holdersv1.Holder)
		var meta driver.DocumentMeta
		meta, err := cursor.ReadDocument(ctx, holder)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("Failed to read book document: %s", err)
		}
		holder.Id = meta.Key
		allHolders = append(allHolders, holder)
	}

	return &holdersv1.ListHoldersResponse{Holders: allHolders}, nil
}

func (s *Server) GetHolder(ctx context.Context, in *holdersv1.GetHolderRequest) (*holdersv1.GetHolderResponse, error) {
	if in == nil || in.Id == "" {
		return nil, fmt.Errorf("Holder id is not provided")
	}

	holder := new(holdersv1.Holder)
	meta, err := s.holdersCollection.ReadDocument(ctx, in.Id, holder)
	if err != nil {
		if driver.IsNotFound(err) {
			err = fmt.Errorf("Holder with id '%s' not found", in.Id)
		} else {
			err = fmt.Errorf("Failed to get holder with id '%s': %s", in.Id, err)
		}
		return nil, err
	}
	holder.Id = meta.Key

	return &holdersv1.GetHolderResponse{Holder: holder}, nil
}

func (s *Server) GetHolderByBookId(ctx context.Context, in *holdersv1.GetHolderByBookIdRequest) (*holdersv1.GetHolderByBookIdResponse, error) {
	if in == nil || in.Id == "" {
		return nil, fmt.Errorf("Book id is not provided")
	}

	const queryHolderByBookId = `
	FOR holder IN %s
		FOR bookId IN holder.held_books
			FILTER bookId == @bookId
				RETURN holder`

	query := fmt.Sprintf(queryHolderByBookId, holdersCollectionName)
	bindVars := map[string]interface{}{"bookId": in.Id}

	cursor, err := s.database.Query(ctx, query, bindVars)
	if err != nil {
		return nil, fmt.Errorf("Failed to iterate over holder documents with query '%s': %s", queryHolderByBookId, err)
	}
	defer cursor.Close()

	h := new(holdersv1.Holder)
	meta, err := cursor.ReadDocument(ctx, h)
	if driver.IsNoMoreDocuments(err) {
		return nil, fmt.Errorf("Holder that held book with id %s not found: %s", in.Id, err)
	} else if err != nil {
		return nil, fmt.Errorf("Failed to read holder document: %s", err)
	}

	h.Id = meta.Key

	return &holdersv1.GetHolderByBookIdResponse{Holder: h}, nil
}

func (s *Server) AddHolder(ctx context.Context, in *holdersv1.AddHolderRequest) (*holdersv1.AddHolderResponse, error) {
	if in == nil || in.Holder == nil {
		return nil, fmt.Errorf("Book is not provided")
	}

	meta, err := s.holdersCollection.CreateDocument(ctx, in.Holder)

	if err != nil {
		return nil, fmt.Errorf("Failed to create book: %s", err)
	}

	in.Holder.Id = meta.Key
	return &holdersv1.AddHolderResponse{Holder: in.Holder}, nil
}

func (s *Server) UpdateHolder(ctx context.Context, in *holdersv1.UpdateHolderRequest) (*holdersv1.UpdateHolderResponse, error) {
	if in == nil || in.Holder == nil || in.Holder.Id == "" {
		return nil, fmt.Errorf("Existing holder is not provided")
	}

	_, err := s.holdersCollection.ReplaceDocument(ctx, in.Holder.Id, in.Holder)
	if err != nil {
		return nil, fmt.Errorf("Failed to update holder with id %s: %s", in.Holder.Id, err)
	}

	return &holdersv1.UpdateHolderResponse{Holder: in.Holder}, nil
}
