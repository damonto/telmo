package notification

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	sgp22 "github.com/damonto/euicc-go/v2"
	"github.com/labstack/echo/v4"

	"github.com/damonto/sigmo/internal/app/handler"
	"github.com/damonto/sigmo/internal/pkg/config"
	"github.com/damonto/sigmo/internal/pkg/lpa"
	mmodem "github.com/damonto/sigmo/internal/pkg/modem"
)

type Handler struct {
	handler.Handler
	manager *mmodem.Manager
	service *Service
}

func New(cfg *config.Config, manager *mmodem.Manager) *Handler {
	return &Handler{
		manager: manager,
		service: NewService(cfg),
	}
}

func (h *Handler) List(c echo.Context) error {
	modem, err := h.FindModem(h.manager, c.Param("id"))
	if err != nil {
		return h.NotFound(c, err)
	}
	response, err := h.service.List(modem)
	if err != nil {
		if errors.Is(err, lpa.ErrNoSupportedAID) {
			return h.NotFound(c, err)
		}
		return h.InternalServerError(c, err)
	}
	return h.Respond(c, response)
}

func (h *Handler) Resend(c echo.Context) error {
	modem, err := h.FindModem(h.manager, c.Param("id"))
	if err != nil {
		return h.NotFound(c, err)
	}
	sequence, err := sequenceFromParam(c)
	if err != nil {
		return h.BadRequest(c, err)
	}
	if err := h.service.Resend(modem, sequence); err != nil {
		if errors.Is(err, errNotificationNotFound) {
			return h.NotFound(c, err)
		}
		if errors.Is(err, lpa.ErrNoSupportedAID) {
			return h.NotFound(c, err)
		}
		return h.InternalServerError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) Delete(c echo.Context) error {
	modem, err := h.FindModem(h.manager, c.Param("id"))
	if err != nil {
		return h.NotFound(c, err)
	}
	sequence, err := sequenceFromParam(c)
	if err != nil {
		return h.BadRequest(c, err)
	}
	if err := h.service.Delete(modem, sequence); err != nil {
		if errors.Is(err, errNotificationNotFound) {
			return h.NotFound(c, err)
		}
		if errors.Is(err, lpa.ErrNoSupportedAID) {
			return h.NotFound(c, err)
		}
		return h.InternalServerError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func sequenceFromParam(c echo.Context) (sgp22.SequenceNumber, error) {
	raw := strings.TrimSpace(c.Param("sequence"))
	if raw == "" {
		return 0, errors.New("sequence number is required")
	}
	value, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid sequence number %q: %w", raw, err)
	}
	return sgp22.SequenceNumber(value), nil
}
