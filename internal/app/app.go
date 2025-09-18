package app

import (
	"context"
	"time"

	"github.com/damonto/telmo/internal/app/router"
	"github.com/damonto/telmo/internal/pkg/config"
	"github.com/damonto/telmo/internal/pkg/modem"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

type application struct {
	Bot     *telego.Bot
	m       *modem.Manager
	config  *config.Config
	handler *th.BotHandler
	updates <-chan telego.Update
	ctx     context.Context
}

func New(ctx context.Context, bot *telego.Bot, m *modem.Manager, config *config.Config) (*application, error) {
	app := &application{Bot: bot, m: m, config: config, ctx: ctx}
	var err error
	app.updates, err = bot.UpdatesViaLongPolling(ctx, nil)
	if err != nil {
		return nil, err
	}
	app.handler, err = th.NewBotHandler(bot, app.updates)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (app *application) Start() error {
	app.handler.Use(th.PanicRecovery())
	router.New(app.Bot, app.handler, app.m, app.config).Register()
	return app.handler.Start()
}

func (app *application) Shutdown() {
	stopCtx, stopCancel := context.WithTimeout(context.Background(), time.Second*30)
	defer stopCancel()

outer:
	for len(app.updates) > 0 {
		select {
		case <-stopCtx.Done():
			break outer
		case <-time.After(100 * time.Microsecond):
			//
		}
	}
	app.handler.StopWithContext(stopCtx)
}
