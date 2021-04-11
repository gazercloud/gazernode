package units_common

import "encoding/json"

type UnitConfigItem struct {
	Name              string            `json:"name"`
	DisplayName       string            `json:"display_name"`
	Type              string            `json:"type"`
	MinValue          string            `json:"min_value"`
	MaxValue          string            `json:"max_value"`
	Format            string            `json:"format"`
	Children          []*UnitConfigItem `json:"children"`
	DefaultValue      string            `json:"default_value"`
	ItemIsDisplayName bool              `json:"item_is_display_name"`
}

func NewUnitConfigItem(name string, displayName string, defaultValue string, tp string, minValue string, maxValue string, format string) *UnitConfigItem {
	var c UnitConfigItem
	c.Name = name
	c.DefaultValue = defaultValue
	c.DisplayName = displayName
	c.Type = tp
	c.MinValue = minValue
	c.MaxValue = maxValue
	c.Format = format
	c.Children = make([]*UnitConfigItem, 0)
	return &c
}

func LoadUnitConfigItems(data string) []*UnitConfigItem {
	res := make([]*UnitConfigItem, 0)

	json.Unmarshal([]byte(data), &res)

	return res
}

func (c *UnitConfigItem) Add(name string, displayName string, defaultValue string, tp string, minValue string, maxValue string, format string) *UnitConfigItem {
	newItem := NewUnitConfigItem(name, displayName, defaultValue, tp, minValue, maxValue, format)
	c.Children = append(c.Children, newItem)
	return newItem
}

func (c *UnitConfigItem) Marshal() string {
	bs, _ := json.MarshalIndent(c.Children, "", " ")
	return string(bs)
}
