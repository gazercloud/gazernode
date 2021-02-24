package forms

import (
	"encoding/base64"
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/gazer_dictionary"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazernode/settings"
	"github.com/gazercloud/gazernode/widgets/widget_chart"
	"github.com/gazercloud/gazerui/canvas"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
	"github.com/go-gl/glfw/v3.3/glfw"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"strconv"
	"strings"
	"time"
)

type PanelUnits struct {
	uicontrols.Panel
	client  *client.Client
	lvUnits *uicontrols.ListView

	btnAdd    *uicontrols.Button
	btnEdit   *uicontrols.Button
	btnRemove *uicontrols.Button

	btnStart *uicontrols.Button
	btnStop  *uicontrols.Button

	btnShowFullScreen  *uicontrols.Button
	btnAddToCloud      *uicontrols.Button
	btnRemoveFromCloud *uicontrols.Button
	btnOpenInBrowser   *uicontrols.Button

	currentUnitId   string
	currentUnitName string
	currentMainItem string

	lvItems      *uicontrols.ListView
	timer        *uievents.FormTimer
	wItemDetails *widget_chart.WidgetCharts

	menuUnits *uicontrols.PopupMenu
}

func NewPanelUnits(parent uiinterfaces.Widget, client *client.Client) *PanelUnits {
	var c PanelUnits
	c.client = client
	c.InitControl(parent, &c)
	return &c
}

