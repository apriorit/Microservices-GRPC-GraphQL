package server

import (
	"context"
	"log"
	"os"
	"testing"
	"time"
	"tutorial/db"

	"github.com/arangodb/go-driver"
	"github.com/stretchr/testify/assert"

	booksv1 "tutorial/gen/go/proto/books"
)

const (
	dbNameForTests                  = "BooksTest"
	deadlinePerTest                 = time.Duration(5 * time.Second)
	deadlineOnStartContanerForTests = time.Duration(60 * time.Second)
)

var dbConf = db.CreateTestDbConfig()

func TestMain(m *testing.M) {
	ctx, cancelFn := context.WithTimeout(context.Background(), deadlineOnStartContanerForTests)
	defer cancelFn()

	testContainer, err := db.RunContainerForTest(ctx, dbConf)
	if err != nil {
		log.Printf("Failed to create test container: %s", err)
		os.Exit(1)
	}
	log.Printf("success.\n")

	defer testContainer.Terminate(ctx)
	os.Exit(m.Run())
}

func newTestServer(t *testing.T, operationContext context.Context) *Server {
	db, err := db.Connect(operationContext, dbConf)
	assert.NoError(t, err)

	col, err := db.Collection(operationContext, booksCollectionName)
	if err != nil {
		assert.True(t, driver.IsNotFound(err))
	} else {
		err = col.Remove(operationContext)
		assert.NoError(t, err)
	}

	srv, err := NewServer(operationContext, db)
	assert.NoError(t, err)
	return srv
}

func assertBooksEqual(t *testing.T, expected *booksv1.Book, actual *booksv1.Book) {
	assert.Equal(t, expected.Id, actual.Id)
	assert.Equal(t, expected.Author, actual.Author)
	assert.Equal(t, expected.Title, actual.Title)
	assert.Equal(t, expected.Isbn, actual.Isbn)
}

func createTestBook(t *testing.T, operationContext context.Context, s *Server, testBook *booksv1.Book) *booksv1.Book {
	addResponse, err := s.AddBook(operationContext, &booksv1.AddBookRequest{Book: testBook})
	assert.NoError(t, err)
	assert.NotNil(t, addResponse)
	assert.NotEmpty(t, addResponse.Book.Id)

	return addResponse.Book
}

func TestAddBook(t *testing.T) {
	contextWithTimeOut, cancelFn := context.WithTimeout(context.Background(), deadlinePerTest)
	defer cancelFn()
	srv := newTestServer(t, contextWithTimeOut)

	createTestBook(t, contextWithTimeOut, srv, testBook1)

	// Cannot create book with the same isbn because its unique
	testBook2.Isbn = testBook1.Isbn
	addResponse, err := srv.AddBook(contextWithTimeOut, &booksv1.AddBookRequest{Book: testBook2})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unique constraint violated")
	assert.Nil(t, addResponse)
}

func TestDeleteBook(t *testing.T) {
	contextWithTimeOut, cancelFn := context.WithTimeout(context.Background(), deadlinePerTest)
	defer cancelFn()
	srv := newTestServer(t, contextWithTimeOut)

	testBookWithId := createTestBook(t, contextWithTimeOut, srv, testBook1)
	deleteResponse, err := srv.DeleteBook(contextWithTimeOut, &booksv1.DeleteBookRequest{Id: testBookWithId.Id})
	assert.NoError(t, err)
	assert.NotNil(t, deleteResponse)

	deleteResponse, err = srv.DeleteBook(contextWithTimeOut, &booksv1.DeleteBookRequest{Id: ""})
	assert.Error(t, err, "because id is empty")
	assert.Nil(t, deleteResponse)

	deleteResponse, err = srv.DeleteBook(contextWithTimeOut, nil)
	assert.Error(t, err, "because request is empty")
	assert.Nil(t, deleteResponse)

	deleteResponse, err = srv.DeleteBook(contextWithTimeOut, &booksv1.DeleteBookRequest{Id: "unknown"})
	assert.Error(t, err, "because id is unknown")
	assert.Nil(t, deleteResponse)
}

