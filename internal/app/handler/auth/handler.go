package auth

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/damonto/sigmo/internal/app/auth"
	"github.com/damonto/sigmo/internal/app/handler"
	"github.com/damonto/sigmo/internal/pkg/config"
)

type Handler struct {
	handler.Handler
	service *Service
}

func New(cfg *config.Config, store *auth.Store) *Handler {
	return &Handler{
		service: NewService(cfg, store),
	}
}

func (h *Handler) SendOTP(c echo.Context) error {
	if err := h.service.SendOTP(); err != nil {
		if errors.Is(err, auth.ErrOTPCooldown) {
			return h.Error(c, http.StatusTooManyRequests, err)
		}
		return h.InternalServerError(c, err)
	}
	return c.NoContent(http.StatusCreated)
}

func (h *Handler) VerifyOTP(c echo.Context) error {
	var req VerifyOTPRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}
	token, err := h.service.VerifyOTP(req.Code)
	if err != nil {
		if errors.Is(err, errInvalidOTP) {
			return h.Unauthorized(c, err)
		}
		return h.InternalServerError(c, err)
	}
	return h.Respond(c, VerifyOTPResponse{Token: token})
}
