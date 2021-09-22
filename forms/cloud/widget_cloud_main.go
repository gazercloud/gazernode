package cloud

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type WidgetCloudMain struct {
	uicontrols.Panel
	client *client.Client

	//wState    *WidgetCloudState
	wNodes    *WidgetCloudNodes
	wSettings *WidgetCloudSettings

	OnNeedToConnect func(userName string, nodeId string, sessionKey string)
}

func NewWidgetCloudMain(parent uiinterfaces.Widget, client *client.Client) *WidgetCloudMain {
	var c WidgetCloudMain
	c.client = client
	c.InitControl(parent, &c)
	c.SetPanelPadding(0)
	return &c
}

func (c *WidgetCloudMain) OnInit() {
	c.wNodes = NewWidgetCloudNodes(c, c.client)
	c.wNodes.OnNeedToConnect = func(userName string, nodeId string, sessionKey string) {
		if c.OnNeedToConnect != nil {
			c.OnNeedToConnect(userName, nodeId, sessionKey)
		}
	}
	c.AddWidgetOnGrid(c.wNodes, 0, 0)
	c.wSettings = NewWidgetCloudSettings(c, c.client)
	c.AddWidgetOnGrid(c.wSettings, 1, 0)
	/*c.wState = NewWidgetCloudState(c, c.client)
	c.AddWidgetOnGrid(c.wState, 2, 0)*/
	c.UpdateStyle()
}

func (c *WidgetCloudMain) Dispose() {
	c.client = nil

	c.Panel.Dispose()
}

func (c *WidgetCloudMain) SetState(response nodeinterface.CloudStateResponse) {
	//c.wState.SetState(response)
	c.wNodes.SetState(response)
	c.wSettings.SetState(response)
}
