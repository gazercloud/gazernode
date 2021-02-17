package simplemap

import "github.com/gazercloud/gazerui/uiproperties"

type Map struct {
	RootItem *MapItem    `json:"root_item"`
	Layers   []*MapLayer `json:"layers"`
}

type MapItem struct {
	Props []*uiproperties.PropertyStruct `json:"props"`
}

type MapLayer struct {
	Name    string `json:"name"`
	Visible bool   `json:"visible"`
	Items   []*MapItem
}

func NewMapLayer() *MapLayer {
	var c MapLayer
	return &c
}

func NewMapItem() *MapItem {
	var c MapItem
	c.Props = make([]*uiproperties.PropertyStruct, 0)
	return &c
}

func NewMap() *Map {
	var c Map
	c.Layers = make([]*MapLayer, 0)
	return &c
}
