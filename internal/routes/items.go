package routes

import (
	"net/http"

	"foodbank/internal/db"
	"foodbank/internal/model"

	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type ItemsHandler struct {
	DB *db.FirestoreDB
}

func (h *ItemsHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/item", h.PutItem)
}

func (h *ItemsHandler) PutItem(c echo.Context) error {
	var itemInput model.Item
	if err := c.Bind(&itemInput); err != nil {
		log.Error().Err(err).Msg("Invalid JSON format")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
	}

	// Validate the item input
	errors := itemInput.Validate()
	if errors.HasErrors() {
		log.Error().Msg("Validation errors")
		return c.JSON(http.StatusBadRequest, errors)
	}

	// Set ULID if Id is not set or not a valid ULID
	if itemInput.Id == "" || len(itemInput.Id) != 26 {
		itemInput.Id = ulid.Make().String()
	}

	// Save the item to the database
	if err := h.DB.PutItem(c.Request().Context(), itemInput); err != nil {
		log.Error().Err(err).Msg("Failed to save item")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, itemInput)
}
