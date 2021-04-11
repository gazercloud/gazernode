package forms

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/gazer_dictionary"
	"github.com/gazercloud/gazernode/settings"
	"github.com/gazercloud/gazernode/system/cloud"
	"github.com/gazercloud/gazerui/canvas"
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
	"github.com/go-gl/glfw/v3.3/glfw"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"os"
	"strconv"
	"strings"
	"time"
)

type PanelCloud struct {
	uicontrols.Panel
	client *client.Client

	lvChannels *uicontrols.ListView

	btnAdd    *uicontrols.Button
	btnEdit   *uicontrols.Button
	btnRemove *uicontrols.Button

	btnRefresh *uicontrols.Button

	txtHeaderChartGroup *uicontrols.TextBlock

	btnStart         *uicontrols.Button
	btnStop          *uicontrols.Button
	btnOpenInBrowser *uicontrols.Button
	btnCopyLink      *uicontrols.Button

	btnShowFullScreen  *uicontrols.Button
	btnRemoveFromCloud *uicontrols.Button

	txtLink        *uicontrols.TextBox
	lblChannelName *uicontrols.TextBlock

	lvItems          *uicontrols.ListView
	timer            *uievents.FormTimer
	unitId           string
	currentChannelId string
}

func NewPanelCloud(parent uiinterfaces.Widget, client *client.Client) *PanelCloud {
	var c PanelCloud
	c.client = client
	c.InitControl(parent, &c)

	return &c
}