func TestGetBookByISBN(t *testing.T) {
	contextWithTimeOut, cancelFn := context.WithTimeout(context.Background(), deadlinePerTest)
	defer cancelFn()
	srv := newTestServer(t, contextWithTimeOut)

	testBookWithId := createTestBook(t, contextWithTimeOut, srv, testBook1)
	getResponse, err := srv.GetBookByISBN(contextWithTimeOut, &booksv1.GetBookByISBNRequest{Isbn: testBookWithId.Isbn})
	assert.NoError(t, err)
	assert.NotNil(t, getResponse)
	assertBooksEqual(t, testBookWithId, getResponse.Book)

	getResponse, err = srv.GetBookByISBN(contextWithTimeOut, &booksv1.GetBookByISBNRequest{Isbn: ""})
	assert.Error(t, err, "because isbn is empty")
	assert.Contains(t, err.Error(), "not provided")
	assert.Nil(t, getResponse)

	getResponse, err = srv.GetBookByISBN(contextWithTimeOut, nil)
	assert.Error(t, err, "because request is empty")
	assert.Contains(t, err.Error(), "not provided")
	assert.Nil(t, getResponse)

	getResponse, err = srv.GetBookByISBN(contextWithTimeOut, &booksv1.GetBookByISBNRequest{Isbn: "unknown"})
	assert.Error(t, err, "because isbn is unknown")
	assert.Contains(t, err.Error(), "not found")
	assert.Nil(t, getResponse)
}

func TestGetBook(t *testing.T) {
	contextWithTimeOut, cancelFn := context.WithTimeout(context.Background(), deadlinePerTest)
	defer cancelFn()
	srv := newTestServer(t, contextWithTimeOut)

	testBookWithId := createTestBook(t, contextWithTimeOut, srv, testBook1)
	getResponse, err := srv.GetBook(contextWithTimeOut, &booksv1.GetBookRequest{Id: testBookWithId.Id})
	assert.NoError(t, err)
	assert.NotNil(t, getResponse)
	assertBooksEqual(t, testBookWithId, getResponse.Book)

	getResponse, err = srv.GetBook(contextWithTimeOut, &booksv1.GetBookRequest{Id: ""})
	assert.Error(t, err, "because id is empty")
	assert.Contains(t, err.Error(), "not provided")
	assert.Nil(t, getResponse)

	getResponse, err = srv.GetBook(contextWithTimeOut, nil)
	assert.Error(t, err, "because request is empty")
	assert.Contains(t, err.Error(), "not provided")
	assert.Nil(t, getResponse)

	getResponse, err = srv.GetBook(contextWithTimeOut, &booksv1.GetBookRequest{Id: "unknown"})
	assert.Error(t, err, "because id is unknown")
	assert.Contains(t, err.Error(), "not found")
	assert.Nil(t, getResponse)
}

func TestGetBooks(t *testing.T) {
	contextWithTimeOut, cancelFn := context.WithTimeout(context.Background(), deadlinePerTest)
	defer cancelFn()
	srv := newTestServer(t, contextWithTimeOut)

	testBookWithId1 := createTestBook(t, contextWithTimeOut, srv, testBook1)
	testBookWithId2 := createTestBook(t, contextWithTimeOut, srv, testBook2)
	getBooksResponse, err := srv.GetBooks(contextWithTimeOut, &booksv1.GetBooksRequest{Ids: []string{testBookWithId1.Id, testBookWithId2.Id}})
	assert.NoError(t, err)
	assert.NotNil(t, getBooksResponse)
	assertBooksEqual(t, testBookWithId1, getBooksResponse.Books[0])
	assertBooksEqual(t, testBookWithId2, getBooksResponse.Books[1])

	getBooksResponse, err = srv.GetBooks(contextWithTimeOut, &booksv1.GetBooksRequest{Ids: nil})
	assert.Error(t, err, "because ids are not provided")
	assert.Contains(t, err.Error(), "not provided")
	assert.Nil(t, getBooksResponse)

	getBooksResponse, err = srv.GetBooks(contextWithTimeOut, nil)
	assert.Error(t, err, "because request is empty")
	assert.Contains(t, err.Error(), "not provided")
	assert.Nil(t, getBooksResponse)
}

func TestListBooks(t *testing.T) {
	contextWithTimeOut, cancelFn := context.WithTimeout(context.Background(), deadlinePerTest)
	defer cancelFn()
	srv := newTestServer(t, contextWithTimeOut)

	testbook1 := createTestBook(t, contextWithTimeOut, srv, testBook1)
	testbook2 := createTestBook(t, contextWithTimeOut, srv, testBook2)

	listResponse, err := srv.ListBooks(contextWithTimeOut, &booksv1.ListBooksRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listResponse)
	assert.Len(t, listResponse.Books, 2)
	assertBooksEqual(t, testbook1, listResponse.Books[0])
	assertBooksEqual(t, testbook2, listResponse.Books[1])

	listResponse, err = srv.ListBooks(contextWithTimeOut, nil)
	assert.Error(t, err)
	assert.Nil(t, listResponse)
}
