package simplemap

import (
	"bytes"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
	"github.com/gazercloud/gazerui/uistyles"
	"image"
	"image/png"
)

type MapToolbox struct {
	uicontrols.Panel
	mapDataSource IMapDataSource

	currentTool        string
	itemsCount         int
	defaultCurrentTool string

	panelFilter          *uicontrols.Panel
	panelStandardButtons *uicontrols.Panel
	panelButtons         *uicontrols.Panel
	panelNavButtons      *uicontrols.Panel
	btnNavLeft           *uicontrols.Button
	btnNavRight          *uicontrols.Button

	txtFilter *uicontrols.TextBox
	btnSearch *uicontrols.Button

	btnNoItem *uicontrols.Button
	btnItem   *uicontrols.Button
	btnCircle *uicontrols.Button
	btnLine   *uicontrols.Button

	standardButtons []*uicontrols.Button
	buttons         []*uicontrols.Button

	offset             int
	maxOffset          int
	countOfItemsOnPage int
}

func NewMapToolbox(parent uiinterfaces.Widget) *MapToolbox {
	var c MapToolbox
	c.Panel.InitControl(parent, &c)
	c.Panel.SetName("MapToolBox")
	c.Panel.SetPanelPadding(0)
	c.panelStandardButtons = c.Panel.AddPanelOnGrid(0, 0)
	c.panelStandardButtons.SetPanelPadding(0)
	c.Panel.AddTextBlockOnGrid(0, 1, " ")
	c.panelFilter = c.Panel.AddPanelOnGrid(0, 2)
	c.panelFilter.SetPanelPadding(0)
	c.panelButtons = c.Panel.AddPanelOnGrid(0, 3)
	c.panelButtons.SetPanelPadding(0)
	c.panelNavButtons = c.Panel.AddPanelOnGrid(0, 4)
	c.panelNavButtons.SetPanelPadding(0)
	c.Panel.AddVSpacerOnGrid(0, 5)

	c.countOfItemsOnPage = 7
	c.defaultCurrentTool = ""
	c.currentTool = c.defaultCurrentTool

	c.panelButtons.SetMinHeight(c.countOfItemsOnPage * 64)

	c.btnNavLeft = c.panelNavButtons.AddButtonOnGrid(0, 0, "<", func(event *uievents.Event) {
		c.offset -= c.countOfItemsOnPage
		if c.offset < 0 {
			c.offset = 0
		}
		c.loadItems()
	})
	c.btnNavRight = c.panelNavButtons.AddButtonOnGrid(1, 0, ">", func(event *uievents.Event) {
		c.offset += c.countOfItemsOnPage
		if c.offset > c.maxOffset {
			c.offset = c.maxOffset
		}
		c.loadItems()
	})

	c.txtFilter = c.panelFilter.AddTextBoxOnGrid(0, 0)
	c.txtFilter.SetEmptyText("Search ...")
	c.txtFilter.OnTextChanged = func(txtBox *uicontrols.TextBox, oldValue string, newValue string) {
		c.offset = 0
		c.loadItems()
	}
	c.btnSearch = c.panelFilter.AddButtonOnGrid(1, 0, "", func(event *uievents.Event) {
		c.offset = 0
		c.loadItems()
	})

	c.btnNoItem = c.panelStandardButtons.AddButtonOnGrid(0, 0, "", c.onBtnClicked)
	c.btnNoItem.SetImageSize(24, 24)
	c.standardButtons = append(c.standardButtons, c.btnNoItem)
	c.btnNoItem.SetUserData("codeName", "")

	//pStItems := c.panelStandardButtons.AddPanelOnGrid(0, 1)
	//pStItems.SetPanelPadding(0)

	c.btnItem = c.panelStandardButtons.AddButtonOnGrid(1, 0, "", c.onBtnClicked)
	c.standardButtons = append(c.standardButtons, c.btnItem)
	c.btnItem.SetImageSize(24, 24)
	c.btnItem.SetUserData("codeName", "text")

	c.btnCircle = c.panelStandardButtons.AddButtonOnGrid(0, 1, "", c.onBtnClicked)
	c.standardButtons = append(c.standardButtons, c.btnCircle)
	c.btnCircle.SetImageSize(24, 24)
	c.btnCircle.SetUserData("codeName", "circle")

	c.btnLine = c.panelStandardButtons.AddButtonOnGrid(1, 1, "", c.onBtnClicked)
	c.btnLine.SetImageSize(24, 24)
	c.standardButtons = append(c.standardButtons, c.btnLine)
	c.btnLine.SetUserData("codeName", "line")

	c.loadItems()
	return &c
}

func (c *MapToolbox) Dispose() {
	c.buttons = nil
	c.standardButtons = nil
	c.Panel.Dispose()
}

func (c *MapToolbox) Reset() {
	c.currentTool = c.defaultCurrentTool
	c.txtFilter.SetText("")
	c.loadItems()
}

