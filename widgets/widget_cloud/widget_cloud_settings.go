package widget_cloud

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
)

type WidgetCloudSettings struct {
	uicontrols.Panel
	client *client.Client

	btnAllow   *uicontrols.Button
	btnDeny    *uicontrols.Button
	btnRefresh *uicontrols.Button

	lvItems *uicontrols.ListView

	btnApply *uicontrols.Button
	//wState   *WidgetCloudState
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

	pButtons := pContent.AddPanelOnGrid(0, 0)
	pButtons.SetPanelPadding(0)

	c.btnAllow = pButtons.AddButtonOnGrid(0, 0, "", func(event *uievents.Event) {
	})
	c.btnAllow.SetTooltip("Allow")

	c.btnDeny = pButtons.AddButtonOnGrid(1, 0, "", func(event *uievents.Event) {
	})
	c.btnDeny.SetTooltip("Deny")

	c.btnRefresh = pButtons.AddButtonOnGrid(2, 0, "", func(event *uievents.Event) {
		c.loadFunctions()
	})
	c.btnRefresh.SetTooltip("Refresh")

	pButtons.AddHSpacerOnGrid(10, 0)

	c.lvItems = pContent.AddListViewOnGrid(0, 1)
	c.lvItems.AddColumn("Function", 200)
	c.lvItems.AddColumn("Enabled", 50)

	pBottom := pContent.AddPanelOnGrid(0, 2)

	pBottom.AddHSpacerOnGrid(0, 0)
	c.btnApply = pBottom.AddButtonOnGrid(1, 0, "Apply", func(event *uievents.Event) {
		var req nodeinterface.CloudSetSettingsRequest
		c.client.CloudSetSettings(req, func(response nodeinterface.CloudSetSettingsResponse, err error) {
		})
	})

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
	//c.wState.SetState(response)
}

func (c *WidgetCloudSettings) UpdateStyle() {
	c.Panel.UpdateStyle()

	activeColor := c.AccentColor()
	inactiveColor := c.InactiveColor()

	c.btnAllow.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_add_materialicons_48dp_1x_baseline_add_black_48dp_png, activeColor))
	c.btnDeny.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_create_materialicons_48dp_1x_baseline_create_black_48dp_png, activeColor))
	c.btnRefresh.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_clear_materialicons_48dp_1x_baseline_clear_black_48dp_png, activeColor))

	c.btnAllow.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_add_materialicons_48dp_1x_baseline_add_black_48dp_png, inactiveColor))
	c.btnDeny.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_create_materialicons_48dp_1x_baseline_create_black_48dp_png, inactiveColor))
	c.btnRefresh.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_clear_materialicons_48dp_1x_baseline_clear_black_48dp_png, inactiveColor))
}

func (c *WidgetCloudSettings) loadFunctions() {
	c.client.CloudGetSettings(func(response nodeinterface.CloudGetSettingsResponse, err error) {
		c.lvItems.RemoveItems()
		for _, item := range response.Items {
			c.lvItems.AddItem(item.Function)
		}
	})
}
