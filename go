package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"cupboard/internal/db"
	"cupboard/internal/model"
)

type PersonsHandler struct {
	DB *db.FirestoreDB
}

func (h *PersonsHandler) ValidatePerson(c echo.Context) error {
	return c.String(http.StatusOK, "Validate Person")
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
package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"cupboard/internal/db"
	"cupboard/internal/model"
)

type FoodBanksHandler struct {
	DB *db.FirestoreDB
}

func (h *FoodBanksHandler) ValidateFoodBank(c echo.Context) error {
	return c.String(http.StatusOK, "Validate Food Bank")
}

func (h *FoodBanksHandler) LoadFoodBanks(c echo.Context) error {
	return c.String(http.StatusOK, "Load Food Banks")
}

func (h *FoodBanksHandler) AssignPersonPermissions(c echo.Context) error {
	return c.String(http.StatusOK, "Assign Person Permissions")
}
package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"cupboard/internal/db"
	"cupboard/internal/model"
)

type ItemsHandler struct {
	DB *db.FirestoreDB
}

func (h *ItemsHandler) ValidateItem(c echo.Context) error {
	return c.String(http.StatusOK, "Validate Item")
}
package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"cupboard/internal/db"
	"cupboard/internal/model"
)

type VisitsHandler struct {
	DB *db.FirestoreDB
}

func (h *VisitsHandler) LoadHouseholdVisits(c echo.Context) error {
	return c.String(http.StatusOK, "Load Household Visits")
}

func (h *VisitsHandler) ComputeItemLimits(c echo.Context) error {
	return c.String(http.StatusOK, "Compute Item Limits")
}

func (h *VisitsHandler) AddItemToVisit(c echo.Context) error {
	return c.String(http.StatusOK, "Add Item to Visit")
}

func (h *VisitsHandler) ValidateVisit(c echo.Context) error {
	return c.String(http.StatusOK, "Validate Visit")
}
package routes

import (
	"cupboard/internal/db"
)

func NewRoutes(dbInstance *db.FirestoreDB) (*PersonsHandler, *FoodBanksHandler, *ItemsHandler, *VisitsHandler) {
	return &PersonsHandler{DB: dbInstance}, &FoodBanksHandler{DB: dbInstance}, &ItemsHandler{DB: dbInstance}, &VisitsHandler{DB: dbInstance}
}
