package router

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/damonto/sigmo/internal/app/handler/esim"
	"github.com/damonto/sigmo/internal/app/handler/euicc"
	hmessage "github.com/damonto/sigmo/internal/app/handler/message"
	hmodem "github.com/damonto/sigmo/internal/app/handler/modem"
	hnetwork "github.com/damonto/sigmo/internal/app/handler/network"
	hussd "github.com/damonto/sigmo/internal/app/handler/ussd"
	"github.com/damonto/sigmo/internal/pkg/config"
	"github.com/damonto/sigmo/internal/pkg/modem"
	"github.com/damonto/sigmo/web"
)

func Register(e *echo.Echo, cfg *config.Config, manager *modem.Manager) {
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: http.FS(web.Root()),
		Index:      "index.html",
		HTML5:      true,
		Skipper: func(c echo.Context) bool {
			path := c.Request().URL.Path
			return strings.HasPrefix(path, "/api/")
		},
	}))

	v1 := e.Group("/api/v1")

	{
		h := hmodem.New(cfg, manager)
		v1.GET("/modems", h.List)
		v1.GET("/modems/:id", h.Get)
		v1.PUT("/modems/:id/sim-slots/:identifier", h.SwitchSimSlot)
		v1.PUT("/modems/:id/msisdn", h.UpdateMSISDN)
		v1.GET("/modems/:id/settings", h.GetSettings)
		v1.PUT("/modems/:id/settings", h.UpdateSettings)

		{
			h := hmessage.New(manager)
			v1.GET("/modems/:id/messages", h.List)
			v1.GET("/modems/:id/messages/:participant", h.ListByParticipant)
			v1.POST("/modems/:id/messages", h.Send)
			v1.DELETE("/modems/:id/messages/:participant", h.DeleteByParticipant)
		}

		{
			h := hussd.New(manager)
			v1.POST("/modems/:id/ussd", h.Execute)
		}

		{
			h := hnetwork.New(manager)
			v1.GET("/modems/:id/networks", h.List)
			v1.PUT("/modems/:id/networks/:operatorCode", h.Register)
		}

		{
			h := euicc.New(cfg, manager)
			v1.GET("/modems/:id/euicc", h.Get)
		}

		{
			h := esim.New(cfg, manager)
			v1.GET("/modems/:id/esims", h.List)
			v1.GET("/modems/:id/esims/download", h.Download)
			v1.POST("/modems/:id/esims/:iccid/enabling", h.Enable)
			v1.PUT("/modems/:id/esims/:iccid/nickname", h.UpdateNickname)
			v1.DELETE("/modems/:id/esims/:iccid", h.Delete)
		}
	}
}
