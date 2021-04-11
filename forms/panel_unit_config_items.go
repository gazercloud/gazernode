package forms

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type PanelUnitConfigItems struct {
	uicontrols.Panel
	client     *client.Client
	configMeta []*units_common.UnitConfigItem
	config     interface{}
	OnChanged  func()

	OnConfigChanged func()
}

func NewPanelUnitConfigItems(parent uiinterfaces.Widget, configMeta []*units_common.UnitConfigItem, config interface{}, client *client.Client) *PanelUnitConfigItems {
	var c PanelUnitConfigItems
	c.InitControl(parent, &c)
	c.configMeta = configMeta
	c.config = config
	c.client = client

	hasTable := false

	labelWidth := 0

	items := make([]*PanelUnitConfigItem, 0)

	for i, configMetaItem := range c.configMeta {
		p := NewPanelUnitConfigItem(&c, configMetaItem, config, client)
		p.OnChanged = func() {
			c.Changed()
		}
		c.AddWidgetOnGrid(p, 0, i)
		if configMetaItem.Type == "table" {
			hasTable = true
		}

		if p.LabelWidth() > labelWidth {
			labelWidth = p.LabelWidth()
		}

		items = append(items, p)
	}

	for _, p := range items {
		p.SetLabelWidth(labelWidth)
	}

	if !hasTable {
		c.AddVSpacerOnGrid(0, len(c.configMeta)+1)
	}

	return &c
}

func (c *PanelUnitConfigItems) Dispose() {
	c.configMeta = nil
	c.config = nil
	c.Panel.Dispose()
}

func (c *PanelUnitConfigItems) Changed() {
	if c.OnChanged != nil {
		c.OnChanged()
	}

	if c.OnConfigChanged != nil {
		c.OnConfigChanged()
	}
}
