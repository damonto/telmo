package euicc

import (
	"errors"

	"github.com/labstack/echo/v4"

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
		service: NewService(cfg),
	}
}

func (h *Handler) Get(c echo.Context) error {
	modem, err := h.FindModem(h.manager, c.Param("id"))
	if err != nil {
		return h.NotFound(c, err)
	}
	response, err := h.service.Get(modem)
	if err != nil {
		if errors.Is(err, lpa.ErrNoSupportedAID) {
			return h.NotFound(c, err)
		}
		return h.InternalServerError(c, err)
	}
	return h.Respond(c, response)
}
