package forms

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/widgets/widget_chart"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type FormAddToChartGroup struct {
	uicontrols.Dialog
	client             *client.Client
	txtUnitName        *uicontrols.TextBox
	txtError           *uicontrols.TextBlock
	lvChartGroups      *uicontrols.ListView
	items              []string
	allItems           []string
	currentChartGroups []string

	timer *uievents.FormTimer

	addedToChartGroup map[string]error
}

func NewFormAddToChartGroup(parent uiinterfaces.Widget, client *client.Client, items []string, allItems []string) *FormAddToChartGroup {
	var c FormAddToChartGroup
	c.client = client
	c.items = items
	c.allItems = allItems
	c.addedToChartGroup = make(map[string]error)
	c.currentChartGroups = make([]string, 0)
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
	c.lvChartGroups.AddColumn("Result", 100)
	c.lvChartGroups.OnSelectionChanged = func() {
		c.currentChartGroups = make([]string, 0)
		if len(c.lvChartGroups.SelectedItems()) > 0 {
			for _, ch := range c.lvChartGroups.SelectedItems() {
				channelId := ch.TempData
				c.currentChartGroups = append(c.currentChartGroups, channelId)
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
		c.TryAccept = nil
		btnOK.SetEnabled(false)
		c.currentChartGroups = make([]string, 0)
		if len(c.lvChartGroups.SelectedItems()) > 0 {
			for _, ch := range c.lvChartGroups.SelectedItems() {
				chartGroupId := ch.TempData
				c.currentChartGroups = append(c.currentChartGroups, chartGroupId)
			}
		}
		if len(c.currentChartGroups) > 0 {
			c.addedToChartGroup = make(map[string]error)
			for _, chartGroup := range c.currentChartGroups {
				c.addItemToChartGroup(chartGroup, items)
			}
		}
		return false
	}

	c.timer = c.Window().NewTimer(500, c.timerUpdate)
	c.timer.StartTimer()

	c.SetAcceptButton(btnOK)
	c.SetRejectButton(btnCancel)

	return &c
}

func (c *FormAddToChartGroup) Dispose() {
	c.lvChartGroups = nil
	if c.timer != nil {
		c.timer.StopTimer()
		c.Window().RemoveTimer(c.timer)
	}
	c.timer = nil
	c.Dialog.Dispose()
}

func (c *FormAddToChartGroup) timerUpdate() {
	count := 0
	for i := 0; i < c.lvChartGroups.ItemsCount(); i++ {
		if res, ok := c.addedToChartGroup[c.lvChartGroups.Item(i).TempData]; ok {
			if res == nil {
				c.lvChartGroups.SetItemValue(i, 2, "ok")
				count++
			} else {
				c.lvChartGroups.SetItemValue(i, 2, res.Error())
			}
		} else {
			c.lvChartGroups.SetItemValue(i, 2, "-")
		}
	}

	if count > 0 && count == len(c.currentChartGroups) {
		c.Accept()
	}
}

func (c *FormAddToChartGroup) addItemToChartGroup(chartGroup string, items []string) {
	c.client.ResGet(chartGroup, func(item *common_interfaces.ResourcesItem, err error) {
		if err != nil {
			c.addedToChartGroup[chartGroup] = err
			return
		}

		var res widget_chart.ChartSettings
		err = json.Unmarshal(item.Content, &res)
		if err != nil {
			c.addedToChartGroup[chartGroup] = err
			return
		}

		for _, it := range items {
			var area widget_chart.ChartSettingsArea
			var ser widget_chart.ChartSettingsSeries
			ser.Item = it
			ser.Color = "#00FF00"
			area.Series = make([]*widget_chart.ChartSettingsSeries, 0)
			area.Series = append(area.Series, &ser)
			if res.Areas == nil {
				res.Areas = make([]*widget_chart.ChartSettingsArea, 0)
			}
			res.Areas = append(res.Areas, &area)
		}

		var bs []byte
		bs, err = json.MarshalIndent(res, "", " ")

		c.client.ResSet(chartGroup, []byte{}, bs, func(err error) {
			c.addedToChartGroup[chartGroup] = err
		})
	})
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