func (c *PanelUnits) OnInit() {
	//pHeader := c.AddPanelOnGrid(0, 0)

	pContent := c.AddPanelOnGrid(0, 0)
	pContent.SetPanelPadding(0)
	splitter := pContent.AddSplitContainerOnGrid(0, 0)
	splitter.SetYExpandable(true)
	splitter.SetPosition(420)

	pUnitsList := splitter.Panel1.AddPanelOnGrid(0, 0)
	pUnitsList.SetPanelPadding(0)
	txtHeader := pUnitsList.AddTextBlockOnGrid(0, 0, "Units")
	txtHeader.SetFontSize(24)

	pButtons := pUnitsList.AddPanelOnGrid(0, 1)
	pButtons.SetPanelPadding(0)

	c.btnAdd = pButtons.AddButtonOnGrid(0, 0, "", func(event *uievents.Event) {
		c.addUnit()
	})
	c.btnAdd.SetTooltip("Add unit ...")
	c.btnAdd.SetMinWidth(60)

	c.btnEdit = pButtons.AddButtonOnGrid(2, 0, "", func(event *uievents.Event) {
		c.editUnit()
	})
	c.btnEdit.SetTooltip("Edit unit ...")

	c.btnRemove = pButtons.AddButtonOnGrid(3, 0, "", func(event *uievents.Event) {
		c.removeUnit()
	})
	c.btnRemove.SetTooltip("Remove selected units")

	pButtons.AddTextBlockOnGrid(4, 0, " | ")

	c.btnStart = pButtons.AddButtonOnGrid(5, 0, "", func(event *uievents.Event) {
		c.startUnit()
	})
	c.btnStart.SetTooltip("Start selected units")

	c.btnStop = pButtons.AddButtonOnGrid(6, 0, "", func(event *uievents.Event) {
		c.stopUnit()
	})
	c.btnStop.SetTooltip("Stop selected units")

	pButtons.AddHSpacerOnGrid(7, 0)

	c.lvUnits = pUnitsList.AddListViewOnGrid(0, 2)
	c.lvUnits.AddColumn("Name", 150)
	c.lvUnits.AddColumn("Type", 100)
	c.lvUnits.AddColumn("Value", 150)
	c.lvUnits.OnSelectionChanged = func() {
		selectedItem := c.lvUnits.SelectedItem()
		if selectedItem != nil {
			unitId := selectedItem.UserData("id").(string)
			unitName := selectedItem.UserData("name").(string)
			unitState, ok := selectedItem.UserData("state").(nodeinterface.UnitStateResponse)
			if ok {
				//if unitState != nil {
				c.currentMainItem = unitState.MainItem
				//}
			}
			c.SetCurrentUnit(unitId, unitName, c.currentMainItem)
		} else {
			c.SetCurrentUnit("", "", "")
		}
	}

	c.menuUnits = uicontrols.NewPopupMenu(c.lvUnits)
	c.menuUnits.AddItemWithUiResImage("Add unit ...", func(event *uievents.Event) {
		c.addUnit()
	}, uiresources.R_icons_material4_png_content_add_materialiconsoutlined_48dp_1x_outline_add_black_48dp_png, "")
	c.menuUnits.AddItemWithUiResImage("Edit unit ...", func(event *uievents.Event) {
		c.editUnit()
	}, uiresources.R_icons_material4_png_content_create_materialiconsoutlined_48dp_1x_outline_create_black_48dp_png, "")
	c.menuUnits.AddItemWithUiResImage("Remove unit", func(event *uievents.Event) {
		c.removeUnit()
	}, uiresources.R_icons_material4_png_content_clear_materialiconsoutlined_48dp_1x_outline_clear_black_48dp_png, "")
	c.menuUnits.AddItemWithUiResImage("Start unit", func(event *uievents.Event) {
		c.startUnit()
	}, uiresources.R_icons_material4_png_av_play_arrow_materialicons_48dp_1x_baseline_play_arrow_black_48dp_png, "")
	c.menuUnits.AddItemWithUiResImage("Stop unit", func(event *uievents.Event) {
		c.stopUnit()
	}, uiresources.R_icons_material4_png_av_pause_materialicons_48dp_1x_baseline_pause_black_48dp_png, "")
	c.menuUnits.AddItemWithUiResImage("View Log", func(event *uievents.Event) {
		c.viewLog()
	}, uiresources.R_icons_material4_png_action_view_headline_materialiconsoutlined_48dp_1x_outline_view_headline_black_48dp_png, "")
	c.lvUnits.SetContextMenu(c.menuUnits)

	pItems := splitter.Panel2.AddPanelOnGrid(1, 0)
	pItems.SetPanelPadding(0)

	txtHeaderItems := pItems.AddTextBlockOnGrid(0, 0, "Data Items")
	txtHeaderItems.SetFontSize(24)

	pItems.SetOnKeyDown(func(event *uievents.KeyDownEvent) bool {
		if event.Key == glfw.KeyEnter || event.Key == glfw.KeyKPEnter {
			items := c.SelectedItems()
			if len(items) > 0 {
				MainFormInstance.ShowFullScreenValue(true, items[0])
			}
			return true
		}
		if event.Key == glfw.KeyF4 {
			items := c.SelectedItems()
			if len(items) == 1 {
				NewFormWriteValue(c, c.client, items[0]).ShowDialog()
			}
			return true
		}
		return false
	})

	pButtonsRight := pItems.AddPanelOnGrid(0, 1)
	pButtonsRight.SetPanelPadding(0)

	pButtons.AddHSpacerOnGrid(5, 0)

	c.btnAddToCloud = pButtonsRight.AddButtonOnGrid(0, 0, "", func(event *uievents.Event) {
		items := c.SelectedItems()
		allItems := c.AllItems()
		f := NewFormAddToCloud(c, c.client, items, allItems, nil)
		f.ShowDialog()
		f.OnAccept = func() {
		}
	})
	c.btnAddToCloud.SetTooltip("Add to Gazer-Cloud")
	c.btnAddToCloud.SetMinWidth(60)

	c.btnOpenInBrowser = pButtonsRight.AddButtonOnGrid(2, 0, "", func(event *uievents.Event) {
		channels := make(map[string]string)
		for i := 0; i < c.lvItems.ItemsCount(); i++ {
			item := c.lvItems.Item(i)

			if cloudChannels, ok := item.UserData("cloud").([]string); ok {
				for _, ch := range cloudChannels {
					channels[ch] = "1"
				}
			}
		}
		for channelId, _ := range channels {
			client.OpenBrowser(gazer_dictionary.ChannelUrl(channelId))
		}
	})
	c.btnOpenInBrowser.SetTooltip("Open in browser")
	c.btnOpenInBrowser.SetMinWidth(60)

	pButtonsRight.AddHSpacerOnGrid(3, 0)

	c.btnRemoveFromCloud = pButtonsRight.AddButtonOnGrid(4, 0, "", func(event *uievents.Event) {
		f := NewFormRemoveFromCloud(c, c.client, c.SelectedItems(), c.AllItems(), nil)
		f.ShowDialog()
		f.OnAccept = func() {
		}
	})
	c.btnRemoveFromCloud.SetTooltip("Remove from Gazer-Cloud")

	c.btnShowFullScreen = pButtonsRight.AddButtonOnGrid(5, 0, "", func(event *uievents.Event) {
		items := c.SelectedItems()
		if len(items) > 0 {
			MainFormInstance.ShowFullScreenValue(true, items[0])
		}
	})
	c.btnShowFullScreen.SetTooltip("Full screen")

	c.lvItems = pItems.AddListViewOnGrid(0, 2)
	c.lvItems.AddColumn("Name", 200)
	c.lvItems.AddColumn("Value", 100)
	c.lvItems.AddColumn("UOM", 70)
	c.lvItems.AddColumn("Time", 80)
	c.lvItems.AddColumn("Sharing", 200)

	c.lvItems.SetColumnTextAlign(1, canvas.HAlignRight)

	menuItems := uicontrols.NewPopupMenu(c.lvUnits)
	menuItems.AddItem("Add to cloud ...", func(event *uievents.Event) {
		c.addSelectedItemsToCloud()
	}, uiresources.ResImgCol(uiresources.R_icons_material4_png_file_cloud_upload_materialiconsoutlined_48dp_1x_outline_cloud_upload_black_48dp_png, c.ForeColor()), "")
	menuItems.AddItem("Remove from cloud ...", func(event *uievents.Event) {
		c.removeSelectedItemsFromCloud()
	}, uiresources.ResImgCol(uiresources.R_icons_material4_png_file_cloud_off_materialiconsoutlined_48dp_1x_outline_cloud_off_black_48dp_png, c.ForeColor()), "")
	menuItems.AddItem("Open in browser", func(event *uievents.Event) {
		c.openSelectedItemInBrowser()
	}, uiresources.ResImgCol(uiresources.R_icons_material4_png_action_open_in_browser_materialiconsoutlined_48dp_1x_outline_open_in_browser_black_48dp_png, c.ForeColor()), "")
	menuItems.AddItem("Big view ...", func(event *uievents.Event) {
		items := c.SelectedItems()
		if len(items) > 0 {
			MainFormInstance.ShowFullScreenValue(true, items[0])
		}
	}, uiresources.ResImgCol(uiresources.R_icons_material4_png_navigation_fullscreen_materialiconsoutlined_48dp_1x_outline_fullscreen_black_48dp_png, c.ForeColor()), "")
	menuItems.AddItem("History ...", func(event *uievents.Event) {
		items := c.SelectedItems()
		if len(items) == 1 {
			NewFormItemHistory(c, c.client, items[0]).ShowDialog()
		}
	}, uiresources.ResImgCol(uiresources.R_icons_material4_png_action_view_headline_materialiconsoutlined_48dp_1x_outline_view_headline_black_48dp_png, c.ForeColor()), "")
	menuItems.AddItem("Properties ...", func(event *uievents.Event) {
		items := c.SelectedItems()
		if len(items) == 1 {
			NewFormItemProperties(c, c.client, items[0]).ShowDialog()
		}
	}, uiresources.ResImgCol(uiresources.R_icons_material4_png_action_info_materialiconsoutlined_48dp_1x_outline_info_black_48dp_png, c.ForeColor()), "")
	menuItems.AddItem("Write ...", func(event *uievents.Event) {
		items := c.SelectedItems()
		if len(items) == 1 {
			NewFormWriteValue(c, c.client, items[0]).ShowDialog()
		}
	}, uiresources.ResImgCol(uiresources.R_icons_material4_png_action_info_materialiconsoutlined_48dp_1x_outline_info_black_48dp_png, c.ForeColor()), "")
	menuItems.AddItem("Copy full item name", func(event *uievents.Event) {
		items := c.SelectedItems()
		if len(items) == 1 {
			glfw.SetClipboardString(items[0])
		}
	}, uiresources.ResImgCol(uiresources.R_icons_material4_png_action_info_materialiconsoutlined_48dp_1x_outline_info_black_48dp_png, c.ForeColor()), "")
	c.lvItems.SetContextMenu(menuItems)
	c.lvItems.OnSelectionChanged = func() {
		c.wItemDetails.SetDataItems(c.SelectedItems())
		c.wItemDetails.SetShowQualities(true)
	}

	c.wItemDetails = widget_chart.NewWidgetCharts(pItems, c.client)
	pItems.AddWidgetOnGrid(c.wItemDetails, 0, 3)

	c.timer = c.Window().NewTimer(500, c.timerUpdate)
	c.timer.StartTimer()

	c.loadUnits()
	c.UpdateStyle()
}

