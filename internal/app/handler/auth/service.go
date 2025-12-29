package auth

import (
	"errors"
	"fmt"

	"github.com/damonto/sigmo/internal/app/auth"
	"github.com/damonto/sigmo/internal/pkg/config"
	"github.com/damonto/sigmo/internal/pkg/notify"
)

var (
	errNoChannelsConfigured      = errors.New("no notification channels configured")
	errAuthProviderRequired      = errors.New("auth provider is required")
	errAuthProviderNotConfigured = errors.New("auth provider is not configured")
	errInvalidOTP                = errors.New("invalid otp")
)

type Service struct {
	cfg   *config.Config
	store *auth.Store
}

func NewService(cfg *config.Config, store *auth.Store) *Service {
	return &Service{
		cfg:   cfg,
		store: store,
	}
}

func (s *Service) OTPRequired() bool {
	return s.cfg.App.OTPRequired
}

func (s *Service) SendOTP() error {
	if !s.OTPRequired() {
		return nil
	}
	channels, err := buildChannels(s.cfg)
	if err != nil {
		return err
	}
	if len(channels) == 0 {
		return errNoChannelsConfigured
	}
	code, _, err := s.store.IssueOTP()
	if err != nil {
		return err
	}
	msg := notify.TextMessage{Text: fmt.Sprintf("Your verification code is %s", code)}
	var combined error
	for name, channel := range channels {
		if err := notify.Send(channel, msg); err != nil {
			combined = errors.Join(combined, fmt.Errorf("sending via %s: %w", name, err))
		}
	}
	return combined
}

func (s *Service) VerifyOTP(code string) (string, error) {
	if s.OTPRequired() && !s.store.VerifyOTP(code) {
		return "", errInvalidOTP
	}
	token, _, err := s.store.IssueToken()
	if err != nil {
		return "", err
	}
	return token, nil
}

func buildChannels(cfg *config.Config) (map[string]notify.Sender, error) {
	channels := make(map[string]notify.Sender)
	if len(cfg.App.AuthProviders) == 0 {
		return nil, errAuthProviderRequired
	}
	for _, provider := range cfg.App.AuthProviders {
		selected, ok := cfg.Channels[provider]
		if !ok {
			return nil, fmt.Errorf("%w: %s", errAuthProviderNotConfigured, provider)
		}
		switch provider {
		case "telegram":
			telegram, err := notify.NewTelegram(&selected)
			if err != nil {
				return nil, err
			}
			channels[provider] = telegram
		case "http":
			httpChannel, err := notify.NewHTTP(&selected)
			if err != nil {
				return nil, err
			}
			channels[provider] = httpChannel
		default:
			return nil, fmt.Errorf("unsupported auth provider %s", provider)
		}
	}
	return channels, nil
}
