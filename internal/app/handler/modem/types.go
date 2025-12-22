package modem

type SIMResponse struct {
	OperatorName       string `json:"operatorName"`
	OperatorIdentifier string `json:"operatorIdentifier"`
	RegionCode         string `json:"regionCode"`
}

type RegisteredOperatorResponse struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type InsertedResponse struct {
	Manufacturer       string                     `json:"manufacturer"`
	ID                 string                     `json:"id"`
	FirmwareRevision   string                     `json:"firmwareRevision"`
	HardwareRevision   string                     `json:"hardwareRevision"`
	Name               string                     `json:"name"`
	Number             string                     `json:"number,omitempty"`
	SIM                SIMResponse                `json:"sim"`
	AccessTechnology   string                     `json:"accessTechnology"`
	RegistrationState  string                     `json:"registrationState"`
	RegisteredOperator RegisteredOperatorResponse `json:"registeredOperator"`
	SignalQuality      uint32                     `json:"signalQuality"`
	IsEsim             bool                       `json:"isEsim"`
}
