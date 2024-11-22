package routes

import (
	"net/http"

	"cupboard/internal/db"
	"cupboard/internal/model"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type PersonsHandler struct {
	DB *db.FirestoreDB
}

func (h *PersonsHandler) PutPerson(c echo.Context) error {
	var personInput model.PersonInput
	if err := c.Bind(&personInput); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(personInput.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
	}

	// Create a Person struct from PersonInput
	person := model.Person{
		PersonCommon: personInput.PersonCommon,
		PasswordHash: string(hashedPassword),
	}

	// Save the person to the database
	if err := h.DB.PutPerson(c.Request().Context(), person); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, person)
}

func (h *PersonsHandler) SearchPersons(c echo.Context) error {
	ctx := c.Request().Context()

	// Fetch persons from the database
	persons, err := h.DB.GetPersons(ctx)
	if err != nil {
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
	return c.String(http.StatusOK, "Email Login Link")
}

func (h *PersonsHandler) ResolvePermissions(c echo.Context) error {
	return c.String(http.StatusOK, "Resolve Permissions")
}
