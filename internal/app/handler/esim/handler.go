package esim

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"

	sgp22 "github.com/damonto/euicc-go/v2"
	"github.com/damonto/sigmo/internal/app/handler"
	"github.com/damonto/sigmo/internal/pkg/config"
	"github.com/damonto/sigmo/internal/pkg/lpa"
	mmodem "github.com/damonto/sigmo/internal/pkg/modem"
)

type Handler struct {
	handler.Handler
	cfg     *config.Config
	manager *mmodem.Manager
	service *Service
}

func New(cfg *config.Config, manager *mmodem.Manager) *Handler {
	return &Handler{
		cfg:     cfg,
		manager: manager,
		service: NewService(cfg, manager),
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

func (h *Handler) Enable(c echo.Context) error {
	modem, err := h.FindModem(h.manager, c.Param("id"))
	if err != nil {
		return h.NotFound(c, err)
	}
	iccid, err := iccidFromParam(c)
	if err != nil {
		return h.BadRequest(c, err)
	}
	if err := h.service.Enable(modem, iccid); err != nil {
		if errors.Is(err, lpa.ErrNoSupportedAID) {
			return h.NotFound(c, err)
		}
		return h.InternalServerError(c, err)
	}
	return h.Respond(c, struct{}{})
}

func (h *Handler) Delete(c echo.Context) error {
	modem, err := h.FindModem(h.manager, c.Param("id"))
	if err != nil {
		return h.NotFound(c, err)
	}
	iccid, err := iccidFromParam(c)
	if err != nil {
		return h.BadRequest(c, err)
	}
	if err := h.service.Delete(modem, iccid); err != nil {
		if errors.Is(err, lpa.ErrNoSupportedAID) {
			return h.NotFound(c, err)
		}
		return h.InternalServerError(c, err)
	}
	return h.Respond(c, struct{}{})
}

func (h *Handler) UpdateNickname(c echo.Context) error {
	modem, err := h.FindModem(h.manager, c.Param("id"))
	if err != nil {
		return h.NotFound(c, err)
	}
	iccid, err := iccidFromParam(c)
	if err != nil {
		return h.BadRequest(c, err)
	}
	var req UpdateNicknameRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}
	if err := h.service.UpdateNickname(modem, iccid, req.Nickname); err != nil {
		if errors.Is(err, errInvalidNickname) {
			return h.BadRequest(c, err)
		}
		if errors.Is(err, lpa.ErrNoSupportedAID) {
			return h.NotFound(c, err)
		}
		return h.InternalServerError(c, err)
	}
	return h.Respond(c, struct{}{})
}

func iccidFromParam(c echo.Context) (sgp22.ICCID, error) {
	iccidParam := c.Param("iccid")
	if iccidParam == "" {
		return nil, errors.New("iccid is required")
	}
	iccid, err := sgp22.NewICCID(iccidParam)
	if err != nil {
		return nil, fmt.Errorf("invalid iccid %q: %w", iccidParam, err)
	}
	return iccid, nil
}
