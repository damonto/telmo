package handler

import (
	"fmt"

	"github.com/damonto/telmo/internal/pkg/util"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

type StartHandler struct {
	*Handler
}

func NewStartHandler() *StartHandler {
	h := new(StartHandler)
	return h
}

func (h *StartHandler) Handle() th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		message := fmt.Sprintf(`
Hello, %s,
Welcome to Telmo\!
I'm here to help you manage your SMS and eSIM profiles\.
If you don't know your Telegram Chat ID, it's *%d*\.
`, util.EscapeText(update.Message.Chat.FirstName), update.Message.Chat.ID)
		_, err := h.Reply(ctx, update, message, nil)
		return err
	}
}
