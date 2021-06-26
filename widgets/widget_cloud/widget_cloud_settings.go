package widget_cloud

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type WidgetCloudSettings struct {
	uicontrols.Panel
	client *client.Client

	btnApply *uicontrols.Button
	wState   *WidgetCloudState
}

func NewWidgetCloudSettings(parent uiinterfaces.Widget, client *client.Client) *WidgetCloudSettings {
	var c WidgetCloudSettings
	c.client = client
	c.InitControl(parent, &c)
	c.SetPanelPadding(0)
	return &c
}

func (c *WidgetCloudSettings) OnInit() {
	pHeader := c.AddPanelOnGrid(0, 0)
	pHeader.SetPanelPadding(0)
	txtHeader := pHeader.AddTextBlockOnGrid(0, 0, "Settings")
	txtHeader.SetFontSize(16)
	txtHeader.SetForeColor(c.AccentColor())
	txtHeader.SetFontSize(c.FontSize() * 1.2)
	pContent := c.AddPanelOnGrid(0, 1)
	pContent.SetPanelPadding(0)

	pContent.AddCheckBoxOnGrid(0, 0, "Enable")
	pContent.AddCheckBoxOnGrid(0, 1, "Allow Write Item")

	c.btnApply = pContent.AddButtonOnGrid(0, 3, "Apply", func(event *uievents.Event) {
		var req nodeinterface.CloudSetSettingsRequest
		c.client.CloudSetSettings(req, func(response nodeinterface.CloudSetSettingsResponse, err error) {
		})
	})

	c.wState = NewWidgetCloudState(c, c.client)
	c.AddWidgetOnGrid(c.wState, 0, 4)

	//pContent.AddVSpacerOnGrid(0, 10)
	c.UpdateStyle()
}

func (c *WidgetCloudSettings) Dispose() {
	c.client = nil

	c.Panel.Dispose()
}

func (c *WidgetCloudSettings) timerUpdate() {
	if !c.IsVisible() {
		return
	}

	c.client.CloudState(func(response nodeinterface.CloudStateResponse, err error) {
		if err != nil {
			return
		}
	})
}

func (c *WidgetCloudSettings) SetState(response nodeinterface.CloudStateResponse) {
	c.wState.SetState(response)
}
