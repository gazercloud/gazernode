package common_interfaces

type ItemProperty struct {
	Name  string `json:"n"`
	Value string `json:"v"`
}

type ItemConfiguration struct {
	Id         uint64          `json:"id"`
	Name       string          `json:"name"`
	Properties []*ItemProperty `json:"p"`
}
