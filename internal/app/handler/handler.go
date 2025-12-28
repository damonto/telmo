package handler

import (
	"fmt"
	"net/http"

	"github.com/damonto/sigmo/internal/pkg/modem"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

type DataResponse struct {
	Data any `json:"data"`
}

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (*Handler) Error(c echo.Context, code int, err error) error {
	return c.JSON(code, HTTPError{Code: code, Message: err.Error()})
}

func (*Handler) Respond(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, DataResponse{Data: data})
}

func (*Handler) BindAndValidate(c echo.Context, i any) error {
	if err := c.Bind(i); err != nil {
		return c.JSON(http.StatusBadRequest, HTTPError{Code: http.StatusBadRequest, Message: err.Error()})
	}
	if err := c.Validate(i); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, HTTPError{Code: http.StatusUnprocessableEntity, Message: err.Error()})
	}
	return nil
}

func (*Handler) BadRequest(c echo.Context, err error) error {
	return c.JSON(http.StatusBadRequest, HTTPError{Code: http.StatusBadRequest, Message: err.Error()})
}

func (*Handler) Unauthorized(c echo.Context, err error) error {
	return c.JSON(http.StatusUnauthorized, HTTPError{Code: http.StatusUnauthorized, Message: err.Error()})
}

func (*Handler) NotFound(c echo.Context, err error) error {
	return c.JSON(http.StatusNotFound, HTTPError{Code: http.StatusNotFound, Message: err.Error()})
}

func (*Handler) Conflict(c echo.Context, err error) error {
	return c.JSON(http.StatusConflict, HTTPError{Code: http.StatusConflict, Message: err.Error()})
}

func (*Handler) InternalServerError(c echo.Context, err error) error {
	return c.JSON(http.StatusInternalServerError, HTTPError{Code: http.StatusInternalServerError, Message: err.Error()})
}

func (h *Handler) FindModem(manager *modem.Manager, id string) (*modem.Modem, error) {
	modems, err := manager.Modems()
	if err != nil {
		return nil, fmt.Errorf("failed to list modems: %w", err)
	}
	for _, m := range modems {
		if m.EquipmentIdentifier == id {
			return m, nil
		}
	}
	return nil, fmt.Errorf("modem with ID %s not found", id)
}