func (c *PanelCloud) OnInit() {
	//pHeader := c.AddPanelOnGrid(0, 0)
	//txtHeader := pHeader.AddTextBlockOnGrid(0, 0, "Public channels")
	//txtHeader.SetFontSize(24)

	pContent := c.AddPanelOnGrid(0, 1)
	pContent.SetPanelPadding(0)
	splitter := pContent.AddSplitContainerOnGrid(0, 0)
	splitter.SetPosition(360)
	splitter.SetYExpandable(true)

	pUnitsList := splitter.Panel1.AddPanelOnGrid(0, 0)
	pUnitsList.SetPanelPadding(0)

	txtHeader := pUnitsList.AddTextBlockOnGrid(0, 0, "Public channels")
	txtHeader.SetFontSize(24)

	pButtons := pUnitsList.AddPanelOnGrid(0, 1)
	pButtons.SetPanelPadding(0)

	c.btnAdd = pButtons.AddButtonOnGrid(0, 0, "", func(event *uievents.Event) {
		f := NewFormAddCloudChannel(c, c.client)
		f.ShowDialog()
		f.OnAccept = func() {
			c.loadChannels()
		}
	})
	c.btnAdd.SetTooltip("Add new public channel")

	c.btnEdit = pButtons.AddButtonOnGrid(1, 0, "", func(event *uievents.Event) {
		channelId := c.lvChannels.SelectedItem().TempData
		f := NewFormEditCloudChannel(c, c.client, channelId, c.lvChannels.SelectedItem().UserData("channelName").(string))
		f.ShowDialog()
		f.OnAccept = func() {
			c.loadChannels()
		}
	})
	c.btnEdit.SetTooltip("Edit public channel")

	c.btnRemove = pButtons.AddButtonOnGrid(2, 0, "", func(event *uievents.Event) {
		uicontrols.ShowQuestionMessageOKCancel(c, "Remove public channel?", "Confirmation", func() {
			channelId := c.lvChannels.SelectedItem().TempData
			c.client.RemoveCloudChannel(channelId, func(err error) {
				c.loadChannels()
			})
		}, nil)
	})
	c.btnRemove.SetTooltip("Remove public channel")

	pButtons.AddTextBlockOnGrid(3, 0, " | ")

	c.btnRefresh = pButtons.AddButtonOnGrid(4, 0, "", func(event *uievents.Event) {
		c.loadChannels()
	})
	c.btnRefresh.SetTooltip("Refresh")

	pButtons.AddHSpacerOnGrid(5, 0)
	/*
		c.btnStart = pButtons.AddButtonOnGrid(4, 0, "", func(event *uievents.Event) {
			unitId := c.lvChannels.SelectedItem().TempData
			c.client.StartUnit(unitId, nil)
		})
		c.btnStart.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_av_play_arrow_materialicons_48dp_1x_baseline_play_arrow_black_48dp_png, c.ForeColor()))

		c.btnStop = pButtons.AddButtonOnGrid(5, 0, "", func(event *uievents.Event) {
			unitId := c.lvChannels.SelectedItem().TempData
			c.client.StopUnit(unitId, nil)
		})
		c.btnStop.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_av_pause_materialiconsoutlined_48dp_1x_outline_pause_black_48dp_png, c.ForeColor()))
	*/
	//pButtons.AddHSpacerOnGrid(5, 0)

	c.lvChannels = pUnitsList.AddListViewOnGrid(0, 2)
	c.lvChannels.AddColumn("Name", 200)
	c.lvChannels.AddColumn("Id", 100)
	c.lvChannels.OnSelectionChanged = c.loadSelected

	menu := uicontrols.NewPopupMenu(c.lvChannels)
	menu.AddItem("Open in Browser ...", func(event *uievents.Event) {
		channelId := c.lvChannels.SelectedItem().TempData
		client.OpenBrowser(gazer_dictionary.ChannelUrl(channelId))
	}, uiresources.ResImgCol(uiresources.R_icons_material4_png_action_open_in_browser_materialicons_48dp_1x_baseline_open_in_browser_black_48dp_png, c.ForeColor()), "")
	menu.AddItem("Copy link to clipboard", func(event *uievents.Event) {
		channelId := c.lvChannels.SelectedItem().TempData
		link := gazer_dictionary.ChannelUrl(channelId)
		glfw.SetClipboardString(link)
	}, uiresources.ResImgCol(uiresources.R_icons_material4_png_content_content_copy_materialicons_48dp_1x_baseline_content_copy_black_48dp_png, c.ForeColor()), "")
	c.lvChannels.SetContextMenu(menu)

	pHeaderRight := splitter.Panel2.AddPanelOnGrid(0, 0)
	pHeaderRight.SetPanelPadding(0)
	c.txtHeaderChartGroup = pHeaderRight.AddTextBlockOnGrid(0, 0, "")
	c.txtHeaderChartGroup.SetFontSize(24)

	pItems := splitter.Panel2.AddPanelOnGrid(0, 1)
	pItems.SetPanelPadding(0)

	pItems.SetOnKeyDown(func(event *uievents.KeyDownEvent) bool {
		if event.Key == glfw.KeyEnter || event.Key == glfw.KeyKPEnter {
			items := c.SelectedItems()
			if len(items) > 0 {
				MainFormInstance.ShowFullScreenValue(true, items[0])
			}
			return true
		}
		return false
	})

	pButtonsRight := pItems.AddPanelOnGrid(0, 0)
	pButtonsRight.SetPanelPadding(0)

	// LINK CONTROL
	pLink := pButtonsRight.AddPanelOnGrid(0, 0)
	pLink.SetPanelPadding(0)

	/*pLinkTxt := pLink.AddPanelOnGrid(0, 1)
	pLinkTxt.SetPanelPadding(0)*/

	c.txtLink = pLink.AddTextBoxOnGrid(0, 1)
	c.txtLink.SetReadOnly(true)

	c.btnOpenInBrowser = pLink.AddButtonOnGrid(1, 0, "", func(event *uievents.Event) {
		channelId := c.currentChannelId
		client.OpenBrowser(gazer_dictionary.ChannelUrl(channelId))
	})
	c.btnOpenInBrowser.SetTooltip("Open in browser")

	c.btnCopyLink = pLink.AddButtonOnGrid(1, 1, "", func(event *uievents.Event) {
		channelId := c.currentChannelId
		link := gazer_dictionary.ChannelUrl(channelId)
		glfw.SetClipboardString(link)
	})
	c.btnCopyLink.SetTooltip("Copy hyper-link to clipboard")

	//pLinkButtons := pLink.AddPanelOnGrid(0, 0)
	//pLinkButtons.SetPanelPadding(0)

	c.lblChannelName = pLink.AddTextBlockOnGrid(0, 0, "")
	c.lblChannelName.SetUnderline(true)
	c.lblChannelName.SetFontSize(18)
	c.lblChannelName.OnClick = func(ev *uievents.Event) {
		client.OpenBrowser(gazer_dictionary.ChannelUrl(c.currentChannelId))
	}
	c.lblChannelName.SetMouseCursor(ui.MouseCursorPointer)

	//pLinkButtons.AddHSpacerOnGrid(1, 0)

	//pLinkButtons.AddHSpacerOnGrid(5, 0)
	////////////////LINK CONTROL

	//

	c.btnShowFullScreen = pLink.AddButtonOnGrid(2, 1, "", func(event *uievents.Event) {
		items := c.SelectedItems()
		if len(items) > 0 {
			MainFormInstance.ShowFullScreenValue(true, items[0])
		}
	})
	c.btnShowFullScreen.SetTooltip("Full screen")

	c.btnRemoveFromCloud = pLink.AddButtonOnGrid(3, 1, "", func(event *uievents.Event) {
		uicontrols.ShowQuestionMessageOKCancel(c, "Remove selected items?", "Confirmation", func() {
			items := c.SelectedItems()
			c.client.CloudRemoveItems([]string{c.currentChannelId}, items, nil)
		}, nil)
	})
	c.btnRemoveFromCloud.SetTooltip("Remove selected items from the public channel")

	c.lvItems = pItems.AddListViewOnGrid(0, 1)
	c.lvItems.AddColumn("Name", 300)
	c.lvItems.AddColumn("Value", 100)
	c.lvItems.AddColumn("UOM", 60)
	c.lvItems.AddColumn("Time", 80)
	c.lvItems.SetColumnTextAlign(1, canvas.HAlignRight)
	c.timer = c.Window().NewTimer(500, c.timerUpdate)
	c.timer.StartTimer()

	c.loadChannels()
	c.UpdateStyle()

	c.createCloudChannelIfItDoesntExists()

	c.loadSelected()
}

