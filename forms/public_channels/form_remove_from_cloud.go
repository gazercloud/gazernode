package public_channels

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/system/public_channel"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
)

type FormRemoveFromCloud struct {
	uicontrols.Dialog
	client            *client.Client
	txtUnitName       *uicontrols.TextBox
	txtError          *uicontrols.TextBlock
	lvChannels        *uicontrols.ListView
	chkAllItems       *uicontrols.CheckBox
	items             []string
	allItems          []string
	currentChannelIds []string
	preferredChannels []string
}

func NewFormRemoveFromCloud(parent uiinterfaces.Widget, client *client.Client, items []string, allItems []string, preferredChannels []string) *FormRemoveFromCloud {
	var c FormRemoveFromCloud
	c.client = client
	c.items = items
	c.allItems = allItems
	c.preferredChannels = preferredChannels
	c.InitControl(parent, &c)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	pLeft := pContent.AddPanelOnGrid(0, 0)
	pLeft.SetPanelPadding(0)
	//pLeft.SetBorderRight(1, c.ForeColor())
	pLeft.SetMinWidth(100)
	pRight := pContent.AddPanelOnGrid(1, 0)
	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)

	img := pLeft.AddImageBoxOnGrid(0, 0, uiresources.ResImgCol(uiresources.R_icons_material4_png_file_cloud_off_materialiconsoutlined_48dp_1x_outline_cloud_off_black_48dp_png, c.AccentColor()))
	img.SetScaling(uicontrols.ImageBoxScaleAdjustImageKeepAspectRatio)
	img.SetMinHeight(64)
	img.SetMinWidth(64)
	pLeft.AddVSpacerOnGrid(0, 1)

	pRight.AddTextBlockOnGrid(0, 0, "Select cloud channel:")

	c.lvChannels = pRight.AddListViewOnGrid(0, 1)
	c.lvChannels.AddColumn("Channel Id", 200)
	c.lvChannels.AddColumn("Name", 200)
	c.lvChannels.OnSelectionChanged = func() {
		c.currentChannelIds = make([]string, 0)
		if len(c.lvChannels.SelectedItems()) > 0 {
			for _, acc := range c.lvChannels.SelectedItems() {
				channelId := acc.TempData
				c.currentChannelIds = append(c.currentChannelIds, channelId)
			}
		}
	}
	c.chkAllItems = pRight.AddCheckBoxOnGrid(0, 2, "All the items of the unit")
	c.chkAllItems.SetChecked(true)

	c.loadChannels()

	pButtons.AddHSpacerOnGrid(0, 0)
	btnOK := pButtons.AddButtonOnGrid(1, 0, "OK", nil)
	btnOK.SetMinWidth(70)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Cancel", nil)
	btnCancel.SetMinWidth(70)

	c.TryAccept = func() bool {
		c.currentChannelIds = make([]string, 0)
		if len(c.lvChannels.SelectedItems()) > 0 {
			for _, ch := range c.lvChannels.SelectedItems() {
				channelId := ch.TempData
				c.currentChannelIds = append(c.currentChannelIds, channelId)
			}
		}
		if len(c.currentChannelIds) > 0 {
			itemsToRemove := c.items
			if c.chkAllItems.IsChecked() {
				itemsToRemove = c.allItems
			}
			c.client.CloudRemoveItems(c.currentChannelIds, itemsToRemove, func(err error) {
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

func (c *FormRemoveFromCloud) loadChannels() {
	c.client.GetCloudChannels(func(channels []public_channel.ChannelInfo, err error) {
		c.lvChannels.RemoveItems()
		for _, s := range channels {
			lvItem := c.lvChannels.AddItem(s.Id)
			lvItem.SetValue(1, s.Name)
			lvItem.TempData = s.Id
			lvItem.SetUserData("channelName", s.Name)
		}

		c.lvChannels.ClearSelection()
		if len(c.preferredChannels) > 0 {
			for i := 0; i < c.lvChannels.ItemsCount(); i++ {
				item := c.lvChannels.Item(i)
				for _, pCh := range c.preferredChannels {
					if pCh == item.TempData {
						c.lvChannels.SelectItemSelection(i, true)
					}
				}
			}
		} else {
			if c.lvChannels.ItemsCount() > 0 {
				c.lvChannels.SelectItem(0)
			}
		}
	})
}

func (c *FormRemoveFromCloud) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("Remove from cloud")
	c.Resize(600, 400)
}

func (c *FormRemoveFromCloud) SetAllItemsCheckBox(checked bool) {
	c.chkAllItems.SetChecked(checked)
}
