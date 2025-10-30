package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"

	"github.com/godbus/dbus/v5"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"

	"github.com/damonto/telmo/internal/pkg/config"
	"github.com/damonto/telmo/internal/pkg/lpa"
	"github.com/damonto/telmo/internal/pkg/modem"
	"github.com/damonto/telmo/internal/pkg/util"
)

const CallbackQueryAskModemPrefix = "ask_modem"

type ModemRequiredMiddleware struct {
	mm     *modem.Manager
	modem  chan *modem.Modem
	config *config.Config
}

type session struct {
	ch     chan *modem.Modem
	ctx    *th.Context
	update telego.Update
}

var sessions sync.Map

func NewModemRequiredMiddleware(mm *modem.Manager, handler *th.BotHandler, config *config.Config) *ModemRequiredMiddleware {
	m := &ModemRequiredMiddleware{
		mm:     mm,
		modem:  make(chan *modem.Modem, 1),
		config: config,
	}
	handler.HandleCallbackQuery(m.HandleModemSelectionCallbackQuery, th.CallbackDataPrefix(CallbackQueryAskModemPrefix))
	return m
}

func (m *ModemRequiredMiddleware) Middleware(eUICCRequired bool) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		modems, err := m.mm.Modems()
		if err != nil {
			return err
		}
		if len(modems) == 0 {
			return m.sendErrorModemNotFound(ctx, update)
		}
		if eUICCRequired {
			for path, modem := range modems {
				// lpa.New will open the ISD-R logical channel, if it fails, the modem is not an eUICC.
				l, err := lpa.New(modem, m.config)
				slog.Debug("checking if the SIM card is an eUICC", "objectPath", path)
				if err != nil {
					delete(modems, path)
					slog.Error("failed to create LPA", "error", err)
				}
				slog.Info("the SIM card is an eUICC", "objectPath", path)
				l.Close()
			}
		}
		return m.run(modems, ctx, update)
	}
}

func (m *ModemRequiredMiddleware) run(modems map[dbus.ObjectPath]*modem.Modem, ctx *th.Context, update telego.Update) error {
	if len(modems) == 0 {
		return m.sendErrorModemNotFound(ctx, update)
	}
	// If there is only one modem, select it automatically.
	if len(modems) == 1 {
		for _, modem := range modems {
			ctx = ctx.WithValue("modem", modem)
			return ctx.Next(update)
		}
	}
	if err := m.ask(ctx, update, modems); err != nil {
		return err
	}

	ch := make(chan *modem.Modem, 1)
	sessions.Store(update.Message.GetMessageID(), &session{
		ch:     ch,
		ctx:    ctx,
		update: update,
	})

	select {
	case modem := <-ch:
		sessions.Delete(update.Message.GetMessageID())
		return ctx.WithValue("modem", modem).Next(update)
	case <-ctx.Done():
		sessions.Delete(update.Message.GetMessageID())
		return ctx.Err()
	}
}

func (m *ModemRequiredMiddleware) HandleModemSelectionCallbackQuery(ctx *th.Context, query telego.CallbackQuery) error {
	parts := strings.Split(query.Data[len(CallbackQueryAskModemPrefix)+1:], ":")
	messageId, _ := strconv.Atoi(parts[1])
	slog.Info("using modem", "objectPath", parts[0], "messageId", messageId)

	modems, err := m.mm.Modems()
	if err != nil {
		return err
	}

	s, ok := sessions.Load(messageId)
	if !ok {
		_, err := ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(query.From.ID), "No pending command found"))
		return err
	}
	s.(*session).ch <- modems[dbus.ObjectPath(parts[0])]
	sessions.Delete(messageId)

	if err := ctx.Bot().AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: query.ID,
	}); err != nil {
		return err
	}
	return ctx.Bot().DeleteMessage(ctx, &telego.DeleteMessageParams{
		ChatID:    tu.ID(query.Message.GetChat().ID),
		MessageID: query.Message.GetMessageID(),
	})
}

func (m *ModemRequiredMiddleware) sendErrorModemNotFound(ctx *th.Context, update telego.Update) error {
	_, err := ctx.Bot().SendMessage(
		ctx,
		tu.Message(
			tu.ID(update.Message.From.ID),
			"No modems were found. Please plug in a modem and try again.",
		).WithReplyParameters(&telego.ReplyParameters{
			MessageID: update.Message.MessageID,
		}),
	)
	if err != nil {
		return err
	}
	return errors.New("no modems were found")
}

func (m *ModemRequiredMiddleware) ask(ctx *th.Context, update telego.Update, modems map[dbus.ObjectPath]*modem.Modem) error {
	var err error
	var buttons [][]telego.InlineKeyboardButton
	var message string
	for path, modem := range modems {
		name, ok := m.config.Alias[modem.EquipmentIdentifier]
		modemName := util.If(ok, name, modem.Model)
		buttons = append(buttons, tu.InlineKeyboardRow(telego.InlineKeyboardButton{
			Text:         fmt.Sprintf("%s (%s)", modemName, modem.EquipmentIdentifier[len(modem.EquipmentIdentifier)-4:]),
			CallbackData: fmt.Sprintf("%s:%s:%d", CallbackQueryAskModemPrefix, path, update.Message.MessageID),
		}))
		message += fmt.Sprintf(`
*%s*
Manufacturer: %s
IMEI: %s
Firmware revision: %s
Hardware revision: %s
ICCID: %s
Operator: %s
Number: %s
		`, util.EscapeText(modemName),
			util.EscapeText(modem.Manufacturer),
			modem.EquipmentIdentifier,
			util.EscapeText(modem.FirmwareRevision),
			util.EscapeText(modem.HardwareRevision),
			modem.Sim.Identifier,
			util.EscapeText(util.If(modem.Sim.OperatorName != "", modem.Sim.OperatorName, util.LookupCarrier(modem.Sim.OperatorIdentifier))),
			util.EscapeText(modem.Number),
		)
	}

	_, err = ctx.Bot().SendMessage(ctx, tu.Message(
		tu.ID(update.Message.From.ID),
		strings.TrimRight(message, "\n"),
	).WithReplyMarkup(tu.InlineKeyboard(buttons...)).WithReplyParameters(&telego.ReplyParameters{
		MessageID: update.Message.MessageID,
	}).WithParseMode(telego.ModeMarkdownV2))
	return err
}
