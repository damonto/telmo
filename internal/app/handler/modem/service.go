package modem

import (
	"context"
	"errors"
	"log/slog"
	"regexp"
	"slices"
	"strings"

	"github.com/damonto/sigmo/internal/pkg/carrier"
	"github.com/damonto/sigmo/internal/pkg/config"
	"github.com/damonto/sigmo/internal/pkg/lpa"
	mmodem "github.com/damonto/sigmo/internal/pkg/modem"
	"github.com/damonto/sigmo/internal/pkg/modem/msisdn"
)

type Service struct {
	cfg     *config.Config
	manager *mmodem.Manager
}

var (
	errSimIdentifierRequired = errors.New("identifier is required")
	errSimSlotsUnavailable   = errors.New("sim slots not available")
	errSimSlotNotFound       = errors.New("sim slot not found")
	errSimSlotAlreadyActive  = errors.New("sim slot already active")
	errMSISDNInvalidNumber   = errors.New("invalid phone number")
	errCompatibleRequired    = errors.New("compatible is required")
)

var msisdnPhoneRE = regexp.MustCompile(`^\+?[0-9]{1,15}$`)

func NewService(cfg *config.Config, manager *mmodem.Manager) *Service {
	return &Service{
		cfg:     cfg,
		manager: manager,
	}
}

func (s *Service) List() ([]*ModemResponse, error) {
	modems, err := s.manager.Modems()
	if err != nil {
		slog.Error("failed to list modems", "error", err)
		return nil, err
	}
	response := make([]*ModemResponse, 0, len(modems))
	for _, m := range modems {
		modemResp, err := s.buildModemResponse(m)
		if err != nil {
			return nil, err
		}
		response = append(response, modemResp)
	}
	slices.SortFunc(response, func(a, b *ModemResponse) int {
		return strings.Compare(a.ID, b.ID)
	})
	return response, nil
}

func (s *Service) Get(modem *mmodem.Modem) (*ModemResponse, error) {
	resp, err := s.buildModemResponse(modem)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *Service) SwitchSimSlot(ctx context.Context, modem *mmodem.Modem, identifier string) error {
	if identifier == "" {
		return errSimIdentifierRequired
	}
	slotIndex, err := s.findSimSlotIndex(modem, identifier)
	if err != nil {
		return err
	}
	if err := modem.SetPrimarySimSlot(slotIndex); err != nil {
		slog.Error("failed to set primary SIM slot", "modem", modem.EquipmentIdentifier, "error", err)
		return err
	}
	_, err = s.manager.WaitForModem(ctx, modem.EquipmentIdentifier)
	if err != nil {
		slog.Error("failed to wait for modem", "modem", modem.EquipmentIdentifier, "error", err)
	}
	return err
}

func (s *Service) UpdateMSISDN(ctx context.Context, modem *mmodem.Modem, number string) error {
	number = strings.TrimSpace(number)
	if !msisdnPhoneRE.MatchString(number) {
		return errMSISDNInvalidNumber
	}
	port, err := modem.Port(mmodem.ModemPortTypeAt)
	if err != nil {
		slog.Error("failed to find AT port", "modem", modem.EquipmentIdentifier, "error", err)
		return err
	}
	client, err := msisdn.New(port.Device)
	if err != nil {
		slog.Error("failed to open MSISDN client", "modem", modem.EquipmentIdentifier, "error", err)
		return err
	}
	defer func() {
		if cerr := client.Close(); cerr != nil {
			slog.Warn("failed to close MSISDN client", "error", cerr, "modem", modem.EquipmentIdentifier)
		}
	}()
	if err := client.Update("", number); err != nil {
		slog.Error("failed to update MSISDN", "modem", modem.EquipmentIdentifier, "error", err)
		return err
	}
	if err := modem.Restart(s.cfg.FindModem(modem.EquipmentIdentifier).Compatible); err != nil {
		slog.Error("failed to restart modem", "modem", modem.EquipmentIdentifier, "error", err)
		return err
	}
	_, err = s.manager.WaitForModem(ctx, modem.EquipmentIdentifier)
	if err != nil {
		slog.Error("failed to wait for modem", "modem", modem.EquipmentIdentifier, "error", err)
	}
	return err
}

