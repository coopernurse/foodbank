package routes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"foodbank/internal/db"
	"foodbank/internal/email"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type VisitsHandlerTestSuite struct {
	suite.Suite
	server *httptest.Server
	db     *db.FirestoreDB
}

func (suite *VisitsHandlerTestSuite) SetupSuite() {
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
	handler := &VisitsHandler{suite.db, mockEmailSender}

	// Create the Echo server
	e := echo.New()
	handler.RegisterRoutes(e)

	// Start the test server
	suite.server = httptest.NewServer(e)
}

func (suite *VisitsHandlerTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *VisitsHandlerTestSuite) TestLoadHouseholdVisits() {
	// Test the LoadHouseholdVisits handler
	resp, err := http.Get(suite.server.URL + "/household/123/visits")
	if err != nil {
		suite.FailNow("Failed to make request", err)
	}
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func (suite *VisitsHandlerTestSuite) TestComputeItemLimits() {
	// Test the ComputeItemLimits handler
	resp, err := http.Get(suite.server.URL + "/household/123/visits/limits")
	if err != nil {
		suite.FailNow("Failed to make request", err)
	}
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func (suite *VisitsHandlerTestSuite) TestAddItemToVisit() {
	// Test the AddItemToVisit handler
	resp, err := http.Post(suite.server.URL+"/household/123/visit/456/item", "application/json", nil)
	if err != nil {
		suite.FailNow("Failed to make request", err)
	}
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func (suite *VisitsHandlerTestSuite) TestPutVisit() {
	// Test the PutVisit handler
	resp, err := http.Post(suite.server.URL+"/visit", "application/json", nil)
	if err != nil {
		suite.FailNow("Failed to make request", err)
	}
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func TestVisitsHandlerSuite(t *testing.T) {
	suite.Run(t, new(VisitsHandlerTestSuite))
}
