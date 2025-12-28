package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/damonto/sigmo/internal/app/auth"
	"github.com/damonto/sigmo/internal/app/handler"
)

const bearerPrefix = "Bearer "

func Auth(store *auth.Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(header, bearerPrefix) {
				return c.JSON(http.StatusUnauthorized, handler.HTTPError{
					Code:    http.StatusUnauthorized,
					Message: "missing or invalid token",
				})
			}
			token := strings.TrimSpace(strings.TrimPrefix(header, bearerPrefix))
			if !store.ValidateToken(token) {
				return c.JSON(http.StatusUnauthorized, handler.HTTPError{
					Code:    http.StatusUnauthorized,
					Message: "missing or invalid token",
				})
			}
			return next(c)
		}
	}
}
