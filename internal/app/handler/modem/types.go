package modem

type SlotResponse struct {
	Active             bool   `json:"active"`
	OperatorName       string `json:"operatorName"`
	OperatorIdentifier string `json:"operatorIdentifier"`
	RegionCode         string `json:"regionCode"`
	Identifier         string `json:"identifier"`
}

type RegisteredOperatorResponse struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type ModemResponse struct {
	Manufacturer       string                     `json:"manufacturer"`
	ID                 string                     `json:"id"`
	FirmwareRevision   string                     `json:"firmwareRevision"`
	HardwareRevision   string                     `json:"hardwareRevision"`
	Name               string                     `json:"name"`
	Number             string                     `json:"number,omitempty"`
	SIM                SlotResponse               `json:"sim"`
	Slots              []SlotResponse             `json:"slots"`
	AccessTechnology   string                     `json:"accessTechnology"`
	RegistrationState  string                     `json:"registrationState"`
	RegisteredOperator RegisteredOperatorResponse `json:"registeredOperator"`
	SignalQuality      uint32                     `json:"signalQuality"`
	SupportsEsim       bool                       `json:"supportsEsim"`
}
