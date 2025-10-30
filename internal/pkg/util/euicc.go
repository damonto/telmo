package util

//go:generate curl -o eum.json https://euicc-manual.osmocom.org/docs/pki/eum/manifest-v2.json
//go:generate curl -o ci.json https://euicc-manual.osmocom.org/docs/pki/ci/manifest.json
//go:generate curl -o accredited.json https://euicc-manual.osmocom.org/docs/pki/eum/accredited.json

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

//go:embed eum.json
var eum []byte

//go:embed ci.json
var ci []byte

//go:embed accredited.json
var accredited []byte

type EUM struct {
	EUM            string    `json:"eum"`
	Country        string    `json:"country"`
	Manufacturer   string    `json:"manufacturer"`
	Accreditations []string  `json:"accreditations"`
	Products       []Product `json:"products"`
}

type Product struct {
	Prefix  string  `json:"prefix"`
	Chip    string  `json:"chip,omitempty"`
	Name    string  `json:"name"`
	InRange [][]int `json:"in-range,omitempty"`
}

type Supplier struct {
	Name      string `json:"name"`
	Abbr      string `json:"abbr,omitempty"`
	Location  string
	Locations map[string]string `json:"locations"`
}

type Accredited struct {
	Suppliers []Supplier `json:"suppliers"`
}

type CertificateIssuer struct {
	KeyID   string `json:"key-id"`
	Country string `json:"country"`
	Name    string `json:"name"`
}

var (
	accreditedSites    map[string]Supplier
	certificateIssuers []CertificateIssuer
	EUMs               []EUM
)

func init() {
	if err := json.Unmarshal(eum, &EUMs); err != nil {
		slog.Error("failed to unmarshal EUMs", "error", err)
	}
	if err := json.Unmarshal(ci, &certificateIssuers); err != nil {
		slog.Error("failed to unmarshal certificate issuers", "error", err)
	}
	var sites Accredited
	if err := json.Unmarshal(accredited, &sites); err != nil {
		slog.Error("failed to unmarshal accredited", "error", err)
	}
	accreditedSites = make(map[string]Supplier)
	for _, site := range sites.Suppliers {
		for prefix, location := range site.Locations {
			site.Location = location
			accreditedSites[prefix] = site
		}
	}
}

func LookupCertificateIssuer(keyID string) string {
	for _, ci := range certificateIssuers {
		if strings.HasPrefix(keyID, ci.KeyID) {
			return ci.Name
		}
	}
	return keyID
}

func LookupEUM(eid string) (country string, manufacturer string, brand string) {
	for _, manifest := range EUMs {
		if strings.HasPrefix(eid, manifest.EUM) {
			country = manifest.Country
			manufacturer = manifest.Manufacturer
			for _, product := range manifest.Products {
				if strings.HasPrefix(eid, product.Prefix) {
					if product.InRange != nil {
						eidRange, _ := strconv.Atoi(eid[len(product.Prefix):30])
						for _, inRange := range product.InRange {
							if eidRange >= inRange[0] && eidRange <= inRange[1] {
								return country, manufacturer, product.Name
							}
						}
					}
					brand = product.Name
				}
			}
		}
	}
	return country, manufacturer, brand
}

func LookupAccredited(sasAccreditationNumber string) string {
	if len(sasAccreditationNumber) < 5 {
		return sasAccreditationNumber
	}
	if supplier, ok := accreditedSites[sasAccreditationNumber[:5]]; ok {
		return fmt.Sprintf(
			"%s %s (%s %s)",
			sasAccreditationNumber,
			If(supplier.Abbr != "", supplier.Abbr, supplier.Name),
			string(0x1F1E6+rune(supplier.Location[0])-'A')+string(0x1F1E6+rune(supplier.Location[1])-'A'),
			supplier.Location,
		)
	}
	return sasAccreditationNumber
}
