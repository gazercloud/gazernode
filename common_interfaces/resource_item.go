package common_interfaces

type ResourcesItem struct {
	Info    ResourcesItemInfo `json:"info"`
	Content []byte            `json:"content"`
}
