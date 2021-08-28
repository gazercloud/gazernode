package widget_cloud

import (
	"fmt"
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
)

type WidgetCloudSettings struct {
	uicontrols.Panel
	client          *client.Client
	btnEditSettings *uicontrols.Button
	lvItems         *uicontrols.ListView
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

	c.btnEditSettings = pButtons.AddButtonOnGrid(0, 0, "", func(event *uievents.Event) {
		dialog := NewDialogRemoteAccessSettings(c, c.client)
		dialog.ShowDialog()
	})
	c.btnEditSettings.SetTooltip("Edit settings")

	pButtons.AddHSpacerOnGrid(10, 0)

	c.lvItems = pContent.AddListViewOnGrid(0, 1)
	c.lvItems.AddColumn("Function", 200)
	c.lvItems.AddColumn("Enabled", 100)
	c.lvItems.AddColumn("Counter", 100)

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

	for _, item := range response.Counters {
		found := false

		for i := 0; i < c.lvItems.ItemsCount(); i++ {
			if c.lvItems.Item(i).Value(0) == item.Name {
				c.lvItems.Item(i).SetValue(1, fmt.Sprint(item.Allow))
				c.lvItems.Item(i).SetValue(2, fmt.Sprint(item.Value))
				found = true
				break
			}
		}

		if !found {
			lvItem := c.lvItems.AddItem(item.Name)
			lvItem.SetValue(1, fmt.Sprint(item.Allow))
			lvItem.SetValue(2, fmt.Sprint(item.Value))
		}

	}

}

func (c *WidgetCloudSettings) UpdateStyle() {
	c.Panel.UpdateStyle()

	activeColor := c.AccentColor()
	inactiveColor := c.InactiveColor()

	c.btnEditSettings.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_create_materialiconsoutlined_48dp_1x_outline_create_black_48dp_png, activeColor))

	c.btnEditSettings.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_create_materialiconsoutlined_48dp_1x_outline_create_black_48dp_png, inactiveColor))
}

func (c *WidgetCloudSettings) loadFunctions() {
	c.client.CloudGetSettings(func(response nodeinterface.CloudGetSettingsResponse, err error) {
		if err == nil {
			c.lvItems.RemoveItems()
			for _, item := range response.Items {
				c.lvItems.AddItem(item.Function)
			}
		}
	})
}