func (c *PanelUnits) Dispose() {
	if c.timer != nil {
		c.timer.StopTimer()
	}
	c.timer = nil

	c.client = nil
	c.lvUnits = nil

	c.btnAdd = nil
	c.btnEdit = nil
	c.btnRemove = nil

	c.btnStart = nil
	c.btnStop = nil

	c.btnShowFullScreen = nil
	c.btnAddToCloud = nil
	c.btnRemoveFromCloud = nil
	c.btnOpenInBrowser = nil

	c.lvItems = nil

	c.Panel.Dispose()
}

func (c *PanelUnits) FullRefresh() {
	c.loadUnits()
}

func (c *PanelUnits) Activate() {
	c.lvUnits.Focus()
}

func (c *PanelUnits) addUnit() {
	f := NewFormAddUnit(c, c.client)
	f.SetName("FormAddUnit")
	f.ShowDialog()
	f.OnAccept = func() {
		logger.Println("OnAccept NewFormAddUnit")
		c.loadUnits()
	}
}

func (c *PanelUnits) editUnit() {
	if len(c.lvUnits.SelectedItems()) == 1 {
		unitId := c.lvUnits.SelectedItem().UserData("id").(string)
		f := NewFormUnitEdit(c, c.client, unitId, "")
		f.ShowDialog()
		f.OnAccept = func() {
			c.loadUnits()
		}
	}
}

