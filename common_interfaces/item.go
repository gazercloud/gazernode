package common_interfaces

type ItemValue struct {
	Value string `json:"v"`
	DT    int64  `json:"t"`
	UOM   string `json:"u"`
	//Flags string `json:"f"`
}

type Item struct {
	Id     uint64
	UnitId string
	Name   string
	Value  ItemValue
}

type ItemGetUnitItems struct {
	Item
	CloudChannels      []string
	CloudChannelsNames []string
}

func NewItem() *Item {
	var c Item
	return &c
}
