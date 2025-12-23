package euicc

type EuiccResponse struct {
	EID          string   `json:"eid"`
	FreeSpace    int32    `json:"freeSpace"`
	SASUP        string   `json:"sasUp"`
	Certificates []string `json:"certificates"`
}
