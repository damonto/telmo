package router

import (
	"github.com/labstack/echo/v4"

	"github.com/damonto/sigmo/internal/app/handler/esim"
	"github.com/damonto/sigmo/internal/app/handler/euicc"
	hmodem "github.com/damonto/sigmo/internal/app/handler/modem"
	"github.com/damonto/sigmo/internal/pkg/config"
	"github.com/damonto/sigmo/internal/pkg/modem"
)

func Register(e *echo.Echo, cfg *config.Config, manager *modem.Manager) {
	v1 := e.Group("/api/v1")

	{
		h := hmodem.New(cfg, manager)
		v1.GET("/modems", h.List)
		v1.GET("/modems/:id", h.Get)

		{
			h := euicc.New(cfg, manager)
			v1.GET("/modems/:id/euicc", h.Get)
		}

			{
				h := esim.New(cfg, manager)
				v1.GET("/modems/:id/esims", h.List)
				v1.POST("/modems/:id/esims/:iccid/enabling", h.Enable)
				v1.PUT("/modems/:id/esims/:iccid/nickname", h.UpdateNickname)
				v1.DELETE("/modems/:id/esims/:iccid", h.Delete)
			}
		}
	}
