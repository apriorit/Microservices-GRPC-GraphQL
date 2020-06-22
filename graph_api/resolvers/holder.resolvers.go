package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"
	booksv1 "tutorial/gen/go/proto/books"
	"tutorial/graph_api/gen"
	"tutorial/graph_api/model"
)

func (r *holderResolver) HeldBooks(ctx context.Context, obj *model.Holder) ([]*model.Book, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFn()

	books := []*model.Book{}
	for _, bookId := range obj.HeldBooks {
		getBookResponse, err := r.services.Books().GetBook(ctx, &booksv1.GetBookRequest{Id: bookId})
		if err != nil {
			return nil, err
		}
		books = append(books, &model.Book{
			ID:     &getBookResponse.Book.Id,
			Author: &getBookResponse.Book.Author,
			Title:  &getBookResponse.Book.Title,
			Isbn:   &getBookResponse.Book.Isbn,
		})
	}

	return books, nil
}

// Holder returns gen.HolderResolver implementation.
func (r *Resolver) Holder() gen.HolderResolver { return &holderResolver{r} }

type holderResolver struct{ *Resolver }
