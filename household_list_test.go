package ui

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"cupboard/internal/db"
	"cupboard/internal/email"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HouseholdListPageTestSuite struct {
	suite.Suite
	server *httptest.Server
	db     *db.FirestoreDB
}

func (suite *HouseholdListPageTestSuite) SetupSuite() {
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
	handler := &HouseholdListPage{DB: suite.db}

	// Create the Echo server
	e := echo.New()
	e.GET("/households", handler.GET)

	// Start the test server
	suite.server = httptest.NewServer(e)
}

func (suite *HouseholdListPageTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *HouseholdListPageTestSuite) TestHouseholdListGET() {
	// Test the HouseholdList GET handler
	resp, err := http.Get(suite.server.URL + "/households")
	if err != nil {
		suite.FailNow("Failed to make request", err)
	}
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func TestHouseholdListPageSuite(t *testing.T) {
	suite.Run(t, new(HouseholdListPageTestSuite))
}
