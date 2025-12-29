package modem

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/damonto/sigmo/internal/app/handler"
	"github.com/damonto/sigmo/internal/pkg/config"
	mmodem "github.com/damonto/sigmo/internal/pkg/modem"
)

type Handler struct {
	handler.Handler
	manager *mmodem.Manager
	service *Service
}

const (
	switchSimSlotTimeout = time.Minute
	updateMSISDNTimeout  = time.Minute
)

var (
	errSwitchSimSlotTimeout = errors.New("switching SIM slot timed out, please refresh to confirm the active slot")
	errUpdateMSISDNTimeout  = errors.New("updating MSISDN timed out, please refresh to confirm the active slot")
)

func New(cfg *config.Config, manager *mmodem.Manager) *Handler {
	return &Handler{
		manager: manager,
		service: NewService(cfg, manager),
	}
}

func (h *Handler) List(c echo.Context) error {
	response, err := h.service.List()
	if err != nil {
		return h.InternalServerError(c, err)
	}
	return h.Respond(c, response)
}

func (h *Handler) Get(c echo.Context) error {
	modem, err := h.FindModem(h.manager, c.Param("id"))
	if err != nil {
		return h.NotFound(c, err)
	}
	response, err := h.service.Get(modem)
	if err != nil {
		return h.InternalServerError(c, err)
	}
	return h.Respond(c, response)
}

func (h *Handler) SwitchSimSlot(c echo.Context) error {
	modem, err := h.FindModem(h.manager, c.Param("id"))
	if err != nil {
		return h.NotFound(c, err)
	}
	identifier := c.Param("identifier")
	if identifier == "" {
		return h.BadRequest(c, errSimIdentifierRequired)
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), switchSimSlotTimeout)
	defer cancel()

	if err := h.service.SwitchSimSlot(ctx, modem, identifier); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return h.Error(c, http.StatusRequestTimeout, errSwitchSimSlotTimeout)
		}
		if errors.Is(err, context.Canceled) {
			return nil
		}
		if errors.Is(err, errSimIdentifierRequired) || errors.Is(err, errSimSlotsUnavailable) || errors.Is(err, errSimSlotNotFound) || errors.Is(err, errSimSlotAlreadyActive) {
			return h.BadRequest(c, err)
		}
		return h.InternalServerError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) UpdateMSISDN(c echo.Context) error {
	modem, err := h.FindModem(h.manager, c.Param("id"))
	if err != nil {
		return h.NotFound(c, err)
	}
	var req UpdateMSISDNRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), updateMSISDNTimeout)
	defer cancel()

	if err := h.service.UpdateMSISDN(ctx, modem, req.Number); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return h.Error(c, http.StatusRequestTimeout, errUpdateMSISDNTimeout)
		}
		if errors.Is(err, errMSISDNInvalidNumber) {
			return h.BadRequest(c, err)
		}
		return h.InternalServerError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) UpdateSettings(c echo.Context) error {
	modem, err := h.FindModem(h.manager, c.Param("id"))
	if err != nil {
		return h.NotFound(c, err)
	}
	var req UpdateModemSettingsRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}
	if err := h.service.UpdateSettings(modem.EquipmentIdentifier, req); err != nil {
		if errors.Is(err, errCompatibleRequired) {
			return h.BadRequest(c, err)
		}
		return h.InternalServerError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) GetSettings(c echo.Context) error {
	modem, err := h.FindModem(h.manager, c.Param("id"))
	if err != nil {
		return h.NotFound(c, err)
	}
	response := h.service.GetSettings(modem.EquipmentIdentifier)
	return h.Respond(c, response)
}
