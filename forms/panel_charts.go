package forms

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/simplemap"
	"github.com/gazercloud/gazernode/widgets/widget_chart"
	"github.com/gazercloud/gazernode/widgets/widget_dataitems"
	"github.com/gazercloud/gazerui/timechart"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
	"image"
	"image/color"
)

type PanelCharts struct {
	uicontrols.Panel
	client     *client.Client
	pUnitsList *uicontrols.Panel
	lvItems    *uicontrols.ListView
	timer      *uievents.FormTimer

	btnAdd    *uicontrols.Button
	btnRename *uicontrols.Button
	btnRemove *uicontrols.Button

	btnEdit   *uicontrols.Button
	btnReject *uicontrols.Button
	btnSave   *uicontrols.Button

	splitterEditor *uicontrols.SplitContainer

	timeChart      *widget_chart.WidgetCharts
	itemsPanel     *widget_dataitems.WidgetDataItems
	currentResId   string
	currentResType string

	isEditing_ bool
}

func NewPanelCharts(parent uiinterfaces.Widget, client *client.Client) *PanelCharts {
	var c PanelCharts
	c.client = client
	c.InitControl(parent, &c)
	c.SetName("PanelCharts")
	return &c
}

func (c *PanelCharts) OnInit() {

	pContent := c.AddPanelOnGrid(0, 1)
	pContent.SetPanelPadding(0)
	splitter := pContent.AddSplitContainerOnGrid(0, 0)
	splitter.SetYExpandable(true)
	splitter.SetPosition(250)

	c.pUnitsList = splitter.Panel1.AddPanelOnGrid(0, 0)
	c.pUnitsList.SetPanelPadding(0)
	txtHeader := c.pUnitsList.AddTextBlockOnGrid(0, 0, "Chart groups")
	txtHeader.SetFontSize(24)

	pButtons := c.pUnitsList.AddPanelOnGrid(0, 1)
	pButtons.SetPanelPadding(0)

	c.lvItems = c.pUnitsList.AddListViewOnGrid(0, 2)
	c.lvItems.AddColumn("Name", 230)
	c.lvItems.OnSelectionChanged = func() {
		c.loadSelected()
	}

	c.btnAdd = pButtons.AddButtonOnGrid(0, 0, "", func(event *uievents.Event) {
		d := NewFormAddChartGroup(c, c.client, "chart_group")
		d.ShowDialog()
		d.OnAccept = func() {
			c.loadChartGroups(d.Id)
		}
	})
	c.btnAdd.SetTooltip("Add chart group ...")
	c.btnAdd.SetMinWidth(60)

	c.btnRename = pButtons.AddButtonOnGrid(1, 0, "", func(event *uievents.Event) {
		if len(c.lvItems.SelectedItems()) != 1 {
			return
		}
		item := c.lvItems.SelectedItems()[0]
		dialog := NewDialogEditChartGroupName(c, c.client, item.TempData, item.Value(0))
		dialog.ShowDialog()
		dialog.OnAccept = func() {
			c.loadChartGroups("")
		}
	})
	c.btnRename.SetTooltip("Rename chart group")

	c.btnRemove = pButtons.AddButtonOnGrid(2, 0, "", func(event *uievents.Event) {
		if len(c.lvItems.SelectedItems()) != 1 {
			return
		}
		uicontrols.ShowQuestionMessageOKCancel(c, "Remove chart group?", "Confirmation", func() {
			item := c.lvItems.SelectedItems()[0]
			c.client.ResRemove(item.TempData, func(err error) {
				c.loadChartGroups("")
			})
		}, nil)
	})
	c.btnRemove.SetTooltip("Remove selected chart group")

	pButtons.AddHSpacerOnGrid(3, 0)

	pMapButtons := splitter.Panel2.AddPanelOnGrid(0, 0)
	c.btnEdit = pMapButtons.AddButtonOnGrid(0, 0, "", func(event *uievents.Event) {
		c.SetEdit(true)
		c.updateButtons()
	})
	c.btnEdit.SetTooltip("Switch to Edit")
	c.btnEdit.SetMinWidth(70)
	c.btnReject = pMapButtons.AddButtonOnGrid(1, 0, "", func(event *uievents.Event) {
		c.SetEdit(false)
		c.loadSelected()
		c.updateButtons()
	})
	c.btnReject.SetTooltip("Reject changes")
	c.btnReject.SetMinWidth(70)
	c.btnSave = pMapButtons.AddButtonOnGrid(2, 0, "", func(event *uievents.Event) {
		if c.currentResId != "" {
			c.client.ResSet(c.currentResId, nil, c.Save(), func(err error) {
				c.SetEdit(false)
				c.updateButtons()
			})
		}
	})
	c.btnSave.SetTooltip("Save")
	c.btnSave.SetMinWidth(70)

	pMapButtons.AddHSpacerOnGrid(3, 0)

	c.splitterEditor = splitter.Panel2.AddSplitContainerOnGrid(0, 1)
	c.splitterEditor.SetYExpandable(true)

	c.timeChart = widget_chart.NewWidgetCharts(c, c.client)
	c.splitterEditor.Panel1.AddWidgetOnGrid(c.timeChart, 0, 1)

	c.timeChart.SetOnChartContextMenuNeed(func(timeChart *timechart.TimeChart, area *timechart.Area, areaIndex int) uiinterfaces.Menu {
		var m *uicontrols.PopupMenu
		if c.isEditing_ {
			m = uicontrols.NewPopupMenu(c.Window().CentralWidget())
			if area != nil {
				if area.ShowQualities() {
					m.AddItem("Hide qualities", func(event *uievents.Event) {
						area.SetShowQualities(false)
					}, nil, "")
				} else {
					m.AddItem("Show qualities", func(event *uievents.Event) {
						area.SetShowQualities(true)
					}, nil, "")
				}

				if area.UnitedScale() {
					m.AddItem("Deactivate United Scale", func(event *uievents.Event) {
						area.SetUnitedScale(false)
					}, nil, "")
				} else {
					m.AddItem("Activate United Scale", func(event *uievents.Event) {
						area.SetUnitedScale(true)
					}, nil, "")
				}

				selectColor := func(event *uievents.Event) {
					ser := event.Sender.(*uicontrols.PopupMenuItem).UserData("ser").(*timechart.Series)
					dialog := uicontrols.NewColorSelector(c, ser.Color())
					dialog.OnColorSelected = func(col color.Color) {
						ser.SetColor(col)
					}
					dialog.OnAccept = func() {
						ser.SetColor(dialog.ResColor())
					}
					dialog.OnReject = func() {
						ser.SetColor(dialog.ResColor())
					}
					dialog.ShowDialog()
				}

				serRemove := func(event *uievents.Event) {
					serIndex := event.Sender.(*uicontrols.PopupMenuItem).UserData("serIndex").(int)
					area.RemoveSeriesByIndex(serIndex)
				}

				for serIndex, ser := range area.Series() {
					menuSeries := uicontrols.NewPopupMenu(c.Window().CentralWidget())
					itemSelectColor := menuSeries.AddItem("Change color", selectColor, nil, "")
					itemSelectColor.SetUserData("ser", ser)
					itemRemove := menuSeries.AddItem("Remove series", serRemove, nil, "")
					itemRemove.SetUserData("serIndex", serIndex)

					serItemImage := image.NewRGBA(image.Rect(0, 0, 16, 16))
					for y := 0; y < 16; y++ {
						for x := 0; x < 16; x++ {
							serItemImage.Set(x, y, ser.Color())
						}
					}

					serItem := m.AddItemWithSubmenu(ser.Id(), serItemImage, menuSeries)
					serItem.AdjustColorForImage = false
				}
				m.AddItem("Remove area", func(event *uievents.Event) {
					timeChart.RemoveAreaByIndex(areaIndex)
				}, nil, "")
			}
			m.AddItem("Save changes", func(event *uievents.Event) {
				c.btnSave.Press()
			}, nil, "")
		} else {
			m = uicontrols.NewPopupMenu(c.Window().CentralWidget())
			m.AddItem("Switch to edit mode", func(event *uievents.Event) {
				c.btnEdit.Press()
			}, nil, "")
		}
		return m
	})

	c.itemsPanel = widget_dataitems.NewWidgetDataItems(c, c.client)
	c.splitterEditor.Panel2.AddWidgetOnGrid(c.itemsPanel, 0, 1)

	c.loadChartGroups("")
	c.SetEdit(false)
	c.updateButtons()
}

