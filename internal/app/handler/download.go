package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/damonto/euicc-go/lpa"
	sgp22 "github.com/damonto/euicc-go/v2"
	"github.com/damonto/telmo/internal/app/state"
	"github.com/damonto/telmo/internal/pkg/config"
	tlpa "github.com/damonto/telmo/internal/pkg/lpa"
	"github.com/damonto/telmo/internal/pkg/modem"
	"github.com/damonto/telmo/internal/pkg/util"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type DownloadHandler struct {
	Handler
	state            *state.StateManager
	config           *config.Config
	cancelled        bool
	confirmed        chan bool
	confirmationCode chan string
}

type downloader struct {
	h               *DownloadHandler
	downloadCtx     context.Context
	cancel          context.CancelFunc
	ctx             *th.Context
	message         telego.Message
	progressMessage *telego.Message
}

type DownloadValue struct {
	modem          *modem.Modem
	cancel         context.CancelFunc
	activationCode *lpa.ActivationCode
}

const (
	DownloadAskActivationCode             state.State = "download_ask_activation_code"
	DownloadAskConfirmationCodeFirst      state.State = "download_ask_confirmation_code_first"
	DownloadAskConfirmationCodeInProgress state.State = "download_ask_confirmation_code_in_progress"
	DownloadConfirm                       state.State = "download_confirm"

	DownloadCallbackDataPrefix = "download"
)

func NewDownloadHandler(config *config.Config, state *state.StateManager) state.Handler {
	return &DownloadHandler{
		state:            state,
		config:           config,
		cancelled:        false,
		confirmationCode: make(chan string, 1),
		confirmed:        make(chan bool, 1),
	}
}

func (h *DownloadHandler) Handle() th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		m := h.Modem(ctx)
		value := &DownloadValue{modem: m}
		s := &state.ChatState{
			Handler: h,
			State:   DownloadAskActivationCode,
			Value:   value,
		}
		h.state.Enter(update.Message.Chat.ID, s)
		update.Message.Text = strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/download"))
		if update.Message.Text != "" {
			return h.downloadProfile(ctx, *update.Message, s, value)
		}
		_, err := h.Reply(ctx, update, util.EscapeText("Please send me the activation code."), nil)
		return err
	}
}

func (h *DownloadHandler) HandleMessage(ctx *th.Context, message telego.Message, s *state.ChatState) error {
	value := s.Value.(*DownloadValue)
	if s.State == DownloadAskActivationCode {
		return h.downloadProfile(ctx, message, s, value)
	}
	if s.State == DownloadAskConfirmationCodeFirst {
		value.activationCode.ConfirmationCode = message.Text
		return h.download(ctx, message, s, value)
	}
	if s.State == DownloadAskConfirmationCodeInProgress {
		h.confirmationCode <- message.Text
		return nil
	}
	return nil
}

func (h *DownloadHandler) downloadProfile(ctx *th.Context, message telego.Message, s *state.ChatState, value *DownloadValue) error {
	var err error
	var ccRequired bool
	value.activationCode, ccRequired, err = h.parseActivationCode(value, message.Text)
	if err != nil {
		return err
	}
	if ccRequired {
		h.state.Current(message.From.ID, DownloadAskConfirmationCodeFirst)
		_, err := h.ReplyMessage(ctx, message, util.EscapeText("Please send me the confirmation code."), nil)
		return err
	}
	return h.download(ctx, message, s, value)
}

func (d *downloader) OnProgress(stage lpa.DownloadStage) {
	percent := map[lpa.DownloadStage]int{
		lpa.DownloadStageAuthenticateClient: 3,
		lpa.DownloadStageAuthenticateServer: 5,
		lpa.DownloadStageInstall:            9,
	}
	progressBar := strings.Repeat("⣿", percent[stage]) + strings.Repeat("⣀", 10-percent[stage])
	progressBar = util.EscapeText("Your profile is being downloaded.\n ⏳ " + progressBar + fmt.Sprintf(" %d%%", percent[stage]*10))
	var err error
	if d.progressMessage != nil {
		_, err = d.ctx.Bot().EditMessageText(d.ctx, &telego.EditMessageTextParams{
			ChatID:    d.progressMessage.Chat.ChatID(),
			MessageID: d.progressMessage.GetMessageID(),
			Text:      progressBar,
			ParseMode: telego.ModeMarkdownV2,
		})
	} else {
		d.progressMessage, err = d.h.ReplyMessage(d.ctx, d.message, progressBar, nil)
	}
	if err != nil {
		slog.Error("failed to send progress message", "error", err)
	}
}

