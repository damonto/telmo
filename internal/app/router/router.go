package router

import (
	"github.com/labstack/echo/v4"

	modemhandler "github.com/damonto/sigmo/internal/app/handler/modem"
	"github.com/damonto/sigmo/internal/pkg/config"
	"github.com/damonto/sigmo/internal/pkg/modem"
)

func Register(e *echo.Echo, cfg *config.Config, manager *modem.Manager) {
	v1 := e.Group("/api/v1")

	modemHandler := modemhandler.NewHandler(cfg, manager)
	v1.GET("/modems", modemHandler.ListInserted)
}