func (c *PanelUnits) removeUnit() {
	units := make([]*nodeinterface.UnitListResponseItem, 0)
	for _, selectedItem := range c.lvUnits.SelectedItems() {
		unitInfo := selectedItem.UserData("info").(*nodeinterface.UnitListResponseItem)
		units = append(units, unitInfo)
	}

	f := NewFormRemoveUnits(c, c.client, units)
	f.ShowDialog()
	f.OnAccept = func() {
		c.loadUnits()
	}
}

func (c *PanelUnits) startUnit() {
	ids := make([]string, 0)
	for _, selectedItem := range c.lvUnits.SelectedItems() {
		unitId := selectedItem.UserData("id").(string)
		ids = append(ids, unitId)
	}
	c.client.StartUnits(ids, nil)
}

func (c *PanelUnits) stopUnit() {
	ids := make([]string, 0)
	for _, selectedItem := range c.lvUnits.SelectedItems() {
		unitId := selectedItem.UserData("id").(string)
		ids = append(ids, unitId)
	}
	c.client.StopUnits(ids, nil)
}

func (c *PanelUnits) viewLog() {
	for _, selectedItem := range c.lvUnits.SelectedItems() {
		sens := selectedItem.UserData("info").(*nodeinterface.UnitListResponseItem)
		f := NewFormItemHistory(c, c.client, sens.Name+"/.service/log")
		f.SetWideValue(true)
		f.ShowDialog()
		break
	}
}

