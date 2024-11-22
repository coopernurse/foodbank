package ui

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cupboard/internal/db"
	"cupboard/internal/email"
	"cupboard/internal/model"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SignupPageTestSuite struct {
	suite.Suite
	server *httptest.Server
	db     *db.FirestoreDB
}

func (suite *SignupPageTestSuite) SetupSuite() {
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
	handler := &SignupPage{DB: suite.db}

	// Create the Echo server
	e := echo.New()
	e.GET("/signup", handler.GET)
	e.POST("/signup", handler.POST)

	// Start the test server
	suite.server = httptest.NewServer(e)
}

func (suite *SignupPageTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *SignupPageTestSuite) TestSignupGET() {
	// Test the Signup GET handler
	resp, err := http.Get(suite.server.URL + "/signup")
	if err != nil {
		suite.FailNow("Failed to make request", err)
	}
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func (suite *SignupPageTestSuite) TestSignupPOST() {
	// Test the Signup POST handler with valid data
	jsonData := `{"firstName": "John", "lastName": "Doe", "email": "john.doe@example.com", "password": "password123"}`
	resp, err := http.Post(suite.server.URL+"/signup", "application/json", strings.NewReader(jsonData))
	if err != nil {
		suite.FailNow("Failed to make request", err)
	}
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Test the Signup POST handler with invalid data (missing email)
	jsonDataInvalid := `{"firstName": "John", "lastName": "Doe", "password": "password123"}`
	respInvalid, errInvalid := http.Post(suite.server.URL+"/signup", "application/json", strings.NewReader(jsonDataInvalid))
	if errInvalid != nil {
		suite.FailNow("Failed to make request", errInvalid)
	}
	defer respInvalid.Body.Close()

	assert.Equal(suite.T(), http.StatusBadRequest, respInvalid.StatusCode)
}

func TestSignupPageSuite(t *testing.T) {
	suite.Run(t, new(SignupPageTestSuite))
}
