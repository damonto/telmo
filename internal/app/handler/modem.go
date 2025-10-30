package handler

import (
	"fmt"
	"log/slog"
	"strings"

	sgp22 "github.com/damonto/euicc-go/v2"
	"github.com/damonto/telmo/internal/pkg/config"
	"github.com/damonto/telmo/internal/pkg/lpa"
	"github.com/damonto/telmo/internal/pkg/modem"
	"github.com/damonto/telmo/internal/pkg/util"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

type ListModemHandler struct {
	Handler
	mm     *modem.Manager
	config *config.Config
}

const ModemMessageTemplate = `
*%s*
Firmware Revision: %s
Hardware Revision: %s
IMEI: %s
Network: %s
Operator: %s
Number: %s
Signal: %d%%
ICCID: %s
`

func NewListModemHandler(config *config.Config, mm *modem.Manager) *ListModemHandler {
	return &ListModemHandler{
		mm:     mm,
		config: config,
	}
}

func (h *ListModemHandler) Handle() th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		modems, err := h.mm.Modems()
		if err != nil {
			return err
		}
		if len(modems) == 0 {
			_, err := h.Reply(ctx, update, util.EscapeText("No modems were found."), nil)
			return err
		}
		var message string
		for _, m := range modems {
			message += h.message(m) + "\n"
		}
		message = strings.TrimSuffix(message, "\n")
		_, err = h.Reply(ctx, update, message, nil)
		return err
	}
}

func (h *ListModemHandler) message(m *modem.Modem) string {
	percent, _, _ := m.SignalQuality()
	code, _ := m.OperatorCode()
	state, _ := m.RegistrationState()
	var accessTech []string
	accessTechnologies, _ := m.AccessTechnologies()
	for _, at := range accessTechnologies {
		accessTech = append(accessTech, at.String())
	}
	name, ok := h.config.Alias[m.EquipmentIdentifier]
	modemName := util.If(ok, name, m.Model)
	message := fmt.Sprintf(ModemMessageTemplate,
		util.EscapeText(modemName),
		util.EscapeText(m.FirmwareRevision),
		util.EscapeText(m.HardwareRevision),
		m.EquipmentIdentifier,
		util.EscapeText(
			fmt.Sprintf("%s (%s - %s)", util.LookupCarrier(code), strings.Join(accessTech, ", "), state),
		),
		util.EscapeText(util.If(m.Sim.OperatorName != "", m.Sim.OperatorName, util.LookupCarrier(m.Sim.OperatorIdentifier))),
		util.EscapeText(m.Number),
		percent,
		m.Sim.Identifier)

	eid, profileName, err := h.euiccInfo(m, m.Sim.Identifier)
	if err != nil {
		slog.Warn("unable to get EID and profile name, maybe it's not an eUICC", "error", err)
		return message
	}
	message += fmt.Sprintf("Profile Name: %s\nEID: `%s`", util.EscapeText(profileName), eid)
	return message
}

func (h *ListModemHandler) euiccInfo(m *modem.Modem, iccid string) (eid string, profileName string, err error) {
	lpa, err := lpa.New(m, h.config)
	if err != nil {
		return "", "", err
	}
	defer lpa.Close()
	info, err := lpa.Info()
	if err != nil {
		return "", "", err
	}
	id, _ := sgp22.NewICCID(iccid)
	profiles, err := lpa.ListProfile(id, nil)
	if err != nil {
		return "", "", err
	}
	if len(profiles) == 0 {
		return info.EID, "", nil
	}
	return info.EID, util.If(profiles[0].ProfileNickname != "", profiles[0].ProfileNickname, profiles[0].ProfileName), nil
}
