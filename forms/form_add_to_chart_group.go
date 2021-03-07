package forms

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type FormAddToChartGroup struct {
	uicontrols.Dialog
	client            *client.Client
	txtUnitName       *uicontrols.TextBox
	txtError          *uicontrols.TextBlock
	lvChartGroups     *uicontrols.ListView
	items             []string
	allItems          []string
	currentChannelIds []string
}

func NewFormAddToChartGroup(parent uiinterfaces.Widget, client *client.Client, items []string, allItems []string) *FormAddToChartGroup {
	var c FormAddToChartGroup
	c.client = client
	c.items = items
	c.allItems = allItems
	c.currentChannelIds = make([]string, 0)
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

	pRight.AddTextBlockOnGrid(0, 0, "Select cloud channel:")

	c.lvChartGroups = pRight.AddListViewOnGrid(0, 1)
	c.lvChartGroups.AddColumn("Name", 200)
	c.lvChartGroups.AddColumn("Channel Id", 200)
	c.lvChartGroups.OnSelectionChanged = func() {
		c.currentChannelIds = make([]string, 0)
		if len(c.lvChartGroups.SelectedItems()) > 0 {
			for _, ch := range c.lvChartGroups.SelectedItems() {
				channelId := ch.TempData
				c.currentChannelIds = append(c.currentChannelIds, channelId)
			}
		}
	}

	c.loadChartGroups()

	pButtons.AddHSpacerOnGrid(0, 0)
	btnOK := pButtons.AddButtonOnGrid(1, 0, "OK", nil)
	btnOK.SetMinWidth(70)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Cancel", nil)
	btnCancel.SetMinWidth(70)

	c.TryAccept = func() bool {
		c.currentChannelIds = make([]string, 0)
		if len(c.lvChartGroups.SelectedItems()) > 0 {
			for _, ch := range c.lvChartGroups.SelectedItems() {
				channelId := ch.TempData
				c.currentChannelIds = append(c.currentChannelIds, channelId)
			}
		}
		if len(c.currentChannelIds) > 0 {
			itemsToAdd := c.items
			c.client.CloudAddItems(c.currentChannelIds, itemsToAdd, func(err error) {
				if err == nil {
					c.TryAccept = nil
					c.Accept()
				} else {
					c.txtError.SetText(err.Error())
				}
			})

		}
		return false
	}

	c.SetAcceptButton(btnOK)
	c.SetRejectButton(btnCancel)

	return &c
}

func (c *FormAddToChartGroup) loadChartGroups() {
	c.client.ResList("chart_group", "", 0, 100000, func(info common_interfaces.ResourcesInfo, err error) {
		c.lvChartGroups.RemoveItems()
		for _, s := range info.Items {
			lvItem := c.lvChartGroups.AddItem(s.Name)
			lvItem.SetValue(1, s.Id)
			lvItem.TempData = s.Id
			lvItem.SetUserData("id", s.Name)
		}

		c.lvChartGroups.ClearSelection()
		if c.lvChartGroups.ItemsCount() > 0 {
			c.lvChartGroups.SelectItem(0)
		}
	})
}

func (c *FormAddToChartGroup) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("Add to chart group")
	c.Resize(600, 400)
}
