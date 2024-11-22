package routes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"cupboard/internal/db"
	"cupboard/internal/email"
	"cupboard/internal/model"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PersonsHandlerTestSuite struct {
	suite.Suite
	server *httptest.Server
	db     *db.FirestoreDB
}

func (suite *PersonsHandlerTestSuite) SetupSuite() {
	// Initialize Firestore client for the emulator
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "project-id")
	if err != nil {
		suite.FailNow("Failed to create Firestore client", err)
	}
	suite.db = db.NewFirestoreDB(client)

	// Create a mock email sender
	mockEmailSender := &email.MockEmailSender{}

	// Create the handler
	handler := routes.NewPersonsHandler(suite.db, mockEmailSender)

	// Create the Echo server
	e := echo.New()
	e.POST("/person/:id/email-login-link", handler.EmailLoginLink)

	// Start the test server
	suite.server = httptest.NewServer(e)
}

func (suite *PersonsHandlerTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *PersonsHandlerTestSuite) TestEmailLoginLink() {
	// Test the EmailLoginLink handler
	resp, err := http.Post(suite.server.URL+"/person/123/email-login-link", "application/json", nil)
	if err != nil {
		suite.FailNow("Failed to make request", err)
	}
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func TestPersonsHandlerSuite(t *testing.T) {
	suite.Run(t, new(PersonsHandlerTestSuite))
}
