package ussd

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/damonto/sigmo/internal/app/handler"
	mmodem "github.com/damonto/sigmo/internal/pkg/modem"
)

type Handler struct {
	handler.Handler
	manager *mmodem.Manager
	service *Service
}

const executeTimeout = time.Minute

var errExecuteTimeout = errors.New("ussd request timed out, please retry")

func New(manager *mmodem.Manager) *Handler {
	return &Handler{
		manager: manager,
		service: NewService(),
	}
}

func (h *Handler) Execute(c echo.Context) error {
	modem, err := h.FindModem(h.manager, c.Param("id"))
	if err != nil {
		return h.NotFound(c, err)
	}
	var req ExecuteRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), executeTimeout)
	defer cancel()

	response, err := h.service.Execute(ctx, modem, req.Action, req.Code)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return h.Error(c, http.StatusRequestTimeout, errExecuteTimeout)
		}
		if errors.Is(err, context.Canceled) {
			return nil
		}
		if errors.Is(err, errInvalidAction) || errors.Is(err, errSessionNotReady) || errors.Is(err, errUnknownSessionStatus) {
			return h.BadRequest(c, err)
		}
		return h.InternalServerError(c, err)
	}
	return h.Respond(c, response)
}
