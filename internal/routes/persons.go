package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"cupboard/internal/db"
)

type PersonsHandler struct {
	DB *db.FirestoreDB
}

func (h *PersonsHandler) PutPerson(c echo.Context) error {
	return c.String(http.StatusOK, "Put Person")
}

func (h *PersonsHandler) SearchPersons(c echo.Context) error {
	return c.String(http.StatusOK, "Search Persons")
}

func (h *PersonsHandler) LoadHouseholdPersons(c echo.Context) error {
	return c.String(http.StatusOK, "Load Household Persons")
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
