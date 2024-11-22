package routes

import (
	"context"
	"encoding/json"
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
	handler := &PersonsHandler{suite.db, mockEmailSender}

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

	// Test the PutPerson handler with invalid data (missing lastName)
	jsonDataInvalidPassword := `{"firstName": "John",  "email": "john.doe@example.com"}`
	respInvalidPassword, errInvalidPassword := http.Post(suite.server.URL+"/person", "application/json", strings.NewReader(jsonDataInvalidPassword))
	if errInvalidPassword != nil {
		suite.FailNow("Failed to make request", errInvalidPassword)
	}
	defer respInvalidPassword.Body.Close()

	assert.Equal(suite.T(), http.StatusBadRequest, respInvalidPassword.StatusCode)
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

func TestLogin(t *testing.T) {
	// Create a mock Firestore client
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "test-project")
	if err != nil {
		t.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	// Create a mock email sender
	mockEmailSender := &email.MockEmailSender{}

	// Create the handler
	dbInstance := db.NewFirestoreDB(client)
	handler := &PersonsHandler{DB: dbInstance, EmailSender: mockEmailSender}

	// Create the Echo server
	e := echo.New()
	e.POST("/login", handler.Login)

	// Start the test server
	server := httptest.NewServer(e)
	defer server.Close()

	// Create a test person input
	testPersonInput := model.PersonInput{
		Person: model.Person{
			PersonCommon: model.PersonCommon{
				Id:    "testPersonID",
				Email: "test@example.com",
			},
		},
		Password: "password123",
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPersonInput.Password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Create the test person with the hashed password
	testPerson := model.Person{
		PersonCommon: testPersonInput.PersonCommon,
		PasswordHash: string(hashedPassword),
	}

	// Save the test person to the mock Firestore
	if err := dbInstance.PutPerson(ctx, testPerson); err != nil {
		t.Fatalf("Failed to save test person: %v", err)
	}

	// Create a valid login request
	loginRequest := `{"email": "test@example.com", "password": "password123"}`
	resp, err := http.Post(server.URL+"/login", "application/json", strings.NewReader(loginRequest))
	if err != nil {
		t.Fatalf("Failed to make login request: %v", err)
	}
	defer resp.Body.Close()

	// Verify the response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify the response body contains the sessionToken
	var responseBody map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	assert.Contains(t, responseBody, "sessionToken")
}