func (s *Service) UpdateSettings(modemID string, req UpdateModemSettingsRequest) error {
	if req.Compatible == nil {
		return errCompatibleRequired
	}
	modem := s.cfg.FindModem(modemID)
	modem.Alias = strings.TrimSpace(req.Alias)
	modem.Compatible = *req.Compatible
	modem.MSS = req.MSS
	if s.cfg.Modems == nil {
		s.cfg.Modems = make(map[string]config.Modem)
	}
	s.cfg.Modems[modemID] = modem
	if err := s.cfg.Save(); err != nil {
		slog.Error("failed to save config", "modem", modemID, "error", err)
		return err
	}
	return nil
}

func (s *Service) GetSettings(modemID string) *ModemSettingsResponse {
	modem := s.cfg.FindModem(modemID)
	return &ModemSettingsResponse{
		Alias:      modem.Alias,
		Compatible: modem.Compatible,
		MSS:        modem.MSS,
	}
}

func (s *Service) buildModemResponse(m *mmodem.Modem) (*ModemResponse, error) {
	sim, err := m.SIMs().Primary()
	if err != nil {
		slog.Error("failed to fetch SIM", "modem", m.EquipmentIdentifier, "error", err)
		return nil, err
	}

	percent, _, err := m.SignalQuality()
	if err != nil {
		slog.Error("failed to fetch signal quality", "modem", m.EquipmentIdentifier, "error", err)
		return nil, err
	}

	access, err := m.AccessTechnologies()
	if err != nil {
		slog.Error("failed to fetch access technologies", "modem", m.EquipmentIdentifier, "error", err)
		return nil, err
	}

	threeGpp := m.ThreeGPP()
	registrationState, err := threeGpp.RegistrationState()
	if err != nil {
		slog.Error("failed to fetch registration state", "modem", m.EquipmentIdentifier, "error", err)
		return nil, err
	}

	registeredOperatorName, err := threeGpp.OperatorName()
	if err != nil {
		slog.Error("failed to fetch operator name", "modem", m.EquipmentIdentifier, "error", err)
		return nil, err
	}

	operatorCode, err := threeGpp.OperatorCode()
	if err != nil {
		slog.Error("failed to fetch operator code", "modem", m.EquipmentIdentifier, "error", err)
		return nil, err
	}

	carrierInfo := carrier.Lookup(sim.OperatorIdentifier)
	supportsEsim, err := supportsEsim(m, s.cfg)
	if err != nil {
		slog.Error("failed to detect eSIM support", "modem", m.EquipmentIdentifier, "error", err)
		return nil, err
	}

	simSlots, err := s.buildSimSlotsResponse(m)
	if err != nil {
		slog.Error("failed to fetch SIM slots", "modem", m.EquipmentIdentifier, "error", err)
		return nil, err
	}

	alias := s.cfg.FindModem(m.EquipmentIdentifier).Alias
	name := m.Model
	if alias != "" {
		name = alias
	}
	simOperatorName := carrierInfo.Name
	if sim.OperatorName != "" {
		simOperatorName = sim.OperatorName
	}
	return &ModemResponse{
		Manufacturer:     m.Manufacturer,
		ID:               m.EquipmentIdentifier,
		FirmwareRevision: m.FirmwareRevision,
		HardwareRevision: m.HardwareRevision,
		Name:             name,
		Number:           m.Number,
		SIM: SlotResponse{
			Active:             sim.Active,
			OperatorName:       simOperatorName,
			OperatorIdentifier: sim.OperatorIdentifier,
			RegionCode:         carrierInfo.Region,
			Identifier:         sim.Identifier,
		},
		Slots:             simSlots,
		AccessTechnology:  accessTechnologyString(access),
		RegistrationState: registrationState.String(),
		RegisteredOperator: RegisteredOperatorResponse{
			Name: registeredOperatorName,
			Code: operatorCode,
		},
		SignalQuality: percent,
		SupportsEsim:  supportsEsim,
	}, nil
}

