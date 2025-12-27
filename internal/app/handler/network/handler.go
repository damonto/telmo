package network

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/damonto/sigmo/internal/app/handler"
	mmodem "github.com/damonto/sigmo/internal/pkg/modem"
)

type Handler struct {
	handler.Handler
	manager *mmodem.Manager
	service *Service
}

func New(manager *mmodem.Manager) *Handler {
	return &Handler{
		manager: manager,
		service: NewService(),
	}
}

func (h *Handler) List(c echo.Context) error {
	modem, err := h.FindModem(h.manager, c.Param("id"))
	if err != nil {
		return h.NotFound(c, err)
	}
	response, err := h.service.List(modem)
	if err != nil {
		return h.InternalServerError(c, err)
	}
	return h.Respond(c, response)
}

func (h *Handler) Register(c echo.Context) error {
	modem, err := h.FindModem(h.manager, c.Param("id"))
	if err != nil {
		return h.NotFound(c, err)
	}
	operatorCode := c.Param("operatorCode")
	if err := h.service.Register(modem, operatorCode); err != nil {
		if errors.Is(err, errOperatorCodeRequired) {
			return h.BadRequest(c, err)
		}
		return h.InternalServerError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}
