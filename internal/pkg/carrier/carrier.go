package carrier

//go:generate curl -L -o carrier.json https://mno-list.harded.org/unified.json

import (
	_ "embed"
	"encoding/json"
)

//go:embed carrier.json
var carrier []byte

type CarrierDataset struct {
	Brand       string              `json:"brand,omitempty"`
	Operator    string              `json:"operator,omitempty"`
	MccmncTuple map[string][]string `json:"mccmnc_tuple,omitempty"`
}

type Carrier struct {
	Name   string `json:"name"`
	Region string `json:"region"`
	Mccmnc string `json:"mccmnc"`
}

var dictionary map[string]Carrier

func init() {
	dictionary = make(map[string]Carrier)
	var c []CarrierDataset
	if err := json.Unmarshal(carrier, &c); err != nil {
		panic(err)
	}
	for _, v := range c {
		name := v.Brand
		if name == "" {
			name = v.Operator
		}
		for region, tuple := range v.MccmncTuple {
			for _, mccmnc := range tuple {
				dictionary[mccmnc] = Carrier{
					Name:   name,
					Region: region,
					Mccmnc: mccmnc,
				}
			}
		}
	}
}

func Lookup(mccmnc string) Carrier {
	if operator, ok := dictionary[mccmnc]; ok {
		return operator
	}
	return Carrier{
		Name:   "Unknown",
		Region: "Unknown",
		Mccmnc: mccmnc,
	}
}
