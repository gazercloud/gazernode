package cloud

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type PanelCloud struct {
	uicontrols.Panel
	client *client.Client
	wCloud *WidgetCloud

	OnNeedToConnect func(userName string, nodeId string, sessionKey string)
}

func NewPanelCloud(parent uiinterfaces.Widget, client *client.Client) *PanelCloud {
	var c PanelCloud
	c.client = client
	c.InitControl(parent, &c)
	c.wCloud = NewWidgetCloud(&c, client)
	c.AddWidgetOnGrid(c.wCloud, 0, 0)

	c.wCloud.OnNeedToConnect = func(userName string, nodeId string, sessionKey string) {
		if c.OnNeedToConnect != nil {
			c.OnNeedToConnect(userName, nodeId, sessionKey)
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

func (c *PanelCloud) IsSomethingWrong() bool {
	return c.wCloud.IsSomethingWrong()
}
