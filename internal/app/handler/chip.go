package handler

import (
	"fmt"
	"strings"

	"github.com/damonto/telmo/internal/pkg/config"
	"github.com/damonto/telmo/internal/pkg/lpa"
	"github.com/damonto/telmo/internal/pkg/util"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

type ChipHandler struct {
	Handler
	config *config.Config
}

func NewChipHandler(config *config.Config) *ChipHandler {
	return &ChipHandler{config: config}
}

func (h *ChipHandler) Handle() th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		lpa, err := h.LPA(ctx, h.config)
		if err != nil {
			return err
		}
		defer lpa.Close()
		info, err := lpa.Info()
		if err != nil {
			return err
		}
		_, err = h.Reply(ctx, update, h.message(info), nil)
		return err
	}
}

func (h *ChipHandler) message(info *lpa.Info) string {
	var keys string
	for _, k := range info.Certificates {
		keys += k + "\n"
	}
	return fmt.Sprintf(`
EID: %s
SAS\-UP: %s
Free Space: %d KiB
Signing Keys:
%s
`,
		fmt.Sprintf("`%s`", info.EID),
		util.EscapeText(info.SasAccreditationNumber),
		info.FreeSpace/1024,
		util.EscapeText(strings.TrimRight(keys, "\n")),
	)
}
