package common_interfaces

import (
	"strconv"
	"strings"
)

type ItemValue struct {
	Value string `json:"v"`
	DT    int64  `json:"t"`
	UOM   string `json:"u"`
}

type Item struct {
	Id               uint64    `json:"id"`
	UnitId           string    `json:"unit_id"`
	Name             string    `json:"name"`
	Value            ItemValue `json:"value"`
	Properties       map[string]*ItemProperty
	TranslateToItems map[uint64]*Item

	SourceItemId            uint64
	PostprocessingTrim      bool    `json:"postprocessing_trim"`
	PostprocessingAdjust    bool    `json:"postprocessing_adjust"`
	PostprocessingScale     float64 `json:"postprocessing_scale"`
	PostprocessingOffset    float64 `json:"postprocessing_offset"`
	PostprocessingPrecision int     `json:"postprocessing_precision"`
}

type ItemGetUnitItems struct {
	Item
	CloudChannels      []string `json:"cloud_channels"`
	CloudChannelsNames []string `json:"cloud_channels_names"`
}

type ItemStateInfo struct {
	Id          uint64 `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Value       string `json:"v"`
	DT          int64  `json:"t"`
	UOM         string `json:"u"`
}

func NewItem() *Item {
	var c Item
	c.Properties = make(map[string]*ItemProperty)
	c.TranslateToItems = make(map[uint64]*Item)
	return &c
}

func (c *Item) PostprocessingValue(value string) string {
	if c.PostprocessingTrim {
		value = strings.Trim(value, " \r\n\t")
	}

	if c.PostprocessingAdjust {
		var err error
		var valueFloat float64
		valueFloat, err = strconv.ParseFloat(value, 64)
		if err == nil {
			valueFloat = valueFloat*c.PostprocessingScale + c.PostprocessingOffset
			value = strconv.FormatFloat(valueFloat, 'f', c.PostprocessingPrecision, 64)
			if strings.Index(value, ".") >= 0 {
				value = strings.TrimRight(value, "0")
			}
		}
	}
	return value
}

func (c *Item) SetPropertyIfDoesntExist(propName string, propValue string) {
	if _, ok := c.Properties[propName]; !ok {
		c.Properties[propName] = &ItemProperty{
			Name:  propName,
			Value: propValue,
		}
	}
}

func (c *Item) SetProperty(propName string, propValue string) {
	c.Properties[propName] = &ItemProperty{
		Name:  propName,
		Value: propValue,
	}
}
