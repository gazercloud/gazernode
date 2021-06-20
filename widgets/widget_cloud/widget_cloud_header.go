package widget_cloud

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazernode/settings"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type WidgetCloudHeader struct {
	uicontrols.Panel
	client    *client.Client
	btnLogout *uicontrols.Button

	lblState1 *uicontrols.TextBlock
	lblState2 *uicontrols.TextBlock
	lblState3 *uicontrols.TextBlock

	OnNeedToLoadState func()
}

func NewWidgetCloudHeader(parent uiinterfaces.Widget, client *client.Client) *WidgetCloudHeader {
	var c WidgetCloudHeader
	c.client = client
	c.InitControl(parent, &c)
	c.SetPanelPadding(0)
	return &c
}

func (c *WidgetCloudHeader) OnInit() {
	pHeader := c.AddPanelOnGrid(0, 0)
	pHeader.SetPanelPadding(0)
	txtHeader := pHeader.AddTextBlockOnGrid(0, 0, "Cloud")
	txtHeader.SetFontSize(24)
	pHeader.AddHSpacerOnGrid(1, 0)
	c.btnLogout = pHeader.AddButtonOnGrid(2, 0, "Logout", func(event *uievents.Event) {
		c.client.CloudLogout(func(err error) {
			if c.OnNeedToLoadState != nil {
				c.OnNeedToLoadState()
			}
		})
	})

	pContent := c.AddPanelOnGrid(0, 1)
	pContent.SetPanelPadding(0)
	pContent.AddTextBlockOnGrid(0, 0, "NodeId:")
	pContent.AddTextBlockOnGrid(0, 1, "UserName:")
	pContent.AddTextBlockOnGrid(0, 2, "Connection:")
	c.lblState1 = pContent.AddTextBlockOnGrid(1, 0, "")
	c.lblState2 = pContent.AddTextBlockOnGrid(1, 1, "")
	c.lblState3 = pContent.AddTextBlockOnGrid(1, 2, "")
	pContent.AddHSpacerOnGrid(2, 0)

	c.UpdateStyle()
}

func (c *WidgetCloudHeader) Dispose() {
	c.client = nil
	c.lblState1 = nil
	c.lblState2 = nil
	c.lblState3 = nil
	c.btnLogout = nil

	c.Panel.Dispose()
}

func (c *WidgetCloudHeader) SetState(response nodeinterface.CloudStateResponse) {
	if response.LoggedIn {
		c.btnLogout.SetEnabled(true)
	} else {
		c.btnLogout.SetEnabled(false)
	}

	c.lblState1.SetText(response.NodeId + " / " + response.IAmStatus)
	if response.IAmStatus == "ok" {
		c.lblState1.SetForeColor(settings.GoodColor)
	} else {
		c.lblState1.SetForeColor(settings.BadColor)
	}

	c.lblState2.SetText(response.UserName + " / " + response.LoginStatus)
	if response.LoginStatus == "ok" {
		c.lblState2.SetForeColor(settings.GoodColor)
	} else {
		c.lblState2.SetForeColor(settings.BadColor)
	}

	if response.Connected {
		c.lblState3.SetForeColor(settings.GoodColor)
		c.lblState3.SetText(response.CurrentRepeater + " / ok")
	} else {
		c.lblState3.SetForeColor(settings.BadColor)
		c.lblState3.SetText(response.CurrentRepeater + " / " + response.ConnectionStatus)
	}
}
