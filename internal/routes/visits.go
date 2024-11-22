package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"cupboard/internal/db"
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

func (h *VisitsHandler) PutVisit(c echo.Context) error {
	return c.String(http.StatusOK, "Put Visit")
}