func (c *PanelCharts) Dispose() {
	if c.timer != nil {
		c.timer.StopTimer()
	}
	c.timer = nil
	c.client = nil

	c.Panel.Dispose()
}

func (c *PanelCharts) FullRefresh() {
	if !c.IsEditing() {
		c.loadChartGroups("")
	}
}

func (c *PanelCharts) SetEdit(editing bool) {
	c.isEditing_ = editing
	c.timeChart.SetEdit(editing)
	c.pUnitsList.SetEnabled(!editing)

	if editing {
		c.splitterEditor.SetRightCollapsed(false)
		c.splitterEditor.SetPositionRelative(0.7)
	} else {
		c.splitterEditor.SetRightCollapsed(true)
	}

	c.updateButtons()
	c.UpdateLayout()
}

func (c *PanelCharts) IsEditing() bool {
	return c.isEditing_
}

func (c *PanelCharts) Save() []byte {
	return c.timeChart.Save()
}

func (c *PanelCharts) loadSelected() {
	selectedItem := c.lvItems.SelectedItem()
	if selectedItem != nil {
		resId := selectedItem.TempData
		c.client.ResGet(resId, func(item *common_interfaces.ResourcesItem, err error) {
			if c.SetCurrentRes(resId, item.Info.Type) {
				c.timeChart.Load(item.Content)
			}
		})
	} else {

		if c.SetCurrentRes("", "") {
			c.timeChart.Load([]byte(""))
		}
	}
}

