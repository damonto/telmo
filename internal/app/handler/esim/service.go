package esim

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/godbus/dbus/v5"

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

func (s *Service) Enable(modem *mmodem.Modem, iccid sgp22.ICCID) error {
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
	var lastSeq sgp22.SequenceNumber
	for _, notification := range notifications {
		lastSeq = max(lastSeq, notification.SequenceNumber)
	}

	if err := client.EnableProfile(iccid, true); err != nil {
		return fmt.Errorf("enabling profile %s on modem %s: %w", iccid.String(), modem.EquipmentIdentifier, err)
	}

	if err := modem.Restart(s.cfg.FindModem(modem.EquipmentIdentifier).Compatible); err != nil {
		return fmt.Errorf("restarting modem %s: %w", modem.EquipmentIdentifier, err)
	}

	go s.awaitModemAndSendNotifications(modem.EquipmentIdentifier, lastSeq)
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

func (s *Service) awaitModemAndSendNotifications(modemID string, lastSeq sgp22.SequenceNumber) {
	var once sync.Once
	err := s.manager.Subscribe(func(modems map[dbus.ObjectPath]*mmodem.Modem) error {
		var target *mmodem.Modem
		for _, modem := range modems {
			if modem.EquipmentIdentifier == modemID {
				target = modem
				break
			}
		}
		if target == nil {
			return nil
		}

		var sendErr error
		once.Do(func() {
			sendErr = s.sendPendingNotifications(target, lastSeq)
		})
		return sendErr
	})
	if err != nil {
		slog.Error("failed to subscribe to modem manager", "error", err, "modem", modemID)
	}
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
