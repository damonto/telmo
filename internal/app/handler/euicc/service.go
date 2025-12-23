package euicc

import (
	"fmt"
	"log/slog"

	"github.com/damonto/sigmo/internal/pkg/config"
	"github.com/damonto/sigmo/internal/pkg/lpa"
	mmodem "github.com/damonto/sigmo/internal/pkg/modem"
)

type Service struct {
	cfg *config.Config
}

func NewService(cfg *config.Config) *Service {
	return &Service{
		cfg: cfg,
	}
}

func (s *Service) Get(modem *mmodem.Modem) (*EuiccResponse, error) {
	client, err := lpa.New(modem, s.cfg)
	if err != nil {
		return nil, fmt.Errorf("modem %s does not support eSIM or failed to create LPA client: %w", modem.EquipmentIdentifier, err)
	}
	defer func() {
		if cerr := client.Close(); cerr != nil {
			slog.Warn("failed to close LPA client", "error", cerr)
		}
	}()

	info, err := client.Info()
	if err != nil {
		return nil, fmt.Errorf("fetching eUICC info for modem %s: %w", modem.EquipmentIdentifier, err)
	}
	return &EuiccResponse{
		EID:          info.EID,
		FreeSpace:    info.FreeSpace,
		SASUP:        info.SASUP,
		Certificates: info.Certificates,
	}, nil
}
