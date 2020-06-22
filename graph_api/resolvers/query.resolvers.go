package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"log"
	booksv1 "tutorial/gen/go/proto/books"
	holdersv1 "tutorial/gen/go/proto/holders"
	"tutorial/graph_api/gen"
	"tutorial/graph_api/model"
)

func (r *queryResolver) Books(ctx context.Context, id *string, isbn *string, holderID *string) ([]*model.Book, error) {
	books := []*model.Book{}

	if id != nil {
		log.Printf("Book id: %s", *id)
		getBookResponse, err := r.services.Books().GetBook(ctx, &booksv1.GetBookRequest{Id: *id})
		if err != nil {
			return nil, err
		}
		books = append(books, service2GraphBook(getBookResponse.Book))
	} else if isbn != nil {
		log.Printf("Book isbn: %s", *isbn)
		getBookByISBNResponse, err := r.services.Books().GetBookByISBN(ctx, &booksv1.GetBookByISBNRequest{Isbn: *isbn})
		if err != nil {
			return nil, err
		}
		books = append(books, service2GraphBook(getBookByISBNResponse.Book))
	} else if holderID != nil {
		getHolderResponse, err := r.services.Holders().GetHolder(ctx, &holdersv1.GetHolderRequest{Id: *holderID})
		if err != nil {
			return nil, err
		}

		for _, bookId := range getHolderResponse.Holder.HeldBooks {
			getBookResponse, err := r.services.Books().GetBook(ctx, &booksv1.GetBookRequest{Id: bookId})
			if err != nil {
				return nil, err
			}
			books = append(books, service2GraphBook(getBookResponse.Book))
		}

	} else {
		listBookResponse, err := r.services.Books().ListBooks(ctx, &booksv1.ListBooksRequest{})
		if err != nil {
			return nil, err
		}
		for _, book := range listBookResponse.Books {
			books = append(books, service2GraphBook(book))
		}
	}

	return books, nil
}

func (r *queryResolver) Holders(ctx context.Context, id *string, bookID *string) ([]*model.Holder, error) {
	holders := []*model.Holder{}

	if id != nil {
		log.Printf("Holder id: %s", *id)
		getHolderResponse, err := r.services.Holders().GetHolder(ctx, &holdersv1.GetHolderRequest{Id: *id})
		if err != nil {
			return nil, err
		}
		holders = append(holders, service2GraphHolder(getHolderResponse.Holder))
	} else if bookID != nil {
		log.Printf("Book id: %s", *bookID)
		getHolderByBookIdResponse, err := r.services.Holders().GetHolderByBookId(ctx, &holdersv1.GetHolderByBookIdRequest{Id: *bookID})
		if err != nil {
			return nil, err
		}
		holders = append(holders, service2GraphHolder(getHolderByBookIdResponse.Holder))
	} else {
		log.Printf("ALl books will be retrieved")
		listHoldersResponse, err := r.services.Holders().ListHolders(ctx, &holdersv1.ListHoldersRequest{})
		if err != nil {
			return nil, err
		}
		for _, holder := range listHoldersResponse.Holders {
			holders = append(holders, service2GraphHolder(holder))
		}
	}

	return holders, nil
}

// Query returns gen.QueryResolver implementation.
func (r *Resolver) Query() gen.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
