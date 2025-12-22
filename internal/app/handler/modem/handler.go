package modem

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/damonto/sigmo/internal/pkg/config"
	mmodem "github.com/damonto/sigmo/internal/pkg/modem"
)

type Handler struct {
	service *Service
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewHandler(cfg *config.Config, manager *mmodem.Manager) *Handler {
	return &Handler{
		service: NewService(cfg, manager),
	}
}

func (h *Handler) ListInserted(c echo.Context) error {
	response, err := h.service.ListInserted()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, response)
}