func (c *PanelUnits) UpdateStyle() {
	c.Panel.UpdateStyle()

	activeColor := c.AccentColor()
	inactiveColor := c.InactiveColor()

	c.btnAdd.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_add_materialiconsoutlined_48dp_1x_outline_add_black_48dp_png, activeColor))
	c.btnEdit.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_create_materialiconsoutlined_48dp_1x_outline_create_black_48dp_png, activeColor))
	c.btnRemove.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_clear_materialiconsoutlined_48dp_1x_outline_clear_black_48dp_png, activeColor))
	c.btnStart.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_av_play_arrow_materialicons_48dp_1x_baseline_play_arrow_black_48dp_png, activeColor))
	c.btnStop.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_av_pause_materialiconsoutlined_48dp_1x_outline_pause_black_48dp_png, activeColor))

	c.btnShowFullScreen.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_navigation_fullscreen_materialiconsoutlined_48dp_1x_outline_fullscreen_black_48dp_png, activeColor))
	c.btnAddToCloud.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_file_cloud_upload_materialicons_48dp_1x_baseline_cloud_upload_black_48dp_png, activeColor))
	c.btnRemoveFromCloud.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_file_cloud_off_materialiconsoutlined_48dp_1x_outline_cloud_off_black_48dp_png, activeColor))
	c.btnOpenInBrowser.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_action_open_in_browser_materialicons_48dp_1x_baseline_open_in_browser_black_48dp_png, activeColor))

	c.btnAdd.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_add_materialiconsoutlined_48dp_1x_outline_add_black_48dp_png, inactiveColor))
	c.btnEdit.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_create_materialiconsoutlined_48dp_1x_outline_create_black_48dp_png, inactiveColor))
	c.btnRemove.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_clear_materialiconsoutlined_48dp_1x_outline_clear_black_48dp_png, inactiveColor))
	c.btnStart.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_av_play_arrow_materialicons_48dp_1x_baseline_play_arrow_black_48dp_png, inactiveColor))
	c.btnStop.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_av_pause_materialiconsoutlined_48dp_1x_outline_pause_black_48dp_png, inactiveColor))

	c.btnShowFullScreen.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_navigation_fullscreen_materialiconsoutlined_48dp_1x_outline_fullscreen_black_48dp_png, inactiveColor))
	c.btnAddToCloud.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_file_cloud_upload_materialicons_48dp_1x_baseline_cloud_upload_black_48dp_png, inactiveColor))
	c.btnRemoveFromCloud.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_file_cloud_off_materialiconsoutlined_48dp_1x_outline_cloud_off_black_48dp_png, inactiveColor))
	c.btnOpenInBrowser.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_action_open_in_browser_materialicons_48dp_1x_baseline_open_in_browser_black_48dp_png, inactiveColor))

}

func (c *PanelUnits) SetCurrentUnit(unitId string, unitName string, mainItem string) {
	c.currentUnitId = unitId
	c.currentUnitName = unitName
	c.lvItems.RemoveItems()
}

func (c *PanelUnits) SelectedItems() []string {
	items := make([]string, 0)
	for _, item := range c.lvItems.SelectedItems() {
		name := item.TempData
		items = append(items, name)
	}
	return items
}

func (c *PanelUnits) AllItems() []string {
	items := make([]string, 0)
	for i := 0; i < c.lvItems.ItemsCount(); i++ {
		name := c.lvItems.Item(i).TempData
		items = append(items, name)
	}
	return items
}

func (c *PanelUnits) loadUnits() {
	c.client.ListOfUnits(func(infos []nodeinterface.UnitListResponseItem, err error) {
		c.lvUnits.RemoveItems()
		for _, s := range infos {
			sens := s
			lvItem := c.lvUnits.AddItem(s.Name)
			lvItem.SetValue(1, s.TypeForDisplay)
			lvItem.SetUserData("info", &sens)
			lvItem.SetUserData("id", sens.Id)
			lvItem.SetUserData("name", sens.Name)
		}
	})
}

func (c *PanelUnits) updateUnitsButtons() {
	if len(c.lvUnits.SelectedItems()) > 0 {
		if len(c.lvUnits.SelectedItems()) == 1 {
			c.btnEdit.SetEnabled(true)
		} else {
			c.btnEdit.SetEnabled(false)
		}
		c.btnRemove.SetEnabled(true)
		c.btnStart.SetEnabled(true)
		c.btnStop.SetEnabled(true)
	} else {
		c.btnEdit.SetEnabled(false)
		c.btnRemove.SetEnabled(false)
		c.btnStart.SetEnabled(false)
		c.btnStop.SetEnabled(false)
	}
}

