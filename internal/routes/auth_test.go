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
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandlerTestSuite struct {
	suite.Suite
	server      *httptest.Server
	db          *db.FirestoreDB
	emailSender *email.MockEmailSender
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
	handler := &AuthHandler{suite.db, mockEmailSender, "http://localhost:8080"}

	// Create the Echo server
	e := echo.New()
	handler.RegisterRoutes(e)

	// Start the test server
	suite.server = httptest.NewServer(e)
	suite.emailSender = mockEmailSender
}

func (suite *AuthHandlerTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *AuthHandlerTestSuite) TestLogin() {
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

func (suite *AuthHandlerTestSuite) TestSendResetPasswordEmail() {
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

	// Create a valid send reset password email request
	sendResetPasswordEmailRequest := `{"email": "test@example.com"}`
	resp, err := http.Post(suite.server.URL+"/send-password-reset-email", "application/json", strings.NewReader(sendResetPasswordEmailRequest))
	if err != nil {
		suite.FailNow("Failed to make send reset password email request", err)
	}
	defer resp.Body.Close()

	// Verify the response status code
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Verify the response body contains the resetPasswordId
	var responseBody map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		suite.FailNow("Failed to decode response body", err)
	}
	assert.Contains(suite.T(), responseBody, "resetPasswordId")

	// Verify that a ResetPassword entity was created with a ULID as the ID
	resetPasswordId := responseBody["resetPasswordId"]
	resetPassword, err := suite.db.GetResetPassword(ctx, resetPasswordId)
	if err != nil {
		suite.FailNow("Failed to get ResetPassword entities", err)
	}

	parsedId, err := ulid.Parse(resetPassword.Id)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), parsedId.String())

	// Verify that the MockEmailSender has a SentEmail sent to the Person's email address
	assert.Len(suite.T(), suite.emailSender.SentEmails, 1)
	assert.Equal(suite.T(), testPerson.Email, suite.emailSender.SentEmails[0].To)
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

	// Generate a ULID for the resetPasswordId
	resetPasswordId := ulid.Make().String()

	// Create a test ResetPassword entity
	testResetPassword := model.ResetPassword{
		Id:       resetPasswordId,
		PersonId: testPerson.Id,
	}

	// Save the test ResetPassword entity to the mock Firestore
	if err := suite.db.PutResetPassword(ctx, testResetPassword); err != nil {
		suite.FailNow("Failed to save test ResetPassword", err)
	}

	// Create a valid reset password request
	resetPasswordRequest := `{"resetPasswordId": "` + resetPasswordId + `", "newPassword": "newPassword123"}`
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

	// Create a valid login request with the new password
	loginRequest := `{"email": "test@example.com", "password": "newPassword123"}`
	resp, err = http.Post(suite.server.URL+"/login", "application/json", strings.NewReader(loginRequest))
	if err != nil {
		suite.FailNow("Failed to make login request", err)
	}
	defer resp.Body.Close()

	// Verify the response status code
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Verify the response body contains the sessionToken
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		suite.FailNow("Failed to decode response body", err)
	}
	assert.Contains(suite.T(), responseBody, "sessionToken")
}

func TestAuthHandlerSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerTestSuite))
}
