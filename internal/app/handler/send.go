package handler

import (
	"github.com/damonto/telmo/internal/app/state"
	"github.com/damonto/telmo/internal/pkg/modem"
	"github.com/damonto/telmo/internal/pkg/util"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"log"
)

type SendHandler struct {
	*Handler
}

type SMSValue struct {
	To    string
	Modem *modem.Modem
}

const (
	SendActionAskPhoneNumber state.State = "send_ask_phone_number"
	SendActionAskText        state.State = "send_ask_text"
)

func NewSendHandler() state.Handler {
	h := new(SendHandler)
	return h
}

func (h *SendHandler) Handle() th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		state.M.Enter(update.Message.Chat.ID, &state.ChatState{
			Handler: h,
			State:   SendActionAskPhoneNumber,
			Value:   &SMSValue{Modem: h.Modem(ctx)},
		})
		_, err := h.Reply(ctx, update, util.EscapeText("Enter the phone number you want to send the SMS to."), nil)
		return err
	}
}

func (h *SendHandler) HandleMessage(ctx *th.Context, message telego.Message, s *state.ChatState) error {
	value := s.Value.(*SMSValue)
	if s.State == SendActionAskPhoneNumber {
		value.To = message.Text
		state.M.Current(message.Chat.ID, SendActionAskText)
		_, err := h.ReplyMessage(ctx, message, util.EscapeText("Enter the text of the SMS you want to send."), nil)
		return err
	}
	if s.State == SendActionAskText {
	        defer state.M.Exit(message.Chat.ID) // 只在最后一步退出状态
	        _, err := value.Modem.SendSMS(value.To, message.Text)
       		 if err != nil {
        	    log.Printf("Failed to send SMS to %s: %v", value.To, err)
       	    		  h.ReplyMessage(ctx, message, util.EscapeText("Failed to send SMS: "+err.Error()), nil)
       		 	 return err
        }
    	    _, err = h.ReplyMessage(ctx, message, util.EscapeText("SMS sent successfully."), nil)
    	    if err != nil {
        	    log.Printf("Failed to send Telegram confirmation: %v", err)
        }
        return err
    }

    return nil
}

func (h *SendHandler) HandleCallbackQuery(ctx *th.Context, query telego.CallbackQuery, s *state.ChatState) error {
	return nil
}
