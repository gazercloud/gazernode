package common_interfaces

type Statistics struct {
	CloudReceivedBytes int `json:"cloud_received_bytes"`
	CloudSentBytes     int `json:"cloud_sent_bytes"`
	ApiCalls           int `json:"api_calls"`
}
