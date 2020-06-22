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

	resp, err := r.services.Books().GetBooks(ctx, &booksv1.GetBooksRequest{Ids: obj.HeldBooks})
	if err != nil {
		return nil, err
	}

	books := make([]*model.Book, len(resp.Books))
	for i, book := range resp.Books {
		books[i] = service2GraphBook(book)
	}

	return books, nil
}

// Holder returns gen.HolderResolver implementation.
func (r *Resolver) Holder() gen.HolderResolver { return &holderResolver{r} }

type holderResolver struct{ *Resolver }