func (c *PanelUnits) updateDataItemsButtons() {
	itemsSelected := c.lvItems.SelectedItems()
	if c.lvItems.ItemsCount() > 0 {
		c.btnAddToCloud.SetEnabled(true)
	} else {
		c.btnAddToCloud.SetEnabled(false)
	}

	if len(itemsSelected) > 0 {
		if len(itemsSelected) == 1 {
			c.btnShowFullScreen.SetEnabled(true)
		} else {
			c.btnShowFullScreen.SetEnabled(false)
		}

		itemHasCloud := false

		for _, item := range itemsSelected {
			if cloudChannels, ok := item.UserData("cloud").([]string); ok {
				if len(cloudChannels) > 0 {
					itemHasCloud = true
				}
			}
		}

		if itemHasCloud {
			c.btnRemoveFromCloud.SetEnabled(true)
		} else {
			c.btnRemoveFromCloud.SetEnabled(false)
		}
	} else {
		c.btnShowFullScreen.SetEnabled(false)
		c.btnRemoveFromCloud.SetEnabled(false)
	}
}

func (c *PanelUnits) updateUnitsState() {
	for i := 0; i < c.lvUnits.ItemsCount(); i++ {
		unitId := c.lvUnits.Item(i).UserData("id").(string)
		c.client.GetUnitState(unitId, func(state nodeinterface.UnitStateResponse, err error) {
			if c.lvUnits == nil {
				return
			}
			for i := 0; i < c.lvUnits.ItemsCount(); i++ {
				sId := c.lvUnits.Item(i).UserData("id").(string)
				if sId == state.UnitId {
					value := state.Value
					{
						if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
							p := message.NewPrinter(language.English)
							value = strings.ReplaceAll(p.Sprint(intValue), ",", " ")
						}
					}

					c.lvUnits.Item(i).SetValue(2, value+" "+state.UOM)
					c.lvUnits.Item(i).SetUserData("state", state)

					if state.UOM == "error" {
						c.lvUnits.Item(i).SetForeColorForCell(2, settings.BadColor)
					} else {
						c.lvUnits.Item(i).SetForeColorForCell(2, settings.GoodColor)
					}

					if state.Status == "stopped" {
						c.lvUnits.Item(i).SetForeColorForCell(0, c.InactiveColor())
						c.lvUnits.Item(i).SetForeColorForCell(1, c.InactiveColor())
					} else {
						c.lvUnits.Item(i).SetForeColorForCell(0, nil)
						c.lvUnits.Item(i).SetForeColorForCell(1, nil)
					}
				}
			}
		})
	}

}

func (c *PanelUnits) updateDataItemsState() {
	if len(c.currentUnitName) > 0 {
		c.client.GetUnitValues(c.currentUnitName, func(items []common_interfaces.ItemGetUnitItems, err error) {
			if err != nil {
				return
			}

			itemsToShow := make([]common_interfaces.ItemGetUnitItems, 0)
			for _, di := range items {
				if !strings.Contains(di.Name, "/.service/") {
					itemsToShow = append(itemsToShow, di)
				}
			}

			needToSelectMainItem := false

			if len(itemsToShow) != c.lvItems.ItemsCount() {
				c.lvItems.RemoveItems()
				for i := 0; i < len(itemsToShow); i++ {
					c.lvItems.AddItem("---")
				}
				needToSelectMainItem = true
			}
			for index, di := range itemsToShow {
				shortName := di.Name
				lastIndexOfSlash := strings.LastIndex(shortName, "/")
				if lastIndexOfSlash > -1 {
					shortName = shortName[lastIndexOfSlash+1:]
				}

				value := di.Value.Value
				{
					if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
						p := message.NewPrinter(language.English)
						value = strings.ReplaceAll(p.Sprint(intValue), ",", " ")
					} else {
						value = strings.ReplaceAll(value, "\r", " ")
						value = strings.ReplaceAll(value, "\n", " ")
					}
				}

				dt := time.Unix(0, di.Value.DT*1000).Format("15:04:05")

				c.lvItems.Item(index).TempData = di.Name
				c.lvItems.Item(index).SetUserData("cloud", di.CloudChannels)
				c.lvItems.SetItemValue(index, 0, shortName)
				c.lvItems.SetItemValue(index, 1, value)
				c.lvItems.SetItemValue(index, 2, di.Value.UOM)
				c.lvItems.SetItemValue(index, 3, dt)
				if len(di.CloudChannels) == 0 {
					c.lvItems.Item(index).SetForeColorForCell(4, c.InactiveColor())
					c.lvItems.SetItemValue(index, 4, "local only")
				} else {
					c.lvItems.Item(index).SetForeColorForCell(4, settings.GoodColor)
					channels := ""
					for _, ch := range di.CloudChannelsNames {
						if len(channels) > 0 {
							channels += ", "
						}
						channels += ch
					}

					c.lvItems.SetItemValue(index, 4, "> "+channels)

					//c.lvItems.SetItemValue(c.lvItems.Item(index), 4, "cloud ("+strconv.Itoa(len(di.CloudChannels))+")")
				}

				//c.lvItems.SetItemValue(c.lvItems.Item(index), 4, strings.Repeat("∆", len(di.CloudChannels)))
				//c.lvItems.SetItemValue(c.lvItems.Item(index), 4, fmt.Sprint(di.CloudChannels))

				if di.Value.UOM == "error" {
					c.lvItems.Item(index).SetForeColorForCell(1, settings.BadColor)
					c.lvItems.Item(index).SetForeColorForCell(2, settings.BadColor)
				} else {
					c.lvItems.Item(index).SetForeColorForCell(1, settings.GoodColor)
					c.lvItems.Item(index).SetForeColorForCell(2, nil)
				}

				if needToSelectMainItem && di.Name == c.currentMainItem {
					c.lvItems.SelectItem(index)
				}
			}
		})
	} else {
		c.lvItems.RemoveItems()
	}

}

