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
	"golang.org/x/crypto/bcrypt"
)

type AuthHandlerTestSuite struct {
	suite.Suite
	server *httptest.Server
	db     *db.FirestoreDB
}

func (suite *AuthHandlerTestSuite) SetupSuite() {
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
	handler := &AuthHandler{suite.db, mockEmailSender}

	// Create the Echo server
	e := echo.New()
	e.POST("/login", handler.Login)

	// Start the test server
	suite.server = httptest.NewServer(e)
}

func (suite *AuthHandlerTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *AuthHandlerTestSuite) TestLogin() {
	// Existing test code...
}

func (suite *AuthHandlerTestSuite) TestResetPassword() {
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
		suite.FailNow("Failed to hash password", err)
	}

	// Create the test person with the hashed password
	testPerson := model.Person{
		PersonCommon: testPersonInput.PersonCommon,
		PasswordHash: string(hashedPassword),
	}

	// Save the test person to the mock Firestore
	ctx := context.Background()
	if err := suite.db.PutPerson(ctx, testPerson); err != nil {
		suite.FailNow("Failed to save test person", err)
	}

	// Create a test ResetPassword entity
	testResetPassword := model.ResetPassword{
		Id:       "testResetPasswordID",
		PersonId: testPerson.Id,
	}

	// Save the test ResetPassword entity to the mock Firestore
	if err := suite.db.PutResetPassword(ctx, testResetPassword); err != nil {
		suite.FailNow("Failed to save test ResetPassword", err)
	}

	// Create a valid reset password request
	resetPasswordRequest := `{"resetPasswordId": "testResetPasswordID", "newPassword": "newPassword123"}`
	resp, err := http.Post(suite.server.URL+"/reset-password", "application/json", strings.NewReader(resetPasswordRequest))
	if err != nil {
		suite.FailNow("Failed to make reset password request", err)
	}
	defer resp.Body.Close()

	// Verify the response status code
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Verify the response body contains the success message
	var responseBody map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		suite.FailNow("Failed to decode response body", err)
	}
	assert.Equal(suite.T(), "Password reset successfully", responseBody["message"])

	// Verify the Person's password has been updated
	updatedPerson, err := suite.db.GetPerson(ctx, testPerson.Id)
	if err != nil {
		suite.FailNow("Failed to get updated person", err)
	}
	assert.NoError(suite.T(), bcrypt.CompareHashAndPassword([]byte(updatedPerson.PasswordHash), []byte("newPassword123")))

	// Verify the ResetPassword entity has been deleted
	_, err = suite.db.GetResetPassword(ctx, testResetPassword.Id)
	assert.Error(suite.T(), err)
}
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
		suite.FailNow("Failed to hash password", err)
	}

	// Create the test person with the hashed password
	testPerson := model.Person{
		PersonCommon: testPersonInput.PersonCommon,
		PasswordHash: string(hashedPassword),
	}

	// Save the test person to the mock Firestore
	ctx := context.Background()
	if err := suite.db.PutPerson(ctx, testPerson); err != nil {
		suite.FailNow("Failed to save test person", err)
	}

	// Create a valid login request
	loginRequest := `{"email": "test@example.com", "password": "password123"}`
	resp, err := http.Post(suite.server.URL+"/login", "application/json", strings.NewReader(loginRequest))
	if err != nil {
		suite.FailNow("Failed to make login request", err)
	}
	defer resp.Body.Close()

	// Verify the response status code
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Verify the response body contains the sessionToken
	var responseBody map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		suite.FailNow("Failed to decode response body", err)
	}
	assert.Contains(suite.T(), responseBody, "sessionToken")
}

func TestAuthHandlerSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerTestSuite))
}
