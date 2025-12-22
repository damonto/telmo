package modem

import (
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

func (s *Service) ListInserted() ([]InsertedResponse, error) {
	modems, err := s.manager.Modems()
	if err != nil {
		return nil, fmt.Errorf("listing modems: %w", err)
	}

	response := make([]InsertedResponse, 0, len(modems))
	for _, m := range modems {
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

		operatorName, err := threeGpp.OperatorName()
		if err != nil {
			return nil, fmt.Errorf("fetching operator name for %s: %w", m.EquipmentIdentifier, err)
		}

		operatorCode, err := threeGpp.OperatorCode()
		if err != nil {
			return nil, fmt.Errorf("fetching operator code for %s: %w", m.EquipmentIdentifier, err)
		}

		carrierInfo := carrier.Lookup(sim.OperatorIdentifier)
		isEsim, err := supportsEsim(m, s.cfg)
		if err != nil {
			return nil, fmt.Errorf("detecting eSIM support for %s: %w", m.EquipmentIdentifier, err)
		}

		alias := s.cfg.FindModem(m.EquipmentIdentifier).Alias

		response = append(response, InsertedResponse{
			Manufacturer:     m.Manufacturer,
			ID:               m.EquipmentIdentifier,
			FirmwareRevision: m.FirmwareRevision,
			HardwareRevision: m.HardwareRevision,
			Name:             cond.If(alias != "", alias, m.Model),
			Number:           m.Number,
			SIM: SIMResponse{
				OperatorName:       cond.If(sim.OperatorName != "", sim.OperatorName, carrierInfo.Name),
				OperatorIdentifier: sim.OperatorIdentifier,
				RegionCode:         carrierInfo.Region,
			},
			AccessTechnology:  accessTechnologyString(access),
			RegistrationState: registrationState.String(),
			RegisteredOperator: RegisteredOperatorResponse{
				Name: operatorName,
				Code: operatorCode,
			},
			SignalQuality: percent,
			IsEsim:        isEsim,
		})
	}

	return response, nil
}

func supportsEsim(m *mmodem.Modem, cfg *config.Config) (bool, error) {
	client, err := lpa.New(m, cfg)
	if err != nil {
		return false, nil
	}
	if err := client.Close(); err != nil {
		return false, fmt.Errorf("closing LPA client: %w", err)
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