func (c *PanelUnits) timerUpdate() {
	if c.Disposed() {
		return
	}

	c.updateUnitsButtons()
	c.updateUnitsState()

	c.updateDataItemsButtons()
	c.updateDataItemsState()
}

func (c *PanelUnits) addSelectedItemsToCloud() {
	cloudAccountsForItems := make(map[string]int)

	for i := 0; i < c.lvItems.ItemsCount(); i++ {
		item := c.lvItems.Item(i)

		if cloudChannels, ok := item.UserData("cloud").([]string); ok {
			for _, ch := range cloudChannels {
				if _, ok := cloudAccountsForItems[ch]; ok {
					cloudAccountsForItems[ch]++
				} else {
					cloudAccountsForItems[ch] = 1
				}
			}
		}
	}

	prefChannels := make([]string, 0)
	for ch := range cloudAccountsForItems {
		prefChannels = append(prefChannels, ch)
	}

	f := NewFormAddToCloud(c, c.client, c.SelectedItems(), c.AllItems(), prefChannels)
	f.SetAllItemsCheckBox(false)
	f.ShowDialog()
	f.OnAccept = func() {
	}
}

func (c *PanelUnits) removeSelectedItemsFromCloud() {
	cloudAccountsForItems := make(map[string]int)

	for _, item := range c.lvItems.SelectedItems() {
		if cloudChannels, ok := item.UserData("cloud").([]string); ok {
			for _, ch := range cloudChannels {
				if _, ok := cloudAccountsForItems[ch]; ok {
					cloudAccountsForItems[ch]++
				} else {
					cloudAccountsForItems[ch] = 1
				}
			}
		}
	}

	prefChannels := make([]string, 0)
	for ch := range cloudAccountsForItems {
		prefChannels = append(prefChannels, ch)
	}

	f := NewFormRemoveFromCloud(c, c.client, c.SelectedItems(), c.AllItems(), prefChannels)
	f.SetAllItemsCheckBox(false)
	f.ShowDialog()
	f.OnAccept = func() {
	}
}

func (c *PanelUnits) openSelectedItemInBrowser() {
	for _, item := range c.lvItems.SelectedItems() {
		channel := ""
		if cloudChannels, ok := item.UserData("cloud").([]string); ok {
			if len(cloudChannels) > 0 {
				channel = cloudChannels[0]
			}
		}

		b64 := base64.StdEncoding.EncodeToString([]byte(item.TempData))
		client.OpenBrowser("https://gazer.cloud/item/" + channel + "/" + b64)
	}

	return
	channels := make(map[string]string)
	for _, item := range c.lvItems.SelectedItems() {
		if cloudChannels, ok := item.UserData("cloud").([]string); ok {
			for _, ch := range cloudChannels {
				channels[ch] = "1"
			}
		}
	}
	for channelId, _ := range channels {
		client.OpenBrowser(gazer_dictionary.ChannelUrl(channelId))
	}
}
