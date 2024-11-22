package routes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cupboard/internal/db"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FoodBanksHandlerTestSuite struct {
	suite.Suite
	server *httptest.Server
	db     *db.FirestoreDB
}

func (suite *FoodBanksHandlerTestSuite) SetupSuite() {
	// Initialize Firestore client for the emulator
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "project-id")
	if err != nil {
		suite.FailNow("Failed to create Firestore client", err)
	}
	suite.db = db.NewFirestoreDB(client)

	// Create the handler
	handler := &FoodBanksHandler{suite.db}

	// Create the Echo server
	e := echo.New()
	e.POST("/foodbank", handler.PutFoodBank)
	e.GET("/foodbanks", handler.LoadFoodBanks)
	e.POST("/foodbank/:id/assign-permissions", handler.AssignPersonPermissions)

	// Start the test server
	suite.server = httptest.NewServer(e)
}

func (suite *FoodBanksHandlerTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *FoodBanksHandlerTestSuite) TestPutFoodBank() {
	// Test the PutFoodBank handler with valid data
	jsonData := `{"name": "Food Bank Name", "address": {"street1": "123 Main St", "city": "Anytown", "state": "CA", "zip": "12345", "country": "USA"}}`
	resp, err := http.Post(suite.server.URL+"/foodbank", "application/json", strings.NewReader(jsonData))
	if err != nil {
		suite.FailNow("Failed to make request", err)
	}
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Test the PutFoodBank handler with invalid data (missing address)
	jsonDataInvalid := `{"name": "Food Bank Name"}`
	respInvalid, errInvalid := http.Post(suite.server.URL+"/foodbank", "application/json", strings.NewReader(jsonDataInvalid))
	if errInvalid != nil {
		suite.FailNow("Failed to make request", errInvalid)
	}
	defer respInvalid.Body.Close()

	assert.Equal(suite.T(), http.StatusBadRequest, respInvalid.StatusCode)
}

func (suite *FoodBanksHandlerTestSuite) TestLoadFoodBanks() {
	// Test the LoadFoodBanks handler
	resp, err := http.Get(suite.server.URL + "/foodbanks")
	if err != nil {
		suite.FailNow("Failed to make request", err)
	}
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func (suite *FoodBanksHandlerTestSuite) TestAssignPersonPermissions() {
	// Test the AssignPersonPermissions handler
	resp, err := http.Post(suite.server.URL+"/foodbank/123/assign-permissions", "application/json", nil)
	if err != nil {
		suite.FailNow("Failed to make request", err)
	}
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func TestFoodBanksHandlerSuite(t *testing.T) {
	suite.Run(t, new(FoodBanksHandlerTestSuite))
}
