package router

import (
	"context"
	"log/slog"
	"slices"
	"strings"

	"github.com/damonto/telmo/internal/app/handler"
	"github.com/damonto/telmo/internal/app/middleware"
	"github.com/damonto/telmo/internal/app/state"
	"github.com/damonto/telmo/internal/pkg/config"
	"github.com/damonto/telmo/internal/pkg/modem"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

type router struct {
	*th.BotHandler
	bot          *telego.Bot
	mm           *modem.Manager
	config       *config.Config
	stateManager *state.StateManager
}

func New(bot *telego.Bot, handler *th.BotHandler, mm *modem.Manager, config *config.Config) *router {
	return &router{
		bot:          bot,
		BotHandler:   handler,
		mm:           mm,
		config:       config,
		stateManager: state.New(),
	}
}

func (r *router) Register() {
	r.stateManager.RegisterCallback(r.BotHandler)
	r.registerCommands()
	r.registerHandlers()
	r.stateManager.RegisterMessage(r.BotHandler)
}

func (r *router) registerCommands() {
	commands := []telego.BotCommand{
		{Command: "modem", Description: "List all plugged in modems"},
		{Command: "slot", Description: "List all SIM slots on the modem"},
		{Command: "chip", Description: "Get the eUICC chip information"},
		{Command: "ussd", Description: "Send a USSD command to the carrier"},
		{Command: "send", Description: "Send an SMS to a phone number"},
		{Command: "msisdn", Description: "Update the MSISDN (phone number) on the SIM"},
		{Command: "register_network", Description: "Register the modem to a network"},
		{Command: "profiles", Description: "List all profiles on the eUICC"},
		{Command: "discovery", Description: "Discover profiles from the SM-DS+ server"},
		{Command: "download", Description: "Download a profile into the eUICC"},
	}

	if err := r.bot.SetMyCommands(context.Background(), &telego.SetMyCommandsParams{
		Scope: &telego.BotCommandScopeAllPrivateChats{
			Type: telego.ScopeTypeAllPrivateChats,
		},
		Commands: commands,
	}); err != nil {
		slog.Error("failed to set commands", "error", err)
	}
}

func (r *router) registerHandlers() {
	r.Handle(handler.NewStartHandler().Handle(), th.CommandEqual("start"))

	modemRequiredMiddleware := middleware.NewModemRequiredMiddleware(r.mm, r.BotHandler, r.config)

	admin := r.Group(th.Not(th.CommandEqual("start")))
	admin.Use(middleware.Admin(r.config))
	admin.Handle(handler.NewListModemHandler(r.config, r.mm).Handle(), th.CommandEqual("modem"))

	{
		standard := admin.Group(r.predicate([]string{"/send", "/slot", "/ussd", "/msisdn", "/register_network"}))
		standard.Use(modemRequiredMiddleware.Middleware(false))
		standard.Handle(handler.NewSIMSlotHandler(r.stateManager).Handle(), th.CommandEqual("slot"))
		standard.Handle(handler.NewUSSDHandler(r.stateManager).Handle(), th.CommandEqual("ussd"))
		standard.Handle(handler.NewSendHandler(r.stateManager).Handle(), th.CommandEqual("send"))
		standard.Handle(handler.NewMSISDNHandler(r.config, r.stateManager).Handle(), th.CommandEqual("msisdn"))
		standard.Handle(handler.NewRegisterNetworkHandler(r.stateManager).Handle(), th.CommandEqual("register_network"))
	}

	{
		euicc := admin.Group(r.predicate([]string{"/chip", "/profiles", "/download", "/send_notification", "/discovery"}))
		euicc.Use(modemRequiredMiddleware.Middleware(true))
		euicc.Handle(handler.NewChipHandler(r.config).Handle(), th.CommandEqual("chip"))
		euicc.Handle(handler.NewProfileHandler(r.config, r.stateManager).Handle(), th.CommandEqual("profiles"))
		euicc.Handle(handler.NewDiscoveryHandler(r.config).Handle(), th.CommandEqual("discovery"))
		euicc.Handle(handler.NewDownloadHandler(r.config, r.stateManager).Handle(), th.CommandEqual("download"))
		euicc.Handle(handler.NewSendNotificationHandler(r.config).Handle(), th.CommandEqualArgc("send_notification", 1))
	}
}

func (r *router) predicate(filters []string) th.Predicate {
	return func(ctx context.Context, update telego.Update) bool {
		return slices.Contains(filters, strings.Split(update.Message.Text, " ")[0])
	}
}
