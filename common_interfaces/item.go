package common_interfaces

type ItemValue struct {
	Value string `json:"v"`
	DT    int64  `json:"t"`
	UOM   string `json:"u"`
	//Flags string `json:"f"`
}

type Item struct {
	Id     uint64    `json:"id"`
	UnitId string    `json:"unit_id"`
	Name   string    `json:"name"`
	Value  ItemValue `json:"value"`
}

type ItemGetUnitItems struct {
	Item
	CloudChannels      []string `json:"cloud_channels"`
	CloudChannelsNames []string `json:"cloud_channels_names"`
}

func NewItem() *Item {
	var c Item
	return &c
}
