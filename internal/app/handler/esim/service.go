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
		slog.Error("failed to create LPA client", "modem", modem.EquipmentIdentifier, "error", err)
		return nil, err
	}
	defer func() {
		if cerr := client.Close(); cerr != nil {
			slog.Warn("failed to close LPA client", "error", cerr)
		}
	}()

	profiles, err := client.ListProfile(nil, nil)
	if err != nil {
		slog.Error("failed to list profiles", "modem", modem.EquipmentIdentifier, "error", err)
		return nil, err
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

func (s *Service) Discover(modem *mmodem.Modem) ([]DiscoverResponse, error) {
	client, err := lpa.New(modem, s.cfg)
	if err != nil {
		slog.Error("failed to create LPA client", "modem", modem.EquipmentIdentifier, "error", err)
		return nil, err
	}
	defer func() {
		if cerr := client.Close(); cerr != nil {
			slog.Warn("failed to close LPA client", "error", cerr)
		}
	}()

	imeiValue, err := modem.ThreeGPP().IMEI()
	if err != nil {
		slog.Error("failed to read modem IMEI", "modem", modem.EquipmentIdentifier, "error", err)
		return nil, err
	}
	imei, err := sgp22.NewIMEI(imeiValue)
	if err != nil {
		slog.Error("invalid IMEI", "modem", modem.EquipmentIdentifier, "imei", imeiValue, "error", err)
		return nil, err
	}

	entries, err := client.Discover(imei)
	if err != nil {
		slog.Error("failed to discover profiles", "modem", modem.EquipmentIdentifier, "error", err)
		return nil, err
	}

	response := make([]DiscoverResponse, 0, len(entries))
	for _, entry := range entries {
		response = append(response, DiscoverResponse{
			EventID: entry.EventID,
			Address: entry.Address,
		})
	}
	return response, nil
}

func (s *Service) Enable(ctx context.Context, modem *mmodem.Modem, iccid sgp22.ICCID) error {
	client, err := lpa.New(modem, s.cfg)
	if err != nil {
		slog.Error("failed to create LPA client", "modem", modem.EquipmentIdentifier, "error", err)
		return err
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
		slog.Error("failed to list notifications", "modem", modem.EquipmentIdentifier, "error", err)
		return err
	}
	var lastSeq sgp22.SequenceNumber
	for _, notification := range notifications {
		lastSeq = max(lastSeq, notification.SequenceNumber)
	}

	if err := client.EnableProfile(iccid, true); err != nil {
		slog.Error("failed to enable profile", "modem", modem.EquipmentIdentifier, "iccid", iccid.String(), "error", err)
		return err
	}

	closeClient()

	if err := modem.Restart(s.cfg.FindModem(modem.EquipmentIdentifier).Compatible); err != nil {
		slog.Error("failed to restart modem", "modem", modem.EquipmentIdentifier, "error", err)
		return err
	}

	target, err := s.manager.WaitForModem(ctx, modem.EquipmentIdentifier)
	if err != nil {
		slog.Error("failed to wait for modem", "modem", modem.EquipmentIdentifier, "error", err)
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
		slog.Error("failed to create LPA client", "modem", modem.EquipmentIdentifier, "error", err)
		return err
	}
	defer func() {
		if cerr := client.Close(); cerr != nil {
			slog.Warn("failed to close LPA client", "error", cerr)
		}
	}()

	if err := client.Delete(iccid); err != nil {
		slog.Error("failed to delete profile", "modem", modem.EquipmentIdentifier, "iccid", iccid.String(), "error", err)
		return err
	}
	return nil
}

func (s *Service) Download(ctx context.Context, modem *mmodem.Modem, activationCode *elpa.ActivationCode, opts *elpa.DownloadOptions) error {
	client, err := lpa.New(modem, s.cfg)
	if err != nil {
		slog.Error("failed to create LPA client", "modem", modem.EquipmentIdentifier, "error", err)
		return err
	}
	defer func() {
		if cerr := client.Close(); cerr != nil {
			slog.Warn("failed to close LPA client", "error", cerr)
		}
	}()

	if err := client.Download(ctx, activationCode, opts); err != nil {
		slog.Error("failed to download profile", "modem", modem.EquipmentIdentifier, "error", err)
		return err
	}
	return nil
}

func (s *Service) UpdateNickname(modem *mmodem.Modem, iccid sgp22.ICCID, nickname string) error {
	if err := validateNickname(nickname); err != nil {
		return err
	}
	client, err := lpa.New(modem, s.cfg)
	if err != nil {
		slog.Error("failed to create LPA client", "modem", modem.EquipmentIdentifier, "error", err)
		return err
	}
	defer func() {
		if cerr := client.Close(); cerr != nil {
			slog.Warn("failed to close LPA client", "error", cerr)
		}
	}()

	if err := client.SetNickname(iccid, nickname); err != nil {
		slog.Error("failed to set nickname", "modem", modem.EquipmentIdentifier, "iccid", iccid.String(), "error", err)
		return err
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
		slog.Error("failed to create LPA client", "modem", modem.EquipmentIdentifier, "error", err)
		return err
	}
	defer func() {
		if cerr := client.Close(); cerr != nil {
			slog.Warn("failed to close LPA client", "error", cerr)
		}
	}()
	notifications, err := client.ListNotification()
	if err != nil {
		slog.Error("failed to list notifications", "modem", modem.EquipmentIdentifier, "error", err)
		return err
	}
	var errs error
	for _, notification := range notifications {
		if notification.SequenceNumber <= lastSeq {
			continue
		}
		if err := client.SendNotification(notification.SequenceNumber, true); err != nil {
			slog.Error("failed to send notification", "sequence", notification.SequenceNumber, "error", err)
			errs = errors.Join(errs, err)
		}
	}
	return errs
}
