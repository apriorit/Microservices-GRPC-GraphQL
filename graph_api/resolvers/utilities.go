package resolvers

import (
	"context"
	"log"
	booksv1 "tutorial/gen/go/proto/books"
	holdersv1 "tutorial/gen/go/proto/holders"
	apiHolder "tutorial/graph_api/api_holder"
	"tutorial/graph_api/model"
)

func service2GraphBook(book *booksv1.Book) *model.Book {
	return &model.Book{
		ID:     &book.Id,
		Author: &book.Author,
		Title:  &book.Title,
		Isbn:   &book.Isbn,
	}
}

func graph2ServiceHolderInput(holderInput *model.HolderInput) *holdersv1.Holder {
	return &holdersv1.Holder{
		FirstName: softDeference(holderInput.FirstName),
		LastName:  softDeference(holderInput.LastName),
		Phone:     softDeference(holderInput.Phone),
		Email:     softDeference(holderInput.Email),
	}
}

func graph2ServiceBookInput(bookInput *model.BookInput) *booksv1.Book {
	return &booksv1.Book{
		Id:     softDeference(bookInput.ID),
		Author: softDeference(bookInput.Author),
		Title:  softDeference(bookInput.Title),
		Isbn:   softDeference(bookInput.Isbn),
	}
}
func service2GraphHolder(holder *holdersv1.Holder) *model.Holder {
	return &model.Holder{
		ID:        holder.Id,
		FirstName: holder.FirstName,
		LastName:  holder.LastName,
		Email:     holder.Email,
		Phone:     holder.Phone,
		HeldBooks: holder.HeldBooks,
	}
}

func softDeference(field *string) string {
	if field == nil {
		return ""
	}
	return *field
}

func getBookAndHolder(ctx context.Context, ah apiHolder.Services, holderId, bookId string) (*booksv1.Book, *holdersv1.Holder, error) {
	log.Printf("Request data. HolderID: %s, bookID: %s", holderId, bookId)

	getHolderResponse, err := ah.Holders().GetHolder(ctx, &holdersv1.GetHolderRequest{Id: holderId})
	if err != nil {
		return nil, nil, err
	}
	getBookResponse, err := ah.Books().GetBook(ctx, &booksv1.GetBookRequest{Id: bookId})
	if err != nil {
		return nil, nil, err
	}

	return getBookResponse.Book, getHolderResponse.Holder, nil
}
