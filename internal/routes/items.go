package routes

import (
	"net/http"

	"cupboard/internal/db"

	"github.com/labstack/echo/v4"
)

type ItemsHandler struct {
	DB *db.FirestoreDB
}

func (h *ItemsHandler) ValidateItem(c echo.Context) error {
	return c.String(http.StatusOK, "Validate Item")
}