func (d *downloader) OnConfirm(metadata *sgp22.ProfileInfo) bool {
	d.h.ReplyMessage(
		d.ctx, d.message,
		util.EscapeText(fmt.Sprintf(`
		Are you sure you want to download this profile?
Provider Name: %s
Profile Name: %s
ICCID: %s
`, metadata.ServiceProviderName, metadata.ProfileName, metadata.ICCID)),
		func(message *telego.SendMessageParams) error {
			message.WithReplyMarkup(tu.InlineKeyboard(
				tu.InlineKeyboardRow(
					telego.InlineKeyboardButton{
						Text:         "Yes",
						CallbackData: fmt.Sprintf("%s:%s", DownloadCallbackDataPrefix, "yes"),
					},
					telego.InlineKeyboardButton{
						Text:         "No",
						CallbackData: fmt.Sprintf("%s:%s", DownloadCallbackDataPrefix, "no"),
					},
				),
			))
			return nil
		},
	)
	d.h.state.Current(d.message.From.ID, DownloadConfirm)
	return <-d.h.confirmed
}

func (d *downloader) OnEnterConfirmationCode() string {
	defer func() {
		d.ctx.Bot().DeleteMessage(d.ctx, &telego.DeleteMessageParams{
			MessageID: d.progressMessage.GetMessageID(),
			ChatID:    d.progressMessage.Chat.ChatID(),
		})
		d.progressMessage = nil
	}()
	if _, err := d.h.ReplyMessage(d.ctx, d.message, util.EscapeText("Please send me the confirmation code."), nil); err != nil {
		d.h.state.Exit(d.message.From.ID)
		d.cancel()
		return <-d.h.confirmationCode
	}
	d.h.state.Current(d.message.From.ID, DownloadAskConfirmationCodeInProgress)
	return <-d.h.confirmationCode
}

func (h *DownloadHandler) download(ctx *th.Context, message telego.Message, _ *state.ChatState, value *DownloadValue) error {
	defer h.state.Exit(message.From.ID)
	var downloadCtx context.Context
	downloadCtx, value.cancel = context.WithTimeout(context.Background(), 10*time.Minute)
	defer value.cancel()
	d := &downloader{h: h, ctx: ctx, downloadCtx: downloadCtx, cancel: value.cancel, message: message}
	l, err := tlpa.New(value.modem, h.config)
	if err != nil {
		return err
	}
	defer l.Close()
	if err := l.Download(downloadCtx, value.activationCode, &lpa.DownloadOptions{
		OnProgress:              d.OnProgress,
		OnConfirm:               d.OnConfirm,
		OnEnterConfirmationCode: d.OnEnterConfirmationCode,
	}); err != nil {
		h.ReplyMessage(ctx, message, util.EscapeText(err.Error()), nil)
		if d.progressMessage != nil {
			ctx.Bot().DeleteMessage(ctx, &telego.DeleteMessageParams{
				MessageID: d.progressMessage.GetMessageID(),
				ChatID:    d.message.Chat.ChatID(),
			})
		}
		return err
	}
	result := util.If(h.cancelled, "The download has been canceled. /profiles", "The profile has been downloaded. /profiles")
	_, err = ctx.Bot().EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:    message.Chat.ChatID(),
		MessageID: d.progressMessage.GetMessageID(),
		Text:      util.EscapeText(result),
		ParseMode: telego.ModeMarkdownV2,
	})
	return err
}

func (h *DownloadHandler) parseActivationCode(value *DownloadValue, text string) (ac *lpa.ActivationCode, ccRequired bool, err error) {
	parts := strings.Split(text, "$")
	ac = &lpa.ActivationCode{
		SMDP: &url.URL{Scheme: "https", Host: parts[1]},
		IMEI: value.modem.EquipmentIdentifier,
	}
	if len(parts) >= 3 {
		ac.MatchingID = parts[2]
	}
	if len(parts) == 5 && parts[4] == "1" {
		ccRequired = true
	}
	return ac, ccRequired, nil
}

func (h *DownloadHandler) HandleCallbackQuery(ctx *th.Context, query telego.CallbackQuery, s *state.ChatState) error {
	if s.State != DownloadConfirm {
		return nil
	}
	confirmed := query.Data[len(DownloadCallbackDataPrefix)+1:]
	h.confirmed <- confirmed == "yes"
	if confirmed == "yes" {
		if err := ctx.Bot().DeleteMessage(ctx, &telego.DeleteMessageParams{
			ChatID:    tu.ID(query.From.ID),
			MessageID: query.Message.GetMessageID(),
		}); err != nil {
			slog.Warn("failed to delete message", "error", err)
		}
	}
	if confirmed == "no" {
		h.cancelled = true
		value := s.Value.(*DownloadValue)
		value.cancel()
		slog.Info("download canceled", "activationCode", value.activationCode)
		h.state.Exit(query.From.ID)
		_, err := h.ReplyCallbackQuery(ctx, query, util.EscapeText("Download canceled!"), nil)
		return err
	}
	return nil
}
