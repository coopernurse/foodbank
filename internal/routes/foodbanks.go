package routes

import (
	"net/http"

	"cupboard/internal/db"
	"cupboard/internal/model"

	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type FoodBanksHandler struct {
	DB *db.FirestoreDB
}

func (h *FoodBanksHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/foodbank", h.PutFoodBank)
	e.GET("/foodbanks", h.LoadFoodBanks)
	e.POST("/foodbank/:id/assign-permissions", h.AssignPersonPermissions)
}

func (h *FoodBanksHandler) PutFoodBank(c echo.Context) error {
	var foodBankInput model.FoodBank
	if err := c.Bind(&foodBankInput); err != nil {
		log.Error().Err(err).Msg("Invalid JSON format")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
	}

	// Validate the food bank input
	errors := foodBankInput.Validate()
	if errors.HasErrors() {
		log.Error().Msg("Validation errors")
		return c.JSON(http.StatusBadRequest, errors)
	}

	// Set ULID if Id is not set or not a valid ULID
	if foodBankInput.Id == "" || len(foodBankInput.Id) != 26 {
		foodBankInput.Id = ulid.Make().String()
	}

	// Save the food bank to the database
	if err := h.DB.PutFoodBank(c.Request().Context(), foodBankInput); err != nil {
		log.Error().Err(err).Msg("Failed to save food bank")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, foodBankInput)
}

func (h *FoodBanksHandler) LoadFoodBanks(c echo.Context) error {
	return c.String(http.StatusOK, "Load Food Banks")
}

func (h *FoodBanksHandler) AssignPersonPermissions(c echo.Context) error {
	return c.String(http.StatusOK, "Assign Person Permissions")
}
