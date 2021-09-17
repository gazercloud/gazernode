package common_interfaces

type Requester interface {
	RequestJson(function string, requestText []byte, host string, fromCloud bool) ([]byte, error)
}
