package modem

import (
	"context"
	"errors"
	"fmt"
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
		return nil, fmt.Errorf("listing modems: %w", err)
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
		return fmt.Errorf("setting primary SIM slot for %s: %w", modem.EquipmentIdentifier, err)
	}
	_, err = s.manager.WaitForModem(ctx, modem.EquipmentIdentifier)
	return err
}

func (s *Service) UpdateMSISDN(ctx context.Context, modem *mmodem.Modem, number string) error {
	number = strings.TrimSpace(number)
	if !msisdnPhoneRE.MatchString(number) {
		return errMSISDNInvalidNumber
	}
	port, err := modem.Port(mmodem.ModemPortTypeAt)
	if err != nil {
		return fmt.Errorf("finding AT port for modem %s: %w", modem.EquipmentIdentifier, err)
	}
	client, err := msisdn.New(port.Device)
	if err != nil {
		return fmt.Errorf("opening MSISDN client for modem %s: %w", modem.EquipmentIdentifier, err)
	}
	defer func() {
		if cerr := client.Close(); cerr != nil {
			slog.Warn("failed to close MSISDN client", "error", cerr, "modem", modem.EquipmentIdentifier)
		}
	}()
	if err := client.Update("", number); err != nil {
		return fmt.Errorf("updating MSISDN for modem %s: %w", modem.EquipmentIdentifier, err)
	}
	if err := modem.Restart(s.cfg.FindModem(modem.EquipmentIdentifier).Compatible); err != nil {
		return fmt.Errorf("restarting modem %s: %w", modem.EquipmentIdentifier, err)
	}
	_, err = s.manager.WaitForModem(ctx, modem.EquipmentIdentifier)
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
		return nil, fmt.Errorf("fetching SIM for %s: %w", m.EquipmentIdentifier, err)
	}

	percent, _, err := m.SignalQuality()
	if err != nil {
		return nil, fmt.Errorf("fetching signal quality for %s: %w", m.EquipmentIdentifier, err)
	}

	access, err := m.AccessTechnologies()
	if err != nil {
		return nil, fmt.Errorf("fetching access technologies for %s: %w", m.EquipmentIdentifier, err)
	}

	threeGpp := m.ThreeGPP()
	registrationState, err := threeGpp.RegistrationState()
	if err != nil {
		return nil, fmt.Errorf("fetching registration state for %s: %w", m.EquipmentIdentifier, err)
	}

	registeredOperatorName, err := threeGpp.OperatorName()
	if err != nil {
		return nil, fmt.Errorf("fetching operator name for %s: %w", m.EquipmentIdentifier, err)
	}

	operatorCode, err := threeGpp.OperatorCode()
	if err != nil {
		return nil, fmt.Errorf("fetching operator code for %s: %w", m.EquipmentIdentifier, err)
	}

	carrierInfo := carrier.Lookup(sim.OperatorIdentifier)
	supportsEsim, err := supportsEsim(m, s.cfg)
	if err != nil {
		return nil, fmt.Errorf("detecting eSIM support for %s: %w", m.EquipmentIdentifier, err)
	}

	simSlots, err := s.buildSimSlotsResponse(m)
	if err != nil {
		return nil, fmt.Errorf("fetching SIM slots for %s: %w", m.EquipmentIdentifier, err)
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
			return nil, fmt.Errorf("fetching SIM for slot %s: %w", slotPath, err)
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
			return 0, fmt.Errorf("fetching SIM for slot %s: %w", slotPath, err)
		}
		if sim.Identifier == identifier && !sim.Active {
			return uint32(index + 1), nil
		}
	}
	return 0, fmt.Errorf("sim slot %q not found: %w", identifier, errSimSlotNotFound)
}

func supportsEsim(m *mmodem.Modem, cfg *config.Config) (bool, error) {
	client, err := lpa.New(m, cfg)
	if err != nil {
		if errors.Is(err, lpa.ErrNoSupportedAID) {
			return false, nil
		}
		return false, fmt.Errorf("creating LPA client for %s: %w", m.EquipmentIdentifier, err)
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
