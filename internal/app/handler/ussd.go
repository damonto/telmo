package handler

import (
	"github.com/damonto/telmo/internal/app/state"
	"github.com/damonto/telmo/internal/pkg/modem"
	"github.com/damonto/telmo/internal/pkg/util"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

type USSDHandler struct {
	Handler
	state *state.StateManager
}

type USSDValue struct {
	Command string
	Modem   *modem.Modem
}

const USSDActionRespond state.State = "ussd_respond"

func NewUSSDHandler(state *state.StateManager) state.Handler {
	return &USSDHandler{
		state: state,
	}
}

func (h *USSDHandler) Handle() th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		m := h.Modem(ctx)
		s, err := m.USSDState()
		if err != nil {
			return err
		}
		if s != modem.Modem3gppUssdSessionStateIdle {
			if err := m.CancelUSSD(); err != nil {
				return err
			}
		}
		h.state.Enter(update.Message.Chat.ID, &state.ChatState{
			Handler: h,
			Value:   &USSDValue{Modem: m},
		})
		_, err = h.Reply(ctx, update, util.EscapeText("Okay, Send me the USSD command you want execute."), nil)
		return err
	}
}

func (h *USSDHandler) HandleMessage(ctx *th.Context, message telego.Message, s *state.ChatState) error {
	if s.State != USSDActionRespond {
		return h.initiate(ctx, message, s)
	}
	if s.State == USSDActionRespond {
		return h.respond(ctx, message, s)
	}
	return nil
}

func (h *USSDHandler) respond(ctx *th.Context, message telego.Message, s *state.ChatState) error {
	m := s.Value.(*USSDValue).Modem
	response, err := m.RespondUSSD(message.Text)
	if err != nil {
		h.ReplyMessage(ctx, message, util.EscapeText(err.Error()), nil)
		return err
	}
	_, err = h.ReplyMessage(ctx, message, util.EscapeText(response), nil)
	return err
}

func (h *USSDHandler) initiate(ctx *th.Context, message telego.Message, s *state.ChatState) error {
	m := s.Value.(*USSDValue).Modem
	response, err := m.InitiateUSSD(message.Text)
	if err != nil {
		h.ReplyMessage(ctx, message, util.EscapeText(err.Error()), nil)
		return err
	}
	h.state.Current(message.Chat.ID, USSDActionRespond)
	_, err = h.ReplyMessage(ctx, message, util.EscapeText(response), nil)
	return err
}

func (h *USSDHandler) HandleCallbackQuery(ctx *th.Context, query telego.CallbackQuery, s *state.ChatState) error {
	return nil
}
