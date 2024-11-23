package routes

import (
	"net/http"
	"time"

	"cupboard/internal/auth"
	"cupboard/internal/db"
	"cupboard/internal/email"
	"cupboard/internal/model"

	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB          *db.FirestoreDB
	EmailSender email.EmailSender
}

func NewAuthHandler(dbInstance *db.FirestoreDB, emailSender email.EmailSender) *AuthHandler {
	return &AuthHandler{DB: dbInstance, EmailSender: emailSender}
}

type SendResetPasswordEmailInput struct {
	Email string `json:"email"`
}

func (h *AuthHandler) SendResetPasswordEmail(c echo.Context) error {
	var input SendResetPasswordEmailInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
	}

	// Load the Person by email
	person, err := h.DB.GetPersonByEmail(c.Request().Context(), input.Email)
	if err != nil {
		log.Error().Err(err).Msg("Failed to load person by email")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to load person"})
	}
	if person == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Person not found"})
	}

	// Create a ResetPassword entity
	resetPassword := model.ResetPassword{
		Id:       ulid.Make().String(),
		PersonId: person.Id,
	}

	// Save the ResetPassword entity to Firestore
	if err := h.DB.PutResetPassword(c.Request().Context(), resetPassword); err != nil {
		log.Error().Err(err).Msg("Failed to save ResetPassword")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save ResetPassword"})
	}

	// Send an email with the reset link
	resetLink := "http://yourapp.com/reset-password?resetPasswordId=" + resetPassword.Id
	if err := h.EmailSender.SendEmail(c.Request().Context(), person.Email, "Password Reset", "Click here to reset your password: "+resetLink); err != nil {
		log.Error().Err(err).Msg("Failed to send reset password email")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send reset password email"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Reset password email sent"})
}

type ResetPasswordInput struct {
	ResetPasswordId string `json:"resetPasswordId"`
	NewPassword     string `json:"newPassword"`
}

func (h *AuthHandler) ResetPassword(c echo.Context) error {
	var input ResetPasswordInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
	}

	// Decode the time component of the ULID
	resetPasswordId, err := ulid.Parse(input.ResetPasswordId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid reset password ID"})
	}

	// Check if the reset password request is more than 12 hours old
	if time.Since(time.Unix(int64(resetPasswordId.Time()), 0)) > 12*time.Hour {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Reset password link expired"})
	}

	// Load the ResetPassword entity by resetPasswordId
	resetPassword, err := h.DB.GetResetPassword(c.Request().Context(), input.ResetPasswordId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to load ResetPassword")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to load ResetPassword"})
	}
	if resetPassword == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "ResetPassword not found"})
	}

	// Load the Person by ResetPassword.PersonId
	person, err := h.DB.GetPerson(c.Request().Context(), resetPassword.PersonId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to load person")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to load person"})
	}
	if person == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Person not found"})
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
	}

	// Update the Person's password
	person.PasswordHash = string(hashedPassword)

	// Save the updated Person to Firestore
	if err := h.DB.PutPerson(c.Request().Context(), *person); err != nil {
		log.Error().Err(err).Msg("Failed to save person")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save person"})
	}

	// Delete the ResetPassword entity
	if err := h.DB.DeleteResetPassword(c.Request().Context(), input.ResetPasswordId); err != nil {
		log.Error().Err(err).Msg("Failed to delete ResetPassword")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete ResetPassword"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Password reset successfully"})
}

// LoginInput defines the input for the login request
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login handles the login request
func (h *AuthHandler) Login(c echo.Context) error {
	var loginInput LoginInput
	if err := c.Bind(&loginInput); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
	}

	// Fetch the person from the database by email
	person, err := h.DB.GetPersonByEmail(c.Request().Context(), loginInput.Email)
	if err != nil {
		log.Error().Err(err).Msg("Failed to load person by email")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to load person"})
	}
	if person == nil || person.PasswordHash == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	// Compare the provided password with the stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(person.PasswordHash), []byte(loginInput.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	// Authentication successful
	sessionToken, err := auth.EncryptSessionToken(person.Id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to encrypt session token")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create session token"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Login successful", "sessionToken": sessionToken})
}
