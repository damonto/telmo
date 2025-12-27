package esim

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"unicode/utf8"

	elpa "github.com/damonto/euicc-go/lpa"
	sgp22 "github.com/damonto/euicc-go/v2"
	"github.com/damonto/sigmo/internal/pkg/carrier"
	"github.com/damonto/sigmo/internal/pkg/config"
	"github.com/damonto/sigmo/internal/pkg/lpa"
	mmodem "github.com/damonto/sigmo/internal/pkg/modem"
)

type Service struct {
	cfg     *config.Config
	manager *mmodem.Manager
}

var errInvalidNickname = errors.New("nickname must be valid utf-8 and 64 bytes or fewer")

func NewService(cfg *config.Config, manager *mmodem.Manager) *Service {
	return &Service{
		cfg:     cfg,
		manager: manager,
	}
}

func (s *Service) List(modem *mmodem.Modem) ([]ProfileResponse, error) {
	client, err := lpa.New(modem, s.cfg)
	if err != nil {
		return nil, fmt.Errorf("creating LPA client for modem %s: %w", modem.EquipmentIdentifier, err)
	}
	defer func() {
		if cerr := client.Close(); cerr != nil {
			slog.Warn("failed to close LPA client", "error", cerr)
		}
	}()

	profiles, err := client.ListProfile(nil, nil)
	if err != nil {
		return nil, fmt.Errorf("listing profiles for modem %s: %w", modem.EquipmentIdentifier, err)
	}

	response := make([]ProfileResponse, 0, len(profiles))
	for _, profile := range profiles {
		name := profile.ProfileNickname
		if name == "" {
			name = profile.ProfileName
		}
		carrierInfo := carrier.Lookup(profile.ProfileOwner.MCC() + profile.ProfileOwner.MNC())
		icon := ""
		if fileType := profile.Icon.FileType(); fileType != "" {
			icon = fmt.Sprintf("data:%s;base64,%s", fileType, base64.StdEncoding.EncodeToString(profile.Icon))
		}
		regionCode := carrierInfo.Region
		response = append(response, ProfileResponse{
			Name:                name,
			ServiceProviderName: profile.ServiceProviderName,
			ICCID:               profile.ICCID.String(),
			Icon:                icon,
			ProfileState:        uint8(profile.ProfileState),
			RegionCode:          regionCode,
		})
	}
	return response, nil
}

func (s *Service) Enable(ctx context.Context, modem *mmodem.Modem, iccid sgp22.ICCID) error {
	client, err := lpa.New(modem, s.cfg)
	if err != nil {
		return fmt.Errorf("creating LPA client for modem %s: %w", modem.EquipmentIdentifier, err)
	}
	closeClient := func() {
		if client == nil {
			return
		}
		if cerr := client.Close(); cerr != nil {
			slog.Warn("failed to close LPA client", "error", cerr)
		}
		client = nil
	}
	defer closeClient()

	notifications, err := client.ListNotification()
	if err != nil {
		return fmt.Errorf("listing notifications for modem %s: %w", modem.EquipmentIdentifier, err)
	}
	var lastSeq sgp22.SequenceNumber
	for _, notification := range notifications {
		lastSeq = max(lastSeq, notification.SequenceNumber)
	}

	if err := client.EnableProfile(iccid, true); err != nil {
		slog.Error("enabling profile", "error", err)
		return fmt.Errorf("enabling profile %s on modem %s: %w", iccid.String(), modem.EquipmentIdentifier, err)
	}

	closeClient()

	if err := modem.Restart(s.cfg.FindModem(modem.EquipmentIdentifier).Compatible); err != nil {
		return fmt.Errorf("restarting modem %s: %w", modem.EquipmentIdentifier, err)
	}

	target, err := s.manager.WaitForModem(ctx, modem.EquipmentIdentifier)
	if err != nil {
		return err
	}
	if err := s.sendPendingNotifications(target, lastSeq); err != nil {
		slog.Warn("failed to handle modem notifications", "error", err, "modem", modem.EquipmentIdentifier)
	}
	return nil
}

