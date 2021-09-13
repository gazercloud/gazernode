package cloud

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

	OnNeedToConnect func(nodeId string, sessionKey string)
}

func NewPanelCloud(parent uiinterfaces.Widget, client *client.Client) *PanelCloud {
	var c PanelCloud
	c.client = client
	c.InitControl(parent, &c)
	c.wCloud = widget_cloud.NewWidgetCloud(&c, client)
	c.AddWidgetOnGrid(c.wCloud, 0, 0)

	c.wCloud.OnNeedToConnect = func(nodeId string, sessionKey string) {
		if c.OnNeedToConnect != nil {
			c.OnNeedToConnect(nodeId, sessionKey)
		}
	}

	return &c
}

func (c *PanelCloud) OnInit() {
}

func (c *PanelCloud) FullRefresh() {
}

func (c *PanelCloud) UpdateStyle() {
	c.Panel.UpdateStyle()
}
