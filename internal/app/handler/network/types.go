package network

type NetworkResponse struct {
	Status             string   `json:"status"`
	OperatorName       string   `json:"operatorName"`
	OperatorShortName  string   `json:"operatorShortName"`
	OperatorCode       string   `json:"operatorCode"`
	AccessTechnologies []string `json:"accessTechnologies"`
}
