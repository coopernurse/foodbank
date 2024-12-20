package routes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"foodbank/internal/db"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ItemsHandlerTestSuite struct {
	suite.Suite
	server *httptest.Server
	db     *db.FirestoreDB
}

func (suite *ItemsHandlerTestSuite) SetupSuite() {
	// Initialize Firestore client for the emulator
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "project-id")
	if err != nil {
		suite.FailNow("Failed to create Firestore client", err)
	}
	suite.db = db.NewFirestoreDB(client)

	// Create the handler
	handler := &ItemsHandler{suite.db}

	// Create the Echo server
	e := echo.New()
	handler.RegisterRoutes(e)

	// Start the test server
	suite.server = httptest.NewServer(e)
}

func (suite *ItemsHandlerTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *ItemsHandlerTestSuite) TestPutItem() {
	// Test the PutItem handler with valid data
	jsonData := `{"foodBankId": "123", "name": "Item Name", "points": 10}`
	resp, err := http.Post(suite.server.URL+"/item", "application/json", strings.NewReader(jsonData))
	if err != nil {
		suite.FailNow("Failed to make request", err)
	}
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Test the PutItem handler with invalid data (missing name)
	jsonDataInvalid := `{"foodBankId": "123", "points": 10}`
	respInvalid, errInvalid := http.Post(suite.server.URL+"/item", "application/json", strings.NewReader(jsonDataInvalid))
	if errInvalid != nil {
		suite.FailNow("Failed to make request", errInvalid)
	}
	defer respInvalid.Body.Close()

	assert.Equal(suite.T(), http.StatusBadRequest, respInvalid.StatusCode)
}

func TestItemsHandlerSuite(t *testing.T) {
	suite.Run(t, new(ItemsHandlerTestSuite))
}