func (s *Service) Delete(modem *mmodem.Modem, iccid sgp22.ICCID) error {
	client, err := lpa.New(modem, s.cfg)
	if err != nil {
		return fmt.Errorf("creating LPA client for modem %s: %w", modem.EquipmentIdentifier, err)
	}
	defer func() {
		if cerr := client.Close(); cerr != nil {
			slog.Warn("failed to close LPA client", "error", cerr)
		}
	}()

	if _, err := client.Delete(iccid); err != nil {
		return fmt.Errorf("deleting profile %s on modem %s: %w", iccid.String(), modem.EquipmentIdentifier, err)
	}
	return nil
}

func (s *Service) Download(ctx context.Context, modem *mmodem.Modem, activationCode *elpa.ActivationCode, opts *elpa.DownloadOptions) error {
	client, err := lpa.New(modem, s.cfg)
	if err != nil {
		return fmt.Errorf("creating LPA client for modem %s: %w", modem.EquipmentIdentifier, err)
	}
	defer func() {
		if cerr := client.Close(); cerr != nil {
			slog.Warn("failed to close LPA client", "error", cerr)
		}
	}()

	if err := client.Download(ctx, activationCode, opts); err != nil {
		return fmt.Errorf("downloading profile on modem %s: %w", modem.EquipmentIdentifier, err)
	}
	return nil
}

func (s *Service) UpdateNickname(modem *mmodem.Modem, iccid sgp22.ICCID, nickname string) error {
	if err := validateNickname(nickname); err != nil {
		return err
	}
	client, err := lpa.New(modem, s.cfg)
	if err != nil {
		return fmt.Errorf("creating LPA client for modem %s: %w", modem.EquipmentIdentifier, err)
	}
	defer func() {
		if cerr := client.Close(); cerr != nil {
			slog.Warn("failed to close LPA client", "error", cerr)
		}
	}()

	if err := client.SetNickname(iccid, nickname); err != nil {
		return fmt.Errorf("setting nickname for profile %s on modem %s: %w", iccid.String(), modem.EquipmentIdentifier, err)
	}
	return nil
}

func validateNickname(nickname string) error {
	if !utf8.ValidString(nickname) {
		return errInvalidNickname
	}
	if len(nickname) > 64 {
		return errInvalidNickname
	}
	return nil
}

func (s *Service) sendPendingNotifications(modem *mmodem.Modem, lastSeq sgp22.SequenceNumber) error {
	client, err := lpa.New(modem, s.cfg)
	if err != nil {
		return fmt.Errorf("creating LPA client for modem %s: %w", modem.EquipmentIdentifier, err)
	}
	defer func() {
		if cerr := client.Close(); cerr != nil {
			slog.Warn("failed to close LPA client", "error", cerr)
		}
	}()

	notifications, err := client.ListNotification()
	if err != nil {
		return fmt.Errorf("listing notifications for modem %s: %w", modem.EquipmentIdentifier, err)
	}

	var errs error
	for _, notification := range notifications {
		if notification.SequenceNumber <= lastSeq {
			continue
		}
		pending, err := client.RetrieveNotificationList(notification.SequenceNumber)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("retrieve notification %d: %w", notification.SequenceNumber, err))
			continue
		}
		if len(pending) == 0 {
			continue
		}
		if err := client.HandleNotification(pending[0]); err != nil {
			errs = errors.Join(errs, fmt.Errorf("handle notification %d: %w", notification.SequenceNumber, err))
			continue
		}
		if err := client.RemoveNotificationFromList(notification.SequenceNumber); err != nil {
			errs = errors.Join(errs, fmt.Errorf("remove notification %d: %w", notification.SequenceNumber, err))
			continue
		}
	}
	return errs
}
