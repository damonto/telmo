package esim

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	elpa "github.com/damonto/euicc-go/lpa"
	sgp22 "github.com/damonto/euicc-go/v2"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"github.com/damonto/sigmo/internal/app/handler"
	"github.com/damonto/sigmo/internal/pkg/carrier"
	"github.com/damonto/sigmo/internal/pkg/config"
	"github.com/damonto/sigmo/internal/pkg/lpa"
	mmodem "github.com/damonto/sigmo/internal/pkg/modem"
)

type Handler struct {
	handler.Handler
	cfg     *config.Config
	manager *mmodem.Manager
	service *Service
}

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const enableTimeout = time.Minute

var errEnableTimeout = errors.New("enabling timed out, please refresh to confirm whether the profile is active")

const (
	wsTypeStart                    = "start"
	wsTypeProgress                 = "progress"
	wsTypePreview                  = "preview"
	wsTypeConfirm                  = "confirm"
	wsTypeConfirmationCode         = "confirmation_code"
	wsTypeConfirmationCodeRequired = "confirmation_code_required"
	wsTypeCancel                   = "cancel"
	wsTypeCompleted                = "completed"
	wsTypeError                    = "error"
)

func New(cfg *config.Config, manager *mmodem.Manager) *Handler {
	return &Handler{
		cfg:     cfg,
		manager: manager,
		service: NewService(cfg, manager),
	}
}

func (h *Handler) List(c echo.Context) error {
	modem, err := h.FindModem(h.manager, c.Param("id"))
	if err != nil {
		return h.NotFound(c, err)
	}
	response, err := h.service.List(modem)
	if err != nil {
		if errors.Is(err, lpa.ErrNoSupportedAID) {
			return h.NotFound(c, err)
		}
		return h.InternalServerError(c, err)
	}
	return h.Respond(c, response)
}

func (h *Handler) Enable(c echo.Context) error {
	modem, err := h.FindModem(h.manager, c.Param("id"))
	if err != nil {
		return h.NotFound(c, err)
	}
	iccid, err := iccidFromParam(c)
	if err != nil {
		return h.BadRequest(c, err)
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), enableTimeout)
	defer cancel()
	if err := h.service.Enable(ctx, modem, iccid); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return h.Error(c, http.StatusRequestTimeout, errEnableTimeout)
		}
		if errors.Is(err, context.Canceled) {
			return nil
		}
		if errors.Is(err, lpa.ErrNoSupportedAID) {
			return h.NotFound(c, err)
		}
		return h.InternalServerError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) Delete(c echo.Context) error {
	modem, err := h.FindModem(h.manager, c.Param("id"))
	if err != nil {
		return h.NotFound(c, err)
	}
	iccid, err := iccidFromParam(c)
	if err != nil {
		return h.BadRequest(c, err)
	}
	if err := h.service.Delete(modem, iccid); err != nil {
		if errors.Is(err, lpa.ErrNoSupportedAID) {
			return h.NotFound(c, err)
		}
		return h.InternalServerError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) Download(c echo.Context) error {
	modem, err := h.FindModem(h.manager, c.Param("id"))
	if err != nil {
		return h.NotFound(c, err)
	}

	conn, err := wsUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	start, err := readStartMessage(conn)
	if err != nil {
		_ = conn.WriteJSON(downloadServerMessage{Type: wsTypeError, Message: err.Error()})
		return nil
	}

	activationCode, err := buildActivationCode(modem, start)
	if err != nil {
		_ = conn.WriteJSON(downloadServerMessage{Type: wsTypeError, Message: err.Error()})
		return nil
	}

	downloadCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	session := newDownloadSession(conn, cancel)

	opts := &elpa.DownloadOptions{
		OnProgress: func(stage elpa.DownloadStage) {
			session.sendIfConnected(downloadServerMessage{
				Type:  wsTypeProgress,
				Stage: stage.String(),
			})
		},
		OnConfirm: func(info *sgp22.ProfileInfo) bool {
			preview := profilePreviewFrom(info)
			if err := session.send(downloadServerMessage{
				Type:    wsTypePreview,
				Profile: &preview,
			}); err != nil {
				return false
			}
			return session.waitForConfirm(downloadCtx)
		},
		OnEnterConfirmationCode: func() string {
			session.sendIfConnected(downloadServerMessage{
				Type: wsTypeConfirmationCodeRequired,
			})
			code := session.waitForConfirmationCode(downloadCtx)
			return strings.TrimSpace(code)
		},
	}

	if err := h.service.Download(downloadCtx, modem, activationCode, opts); err != nil {
		_ = session.send(downloadServerMessage{Type: wsTypeError, Message: err.Error()})
		return nil
	}

	_ = session.send(downloadServerMessage{Type: wsTypeCompleted})
	return nil
}

func (h *Handler) UpdateNickname(c echo.Context) error {
	modem, err := h.FindModem(h.manager, c.Param("id"))
	if err != nil {
		return h.NotFound(c, err)
	}
	iccid, err := iccidFromParam(c)
	if err != nil {
		return h.BadRequest(c, err)
	}
	var req UpdateNicknameRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}
	if err := h.service.UpdateNickname(modem, iccid, req.Nickname); err != nil {
		if errors.Is(err, errInvalidNickname) {
			return h.BadRequest(c, err)
		}
		if errors.Is(err, lpa.ErrNoSupportedAID) {
			return h.NotFound(c, err)
		}
		return h.InternalServerError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func iccidFromParam(c echo.Context) (sgp22.ICCID, error) {
	iccidParam := c.Param("iccid")
	if iccidParam == "" {
		return nil, errors.New("iccid is required")
	}
	iccid, err := sgp22.NewICCID(iccidParam)
	if err != nil {
		return nil, fmt.Errorf("invalid iccid %q: %w", iccidParam, err)
	}
	return iccid, nil
}

func readStartMessage(conn *websocket.Conn) (downloadClientMessage, error) {
	var start downloadClientMessage
	if err := conn.ReadJSON(&start); err != nil {
		return downloadClientMessage{}, err
	}
	if start.Type != "" && start.Type != wsTypeStart {
		return downloadClientMessage{}, fmt.Errorf("unexpected message type %q", start.Type)
	}
	if start.SMDP == "" {
		return downloadClientMessage{}, errors.New("smdp is required")
	}
	return start, nil
}

func profilePreviewFrom(info *sgp22.ProfileInfo) downloadProfilePreview {
	carrierInfo := carrier.Lookup(info.ProfileOwner.MCC() + info.ProfileOwner.MNC())
	preview := downloadProfilePreview{
		ICCID:               info.ICCID.String(),
		ServiceProviderName: info.ServiceProviderName,
		ProfileName:         info.ProfileName,
		ProfileNickname:     info.ProfileNickname,
		ProfileState:        info.ProfileState.String(),
		RegionCode:          carrierInfo.Region,
	}
	if info.Icon.Valid() {
		preview.Icon = fmt.Sprintf("data:%s;base64,%s", info.Icon.FileType(), info.Icon.String())
	}
	return preview
}
