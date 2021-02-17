package forms

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type PanelAdmin struct {
	uicontrols.Panel
	client *client.Client
}

func NewPanelAdmin(parent uiinterfaces.Widget, client *client.Client) *PanelAdmin {
	var c PanelAdmin
	c.client = client
	c.InitControl(parent, &c)
	return &c
}

func (c *PanelAdmin) OnInit() {
	tvItems := NewUnitTreeControl(c, c.client)
	txtHeader := c.AddTextBlockOnGrid(0, 0, "Administrator's interface")
	txtHeader.SetFontSize(24)
	c.AddWidgetOnGrid(tvItems, 0, 1)
}
