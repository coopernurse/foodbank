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
