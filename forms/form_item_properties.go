package forms

import (
	"fmt"
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"time"
)

type FormItemProperties struct {
	uicontrols.Dialog
	client       *client.Client
	lvProperties *uicontrols.ListView
	itemName     string
	timer        *uievents.FormTimer
}

func NewFormItemProperties(parent uiinterfaces.Widget, client *client.Client, itemName string) *FormItemProperties {
	var c FormItemProperties
	c.client = client
	c.itemName = itemName
	c.InitControl(parent, &c)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	pLeft := pContent.AddPanelOnGrid(0, 0)
	pLeft.SetPanelPadding(0)
	//pLeft.SetBorderRight(1, c.ForeColor())
	pLeft.SetMinWidth(100)
	pRight := pContent.AddPanelOnGrid(1, 0)
	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)

	/*img := pLeft.AddImageBoxOnGrid(0, 0, uiresources.ResImageAdjusted("icons/material/image/drawable-hdpi/ic_blur_on_black_48dp.png", c.ForeColor()))
	img.SetScaling(uicontrols.ImageBoxScaleAdjustImageKeepAspectRatio)
	img.SetMinHeight(64)
	img.SetMinWidth(64)*/
	pLeft.AddVSpacerOnGrid(0, 1)

	c.lvProperties = pRight.AddListViewOnGrid(0, 1)
	c.lvProperties.AddColumn("Name", 150)
	c.lvProperties.AddColumn("Value", 400)

	c.loadProperties()

	pButtons.AddHSpacerOnGrid(0, 0)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Close", func(event *uievents.Event) {
		c.Reject()
	})
	btnCancel.SetMinWidth(70)

	c.SetRejectButton(btnCancel)

	c.timer = c.Window().NewTimer(1000, func() {
		c.loadProperties()
	})
	c.timer.StartTimer()

	c.CloseEvent = func() {
		if c.timer != nil {
			c.timer.StopTimer()
			c.timer = nil
		}
	}

	return &c
}

func (c *FormItemProperties) loadProperties() {
	items := []string{c.itemName}
	c.client.GetItemsValues(items, func(items []common_interfaces.ItemGetUnitItems, err error) {
		for _, item := range items {
			if item.Name == c.itemName {
				if c.lvProperties.ItemsCount() == 0 {
					c.lvProperties.AddItem2("Name", item.Name)
					c.lvProperties.AddItem2("Id (dec)", fmt.Sprintf("%d", item.Id))
					c.lvProperties.AddItem2("Id (hex)", fmt.Sprintf("%016X", item.Id))
					c.lvProperties.AddItem2("Unit Id", item.UnitId)
					c.lvProperties.AddItem2("Value", item.Value.Value)
					c.lvProperties.AddItem2("UOM", item.Value.UOM)
					c.lvProperties.AddItem2("Date/Time", time.Unix(0, item.Value.DT*1000).Format("2006-01-02 15-03-04.000"))
					c.lvProperties.AddItem2("Flags", item.Value.Flags)
				} else {
					c.lvProperties.SetItemValue(0, 1, item.Name)
					c.lvProperties.SetItemValue(1, 1, fmt.Sprintf("%d", item.Id))
					c.lvProperties.SetItemValue(2, 1, fmt.Sprintf("%016X", item.Id))
					c.lvProperties.SetItemValue(3, 1, item.UnitId)
					c.lvProperties.SetItemValue(4, 1, item.Value.Value)
					c.lvProperties.SetItemValue(5, 1, item.Value.UOM)
					c.lvProperties.SetItemValue(6, 1, time.Unix(0, item.Value.DT*1000).Format("2006-01-02 15-03-04.000"))
					c.lvProperties.SetItemValue(7, 1, item.Value.Flags)
				}
			}
		}
	})
}

func (c *FormItemProperties) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("Item properties")
	c.Resize(800, 400)
}
