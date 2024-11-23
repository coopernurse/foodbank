package routes

import (
	"net/http"
	"strings"

	"foodbank/internal/db"
	"foodbank/internal/email"
	"foodbank/internal/model"

	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type PersonsHandler struct {
	DB          *db.FirestoreDB
	EmailSender email.EmailSender
}

func NewPersonsHandler(dbInstance *db.FirestoreDB, emailSender email.EmailSender) *PersonsHandler {
	return &PersonsHandler{DB: dbInstance, EmailSender: emailSender}
}

func (h *PersonsHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/person", h.PutPerson)
	e.GET("/persons/search", h.SearchPersons)
	e.GET("/household/:id/persons", h.LoadHouseholdPersons)
	e.POST("/person/:id/reset-password", h.ResetPassword)
	e.POST("/person/:id/email-login-link", h.EmailLoginLink)
	e.GET("/person/:id/permissions", h.ResolvePermissions)
}

func (h *PersonsHandler) PutPerson(c echo.Context) error {
	var personInput model.PersonInput
	if err := c.Bind(&personInput); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
	}

	// Validate the person input
	errors := personInput.Validate()
	if errors.HasErrors() {
		return c.JSON(http.StatusBadRequest, errors)
	}

	personInput.Email = strings.ToLower(personInput.Email)

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(personInput.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
	}

	// Create a Person struct from PersonInput
	person := model.Person{
		PersonCommon: personInput.PersonCommon,
		PasswordHash: string(hashedPassword),
	}

	// Set ULID if Id is not set or not a valid ULID
	if person.Id == "" || len(person.Id) != 26 {
		person.Id = ulid.Make().String()
	}

	// Save the person to the database
	if err := h.DB.PutPerson(c.Request().Context(), person); err != nil {
		log.Error().Err(err).Msg("Failed to save person")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, person)
}

func (h *PersonsHandler) SearchPersons(c echo.Context) error {
	ctx := c.Request().Context()

	// Fetch persons from the database
	persons, err := h.DB.GetPersons(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve persons")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Convert persons to PersonOutput
	var personOutputs []model.PersonOutput
	for _, person := range persons {
		personOutputs = append(personOutputs, model.PersonOutput{
			PersonCommon: person.PersonCommon,
		})
	}

	return c.JSON(http.StatusOK, personOutputs)
}

func (h *PersonsHandler) LoadHouseholdPersons(c echo.Context) error {
	ctx := c.Request().Context()
	householdID := c.Param("id")

	// Fetch persons from the database for the given household ID
	persons, err := h.DB.GetHouseholdPersons(ctx, householdID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve household persons")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Convert persons to PersonOutput
	var personOutputs []model.PersonOutput
	for _, person := range persons {
		personOutputs = append(personOutputs, model.PersonOutput{
			PersonCommon: person.PersonCommon,
		})
	}

	return c.JSON(http.StatusOK, personOutputs)
}

func (h *PersonsHandler) ResetPassword(c echo.Context) error {
	return c.String(http.StatusOK, "Reset Password")
}

func (h *PersonsHandler) EmailLoginLink(c echo.Context) error {
	// Use h.EmailSender to send email
	return c.String(http.StatusOK, "Email Login Link")
}

func (h *PersonsHandler) ResolvePermissions(c echo.Context) error {
	return c.String(http.StatusOK, "Resolve Permissions")
}
