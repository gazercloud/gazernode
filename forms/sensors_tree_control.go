package forms

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"time"
)

type UnitsTreeControl struct {
	uicontrols.Panel
	lvItems   *uicontrols.ListView
	client    *client.Client
	timer     *uievents.FormTimer
	unitId    string
	channelId string
}

func NewUnitTreeControl(parent uiinterfaces.Widget, client *client.Client) *UnitsTreeControl {
	var c UnitsTreeControl
	c.client = client
	c.InitControl(parent, &c)
	return &c
}

func (c *UnitsTreeControl) OnInit() {
	c.SetPanelPadding(0)

	pContent := c.AddPanelOnGrid(0, 1)
	pContent.SetPanelPadding(0)

	c.lvItems = pContent.AddListViewOnGrid(0, 1)
	c.lvItems.AddColumn("Name", 250)
	c.lvItems.AddColumn("Value", 120)
	c.lvItems.AddColumn("UOM", 50)
	c.lvItems.AddColumn("Date/Time", 150)
	c.timer = c.Window().NewTimer(500, c.timerUpdate)
	c.timer.StartTimer()
}

func (c *UnitsTreeControl) SelectedItems() []string {
	items := make([]string, 0)
	for _, item := range c.lvItems.SelectedItems() {
		name := item.TempData
		items = append(items, name)
	}
	return items
}

func (c *UnitsTreeControl) SetUnitId(unitId string) {
	c.unitId = unitId
}

func (c *UnitsTreeControl) SetChannelId(channelId string) {
	c.channelId = channelId
}

func (c *UnitsTreeControl) timerUpdate() {
	if len(c.unitId) > 0 {
		c.client.GetUnitValues(c.unitId, func(items []common_interfaces.ItemGetUnitItems, err error) {
			if len(items) != c.lvItems.ItemsCount() {
				c.lvItems.RemoveItems()
				for i := 0; i < len(items); i++ {
					c.lvItems.AddItem("---")
				}
			}
			for index, di := range items {
				c.lvItems.Item(index).TempData = di.Name
				c.lvItems.SetItemValue(index, 0, di.Name)
				c.lvItems.SetItemValue(index, 1, di.Value.Value)
				c.lvItems.SetItemValue(index, 2, di.Value.UOM)
				c.lvItems.SetItemValue(index, 3, time.Unix(0, di.Value.DT*1000).Format("2006-01-02 15-04-05"))
			}
		})
	}
	if len(c.channelId) > 0 {
		c.client.GetCloudChannelValues(c.channelId, func(items []common_interfaces.Item, err error) {
			if len(items) != c.lvItems.ItemsCount() {
				c.lvItems.RemoveItems()
				for i := 0; i < len(items); i++ {
					c.lvItems.AddItem("---")
				}
			}
			for index, di := range items {
				c.lvItems.Item(index).TempData = di.Name
				c.lvItems.SetItemValue(index, 0, di.Name)
				c.lvItems.SetItemValue(index, 1, di.Value.Value)
				c.lvItems.SetItemValue(index, 2, di.Value.UOM)
				c.lvItems.SetItemValue(index, 3, time.Unix(0, di.Value.DT*1000).Format("2006-01-02 15-04-05"))
			}
		})
	}
}
