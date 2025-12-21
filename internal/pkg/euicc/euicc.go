package euicc

//go:generate curl -L -o ci.json https://euicc-manual.osmocom.org/docs/pki/ci/manifest.json
//go:generate curl -L -o accredited.json https://euicc-manual.osmocom.org/docs/pki/eum/accredited.json

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
	"strings"
)

//go:embed ci.json
var ci []byte

//go:embed accredited.json
var accredited []byte

type Accredited struct {
	Version   uint8      `json:"version"`
	Suppliers []Supplier `json:"suppliers"`
}

type Supplier struct {
	Name      string            `json:"name"`
	Abbr      string            `json:"abbr,omitempty"`
	Region    string            `json:"country"`
	EUM       []string          `json:"eum,omitempty"`
	Locations map[string]string `json:"locations"`
}

type CertificateIssuer struct {
	KeyID   string `json:"key-id"`
	Country string `json:"country"`
	Name    string `json:"name"`
}

var (
	issuers []CertificateIssuer
	sites   Accredited
)

func init() {
	if err := json.Unmarshal(ci, &issuers); err != nil {
		slog.Error("failed to unmarshal certificate issuers", "error", err)
	}
	if err := json.Unmarshal(accredited, &sites); err != nil {
		slog.Error("failed to unmarshal accredited", "error", err)
	}
}

func LookupCertificateIssuer(keyID string) string {
	for _, ci := range issuers {
		if strings.HasPrefix(keyID, ci.KeyID) {
			return ci.Name
		}
	}
	return keyID
}

func LookupSASUP(eid, sasAccreditationNumber string) string {
	eum := eid[:8]
	for _, supplier := range sites.Suppliers {
		if slices.Contains(supplier.EUM, eum) {
			fallback := fmt.Sprintf("%s (%s)", supplier.Name, string(0x1F1E6+rune(supplier.Region[0])-'A')+string(0x1F1E6+rune(supplier.Region[1])-'A'))
			if len(sasAccreditationNumber) < 5 {
				return fallback
			}
			if value, ok := supplier.Locations[sasAccreditationNumber[:5]]; ok {
				return fmt.Sprintf("%s (%s)", supplier.Name, string(0x1F1E6+rune(value[0])-'A')+string(0x1F1E6+rune(value[1])-'A'))
			}
			return fallback
		}
	}
	return sasAccreditationNumber
}
