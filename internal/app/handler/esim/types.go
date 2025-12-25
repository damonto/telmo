package esim

type ProfileResponse struct {
	Name                string `json:"name"`
	ServiceProviderName string `json:"serviceProviderName"`
	ICCID               string `json:"iccid"`
	Icon                string `json:"icon"`
	ProfileState        uint8  `json:"profileState"`
	RegionCode          string `json:"regionCode,omitempty"`
}
