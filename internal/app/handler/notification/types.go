package notification

type NotificationResponse struct {
	SequenceNumber string `json:"sequenceNumber"`
	ICCID          string `json:"iccid"`
	SMDP           string `json:"smdp"`
	Operation      string `json:"operation"`
}
