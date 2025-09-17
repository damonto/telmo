package handler

import (
	"fmt"
	"log/slog"

	sgp22 "github.com/damonto/euicc-go/v2"
	"github.com/damonto/telmo/internal/app/state"
	"github.com/damonto/telmo/internal/pkg/config"
	"github.com/damonto/telmo/internal/pkg/lpa"
	"github.com/damonto/telmo/internal/pkg/modem"
	"github.com/damonto/telmo/internal/pkg/util"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type ProfileHandler struct {
	Handler
	state *state.StateManager
}

const (
	ProfileActionCallbackDataPrefix = "profile"

	ProfileActionSetNickname state.State = "Set Nickname"
	ProfileActionDelete      state.State = "Delete"
	ProfileActionEnable      state.State = "Enable"
)

type ProfileValue struct {
	ICCID   sgp22.ICCID
	Action  state.State
	Profile *sgp22.ProfileInfo
	Value   string
	Modem   *modem.Modem
}

func NewProfileHandler(s *state.StateManager) state.Handler {
	return &ProfileHandler{
		state: s,
	}
}

func (h *ProfileHandler) HandleCallbackQuery(ctx *th.Context, query telego.CallbackQuery, state *state.ChatState) error {
	var err error
	value := state.Value.(*ProfileValue)
	value.ICCID, _ = sgp22.NewICCID(query.Data[len(ProfileActionCallbackDataPrefix)+1:])
	l, err := lpa.New(value.Modem)
	if err != nil {
		return err
	}
	ps, err := l.ListProfile(value.ICCID, nil)
	if err != nil {
		return err
	}
	defer l.Close()
	value.Profile = ps[0]
	return h.sendActionMessage(ctx, query, ps[0])
}

func (h *ProfileHandler) sendActionMessage(ctx *th.Context, query telego.CallbackQuery, profile *sgp22.ProfileInfo) error {
	var buttons []telego.KeyboardButton
	if profile.ProfileState == sgp22.ProfileEnabled {
		buttons = tu.KeyboardRow(
			tu.KeyboardButton(string(ProfileActionSetNickname)),
		)
	} else {
		buttons = tu.KeyboardRow(
			tu.KeyboardButton(string(ProfileActionSetNickname)),
			tu.KeyboardButton(string(ProfileActionEnable)),
			tu.KeyboardButton(string(ProfileActionDelete)),
		)
	}

	var message string
	name := fmt.Sprintf("[%s] %s",
		profile.ServiceProviderName,
		util.If(profile.ProfileNickname != "", profile.ProfileNickname, profile.ProfileName),
	)
	message += fmt.Sprintf(`
%s *%s*
%s
`,
		util.If(profile.ProfileState == sgp22.ProfileEnabled, "✅", "🅾️"),
		util.EscapeText(name),
		profile.ICCID,
	)
	message = util.EscapeText("What do you want to do with this profile? \n") + message
	_, err := h.ReplyCallbackQuery(ctx, query, message, func(message *telego.SendMessageParams) error {
		message.WithReplyMarkup(
			tu.Keyboard(buttons).
				WithOneTimeKeyboard().
				WithResizeKeyboard(),
		)
		return nil
	})
	return err
}

func (h *ProfileHandler) HandleMessage(ctx *th.Context, message telego.Message, s *state.ChatState) error {
	if state.State(message.Text) == ProfileActionSetNickname {
		return h.askNickname(ctx, message, s)
	}
	if s.State == ProfileActionSetNickname {
		return h.setNickname(ctx, message, s)
	}
	if state.State(message.Text) == ProfileActionEnable {
		return h.enableProfile(ctx, message, s)
	}
	if state.State(message.Text) == ProfileActionDelete {
		return h.confirmDelete(ctx, message, s)
	}
	if s.State == ProfileActionDelete {
		return h.deleteProfile(ctx, message, s)
	}
	h.state.Exit(message.Chat.ID)
	return nil
}

func (h *ProfileHandler) deleteProfile(ctx *th.Context, message telego.Message, s *state.ChatState) error {
	if message.Text != "Yes" {
		_, err := h.ReplyMessage(ctx, message, util.EscapeText("Okay, the profile will not be deleted. /profiles"), nil)
		return err
	}
	value := s.Value.(*ProfileValue)
	l, err := lpa.New(value.Modem)
	if err != nil {
		return err
	}
	defer l.Close()
	seqs, err := l.Delete(value.ICCID)
	if err != nil {
		return err
	}
	var command string
	for _, seq := range seqs {
		command += fmt.Sprintf("`/send_notification %d`\n", seq)
	}
	_, err = h.ReplyMessage(ctx, message,
		fmt.Sprintf("The profile has been deleted\\. If your operator hasn\\'t received the deletion notification, you can resend it using the following %s: \n\n%s\n/profiles", util.If(len(seqs) == 1, "command", "commands"), command),
		nil)
	return err
}

func (h *ProfileHandler) confirmDelete(ctx *th.Context, message telego.Message, s *state.ChatState) error {
	h.state.Current(message.Chat.ID, ProfileActionDelete)
	value := s.Value.(*ProfileValue)
	_, err := h.ReplyMessage(
		ctx,
		message,
		fmt.Sprintf(
			"Are you sure you want to delete this profile *%s*?",
			util.EscapeText(util.If(value.Profile.ProfileNickname != "", value.Profile.ProfileNickname, value.Profile.ProfileName)),
		),
		func(m *telego.SendMessageParams) error {
			m.WithReplyMarkup(tu.Keyboard(
				tu.KeyboardRow(
					tu.KeyboardButton("Yes"),
					tu.KeyboardButton("No"),
				),
			).WithOneTimeKeyboard().WithResizeKeyboard())
			return nil
		},
	)
	return err
}

func (h *ProfileHandler) enableProfile(ctx *th.Context, message telego.Message, s *state.ChatState) error {
	value := s.Value.(*ProfileValue)
	l, err := lpa.New(value.Modem)
	if err != nil {
		return err
	}
	if err := l.EnableProfile(value.ICCID, true); err != nil {
		return err
	}
	l.Close()
	if config.C.Compatible {
		if err := value.Modem.Restart(); err != nil {
			slog.Warn("Failed to restart the modem", "error", err)
		}
	}
	_, err = h.ReplyMessage(
		ctx,
		message,
		util.EscapeText("The profile has been enabled. It may take a few seconds to activate. /profiles"),
		nil,
	)
	return err
}

func (h *ProfileHandler) setNickname(ctx *th.Context, message telego.Message, s *state.ChatState) error {
	value := s.Value.(*ProfileValue)
	value.Value = message.Text
	l, err := lpa.New(value.Modem)
	if err != nil {
		return err
	}
	defer l.Close()
	if err := l.SetNickname(value.ICCID, value.Value); err != nil {
		return err
	}
	_, err = h.ReplyMessage(
		ctx,
		message,
		util.EscapeText("The nickname has been updated. /profiles"),
		nil,
	)
	return err
}

func (h *ProfileHandler) askNickname(ctx *th.Context, message telego.Message, _ *state.ChatState) error {
	h.state.Current(message.Chat.ID, ProfileActionSetNickname)
	_, err := h.ReplyMessage(
		ctx,
		message,
		util.EscapeText("Okay, please send me the new nickname."),
		nil,
	)
	return err
}

func (h *ProfileHandler) Handle() th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		l, err := h.LPA(ctx)
		if err != nil {
			return err
		}
		defer l.Close()

		h.state.Enter(update.Message.Chat.ID, &state.ChatState{
			Handler: h,
			Value: &ProfileValue{
				Modem: h.Modem(ctx),
			},
		})

		profiles, err := l.ListProfile(nil, nil)
		if err != nil {
			return err
		}
		if len(profiles) == 0 {
			_, err := h.Reply(ctx, update, util.EscapeText("No profiles were found."), nil)
			return err
		}
		buttons, message := h.message(profiles)
		_, err = h.Reply(ctx, update, message, func(message *telego.SendMessageParams) error {
			message.WithReplyMarkup(buttons)
			return nil
		})
		return err
	}
}

func (h *ProfileHandler) message(profiles []*sgp22.ProfileInfo) (*telego.InlineKeyboardMarkup, string) {
	var message string
	var buttons [][]telego.InlineKeyboardButton
	for _, profile := range profiles {
		name := fmt.Sprintf("[%s] %s",
			profile.ServiceProviderName,
			util.If(profile.ProfileNickname != "", profile.ProfileNickname, profile.ProfileName),
		)
		message += fmt.Sprintf(`
%s *%s*
%s
`,
			util.If(profile.ProfileState == sgp22.ProfileEnabled, "✅", "🅾️"),
			util.EscapeText(name),
			profile.ICCID,
		)
		id := profile.ICCID.String()
		id = id[len(id)-4:]
		buttons = append(buttons, tu.InlineKeyboardRow(telego.InlineKeyboardButton{
			Text:         fmt.Sprintf("%s (%s)", name, id),
			CallbackData: fmt.Sprintf("%s:%s", ProfileActionCallbackDataPrefix, profile.ICCID),
		}))
	}
	return tu.InlineKeyboard(buttons...), message
}