func (c *PanelCloud) Dispose() {
	if c.timer != nil {
		c.timer.StopTimer()
	}
	c.timer = nil

	c.client = nil

	c.lvChannels = nil

	c.btnAdd = nil
	c.btnEdit = nil
	c.btnRemove = nil

	c.btnStart = nil
	c.btnStop = nil
	c.btnOpenInBrowser = nil
	c.btnCopyLink = nil

	c.btnShowFullScreen = nil
	c.btnRemoveFromCloud = nil

	c.txtLink = nil
	c.lblChannelName = nil

	c.lvItems = nil
	c.Panel.Dispose()
}

func (c *PanelCloud) loadSelected() {
	selectedItem := c.lvChannels.SelectedItem()
	if selectedItem != nil {
		name := c.lvChannels.SelectedItem().UserData("channelName").(string)
		c.txtHeaderChartGroup.SetText("Channel: " + name)
		unitId := selectedItem.TempData
		c.SetCurrentChannelId(unitId, name)
	} else {
		c.txtHeaderChartGroup.SetText("no channel selected")
		unitId := ""
		c.SetCurrentChannelId(unitId, "")
	}
}

func (c *PanelCloud) FullRefresh() {
	c.loadChannels()
}

func (c *PanelCloud) UpdateStyle() {
	c.Panel.UpdateStyle()

	activeColor := c.AccentColor()
	inactiveColor := c.InactiveColor()

	c.btnAdd.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_add_materialicons_48dp_1x_baseline_add_black_48dp_png, activeColor))
	c.btnEdit.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_create_materialicons_48dp_1x_baseline_create_black_48dp_png, activeColor))
	c.btnRemove.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_clear_materialicons_48dp_1x_baseline_clear_black_48dp_png, activeColor))
	//c.btnStart.SetImage(uiresources.ResImageAdjusted("icons/material/av/drawable-hdpi/ic_play_arrow_black_48dp.png", c.ForeColor()))
	//c.btnStop.SetImage(uiresources.ResImageAdjusted("icons/material/av/drawable-hdpi/ic_pause_black_48dp.png", c.ForeColor()))

	c.btnOpenInBrowser.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_action_open_in_browser_materialicons_48dp_1x_baseline_open_in_browser_black_48dp_png, activeColor))
	c.btnCopyLink.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_content_copy_materialicons_48dp_1x_baseline_content_copy_black_48dp_png, activeColor))

	c.btnShowFullScreen.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_navigation_fullscreen_materialicons_48dp_1x_baseline_fullscreen_black_48dp_png, activeColor))
	c.btnRemoveFromCloud.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_file_cloud_off_materialicons_48dp_1x_baseline_cloud_off_black_48dp_png, activeColor))

	c.btnAdd.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_add_materialicons_48dp_1x_baseline_add_black_48dp_png, inactiveColor))
	c.btnEdit.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_create_materialicons_48dp_1x_baseline_create_black_48dp_png, inactiveColor))
	c.btnRemove.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_clear_materialicons_48dp_1x_baseline_clear_black_48dp_png, inactiveColor))
	//c.btnStart.SetImageDisabled(uiresources.ResImageAdjusted("icons/material/av/drawable-hdpi/ic_play_arrow_black_48dp.png", c.InactiveColor()))
	//c.btnStop.SetImageDisabled(uiresources.ResImageAdjusted("icons/material/av/drawable-hdpi/ic_pause_black_48dp.png", c.InactiveColor()))

	c.btnOpenInBrowser.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_action_open_in_browser_materialicons_48dp_1x_baseline_open_in_browser_black_48dp_png, inactiveColor))
	c.btnCopyLink.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_content_copy_materialicons_48dp_1x_baseline_content_copy_black_48dp_png, inactiveColor))

	c.btnShowFullScreen.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_navigation_fullscreen_materialicons_48dp_1x_baseline_fullscreen_black_48dp_png, inactiveColor))
	c.btnRemoveFromCloud.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_file_cloud_off_materialicons_48dp_1x_baseline_cloud_off_black_48dp_png, inactiveColor))

	c.btnRefresh.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_navigation_refresh_materialicons_48dp_1x_baseline_refresh_black_48dp_png, activeColor))
	c.btnRefresh.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_navigation_refresh_materialicons_48dp_1x_baseline_refresh_black_48dp_png, inactiveColor))
}

