package routes

import (
	"net/http"

	"cupboard/internal/db"

	"github.com/labstack/echo/v4"
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