func (s *Service) buildSimSlotsResponse(m *mmodem.Modem) ([]SlotResponse, error) {
	if len(m.SimSlots) == 0 {
		return []SlotResponse{}, nil
	}
	simSlots := make([]SlotResponse, 0, len(m.SimSlots))
	for _, slotPath := range m.SimSlots {
		sim, err := m.SIMs().Get(slotPath)
		if err != nil {
			slog.Error("failed to fetch SIM for slot", "modem", m.EquipmentIdentifier, "slot", slotPath, "error", err)
			return nil, err
		}
		carrierInfo := carrier.Lookup(sim.OperatorIdentifier)
		operatorName := carrierInfo.Name
		if sim.OperatorName != "" {
			operatorName = sim.OperatorName
		}
		simSlots = append(simSlots, SlotResponse{
			Active:             sim.Active,
			OperatorName:       operatorName,
			OperatorIdentifier: sim.OperatorIdentifier,
			RegionCode:         carrierInfo.Region,
			Identifier:         sim.Identifier,
		})
	}
	return simSlots, nil
}

func (s *Service) findSimSlotIndex(modem *mmodem.Modem, identifier string) (uint32, error) {
	if len(modem.SimSlots) == 0 {
		return 0, errSimSlotsUnavailable
	}
	for index, slotPath := range modem.SimSlots {
		sim, err := modem.SIMs().Get(slotPath)
		if err != nil {
			slog.Error("failed to fetch SIM for slot", "modem", modem.EquipmentIdentifier, "slot", slotPath, "error", err)
			return 0, err
		}
		if sim.Identifier == identifier && !sim.Active {
			return uint32(index + 1), nil
		}
	}
	return 0, errSimSlotNotFound
}

func supportsEsim(m *mmodem.Modem, cfg *config.Config) (bool, error) {
	client, err := lpa.New(m, cfg)
	if err != nil {
		if errors.Is(err, lpa.ErrNoSupportedAID) {
			return false, nil
		}
		slog.Error("failed to create LPA client", "modem", m.EquipmentIdentifier, "error", err)
		return false, err
	}
	defer func() {
		if cerr := client.Close(); cerr != nil {
			slog.Warn("failed to close LPA client", "error", cerr)
		}
	}()
	return true, nil
}

func accessTechnologyString(access []mmodem.ModemAccessTechnology) string {
	if len(access) == 0 {
		return ""
	}
	priority := []mmodem.ModemAccessTechnology{
		mmodem.ModemAccessTechnology5GNR,
		mmodem.ModemAccessTechnologyLte,
		mmodem.ModemAccessTechnologyLteCatM,
		mmodem.ModemAccessTechnologyLteNBIot,
		mmodem.ModemAccessTechnologyHspaPlus,
		mmodem.ModemAccessTechnologyHspa,
		mmodem.ModemAccessTechnologyHsupa,
		mmodem.ModemAccessTechnologyHsdpa,
		mmodem.ModemAccessTechnologyUmts,
		mmodem.ModemAccessTechnologyEdge,
		mmodem.ModemAccessTechnologyGprs,
		mmodem.ModemAccessTechnologyGsm,
		mmodem.ModemAccessTechnologyGsmCompact,
		mmodem.ModemAccessTechnologyEvdob,
		mmodem.ModemAccessTechnologyEvdoa,
		mmodem.ModemAccessTechnologyEvdo0,
		mmodem.ModemAccessTechnology1xrtt,
		mmodem.ModemAccessTechnologyPots,
	}
	for _, tech := range priority {
		if slices.Contains(access, tech) {
			return tech.String()
		}
	}
	return access[0].String()
}
