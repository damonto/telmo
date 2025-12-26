package esim

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	elpa "github.com/damonto/euicc-go/lpa"
	sgp22 "github.com/damonto/euicc-go/v2"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"github.com/damonto/sigmo/internal/app/handler"
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
	if err := h.service.Enable(modem, iccid); err != nil {
		if errors.Is(err, lpa.ErrNoSupportedAID) {
			return h.NotFound(c, err)
		}
		return h.InternalServerError(c, err)
	}
	return h.Respond(c, struct{}{})
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
	return h.Respond(c, struct{}{})
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
		_ = sendMessage(conn, nil, downloadServerMessage{Type: wsTypeError, Message: err.Error()})
		return nil
	}

	activationCode, err := buildActivationCode(modem, start)
	if err != nil {
		_ = sendMessage(conn, nil, downloadServerMessage{Type: wsTypeError, Message: err.Error()})
		return nil
	}

	downloadCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	disconnectCh := make(chan struct{})
	var disconnectOnce sync.Once
	disconnect := func() {
		disconnectOnce.Do(func() {
			close(disconnectCh)
		})
	}

	confirmCh := make(chan bool, 1)
	confirmationCodeCh := make(chan string, 1)

	go func() {
		defer disconnect()
		for {
			var msg downloadClientMessage
			if err := conn.ReadJSON(&msg); err != nil {
				return
			}
			switch msg.Type {
			case wsTypeConfirm:
				if msg.Accept != nil {
					select {
					case confirmCh <- *msg.Accept:
					default:
					}
				}
			case wsTypeConfirmationCode:
				select {
				case confirmationCodeCh <- msg.Code:
				default:
				}
			case wsTypeCancel:
				cancel()
			}
		}
	}()

	opts := &elpa.DownloadOptions{
		OnProgress: func(stage elpa.DownloadStage) {
			sendIfConnected(conn, disconnectCh, disconnect, downloadServerMessage{
				Type:  wsTypeProgress,
				Stage: stage.String(),
			})
		},
		OnConfirm: func(info *sgp22.ProfileInfo) bool {
			preview := profilePreviewFrom(info)
			if err := sendMessage(conn, disconnect, downloadServerMessage{
				Type:    wsTypePreview,
				Profile: &preview,
			}); err != nil {
				return false
			}
			return waitForConfirm(downloadCtx, disconnectCh, confirmCh)
		},
		OnEnterConfirmationCode: func() string {
			sendIfConnected(conn, disconnectCh, disconnect, downloadServerMessage{
				Type: wsTypeConfirmationCodeRequired,
			})
			code := waitForConfirmationCode(downloadCtx, disconnectCh, confirmationCodeCh)
			return strings.TrimSpace(code)
		},
	}

	if err := h.service.Download(downloadCtx, modem, activationCode, opts); err != nil {
		_ = sendMessage(conn, disconnect, downloadServerMessage{Type: wsTypeError, Message: err.Error()})
		return nil
	}

	_ = sendMessage(conn, disconnect, downloadServerMessage{Type: wsTypeCompleted})
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
	return h.Respond(c, struct{}{})
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
	if start.ActivationCode == "" {
		return downloadClientMessage{}, errors.New("activationCode is required")
	}
	return start, nil
}

func buildActivationCode(modem *mmodem.Modem, start downloadClientMessage) (*elpa.ActivationCode, error) {
	smdpURL, err := parseSMDP(start.SMDP)
	if err != nil {
		return nil, err
	}
	matchingID := strings.TrimSpace(start.ActivationCode)
	if matchingID == "" {
		return nil, errors.New("activationCode is required")
	}
	imei, err := modem.ThreeGPP().IMEI()
	if err != nil {
		return nil, fmt.Errorf("reading modem IMEI: %w", err)
	}
	return &elpa.ActivationCode{
		SMDP:             smdpURL,
		MatchingID:       matchingID,
		IMEI:             imei,
		ConfirmationCode: strings.TrimSpace(start.ConfirmationCode),
	}, nil
}

func parseSMDP(raw string) (*url.URL, error) {
	smdp := strings.TrimSpace(raw)
	if smdp == "" {
		return nil, errors.New("smdp is required")
	}
	if !strings.Contains(smdp, "://") {
		smdp = "https://" + smdp
	}
	parsed, err := url.Parse(smdp)
	if err != nil || parsed.Host == "" {
		return nil, fmt.Errorf("invalid smdp %q", raw)
	}
	return parsed, nil
}

func profilePreviewFrom(info *sgp22.ProfileInfo) downloadProfilePreview {
	preview := downloadProfilePreview{
		ICCID:               info.ICCID.String(),
		ServiceProviderName: info.ServiceProviderName,
		ProfileName:         info.ProfileName,
		ProfileNickname:     info.ProfileNickname,
		ProfileState:        info.ProfileState.String(),
		OwnerMCC:            info.ProfileOwner.MCC(),
		OwnerMNC:            info.ProfileOwner.MNC(),
	}
	if info.Icon.Valid() {
		preview.Icon = info.Icon.String()
		preview.IconType = info.Icon.FileType()
	}
	return preview
}

func sendMessage(conn *websocket.Conn, disconnect func(), msg downloadServerMessage) error {
	if err := conn.WriteJSON(msg); err != nil {
		if disconnect != nil {
			disconnect()
		}
		return err
	}
	return nil
}

func sendIfConnected(conn *websocket.Conn, disconnectCh <-chan struct{}, disconnect func(), msg downloadServerMessage) {
	select {
	case <-disconnectCh:
		return
	default:
	}
	_ = sendMessage(conn, disconnect, msg)
}

func waitForConfirm(ctx context.Context, disconnectCh <-chan struct{}, confirmCh <-chan bool) bool {
	select {
	case accept := <-confirmCh:
		return accept
	default:
	}
	select {
	case accept := <-confirmCh:
		return accept
	case <-ctx.Done():
		return false
	case <-disconnectCh:
		return false
	}
}

func waitForConfirmationCode(ctx context.Context, disconnectCh <-chan struct{}, confirmationCodeCh <-chan string) string {
	select {
	case code := <-confirmationCodeCh:
		return code
	default:
	}
	select {
	case code := <-confirmationCodeCh:
		return code
	case <-ctx.Done():
		return ""
	case <-disconnectCh:
		return ""
	}
}
