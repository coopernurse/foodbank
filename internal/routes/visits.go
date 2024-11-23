package routes

import (
	"net/http"

	"foodbank/internal/db"
	"foodbank/internal/email"

	"github.com/labstack/echo/v4"
)

type VisitsHandler struct {
	DB          *db.FirestoreDB
	EmailSender email.EmailSender
}

func NewVisitsHandler(dbInstance *db.FirestoreDB, emailSender email.EmailSender) *VisitsHandler {
	return &VisitsHandler{DB: dbInstance, EmailSender: emailSender}
}

func (h *VisitsHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/household/:id/visits", h.LoadHouseholdVisits)
	e.GET("/household/:id/visits/limits", h.ComputeItemLimits)
	e.POST("/household/:id/visit/:visitId/item", h.AddItemToVisit)
	e.POST("/visit", h.PutVisit)
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
