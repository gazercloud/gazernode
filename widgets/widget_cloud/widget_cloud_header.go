package widget_cloud

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazernode/settings"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
)

type WidgetCloudHeader struct {
	uicontrols.Panel
	client        *client.Client
	btnLogout     *uicontrols.Button
	currentNodeId string

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
	txtHeader := pHeader.AddTextBlockOnGrid(0, 0, "Remote Access")
	txtHeader.SetFontSize(24)
	txtHeader.SetForeColor(c.AccentColor())
	pHeader.AddHSpacerOnGrid(1, 0)

	pContent := c.AddPanelOnGrid(0, 1)
	pContent.SetBorderBottom(1, c.InactiveColor())

	pState := pContent.AddPanelOnGrid(0, 0)
	pState.SetPanelPadding(0)
	pState.AddTextBlockOnGrid(0, 0, "NodeId:")
	pState.AddTextBlockOnGrid(0, 1, "UserName:")
	pState.AddTextBlockOnGrid(0, 2, "Connection:")
	c.lblState1 = pState.AddTextBlockOnGrid(1, 0, "")
	c.lblState2 = pState.AddTextBlockOnGrid(1, 1, "")
	c.lblState3 = pState.AddTextBlockOnGrid(1, 2, "")

	pContent.AddHSpacerOnGrid(1, 0)

	btnOpenInBrowser := pContent.AddButtonOnGrid(2, 0, "Open\r\nin browser", func(event *uievents.Event) {
		if len(c.currentNodeId) > 0 {
			client.OpenBrowser("https://" + c.currentNodeId + "-n.gazer.cloud/")
		} else {
			uicontrols.ShowErrorMessage(c, "Please set NodeId for the current instance.", "Error")
		}
	})
	btnOpenInBrowser.SetMinWidth(120)
	btnOpenInBrowser.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_action_open_in_browser_materialicons_48dp_1x_baseline_open_in_browser_black_48dp_png, c.AccentColor()))

	c.btnLogout = pContent.AddButtonOnGrid(3, 0, "Logout", func(event *uievents.Event) {
		c.client.CloudLogout(func(err error) {
			if c.OnNeedToLoadState != nil {
				c.OnNeedToLoadState()
			}
		})
	})
	c.btnLogout.SetMinWidth(120)
	c.btnLogout.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_action_open_in_browser_materialicons_48dp_1x_baseline_open_in_browser_black_48dp_png, c.AccentColor()))

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

	c.currentNodeId = response.NodeId

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
