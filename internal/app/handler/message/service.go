package message

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	mmodem "github.com/damonto/sigmo/internal/pkg/modem"
)

type Service struct{}

var (
	errParticipantRequired = errors.New("participant is required")
	errRecipientRequired   = errors.New("recipient is required")
	errTextRequired        = errors.New("text is required")
)

func NewService() *Service {
	return &Service{}
}

func (s *Service) ListConversations(modem *mmodem.Modem) ([]MessageResponse, error) {
	messages, err := modem.Messaging().List()
	if err != nil {
		return nil, fmt.Errorf("listing messages for modem %s: %w", modem.EquipmentIdentifier, err)
	}

	latest := make(map[string]*mmodem.SMS, len(messages))
	for _, sms := range messages {
		key := strings.TrimSpace(sms.Number)
		existing, ok := latest[key]
		if !ok || sms.Timestamp.After(existing.Timestamp) {
			latest[key] = sms
		}
	}

	response := make([]MessageResponse, 0, len(latest))
	for _, sms := range latest {
		response = append(response, buildMessageResponse(sms))
	}

	slices.SortFunc(response, func(a, b MessageResponse) int {
		if a.ID == b.ID {
			return 0
		}
		if a.ID > b.ID {
			return -1
		}
		return 1
	})
	return response, nil
}

func (s *Service) ListByParticipant(modem *mmodem.Modem, participant string) ([]MessageResponse, error) {
	if strings.TrimSpace(participant) == "" {
		return nil, errParticipantRequired
	}
	messages, err := modem.Messaging().List()
	if err != nil {
		return nil, fmt.Errorf("listing messages for modem %s: %w", modem.EquipmentIdentifier, err)
	}

	response := make([]MessageResponse, 0, len(messages))
	for _, sms := range messages {
		if strings.TrimSpace(sms.Number) != participant {
			continue
		}
		response = append(response, buildMessageResponse(sms))
	}
	slices.SortFunc(response, func(a, b MessageResponse) int {
		if a.ID == b.ID {
			return 0
		}
		if a.ID < b.ID {
			return -1
		}
		return 1
	})
	return response, nil
}

func (s *Service) Send(modem *mmodem.Modem, to string, text string) error {
	if strings.TrimSpace(to) == "" {
		return errRecipientRequired
	}
	if strings.TrimSpace(text) == "" {
		return errTextRequired
	}
	_, err := modem.Messaging().Send(to, text)
	if err != nil {
		return fmt.Errorf("sending SMS to %s on modem %s: %w", to, modem.EquipmentIdentifier, err)
	}
	return nil
}

func (s *Service) DeleteByParticipant(modem *mmodem.Modem, participant string) error {
	if strings.TrimSpace(participant) == "" {
		return errParticipantRequired
	}
	messages, err := modem.Messaging().List()
	if err != nil {
		return fmt.Errorf("listing messages for modem %s: %w", modem.EquipmentIdentifier, err)
	}
	messaging := modem.Messaging()
	for _, sms := range messages {
		if strings.TrimSpace(sms.Number) != participant {
			continue
		}
		if err := messaging.Delete(sms.Path()); err != nil {
			return fmt.Errorf("deleting message for participant %s on modem %s: %w", participant, modem.EquipmentIdentifier, err)
		}
	}
	return nil
}

func buildMessageResponse(sms *mmodem.SMS) MessageResponse {
	incoming := sms.State == mmodem.SMSStateReceived || sms.State == mmodem.SMSStateReceiving
	remote := strings.TrimSpace(sms.Number)
	return MessageResponse{
		ID:        messageID(sms),
		Sender:    remote,
		Recipient: remote,
		Text:      sms.Text,
		Timestamp: sms.Timestamp,
		Status:    strings.ToLower(sms.State.String()),
		Incoming:  incoming,
	}
}

func messageID(sms *mmodem.SMS) int64 {
	path := string(sms.Path())
	if path == "" {
		return 0
	}
	idx := strings.LastIndex(path, "/")
	if idx == -1 || idx+1 >= len(path) {
		return 0
	}
	id, err := strconv.ParseInt(path[idx+1:], 10, 64)
	if err != nil {
		return 0
	}
	return id
}
