package message

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

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
	response, err := h.service.ListConversations(modem)
	if err != nil {
		return h.InternalServerError(c, err)
	}
	return h.Respond(c, response)
}

func (h *Handler) ListByParticipant(c echo.Context) error {
	modem, err := h.FindModem(h.manager, c.Param("id"))
	if err != nil {
		return h.NotFound(c, err)
	}
	participant, err := participantFromParam(c)
	if err != nil {
		return h.BadRequest(c, err)
	}
	response, err := h.service.ListByParticipant(modem, participant)
	if err != nil {
		if errors.Is(err, errParticipantRequired) {
			return h.BadRequest(c, err)
		}
		return h.InternalServerError(c, err)
	}
	return h.Respond(c, response)
}

func (h *Handler) Send(c echo.Context) error {
	modem, err := h.FindModem(h.manager, c.Param("id"))
	if err != nil {
		return h.NotFound(c, err)
	}
	var req SendMessageRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}
	if err := h.service.Send(modem, req.To, req.Text); err != nil {
		if errors.Is(err, errRecipientRequired) || errors.Is(err, errTextRequired) {
			return h.BadRequest(c, err)
		}
		return h.InternalServerError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) DeleteByParticipant(c echo.Context) error {
	modem, err := h.FindModem(h.manager, c.Param("id"))
	if err != nil {
		return h.NotFound(c, err)
	}
	participant, err := participantFromParam(c)
	if err != nil {
		return h.BadRequest(c, err)
	}
	if err := h.service.DeleteByParticipant(modem, participant); err != nil {
		if errors.Is(err, errParticipantRequired) {
			return h.BadRequest(c, err)
		}
		return h.InternalServerError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func participantFromParam(c echo.Context) (string, error) {
	raw := c.Param("participant")
	if raw == "" {
		return "", errParticipantRequired
	}
	participant, err := url.PathUnescape(raw)
	if err != nil {
		return "", fmt.Errorf("invalid participant %q: %w", raw, err)
	}
	return participant, nil
}
