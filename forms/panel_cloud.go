package forms

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/widgets/widget_cloud"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type PanelCloud struct {
	uicontrols.Panel
	client *client.Client
	wCloud *widget_cloud.WidgetCloud
}

func NewPanelCloud(parent uiinterfaces.Widget, client *client.Client) *PanelCloud {
	var c PanelCloud
	c.client = client
	c.InitControl(parent, &c)
	c.wCloud = widget_cloud.NewWidgetCloud(&c, client)
	c.AddWidgetOnGrid(c.wCloud, 0, 0)
	return &c
}

func (c *PanelCloud) OnInit() {
}

func (c *PanelCloud) FullRefresh() {
}

func (c *PanelCloud) UpdateStyle() {
	c.Panel.UpdateStyle()
}
