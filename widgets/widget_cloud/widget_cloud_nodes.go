package widget_cloud

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type WidgetCloudNodes struct {
	uicontrols.Panel
	client *client.Client
	timer  *uievents.FormTimer

	lvItems *uicontrols.ListView
}

func NewWidgetCloudNodes(parent uiinterfaces.Widget, client *client.Client) *WidgetCloudNodes {
	var c WidgetCloudNodes
	c.client = client
	c.InitControl(parent, &c)
	c.SetPanelPadding(0)
	return &c
}

func (c *WidgetCloudNodes) OnInit() {
	pHeader := c.AddPanelOnGrid(0, 0)
	pHeader.SetPanelPadding(0)
	txtHeader := pHeader.AddTextBlockOnGrid(0, 0, "Nodes")
	txtHeader.SetFontSize(16)

	pContent := c.AddPanelOnGrid(0, 1)
	pContent.SetPanelPadding(0)

	c.lvItems = pContent.AddListViewOnGrid(0, 0)
	c.lvItems.AddColumn("Id", 100)
	c.lvItems.AddColumn("Name", 100)

	c.UpdateStyle()
}

func (c *WidgetCloudNodes) Dispose() {
	if c.timer != nil {
		c.timer.StopTimer()
	}
	c.timer = nil
	c.client = nil

	c.Panel.Dispose()
}

func (c *WidgetCloudNodes) timerUpdate() {
	if !c.IsVisible() {
		return
	}

	c.client.CloudState(func(response nodeinterface.CloudStateResponse, err error) {
		if err != nil {
			return
		}
	})
}

func (c *WidgetCloudNodes) SetState(response nodeinterface.CloudStateResponse) {
}