func (c *PanelCharts) SetCurrentRes(resId string, resType string) bool {
	if c.isEditing_ {
		uicontrols.ShowQuestionMessageOKCancel(c, "Save current chart group?", "Confirmation", func() {
			c.btnSave.Press()
			c.currentResId = resId
			c.currentResType = resType
			c.SetEdit(false)
		}, func() {
			c.currentResId = resId
			c.currentResType = resType
			c.SetEdit(false)
		})
	} else {
		if c.currentResId != resId {
			c.currentResId = resId
			c.currentResType = resType
			c.SetEdit(false)
		}
	}
	return true
}

func (c *PanelCharts) updateButtons() {
	if c.IsEditing() {
		c.btnEdit.SetVisible(false)
		c.btnReject.SetVisible(true)
		c.btnSave.SetVisible(true)
	} else {
		c.btnEdit.SetVisible(true)
		c.btnReject.SetVisible(false)
		c.btnSave.SetVisible(false)
	}
}

func (c *PanelCharts) UpdateStyle() {
	c.Panel.UpdateStyle()

	activeColor := c.AccentColor()
	inactiveColor := c.InactiveColor()

	c.btnAdd.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_add_materialicons_48dp_1x_baseline_add_black_48dp_png, activeColor))
	c.btnRename.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_create_materialicons_48dp_1x_baseline_create_black_48dp_png, activeColor))
	c.btnEdit.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_create_materialicons_48dp_1x_baseline_create_black_48dp_png, activeColor))
	c.btnSave.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_save_alt_materialicons_48dp_1x_baseline_save_alt_black_48dp_png, activeColor))
	c.btnReject.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_communication_cancel_presentation_materialicons_48dp_1x_baseline_cancel_presentation_black_48dp_png, activeColor))
	c.btnRemove.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_clear_materialicons_48dp_1x_baseline_clear_black_48dp_png, activeColor))

	c.btnAdd.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_add_materialicons_48dp_1x_baseline_add_black_48dp_png, inactiveColor))
	c.btnRename.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_create_materialicons_48dp_1x_baseline_create_black_48dp_png, inactiveColor))
	c.btnEdit.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_create_materialicons_48dp_1x_baseline_create_black_48dp_png, inactiveColor))
	c.btnSave.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_save_alt_materialicons_48dp_1x_baseline_save_alt_black_48dp_png, inactiveColor))
	c.btnReject.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_communication_cancel_presentation_materialicons_48dp_1x_baseline_cancel_presentation_black_48dp_png, inactiveColor))
	c.btnRemove.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_clear_materialicons_48dp_1x_baseline_clear_black_48dp_png, inactiveColor))

}

func (c *PanelCharts) GetDataItemValue(path string, control simplemap.IMapControl) {
	val := c.client.GetItemValue(path)
	control.UpdateValue(val)
}

func (c *PanelCharts) loadChartGroups(selectAfterLoadingId string) {
	c.client.ResList("chart_group", "", 0, 1000000, func(infos common_interfaces.ResourcesInfo, err error) {

		if selectAfterLoadingId == "" {
			if c.lvItems.SelectedItem() != nil {
				selectAfterLoadingId = c.lvItems.SelectedItem().TempData
			}
		}

		c.lvItems.RemoveItems()
		indexForSelect := -1
		for i, info := range infos.Items {
			item := c.lvItems.AddItem(info.Name)
			item.TempData = info.Id
			if info.Id == selectAfterLoadingId {
				indexForSelect = i
			}
		}

		if indexForSelect > -1 {
			c.lvItems.SelectItem(indexForSelect)
		}
	})
}
