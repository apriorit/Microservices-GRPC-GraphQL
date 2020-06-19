package server

import (
	"context"
	"log"
	"os"
	"testing"
	"time"
	"tutorial/db"
	holdersv1 "tutorial/gen/go/proto/holders"

	"github.com/arangodb/go-driver"
	"github.com/stretchr/testify/assert"
)

const (
	dbNameForTests                  = "HoldersTest"
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

	col, err := db.Collection(operationContext, holdersCollectionName)
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

func createTestHolder(t *testing.T, operationContext context.Context, s *Server, testHolder *holdersv1.Holder) *holdersv1.Holder {
	addResponse, err := s.AddHolder(operationContext, &holdersv1.AddHolderRequest{Holder: testHolder})
	assert.NoError(t, err)
	assert.NotNil(t, addResponse)
	assert.NotEmpty(t, addResponse.Holder.Id)

	return addResponse.Holder
}

func assertHoldersEqual(t *testing.T, expected *holdersv1.Holder, actual *holdersv1.Holder) {
	assert.Equal(t, expected.Id, actual.Id)
	assert.Equal(t, expected.FirstName, actual.FirstName)
	assert.Equal(t, expected.LastName, actual.LastName)
	assert.Equal(t, expected.Email, actual.Email)
	assert.Equal(t, expected.Phone, actual.Phone)
	assert.Equal(t, expected.HeldBooks, actual.HeldBooks)
}

func TestAddHolder(t *testing.T) {
	contextWithTimeOut, cancelFn := context.WithTimeout(context.Background(), time.Second*deadlinePerTest)
	defer cancelFn()
	srv := newTestServer(t, contextWithTimeOut)

	createTestHolder(t, contextWithTimeOut, srv, testholder1)
}

func TestGetBook(t *testing.T) {
	contextWithTimeOut, cancelFn := context.WithTimeout(context.Background(), time.Second*deadlinePerTest)
	defer cancelFn()
	srv := newTestServer(t, contextWithTimeOut)

	testHolderWithId := createTestHolder(t, contextWithTimeOut, srv, testholder1)
	getResponse, err := srv.GetHolder(contextWithTimeOut, &holdersv1.GetHolderRequest{Id: testHolderWithId.Id})
	assert.NoError(t, err)
	assert.NotNil(t, getResponse)
	assertHoldersEqual(t, testHolderWithId, getResponse.Holder)

	getResponse, err = srv.GetHolder(contextWithTimeOut, &holdersv1.GetHolderRequest{Id: ""})
	assert.Error(t, err, "because id is empty")
	assert.Contains(t, err.Error(), "not provided")
	assert.Nil(t, getResponse)

	getResponse, err = srv.GetHolder(contextWithTimeOut, nil)
	assert.Error(t, err, "because request is empty")
	assert.Contains(t, err.Error(), "not provided")
	assert.Nil(t, getResponse)

	getResponse, err = srv.GetHolder(contextWithTimeOut, &holdersv1.GetHolderRequest{Id: "unknown"})
	assert.Error(t, err, "because id is unknown")
	assert.Contains(t, err.Error(), "not found")
	assert.Nil(t, getResponse)
}

func TestListHolders(t *testing.T) {
	contextWithTimeOut, cancelFn := context.WithTimeout(context.Background(), time.Second*deadlinePerTest)
	defer cancelFn()
	srv := newTestServer(t, contextWithTimeOut)

	testHolder1 := createTestHolder(t, contextWithTimeOut, srv, testholder1)
	testHolder2 := createTestHolder(t, contextWithTimeOut, srv, testholder2)

	listResponse, err := srv.ListHolders(contextWithTimeOut, &holdersv1.ListHoldersRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listResponse)
	assert.Len(t, listResponse.Holders, 2)
	assertHoldersEqual(t, testHolder1, listResponse.Holders[0])
	assertHoldersEqual(t, testHolder2, listResponse.Holders[1])
}

func TestUpdateHolder(t *testing.T) {
	contextWithTimeOut, cancelFn := context.WithTimeout(context.Background(), time.Second*deadlinePerTest)
	defer cancelFn()
	srv := newTestServer(t, contextWithTimeOut)

	testHolder1 := createTestHolder(t, contextWithTimeOut, srv, testholder1)
	testHolder1.Email = "updated@io.com"

	updateResponse, err := srv.UpdateHolder(contextWithTimeOut, &holdersv1.UpdateHolderRequest{Holder: testHolder1})
	assert.NoError(t, err)
	assert.NotNil(t, updateResponse)

	getResponse, err := srv.GetHolder(contextWithTimeOut, &holdersv1.GetHolderRequest{Id: testHolder1.Id})
	assert.NoError(t, err)
	assert.NotNil(t, getResponse)

	assertHoldersEqual(t, testHolder1, getResponse.Holder)
}

func TestGetHolderByBookId(t *testing.T) {
	contextWithTimeOut, cancelFn := context.WithTimeout(context.Background(), time.Second*deadlinePerTest)
	defer cancelFn()
	srv := newTestServer(t, contextWithTimeOut)

	testholder1.HeldBooks = []string{"bookId1", "bookId2"}
	testHolder1 := createTestHolder(t, contextWithTimeOut, srv, testholder1)

	getResponse, err := srv.GetHolderByBookId(contextWithTimeOut, &holdersv1.GetHolderByBookIdRequest{Id: "bookId1"})
	assert.NoError(t, err)
	assert.NotNil(t, getResponse)
	assertHoldersEqual(t, testHolder1, getResponse.Holder)
}
