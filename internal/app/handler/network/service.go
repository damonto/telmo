package network

import (
	"errors"
	"log/slog"
	"strings"

	mmodem "github.com/damonto/sigmo/internal/pkg/modem"
)

type Service struct{}

var errOperatorCodeRequired = errors.New("operator code is required")

func NewService() *Service {
	return &Service{}
}

func (s *Service) List(modem *mmodem.Modem) ([]NetworkResponse, error) {
	networks, err := modem.ThreeGPP().ScanNetworks()
	if err != nil {
		slog.Error("failed to scan networks", "modem", modem.EquipmentIdentifier, "error", err)
		return nil, err
	}

	response := make([]NetworkResponse, 0, len(networks))
	for _, network := range networks {
		response = append(response, NetworkResponse{
			Status:             network.Status.String(),
			OperatorName:       network.OperatorName,
			OperatorShortName:  network.OperatorShortName,
			OperatorCode:       network.OperatorCode,
			AccessTechnologies: accessTechnologyStrings(network.AccessTechnology),
		})
	}
	return response, nil
}

func (s *Service) Register(modem *mmodem.Modem, operatorCode string) error {
	operatorCode = strings.TrimSpace(operatorCode)
	if operatorCode == "" {
		return errOperatorCodeRequired
	}
	if err := modem.ThreeGPP().RegisterNetwork(operatorCode); err != nil {
		slog.Error("failed to register network", "modem", modem.EquipmentIdentifier, "operator", operatorCode, "error", err)
		return err
	}
	return nil
}

func accessTechnologyStrings(access []mmodem.ModemAccessTechnology) []string {
	names := make([]string, 0, len(access))
	for _, tech := range access {
		names = append(names, tech.String())
	}
	return names
}
