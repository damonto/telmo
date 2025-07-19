package handler

import (
	"fmt"

	sgp22 "github.com/damonto/euicc-go/v2"
	"github.com/damonto/telmo/internal/pkg/util"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

type DiscoveryHandler struct {
	*Handler
}

func NewDiscoveryHandler() *DiscoveryHandler {
	h := new(DiscoveryHandler)
	return h
}

func (h *DiscoveryHandler) Handle() th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		lpa, err := h.LPA(ctx)
		if err != nil {
			return err
		}
		defer lpa.Close()
		imei, _ := sgp22.NewIMEI(h.Modem(ctx).EquipmentIdentifier)
		entries, err := lpa.Discover(imei)
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			_, err := h.Reply(ctx, update, "No profiles found", nil)
			return err
		}
		message := "You have the following profiles available for download:\n\n"
		for _, entry := range entries {
			message += fmt.Sprintf("*%s*\n`/download LPA:1$%s$`\n\n", util.EscapeText(entry.EventID), util.EscapeText(entry.Address))
		}
		message += "You can copy the download command and paste it into the chat to download the profile\\."
		_, err = h.Reply(ctx, update, message, func(m *telego.SendMessageParams) error {
			m.WithParseMode(telego.ModeMarkdownV2)
			return nil
		})
		return err
	}
}
