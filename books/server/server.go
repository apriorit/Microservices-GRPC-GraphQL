package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"tutorial/db"
	booksv1 "tutorial/gen/go/proto/books"

	"github.com/arangodb/go-driver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	booksCollectionName = "Books"
	defaultPort         = "60001"
)

type Server struct {
	database        driver.Database
	booksCollection driver.Collection
}

func NewServer(ctx context.Context, database driver.Database) (*Server, error) {
	collection, err := db.AttachCollection(ctx, database, booksCollectionName)
	if err != nil {
		return nil, err
	}

	_, _, err = collection.EnsurePersistentIndex(ctx, []string{"isbn"}, &driver.EnsurePersistentIndexOptions{Unique: true})
	if err != nil {
		return nil, err
	}

	return &Server{
		database:        database,
		booksCollection: collection,
	}, nil
}

func (s *Server) Run() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = defaultPort
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Fatal("net.Listen failed")
		return
	}
	grpcServer := grpc.NewServer()

	booksv1.RegisterBooksAPIServer(grpcServer, s) // use authogenerated code to register the server
	reflection.Register(grpcServer)

	log.Printf("Starting Books server on port %s", port)
	go func() {
		grpcServer.Serve(listener)
	}()
}

func (s *Server) ListBooks(ctx context.Context, in *booksv1.ListBooksRequest) (*booksv1.ListBooksResponse, error) {
	if in == nil {
		return nil, fmt.Errorf("Request is empty")
	}

	cursor, err := s.database.Query(ctx, db.ListRecords(s.booksCollection.Name()), nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to iterate over documents: %s", err)
	}
	defer cursor.Close()

	allBooks := []*booksv1.Book{}
	for {
		book := new(booksv1.Book)
		var meta driver.DocumentMeta
		meta, err := cursor.ReadDocument(ctx, book)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("Failed to read book document: %s", err)
		}
		book.Id = meta.Key
		allBooks = append(allBooks, book)
	}

	return &booksv1.ListBooksResponse{Books: allBooks}, nil
}

func (s *Server) GetBook(ctx context.Context, in *booksv1.GetBookRequest) (*booksv1.GetBookResponse, error) {
	if in == nil || in.Id == "" {
		return nil, fmt.Errorf("Book id is not provided")
	}

	book := new(booksv1.Book)
	meta, err := s.booksCollection.ReadDocument(ctx, in.Id, book)
	if err != nil {
		if driver.IsNotFound(err) {
			err = fmt.Errorf("Book with id '%s' not found", in.Id)
		} else {
			err = fmt.Errorf("Failed to get book with id '%s': %s", in.Id, err)
		}
		return nil, err
	}
	book.Id = meta.Key

	return &booksv1.GetBookResponse{Book: book}, nil
}

func (s *Server) GetBooks(ctx context.Context, in *booksv1.GetBooksRequest) (*booksv1.GetBooksResponse, error) {
	if in == nil || len(in.Ids) == 0 {
		return nil, fmt.Errorf("Book ids are not provided")
	}

	const queryBooksByIds = `
	FOR book IN %s
		FILTER book._key in @ids
			RETURN book`

	query := fmt.Sprintf(queryBooksByIds, booksCollectionName)
	bindVars := map[string]interface{}{"ids": in.Ids}

	cursor, err := s.database.Query(ctx, query, bindVars)
	if err != nil {
		return nil, fmt.Errorf("Failed to iterate over book documents with query '%s': %s", queryBooksByIds, err)
	}
	defer cursor.Close()

	books := []*booksv1.Book{}
	for {
		book := new(booksv1.Book)
		meta, err := cursor.ReadDocument(ctx, book)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			log.Print(err)
			return nil, fmt.Errorf("Failed to read book document: %s", err)
		}

		book.Id = meta.Key
		books = append(books, book)
	}

	return &booksv1.GetBooksResponse{Books: books}, nil
}

func (s *Server) GetBookByISBN(ctx context.Context, in *booksv1.GetBookByISBNRequest) (*booksv1.GetBookByISBNResponse, error) {
	if in == nil || in.Isbn == "" {
		return nil, fmt.Errorf("Book isbn is not provided")
	}

	const queryBookByISBN = `
	FOR book IN %s
		FILTER book.isbn == @isbn
			RETURN book`

	query := fmt.Sprintf(queryBookByISBN, booksCollectionName)
	bindVars := map[string]interface{}{"isbn": in.Isbn}

	cursor, err := s.database.Query(ctx, query, bindVars)
	if err != nil {
		return nil, fmt.Errorf("Failed to iterate over book documents with query '%s': %s", queryBookByISBN, err)
	}
	defer cursor.Close()

	b := new(booksv1.Book)
	meta, err := cursor.ReadDocument(ctx, b)
	if driver.IsNoMoreDocuments(err) {
		return nil, fmt.Errorf("Book with ISBN '%s' not found: %s", in.Isbn, err)
	} else if err != nil {
		return nil, fmt.Errorf("Failed to read book document: %s", err)
	}

	b.Id = meta.Key

	return &booksv1.GetBookByISBNResponse{Book: b}, nil
}

func (s *Server) AddBook(ctx context.Context, in *booksv1.AddBookRequest) (*booksv1.AddBookResponse, error) {
	if in == nil || in.Book == nil {
		return nil, fmt.Errorf("Book is not provided")
	}

	meta, err := s.booksCollection.CreateDocument(ctx, in.Book)

	if err != nil {
		return nil, fmt.Errorf("Failed to create book: %s", err)
	}

	in.Book.Id = meta.Key
	return &booksv1.AddBookResponse{Book: in.Book}, nil
}

func (s *Server) DeleteBook(ctx context.Context, in *booksv1.DeleteBookRequest) (*booksv1.DeleteBookResponse, error) {
	if in == nil || in.Id == "" {
		return nil, fmt.Errorf("Book id is not provided")
	}

	_, err := s.booksCollection.RemoveDocument(ctx, in.Id)
	if err != nil {
		return nil, fmt.Errorf("Failed to remove existing book: %s", err)
	}

	return &booksv1.DeleteBookResponse{}, nil
}
