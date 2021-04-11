package widget_item_history

import (
	"fmt"
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/history"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/widgets/widget_time_filter"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"sort"
	"strings"
	"time"
)

type WidgetItemHistory struct {
	uicontrols.Panel
	client         *client.Client
	itemName       string
	wideValue      bool
	timer          *uievents.FormTimer
	lastLoadedDT   int64
	loadedItems    []*common_interfaces.ItemValue
	loadedItemsMap map[int64]*common_interfaces.ItemValue

	timeFilter    *widget_time_filter.TimeFilterWidget
	lvItems       *uicontrols.ListView
	chkAutoscroll *uicontrols.CheckBox
	lblStatistics *uicontrols.TextBlock
	loading       bool
}

func NewWidgetItemHistory(parent uiinterfaces.Widget, client *client.Client) *WidgetItemHistory {
	var c WidgetItemHistory
	c.client = client
	c.InitControl(parent, &c)

	c.SetPanelPadding(0)

	c.loadedItems = make([]*common_interfaces.ItemValue, 0)
	c.loadedItemsMap = make(map[int64]*common_interfaces.ItemValue)

	pContent := c.AddPanelOnGrid(0, 0)
	pContent.SetPanelPadding(0)
	pRight := pContent.AddPanelOnGrid(1, 0)
	pRight.SetPanelPadding(0)
	pButtons := c.AddPanelOnGrid(0, 1)

	c.timeFilter = widget_time_filter.NewTimeFilterWidget(pRight)
	c.timeFilter.OnEdited = c.timeFilterChanged
	c.timeFilter.SetGridX(0)
	c.timeFilter.SetGridY(0)
	pRight.AddWidgetOnGrid(c.timeFilter, 0, 0)

	c.lvItems = pRight.AddListViewOnGrid(0, 1)
	c.lvItems.AddColumn("Date/Time", 200)
	c.lvItems.AddColumn("Value", 150)
	c.lvItems.AddColumn("UOM", 50)
	c.lvItems.OnSelectionChanged = func() {
		selectedItem := c.lvItems.SelectedItem()
		if selectedItem == c.lvItems.Item(c.lvItems.ItemsCount()-1) {
			c.chkAutoscroll.SetChecked(true)
		} else {
			c.chkAutoscroll.SetChecked(false)
		}
	}

	c.lblStatistics = pRight.AddTextBlockOnGrid(0, 2, "")

	c.timeFilterChanged()

	c.chkAutoscroll = pButtons.AddCheckBoxOnGrid(0, 0, "Autoscroll")
	c.chkAutoscroll.SetChecked(true)
	pButtons.AddHSpacerOnGrid(1, 0)

	c.timer = c.Window().NewTimer(1000, func() {
		c.loadHistory()
	})
	c.timer.StartTimer()

	c.SetWideValue(c.wideValue)
	c.SetMinHeight(300)

	return &c
}

func (c *WidgetItemHistory) Dispose() {
	c.client = nil
	if c.timer != nil {
		c.timer.StopTimer()
		c.Window().RemoveTimer(c.timer)
		c.timer = nil
	}

	c.loadedItems = nil
	c.loadedItemsMap = nil
	c.timeFilter = nil
	c.lvItems = nil
	c.chkAutoscroll = nil
	c.lblStatistics = nil

	c.Panel.Dispose()
}

func (c *WidgetItemHistory) SetItem(itemName string) {
	c.itemName = itemName
	c.timeFilterChanged()
}

func (c *WidgetItemHistory) loadHistory() {
	if c.itemName == "" {
		c.lvItems.RemoveItems()
		return
	}

	if c.loading {
		return
	}
	c.loading = true
	c.client.ReadHistory(c.itemName, c.lastLoadedDT+1, c.timeFilter.TimeTo(), func(result *history.ReadResult, err error) {
		c.loading = false

		if err != nil {
			return
		}

		if result == nil {
			return
		}

		if c.lvItems == nil {
			return
		}

		if len(c.loadedItems) == 0 && c.lvItems.ItemsCount() > 0 {
			c.lvItems.RemoveItems()
		}
		lastItemsCount := c.lvItems.ItemsCount()

		lastLoadedItemsDT := int64(0)
		if len(c.loadedItems) > 0 {
			lastLoadedItemsDT = c.loadedItems[len(c.loadedItems)-1].DT
		}
		needToRestructure := false

		for _, item := range result.Items {
			if item.DT >= c.timeFilter.TimeFrom() && item.DT < c.timeFilter.TimeTo() {
				//if _, ok := c.loadedItemsMap[item.DT]; !ok {
				{
					c.loadedItemsMap[item.DT] = item
					value := item.Value
					value = strings.ReplaceAll(value, "\r", " ")
					value = strings.ReplaceAll(value, "\n", " ")
					c.lvItems.AddItem3(time.Unix(0, item.DT*1000).Local().Format("2006-01-02 15:04:05.000"), value, item.UOM)
					c.loadedItems = append(c.loadedItems, item)
					if item.DT < lastLoadedItemsDT {
						needToRestructure = true
					}
					lastLoadedItemsDT = item.DT
					c.lastLoadedDT = item.DT
					if item.DT > c.lastLoadedDT {
						c.lastLoadedDT = item.DT
					}
				}
			}
		}

		if needToRestructure {
			sort.Slice(c.loadedItems, func(i, j int) bool {
				return c.loadedItems[i].DT < c.loadedItems[j].DT
			})

			logger.Println("Restructure!")
			for index, loadedItem := range c.loadedItems {
				c.lvItems.SetItemValue(index, 0, time.Unix(0, loadedItem.DT*1000).Local().Format("2006-01-02 15:04:05.000"))
				c.lvItems.SetItemValue(index, 1, loadedItem.Value)
				c.lvItems.SetItemValue(index, 2, loadedItem.UOM)
			}
		}

		if lastItemsCount == 0 {
			c.chkAutoscroll.SetChecked(true)
		}
		if c.chkAutoscroll.IsChecked() {
			c.lvItems.EnsureVisibleItem(c.lvItems.ItemsCount() - 1)
			c.lvItems.SelectItem(c.lvItems.ItemsCount() - 1)
		}

		c.lblStatistics.SetText("Items count: " + fmt.Sprint(len(c.loadedItems)))
	})
}

func (c *WidgetItemHistory) timeFilterChanged() {
	c.lastLoadedDT = c.timeFilter.TimeFrom() - 1
	c.loadedItems = make([]*common_interfaces.ItemValue, 0)
	c.loadedItemsMap = make(map[int64]*common_interfaces.ItemValue)
	c.lvItems.RemoveItems()
	c.lblStatistics.SetText("loading ...")
	//c.lvItems.AddItem("loading ...")
	c.loadHistory()
	//logger.Println("Filter: ", time.Unix(0, c.lastLoadedDT * 1000).Local().Format("2006-01-02 15-04-05.000"), c.lastLoadedDT)
	//logger.Println("Filter: ", c.lastLoadedDT)
}

func (c *WidgetItemHistory) SetWideValue(wideValue bool) {
	c.wideValue = wideValue
	if c.lvItems != nil {
		if c.wideValue {
			c.lvItems.SetColumnWidth(1, 300)
		} else {
			c.lvItems.SetColumnWidth(1, 100)
		}
	}
}
