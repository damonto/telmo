package modem

import (
	"errors"
	"fmt"
	"slices"

	"github.com/damonto/sigmo/internal/pkg/carrier"
	"github.com/damonto/sigmo/internal/pkg/cond"
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

func (s *Service) List() ([]ModemResponse, error) {
	modems, err := s.manager.Modems()
	if err != nil {
		return nil, fmt.Errorf("listing modems: %w", err)
	}
	response := make([]ModemResponse, 0, len(modems))
	for _, m := range modems {
		modemResp, err := s.buildModemResponse(m)
		if err != nil {
			return nil, err
		}
		response = append(response, modemResp)
	}
	return response, nil
}

func (s *Service) Get(modem *mmodem.Modem) (*ModemResponse, error) {
	resp, err := s.buildModemResponse(modem)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *Service) buildModemResponse(m *mmodem.Modem) (ModemResponse, error) {
	sim, err := m.SIMs().Primary()
	if err != nil {
		return ModemResponse{}, fmt.Errorf("fetching SIM for %s: %w", m.EquipmentIdentifier, err)
	}

	percent, _, err := m.SignalQuality()
	if err != nil {
		return ModemResponse{}, fmt.Errorf("fetching signal quality for %s: %w", m.EquipmentIdentifier, err)
	}

	access, err := m.AccessTechnologies()
	if err != nil {
		return ModemResponse{}, fmt.Errorf("fetching access technologies for %s: %w", m.EquipmentIdentifier, err)
	}

	threeGpp := m.ThreeGPP()
	registrationState, err := threeGpp.RegistrationState()
	if err != nil {
		return ModemResponse{}, fmt.Errorf("fetching registration state for %s: %w", m.EquipmentIdentifier, err)
	}

	operatorName, err := threeGpp.OperatorName()
	if err != nil {
		return ModemResponse{}, fmt.Errorf("fetching operator name for %s: %w", m.EquipmentIdentifier, err)
	}

	operatorCode, err := threeGpp.OperatorCode()
	if err != nil {
		return ModemResponse{}, fmt.Errorf("fetching operator code for %s: %w", m.EquipmentIdentifier, err)
	}

	carrierInfo := carrier.Lookup(sim.OperatorIdentifier)
	supportsEsim, err := supportsEsim(m, s.cfg)
	if err != nil {
		return ModemResponse{}, fmt.Errorf("detecting eSIM support for %s: %w", m.EquipmentIdentifier, err)
	}

	simSlots, err := s.buildSimSlotsResponse(m)
	if err != nil {
		return ModemResponse{}, fmt.Errorf("fetching SIM slots for %s: %w", m.EquipmentIdentifier, err)
	}

	alias := s.cfg.FindModem(m.EquipmentIdentifier).Alias
	return ModemResponse{
		Manufacturer:     m.Manufacturer,
		ID:               m.EquipmentIdentifier,
		FirmwareRevision: m.FirmwareRevision,
		HardwareRevision: m.HardwareRevision,
		Name:             cond.If(alias != "", alias, m.Model),
		Number:           m.Number,
		SIM: SlotResponse{
			Active:             sim.Active,
			OperatorName:       cond.If(sim.OperatorName != "", sim.OperatorName, carrierInfo.Name),
			OperatorIdentifier: sim.OperatorIdentifier,
			RegionCode:         carrierInfo.Region,
			Identifier:         sim.Identifier,
		},
		Slots:             simSlots,
		AccessTechnology:  accessTechnologyString(access),
		RegistrationState: registrationState.String(),
		RegisteredOperator: RegisteredOperatorResponse{
			Name: operatorName,
			Code: operatorCode,
		},
		SignalQuality: percent,
		SupportsEsim:  supportsEsim,
	}, nil
}

func (s *Service) buildSimSlotsResponse(m *mmodem.Modem) ([]SlotResponse, error) {
	if len(m.SimSlots) == 0 {
		return nil, nil
	}
	simSlots := make([]SlotResponse, 0, len(m.SimSlots))
	for _, slotPath := range m.SimSlots {
		sim, err := m.SIMs().Get(slotPath)
		if err != nil {
			return nil, fmt.Errorf("fetching SIM for slot %s: %w", slotPath, err)
		}
		carrierInfo := carrier.Lookup(sim.OperatorIdentifier)
		simSlots = append(simSlots, SlotResponse{
			Active:             sim.Active,
			OperatorName:       cond.If(sim.OperatorName != "", sim.OperatorName, carrierInfo.Name),
			OperatorIdentifier: sim.OperatorIdentifier,
			RegionCode:         carrierInfo.Region,
			Identifier:         sim.Identifier,
		})
	}
	return simSlots, nil
}

func supportsEsim(m *mmodem.Modem, cfg *config.Config) (bool, error) {
	client, err := lpa.New(m, cfg)
	if err != nil {
		if errors.Is(err, lpa.ErrNoSupportedAID) {
			return false, nil
		}
		return false, fmt.Errorf("creating LPA client for %s: %w", m.EquipmentIdentifier, err)
	}
	if err := client.Close(); err != nil {
		return false, fmt.Errorf("closing LPA client for %s: %w", m.EquipmentIdentifier, err)
	}
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
