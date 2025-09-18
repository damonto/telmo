package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/damonto/telmo/internal/app"
	"github.com/damonto/telmo/internal/pkg/config"
	"github.com/damonto/telmo/internal/pkg/modem"
	"github.com/damonto/telmo/internal/pkg/util"
	"github.com/godbus/dbus/v5"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

var BuildVersion string

type Subscriber struct {
	ctx    context.Context
	cancel context.CancelFunc
}

var cfg = config.New()

func init() {
	flag.StringVar(&cfg.BotToken, "bot-token", "", "Telegram bot token")
	flag.Var(&cfg.AdminId, "admin-id", "Admin user ID with bot management privileges (multiple allowed)")
	flag.BoolVar(&cfg.Slowdown, "slowdown", false, "Enable slowdown mode (MSS: 120)")
	flag.BoolVar(&cfg.ForceAT, "force-at", false, "Force the use of AT commands as the LPA driver")
	flag.BoolVar(&cfg.Compatible, "compatible", false, "Enable if your modem does not support proactive refresh")
	flag.Var(&cfg.ModemName, "modem-name", "Modem name IMEI:name (multiple allowed)")
	flag.StringVar(&cfg.Endpoint, "endpoint", "https://api.telegram.org", "Telegram Bot API endpoint")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "Enable verbose logging")

	flag.Parse()
}

func main() {
	if cfg.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	if os.Geteuid() != 0 {
		slog.Error("Please run as root")
		os.Exit(1)
	}
	if err := cfg.IsValid(); err != nil {
		slog.Error("Config is invalid", "error", err)
		os.Exit(1)
	}

	slog.Info("Starting Telmo", "version", BuildVersion)

	bot, err := telego.NewBot(cfg.BotToken,
		telego.WithAPIServer(cfg.Endpoint),
		telego.WithDefaultLogger(cfg.Verbose, true),
	)
	if err != nil {
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	me, err := bot.GetMe(ctx)
	if err != nil {
		panic(err)
	}
	slog.Info("Bot started", "username", me.Username, "id", me.ID)

	mm, err := modem.NewManager()
	if err != nil {
		panic(err)
	}

	go subscribe(bot, mm)

	app, err := app.New(ctx, bot, mm, cfg)
	if err != nil {
		panic(err)
	}
	go func() {
		if err := app.Start(); err != nil {
			panic(err)
		}
	}()
	<-ctx.Done()
	slog.Info("Stopping Telmo")
	app.Shutdown()
	slog.Info("Goodbye!")
}

func subscribe(bot *telego.Bot, mm *modem.Manager) {
	subscribers := make(map[dbus.ObjectPath]*Subscriber)
	modems, err := mm.Modems()
	if err != nil {
		panic(err)
	}

	go subscribeMessaging(bot, modems, subscribers)

	if err := mm.Subscribe(func(modems map[dbus.ObjectPath]*modem.Modem) error {
		for path, s := range subscribers {
			slog.Debug("Canceling subscriber", "path", path)
			s.cancel()
		}
		go subscribeMessaging(bot, modems, subscribers)
		return nil
	}); err != nil {
		panic(err)
	}
}

func subscribeMessaging(bot *telego.Bot, modems map[dbus.ObjectPath]*modem.Modem, subscribers map[dbus.ObjectPath]*Subscriber) {
	for path, m := range modems {
		slog.Info("Subscribing to modem messaging", "path", path)
		ctx, cancel := context.WithCancel(context.Background())
		go func(ctx context.Context, m *modem.Modem) {
			if err := m.SubscribeMessaging(ctx, func(message *modem.SMS) error {
				if err := send(bot, m, message); err != nil {
					slog.Error("Failed to send message", "error", err)
				}
				return nil
			}); err != nil {
				slog.Error("Failed to subscribe to modem messaging", "error", err)
			}
		}(ctx, m)
		subscribers[path] = &Subscriber{ctx: ctx, cancel: cancel}
	}
}

func send(bot *telego.Bot, modem *modem.Modem, messsage *modem.SMS) error {
	template := `
☎️ *%s*
👤 *\[%s\] \- %s*
>%s
`
	operatorName, err := modem.OperatorName()
	if err != nil {
		slog.Error("Failed to get operator name", "error", err)
		operatorName = "unknown"
	}
	name, ok := cfg.ModemName[modem.EquipmentIdentifier]
	modemName := util.If(ok, name, modem.Model)
	message := fmt.Sprintf(
		template,
		util.EscapeText(modemName),
		util.EscapeText(operatorName),
		util.EscapeText(messsage.Number),
		fmt.Sprintf("`%s`", util.EscapeText(messsage.Text)),
	)
	for _, adminId := range cfg.AdminId {
		msg, err := bot.SendMessage(context.Background(), tu.Message(tu.ID(adminId), message).WithParseMode(telego.ModeMarkdownV2))
		if err != nil {
			slog.Error("Failed to send message", "error", err, "to", adminId, "message", message)
		} else {
			slog.Info("Message sent", "id", msg.MessageID, "to", adminId)
		}
	}
	return nil
}
