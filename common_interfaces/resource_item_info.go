package common_interfaces

type ResourcesItemInfo struct {
	Id         string          `json:"id"`
	Type       string          `json:"type"`
	Properties []*ItemProperty `json:"p"`
}
