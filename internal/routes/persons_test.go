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
package routes

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
	e.POST("/person", handler.PutPerson)
	e.GET("/persons/search", handler.SearchPersons)
	e.GET("/household/:id/persons", handler.LoadHouseholdPersons)
	e.POST("/person/:id/reset-password", handler.ResetPassword)
	e.POST("/person/:id/email-login-link", handler.EmailLoginLink)
	e.GET("/person/:id/permissions", handler.ResolvePermissions)

	// Start the test server
	suite.server = httptest.NewServer(e)
}

func (suite *PersonsHandlerTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *PersonsHandlerTestSuite) TestPutPerson() {
	// Test the PutPerson handler with valid data
	jsonData := `{"firstName": "John", "lastName": "Doe", "email": "john.doe@example.com", "password": "password123"}`
	resp, err := http.Post(suite.server.URL+"/person", "application/json", strings.NewReader(jsonData))
	if err != nil {
		suite.FailNow("Failed to make request", err)
	}
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Test the PutPerson handler with invalid data (missing email)
	jsonDataInvalid := `{"firstName": "John", "lastName": "Doe", "password": "password123"}`
	respInvalid, errInvalid := http.Post(suite.server.URL+"/person", "application/json", strings.NewReader(jsonDataInvalid))
	if errInvalid != nil {
		suite.FailNow("Failed to make request", errInvalid)
	}
	defer respInvalid.Body.Close()

	assert.Equal(suite.T(), http.StatusBadRequest, respInvalid.StatusCode)
}

func (suite *PersonsHandlerTestSuite) TestSearchPersons() {
	// Test the SearchPersons handler
	resp, err := http.Get(suite.server.URL + "/persons/search")
	if err != nil {
		suite.FailNow("Failed to make request", err)
	}
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func (suite *PersonsHandlerTestSuite) TestLoadHouseholdPersons() {
	// Test the LoadHouseholdPersons handler
	resp, err := http.Get(suite.server.URL + "/household/123/persons")
	if err != nil {
		suite.FailNow("Failed to make request", err)
	}
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func (suite *PersonsHandlerTestSuite) TestResetPassword() {
	// Test the ResetPassword handler
	resp, err := http.Post(suite.server.URL+"/person/123/reset-password", "application/json", nil)
	if err != nil {
		suite.FailNow("Failed to make request", err)
	}
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
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

func (suite *PersonsHandlerTestSuite) TestResolvePermissions() {
	// Test the ResolvePermissions handler
	resp, err := http.Get(suite.server.URL + "/person/123/permissions")
	if err != nil {
		suite.FailNow("Failed to make request", err)
	}
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func TestPersonsHandlerSuite(t *testing.T) {
	suite.Run(t, new(PersonsHandlerTestSuite))
}
