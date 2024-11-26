package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"foodbank/internal/db"
	"foodbank/internal/email"
	"foodbank/internal/model"

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
	handler.RegisterRoutes(e)

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

func (suite *PersonsHandlerTestSuite) TestPostHousehold() {
	t := suite.T()

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		validateDB     bool
	}{
		{
			name: "valid household",
			requestBody: `{
				"head": {
					"firstName": "John",
					"lastName": "Doe",
					"dob": "1980-01-01"
				},
				"members": []
			}`,
			expectedStatus: http.StatusOK,
			validateDB:     true,
		},
		{
			name: "missing required fields",
			requestBody: `{
				"head": {
					"firstName": "",
					"lastName": "Doe",
					"dob": ""
				},
				"members": []
			}`,
			expectedStatus: http.StatusBadRequest,
			validateDB:     false,
		},
		{
			name:           "invalid json",
			requestBody:    `{"invalid": json}`,
			expectedStatus: http.StatusBadRequest,
			validateDB:     false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Make the request
			resp, err := http.Post(suite.server.URL+"/household", "application/json", strings.NewReader(tc.requestBody))
			assert.NoError(t, err)
			defer resp.Body.Close()

			// Check status code
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			// For successful cases, verify the household was stored in DB
			if tc.validateDB {
				// Get the household ID from response
				var household model.Household
				err = json.NewDecoder(resp.Body).Decode(&household)
				assert.NoError(t, err)
				assert.NotEmpty(t, household.Id)

				// Verify household exists in DB
				ctx := context.Background()
				stored, err := suite.db.GetHouseholdByID(ctx, household.Id)
				assert.NoError(t, err)
				assert.Equal(t, household.Head.FirstName, stored.Head.FirstName)
				assert.Equal(t, household.Head.LastName, stored.Head.LastName)
				assert.Equal(t, household.Head.DOB, stored.Head.DOB)
			}
		})
	}
}

func TestPersonsHandlerSuite(t *testing.T) {
	suite.Run(t, new(PersonsHandlerTestSuite))
}