func (c *PanelCloud) SetCurrentChannelId(channelId string, name string) {
	c.lvItems.RemoveItems()
	if len(channelId) > 0 {
		c.currentChannelId = channelId
		c.txtLink.SetText(gazer_dictionary.ChannelUrl(c.currentChannelId))
		c.lblChannelName.SetText(name + " (" + channelId + ")")
	} else {
		c.currentChannelId = channelId
		c.txtLink.SetText(gazer_dictionary.ChannelUrl(c.currentChannelId))
		c.lblChannelName.SetText("")
	}
}

func (c *PanelCloud) loadChannels() {
	if c.client == nil {
		return
	}
	c.client.GetCloudChannels(func(channels []cloud.ChannelInfo, err error) {
		if c.lvChannels == nil {
			return
		}
		c.lvChannels.RemoveItems()
		for _, s := range channels {
			lvItem := c.lvChannels.AddItem(s.Name)
			lvItem.SetValue(1, s.Id)
			lvItem.TempData = s.Id
			lvItem.SetUserData("channelName", s.Name)
		}
	})
}

func (c *PanelCloud) SelectedItems() []string {
	items := make([]string, 0)
	for _, item := range c.lvItems.SelectedItems() {
		name := item.TempData
		items = append(items, name)
	}
	return items
}

func (c *PanelCloud) createCloudChannelIfItDoesntExists() {
	c.client.GetCloudChannels(func(channels []cloud.ChannelInfo, err error) {
		if err == nil && len(channels) == 0 {
			var hostName string
			hostName, err = os.Hostname()
			if err != nil {
				hostName = "Default Channel"
			}
			c.client.AddCloudChannel(hostName, func(err error) {
				c.loadChannels()
			})
		}
	})
}

func (c *PanelCloud) timerUpdate() {
	if !c.IsVisible() {
		return
	}

	if len(c.lvChannels.SelectedItems()) > 0 {
		if len(c.lvChannels.SelectedItems()) == 1 {
			c.btnEdit.SetEnabled(true)
			c.btnRemove.SetEnabled(true)
		} else {
			c.btnEdit.SetEnabled(false)
			c.btnRemove.SetEnabled(false)
		}
		/*c.btnStart.SetEnabled(true)
		c.btnStop.SetEnabled(true)*/
	} else {
		c.btnEdit.SetEnabled(false)
		c.btnRemove.SetEnabled(false)
		/*c.btnStart.SetEnabled(false)
		c.btnStop.SetEnabled(false)*/
	}

	itemsSelected := c.lvItems.SelectedItems()

	if len(c.currentChannelId) > 0 {
		c.btnOpenInBrowser.SetEnabled(true)
		c.btnCopyLink.SetEnabled(true)
	} else {
		c.btnOpenInBrowser.SetEnabled(false)
		c.btnCopyLink.SetEnabled(false)
	}

	if len(itemsSelected) > 0 {
		if len(itemsSelected) == 1 {
			c.btnShowFullScreen.SetEnabled(true)
		} else {
			c.btnShowFullScreen.SetEnabled(false)
		}

		c.btnRemoveFromCloud.SetEnabled(true)

	} else {
		c.btnShowFullScreen.SetEnabled(false)
		c.btnRemoveFromCloud.SetEnabled(false)
	}

	if len(c.currentChannelId) > 0 {
		c.client.GetCloudChannelValues(c.currentChannelId, func(items []common_interfaces.Item, err error) {
			if len(items) != c.lvItems.ItemsCount() {
				c.lvItems.RemoveItems()
				for i := 0; i < len(items); i++ {
					c.lvItems.AddItem("---")
				}
			}
			for index, di := range items {

				value := di.Value.Value
				{
					if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
						p := message.NewPrinter(language.English)
						value = strings.ReplaceAll(p.Sprint(intValue), ",", " ")
					}
				}

				c.lvItems.Item(index).TempData = di.Name
				c.lvItems.SetItemValue(index, 0, di.Name)
				c.lvItems.SetItemValue(index, 1, value)
				c.lvItems.SetItemValue(index, 2, di.Value.UOM)
				c.lvItems.SetItemValue(index, 3, time.Unix(0, di.Value.DT*1000).Format("15:04:05"))

				if di.Value.UOM == "error" {
					c.lvItems.Item(index).SetForeColorForCell(1, settings.BadColor)
					c.lvItems.Item(index).SetForeColorForCell(2, settings.BadColor)
				} else {
					c.lvItems.Item(index).SetForeColorForCell(1, settings.GoodColor)
					c.lvItems.Item(index).SetForeColorForCell(2, nil)
				}
			}
		})
	}
}
