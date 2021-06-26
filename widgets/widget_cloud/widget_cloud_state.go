package widget_cloud

import (
	"fmt"
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"sort"
)

type WidgetCloudState struct {
	uicontrols.Panel
	client *client.Client

	lvItems *uicontrols.ListView
}

func NewWidgetCloudState(parent uiinterfaces.Widget, client *client.Client) *WidgetCloudState {
	var c WidgetCloudState
	c.client = client
	c.InitControl(parent, &c)
	c.SetPanelPadding(0)
	return &c
}

func (c *WidgetCloudState) OnInit() {
	pHeader := c.AddPanelOnGrid(0, 0)
	pHeader.SetPanelPadding(0)
	txtHeader := pHeader.AddTextBlockOnGrid(0, 0, "Frames from the cloud")
	txtHeader.SetFontSize(16)
	txtHeader.SetForeColor(c.AccentColor())
	txtHeader.SetFontSize(c.FontSize() * 1.2)

	pContent := c.AddPanelOnGrid(0, 1)
	pContent.SetPanelPadding(0)

	c.lvItems = pContent.AddListViewOnGrid(0, 0)
	c.lvItems.AddColumn("Function", 200)
	c.lvItems.AddColumn("Count", 100)

	c.UpdateStyle()
}

func (c *WidgetCloudState) Dispose() {
	c.client = nil
	c.lvItems = nil
	c.Panel.Dispose()
}

func (c *WidgetCloudState) SetState(response nodeinterface.CloudStateResponse) {
	type CallItem struct {
		Name  string
		Count int64
	}
	sort.Slice(response.Counters, func(i, j int) bool {
		return response.Counters[i].Name < response.Counters[j].Name
	})

	if c.lvItems.ItemsCount() != len(response.Counters) {
		c.lvItems.RemoveItems()
		for i := 0; i < len(response.Counters); i++ {
			c.lvItems.AddItem("")
		}
	}
	for index, item := range response.Counters {
		c.lvItems.Item(index).SetValue(0, item.Name)
		c.lvItems.Item(index).SetValue(1, fmt.Sprint(item.Value))
	}
}
