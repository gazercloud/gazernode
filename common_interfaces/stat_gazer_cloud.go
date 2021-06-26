package common_interfaces

type StatGazerCloud struct {
	CallsPerSecond float64 `json:"calls_per_second"`
	ReceiveSpeed   float64 `json:"receive_speed"`
	SendSpeed      float64 `json:"send_speed"`
}
