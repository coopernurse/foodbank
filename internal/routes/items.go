package routes

import (
	"net/http"

	"cupboard/internal/db"

	"github.com/labstack/echo/v4"
)

type ItemsHandler struct {
	DB *db.FirestoreDB
}

func (h *ItemsHandler) PutItem(c echo.Context) error {
	return c.String(http.StatusOK, "Put Item")
}
