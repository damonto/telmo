package handler

import (
	"fmt"

	"github.com/damonto/telmo/internal/app/state"
	"github.com/damonto/telmo/internal/pkg/modem"
	"github.com/damonto/telmo/internal/pkg/util"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type RegisterNetworkHandler struct {
	*Handler
}

type RegisterNetworkValue struct {
	Modem *modem.Modem
}

const RegisterNetworkActionCallbackDataPrefix = "register_network"

func NewRegisterNetworkHandler() *RegisterNetworkHandler {
	h := new(RegisterNetworkHandler)
	return h
}

func (h *RegisterNetworkHandler) Handle() th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		m := h.Modem(ctx)
		state.M.Enter(update.Message.Chat.ID, &state.ChatState{
			Handler: h,
			Value:   &RegisterNetworkValue{Modem: m},
		})
		message, err := h.Reply(ctx, update, util.EscapeText("Just a moment, I'm scanning the available networks."), nil)
		if err != nil {
			return err
		}
		networks, err := m.ScanNetworks()
		if err != nil {
			return err
		}
		if len(networks) == 0 {
			h.ReplyMessage(ctx, *message, util.EscapeText("No networks found."), nil)
			return nil
		}
		return h.sendAvailableNetworks(ctx, *message, networks)
	}
}

func (h *RegisterNetworkHandler) sendAvailableNetworks(ctx *th.Context, message telego.Message, networks []*modem.ThreeGPPNetwork) error {
	var buttons [][]telego.InlineKeyboardButton
	for _, network := range networks {
		var accessTechnology string
		for _, tech := range network.AccessTechnology {
			accessTechnology += tech.String() + ", "
		}
		buttons = append(buttons, tu.InlineKeyboardRow(telego.InlineKeyboardButton{
			Text:         fmt.Sprintf("%s - %s (%s)", network.OperatorCode, network.OperatorName, accessTechnology[:len(accessTechnology)-2]),
			CallbackData: fmt.Sprintf("%s:%s", RegisterNetworkActionCallbackDataPrefix, network.OperatorCode),
		}))
	}
	_, err := h.ReplyMessage(ctx, message, util.EscapeText("Please select a network."), func(message *telego.SendMessageParams) error {
		message.WithReplyMarkup(tu.InlineKeyboard(buttons...))
		return nil
	})
	return err
}

func (h *RegisterNetworkHandler) HandleCallbackQuery(ctx *th.Context, query telego.CallbackQuery, s *state.ChatState) error {
	defer state.M.Exit(query.From.ID)
	value := s.Value.(*RegisterNetworkValue)
	ctx.Bot().DeleteMessage(ctx, &telego.DeleteMessageParams{
		ChatID:    tu.ID(query.From.ID),
		MessageID: query.Message.GetMessageID(),
	})
	message, err := ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
		ChatID: tu.ID(query.From.ID),
		Text:   fmt.Sprintf("Registering network to %s...", query.Data[len(RegisterNetworkActionCallbackDataPrefix)+1:]),
	})
	if err != nil {
		return err
	}
	params := &telego.EditMessageTextParams{
		ChatID:    tu.ID(query.From.ID),
		MessageID: message.MessageID,
	}
	if err := value.Modem.RegisterNetwork(query.Data[len(RegisterNetworkActionCallbackDataPrefix)+1:]); err != nil {
		params.Text = fmt.Sprintf("Failed to register network: %s", err.Error())
		_, err := ctx.Bot().EditMessageText(ctx, params)
		return err
	}
	params.Text = "Network registered successfully."
	_, err = ctx.Bot().EditMessageText(ctx, params)
	return err
}

func (h *RegisterNetworkHandler) HandleMessage(ctx *th.Context, message telego.Message, s *state.ChatState) error {
	return nil
}
