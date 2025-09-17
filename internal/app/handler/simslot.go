package handler

import (
	"fmt"
	"strconv"

	"github.com/damonto/telmo/internal/app/state"
	"github.com/damonto/telmo/internal/pkg/modem"
	"github.com/damonto/telmo/internal/pkg/util"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type SIMSlotHandler struct {
	Handler
	state *state.StateManager
}

type SIMValue struct {
	Modem *modem.Modem
}

const CallbackQuerySIMSlotPrefix = "simslot"

func NewSIMSlotHandler(s *state.StateManager) state.Handler {
	return &SIMSlotHandler{
		state: s,
	}
}

func (h *SIMSlotHandler) Handle() th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		var message string
		var buttons [][]telego.InlineKeyboardButton
		modem := h.Modem(ctx)
		h.state.Enter(update.Message.Chat.ID, &state.ChatState{Handler: h, Value: &SIMValue{Modem: modem}})
		for idx, slot := range modem.SimSlots {
			sim, err := modem.SIM(slot)
			if err != nil {
				return err
			}
			button, text := h.message(idx+1, sim)
			message += text
			buttons = append(buttons, button)
		}
		_, err := h.Reply(ctx, update, message, func(message *telego.SendMessageParams) error {
			message.ReplyMarkup = tu.InlineKeyboard(buttons...)
			return nil
		})
		return err
	}
}

func (h *SIMSlotHandler) message(slot int, sim *modem.SIM) ([]telego.InlineKeyboardButton, string) {
	message := fmt.Sprintf(`
*\[Slot %d\]* %s
Operator: %s
IMSI: %s
ICCID: %s
`,
		slot,
		util.If(sim.Active, "🟢", "🔴"),
		util.EscapeText(util.If(sim.OperatorName != "", sim.OperatorName, util.LookupCarrier(sim.OperatorIdentifier))),
		sim.Imsi,
		sim.Identifier,
	)
	return tu.InlineKeyboardRow(telego.InlineKeyboardButton{
		Text:         fmt.Sprintf("%s [Slot %d] %s", util.If(sim.Active, "🟢", "🔴"), slot, sim.Identifier),
		CallbackData: fmt.Sprintf("%s:%d", CallbackQuerySIMSlotPrefix, slot),
	}), message
}

func (h *SIMSlotHandler) HandleCallbackQuery(ctx *th.Context, query telego.CallbackQuery, s *state.ChatState) error {
	defer h.state.Exit(query.From.ID)
	v, err := strconv.Atoi(query.Data[len(CallbackQuerySIMSlotPrefix)+1:])
	if err != nil {
		return err
	}
	if err := s.Value.(*SIMValue).Modem.SetPrimarySimSlot(uint32(v)); err != nil {
		return err
	}
	_, err = h.ReplyCallbackQuery(ctx, query, util.EscapeText(fmt.Sprintf("Primary SIM slot set to %d", v)), nil)
	return err
}

func (h *SIMSlotHandler) HandleMessage(ctx *th.Context, message telego.Message, s *state.ChatState) error {
	return nil
}