func (c *MapToolbox) SetMapDataSource(mapDataSource IMapDataSource) {
	c.mapDataSource = mapDataSource
	c.currentTool = c.defaultCurrentTool
	c.loadItems()
}

func (c *MapToolbox) loadItems() {
	c.buttons = make([]*uicontrols.Button, 0)
	c.panelButtons.RemoveAllWidgets()
	c.itemsCount = 0
	c.currentTool = c.defaultCurrentTool
	c.updateButtonsState()
	c.panelButtons.AddTextBlockOnGrid(0, len(c.panelButtons.Widgets()), "loading ...")

	if c.mapDataSource != nil {
		c.mapDataSource.GetWidgets(c.txtFilter.Text(), c.offset, c.countOfItemsOnPage, c)
	}
}

func (c *MapToolbox) AddItem(displayName string, codeName string, image image.Image, item common_interfaces.ResourcesItemInfo) {
	btn := uicontrols.NewButton(c, " "+displayName, c.onBtnClicked)
	btn.SetUserData("codeName", codeName)
	btn.SetUserData("item", item)
	btn.SetImageSize(48, 48)
	btn.SetTextImageVerticalOrientation(false)
	btn.SetImage(image)
	c.panelButtons.AddWidgetOnGrid(btn, 0, c.itemsCount)
	c.buttons = append(c.buttons, btn)
	c.updateButtonsState()
	c.itemsCount++
}

func (c *MapToolbox) AddStopItem() {
	c.panelButtons.AddVSpacerOnGrid(0, c.itemsCount)
	c.updateButtonsState()
	c.itemsCount++
}

func (c *MapToolbox) onBtnClicked(ev *uievents.Event) {
	c.currentTool = ev.Sender.(*uicontrols.Button).UserData("codeName").(string)
	c.updateButtonsState()
}

func (c *MapToolbox) updateButtonsState() {
	allButtons := make([]*uicontrols.Button, 0)
	allButtons = append(allButtons, c.standardButtons...)
	allButtons = append(allButtons, c.buttons...)
	for _, btn := range allButtons {
		codeName, ok := btn.UserData("codeName").(string)
		if ok {
			if codeName == c.currentTool {
				btn.SetForeColor(uistyles.DefaultBackColor)
				btn.SetBackColor(c.ForeColor())
			} else {
				btn.SetForeColor(nil)
				btn.SetBackColor(nil)
			}
		}

		resItemInfo, ok := btn.UserData("item").(common_interfaces.ResourcesItemInfo)
		if ok {
			var thumbnail image.Image
			if resItemInfo.Thumbnail != nil {
				thumbnail, _ = png.Decode(bytes.NewBuffer(resItemInfo.Thumbnail))
			} else {
				thumbnail = image.NewAlpha(image.Rect(0, 0, 32, 32))
			}
			btn.SetImage(thumbnail)
			_ = resItemInfo
		}

		if btn == c.btnNoItem {
			c.btnNoItem.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_action_highlight_alt_materialicons_48dp_1x_baseline_highlight_alt_black_48dp_png, btn.ForeColor()))
		}
		if btn == c.btnItem {
			c.btnItem.SetImage(uiresources.ResImgCol(uiresources.R_icons_custom_rect_item_png, btn.ForeColor()))
		}
		if btn == c.btnCircle {
			c.btnCircle.SetImage(uiresources.ResImgCol(uiresources.R_icons_custom_circle_item_png, btn.ForeColor()))
		}
		if btn == c.btnLine {
			c.btnLine.SetImage(uiresources.ResImgCol(uiresources.R_icons_custom_line_png, btn.ForeColor()))
		}
	}
}

func (c *MapToolbox) CurrentTool() string {
	return c.currentTool
}

func (c *MapToolbox) ResetCurrentTool() {
	c.currentTool = c.defaultCurrentTool
	c.updateButtonsState()
}

func (c *MapToolbox) UpdateStyle() {
	c.Panel.UpdateStyle()
	c.updateButtonsState()

	c.btnSearch.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_action_youtube_searched_for_materialicons_48dp_1x_baseline_youtube_searched_for_black_48dp_png, c.ForeColor()))
}

func (c *MapToolbox) SetItems(widgets common_interfaces.ResourcesInfo) {
	c.panelButtons.RemoveAllWidgets()

	for _, w := range widgets.Items {
		var thumbnail image.Image
		if w.Thumbnail != nil {
			thumbnail, _ = png.Decode(bytes.NewBuffer(w.Thumbnail))
		}

		c.AddItem(w.Name, w.Id, thumbnail, w)
	}

	c.AddStopItem()

	if c.offset == 0 {
		c.btnNavLeft.SetEnabled(false)
	} else {
		c.btnNavLeft.SetEnabled(true)
	}

	c.maxOffset = widgets.InFilterCount / c.countOfItemsOnPage
	c.maxOffset *= c.countOfItemsOnPage

	if c.offset >= c.maxOffset {
		c.btnNavRight.SetEnabled(false)
	} else {
		c.btnNavRight.SetEnabled(true)
	}

}
