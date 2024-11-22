package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// AuthMiddleware checks for a valid secure cookie
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_token")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		// Validate the session token (this is a placeholder for actual validation logic)
		if cookie.Value != "valid_session_token" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		return next(c)
	}
}
