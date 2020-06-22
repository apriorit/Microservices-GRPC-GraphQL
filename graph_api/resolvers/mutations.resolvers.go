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

func (r *mutationResolver) CreateBook(ctx context.Context, inputData model.BookInput) (*model.Book, error) {
	log.Printf("Book: %+v", inputData)

	book := graph2ServiceBookInput(&inputData)

	addBookResponse, err := r.services.Books().AddBook(ctx, &booksv1.AddBookRequest{Book: book})
	if err != nil {
		return nil, err
	}

	return service2GraphBook(addBookResponse.Book), nil
}

func (r *mutationResolver) DeleteBook(ctx context.Context, id string) (bool, error) {
	_, err := r.services.Books().DeleteBook(ctx, &booksv1.DeleteBookRequest{Id: id})
	return err == nil, err
}

func (r *mutationResolver) TakeBookInUse(ctx context.Context, holderID string, bookID string) (bool, error) {
	book, holder, err := getBookAndHolder(ctx, r.services, holderID, bookID)
	if err != nil {
		return false, err
	}
	log.Printf("Found book: %+v", book)

	// Add book id to holder
	holder.HeldBooks = append(holder.HeldBooks, bookID)
	_, err = r.services.Holders().UpdateHolder(ctx, &holdersv1.UpdateHolderRequest{Holder: holder})
	if err != nil {
		return false, err
	}
	log.Print("Holder updated successfully")

	return err == nil, err
}

func (r *mutationResolver) ReturnBook(ctx context.Context, holderID string, bookID string) (bool, error) {
	book, holder, err := getBookAndHolder(ctx, r.services, holderID, bookID)
	if err != nil {
		return false, err
	}
	log.Printf("Found book: %+v", book)

	for pos, id := range holder.HeldBooks {
		if bookID == id {
			holder.HeldBooks = append(holder.HeldBooks[:pos], holder.HeldBooks[pos+1:]...)
			break
		}
	}

	_, err = r.services.Holders().UpdateHolder(ctx, &holdersv1.UpdateHolderRequest{Holder: holder})
	if err != nil {
		return false, err
	}
	log.Print("Holder updated successfully")

	return err == nil, err
}

func (r *mutationResolver) CreateHolder(ctx context.Context, inputData model.HolderInput) (*model.Holder, error) {
	log.Printf("Holder: %+v", inputData)

	holder := graph2ServiceHolderInput(&inputData)

	addHolderResponse, err := r.services.Holders().AddHolder(ctx, &holdersv1.AddHolderRequest{Holder: holder})
	if err != nil {
		return nil, err
	}

	return service2GraphHolder(addHolderResponse.Holder), nil
}

// Mutation returns gen.MutationResolver implementation.
func (r *Resolver) Mutation() gen.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
