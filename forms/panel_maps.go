package forms

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/simplemap"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiproperties"
	"github.com/gazercloud/gazerui/uiresources"
	"image"
	"image/png"
)

type PanelMaps struct {
	uicontrols.Panel
	client  *client.Client
	lvItems *uicontrols.ListView
	timer   *uievents.FormTimer

	splitter *uicontrols.SplitContainer

	btnAdd    *uicontrols.Button
	btnRename *uicontrols.Button
	btnRemove *uicontrols.Button

	btnEdit   *uicontrols.Button
	btnReject *uicontrols.Button
	btnSave   *uicontrols.Button

	lblScale *uicontrols.TextBlock

	btnItemUpMax   *uicontrols.Button
	btnItemUp      *uicontrols.Button
	btnItemDown    *uicontrols.Button
	btnItemDownMax *uicontrols.Button

	btnZoom100         *uicontrols.Button
	btnZoomInContainer *uicontrols.Button
	btnZoomIn          *uicontrols.Button
	btnZoomOut         *uicontrols.Button

	//mapWidget      *simplemap.MapDocumentWidget
	panelMapDocument *uicontrols.Panel
	currentResId     string
	currentResType   string

	ToolBox *simplemap.MapToolbox

	mapWidget_ *simplemap.MapWidget

	panelRight       *uicontrols.Panel
	propertiesEditor *simplemap.PropertiesEditor

	OnMouseDrop       func(droppedValue interface{}, control simplemap.IMapControl, x int32, y int32)
	OnDocumentChanged func()
	OnScaleChanged    func(scale float64)
	OnActionOpenMap   func(resId string)
	OnActionWriteItem func(item string, value string)

	saving bool
}

func NewPanelMaps(parent uiinterfaces.Widget, client *client.Client) *PanelMaps {
	var c PanelMaps
	c.client = client
	c.InitControl(parent, &c)
	return &c
}

func (c *PanelMaps) OnInit() {

	pContent := c.AddPanelOnGrid(0, 1)
	pContent.SetPanelPadding(0)
	c.splitter = pContent.AddSplitContainerOnGrid(0, 0)
	c.splitter.SetYExpandable(true)
	c.splitter.SetPosition(250)

	pUnitsList := c.splitter.Panel1.AddPanelOnGrid(0, 0)
	pUnitsList.SetPanelPadding(0)
	txtHeader := pUnitsList.AddTextBlockOnGrid(0, 0, "Maps")
	txtHeader.SetFontSize(24)

	pButtons := pUnitsList.AddPanelOnGrid(0, 1)
	pButtons.SetPanelPadding(0)

	c.lvItems = pUnitsList.AddListViewOnGrid(0, 2)
	c.lvItems.AddColumn("Name", 230)
	c.lvItems.OnSelectionChanged = func() {
		c.loadSelected()
	}

	c.btnAdd = pButtons.AddButtonOnGrid(0, 0, "", func(event *uievents.Event) {
		d := NewFormAddMap(c, c.client, "simple_map")
		d.ShowDialog()
		d.OnAccept = func() {
			c.loadMaps()
		}
	})
	c.btnAdd.SetTooltip("Add map ...")
	c.btnAdd.SetMinWidth(60)

	c.btnRename = pButtons.AddButtonOnGrid(1, 0, "", func(event *uievents.Event) {
		if len(c.lvItems.SelectedItems()) != 1 {
			return
		}
		item := c.lvItems.SelectedItems()[0]
		dialog := NewDialogEditMapName(c, c.client, item.TempData, item.Value(0))
		dialog.ShowDialog()
		dialog.OnAccept = func() {
			c.loadMaps()
		}
	})
	c.btnRename.SetTooltip("Rename chart group")

	c.btnRemove = pButtons.AddButtonOnGrid(2, 0, "", func(event *uievents.Event) {
		if len(c.lvItems.SelectedItems()) != 1 {
			return
		}
		uicontrols.ShowQuestionMessageOKCancel(c, "Remove selected map?", "Confirmation", func() {
			item := c.lvItems.SelectedItems()[0]
			c.client.ResRemove(item.TempData, func(err error) {
				c.loadMaps()
			})
		}, nil)
	})
	c.btnRemove.SetTooltip("Remove selected map")

	pButtons.AddHSpacerOnGrid(3, 0)

	pMapButtons := c.splitter.Panel2.AddPanelOnGrid(0, 0)
	c.btnEdit = pMapButtons.AddButtonOnGrid(0, 0, "", func(event *uievents.Event) {
		if c.IsLoaded() {
			c.splitter.SetLeftCollapsed(true)
			c.UpdateLayout()
			c.SetEdit(true)
			c.ZoomToDefault()
			c.updateButtons()
		}
	})
	c.btnEdit.SetMinWidth(70)
	c.btnEdit.SetTooltip("Switch to Edit")
	c.btnReject = pMapButtons.AddButtonOnGrid(1, 0, "", func(event *uievents.Event) {
		funcReject := func() {
			c.splitter.SetLeftCollapsed(false)
			c.UpdateLayout()
			c.SetEdit(false)
			c.ZoomToDefault()
			c.loadSelected()
			c.updateButtons()
		}

		if c.HasChanges() {
			uicontrols.ShowQuestionMessageYesNoCancel(c, "Save changes?", "Confirmation", func() {
				// Yes
				c.btnSave.Press()
			}, func() {
				// No
				funcReject()
			}, func() {
				//Cancel
			})
		} else {
			funcReject()
		}
	})
	c.btnReject.SetTooltip("Reject changes")
	c.btnReject.SetMinWidth(70)
	c.btnSave = pMapButtons.AddButtonOnGrid(2, 0, "", func(event *uievents.Event) {
		c.SaveMap()
	})
	c.btnSave.SetTooltip("Save")
	c.btnSave.SetMinWidth(70)

	pMapButtons.AddHSpacerOnGrid(9, 0)

	c.btnItemUpMax = pMapButtons.AddButtonOnGrid(15, 0, "", func(event *uievents.Event) {
		if c.IsLoaded() {
			c.MoveItemUpMax()
			c.updateButtons()
		}
	})
	c.btnItemUpMax.SetTooltip("Bring to front")

	c.btnItemUp = pMapButtons.AddButtonOnGrid(16, 0, "", func(event *uievents.Event) {
		if c.IsLoaded() {
			c.MoveItemUp()
			c.updateButtons()
		}
	})
	c.btnItemUp.SetTooltip("Bring forward")

	c.btnItemDown = pMapButtons.AddButtonOnGrid(17, 0, "", func(event *uievents.Event) {
		if c.IsLoaded() {
			c.MoveItemDown()
			c.updateButtons()
		}
	})
	c.btnItemDown.SetTooltip("Send backward")

	c.btnItemDownMax = pMapButtons.AddButtonOnGrid(18, 0, "", func(event *uievents.Event) {
		if c.IsLoaded() {
			c.MoveItemDownMax()
			c.updateButtons()
		}
	})
	c.btnItemDownMax.SetTooltip("Send to back")

	c.lblScale = pMapButtons.AddTextBlockOnGrid(19, 0, "")

	c.btnZoomIn = pMapButtons.AddButtonOnGrid(20, 0, "", func(event *uievents.Event) {
		if c.IsLoaded() {
			c.ZoomIn()
			c.updateButtons()
		}
	})
	c.btnZoomIn.SetTooltip("Zoom In")

	c.btnZoomOut = pMapButtons.AddButtonOnGrid(21, 0, "", func(event *uievents.Event) {
		if c.IsLoaded() {
			c.ZoomOut()
			c.updateButtons()
		}
	})
	c.btnZoomOut.SetTooltip("Zoom Out")

	c.btnZoom100 = pMapButtons.AddButtonOnGrid(22, 0, "", func(event *uievents.Event) {
		if c.IsLoaded() {
			c.Zoom100()
			c.updateButtons()
		}
	})
	c.btnZoom100.SetTooltip("Zoom 100%")

	c.btnZoomInContainer = pMapButtons.AddButtonOnGrid(23, 0, "", func(event *uievents.Event) {
		if c.IsLoaded() {
			c.ZoomInContainer()
			c.updateButtons()
		}
	})
	c.btnZoomInContainer.SetTooltip("Show All")

	c.panelMapDocument = c.splitter.Panel2.AddPanelOnGrid(0, 1)
	//c.splitter.Panel2.AddWidgetOnGrid(c.mapWidget, 0, 1)

	///////////////////

	c.ToolBox = simplemap.NewMapToolbox(c.panelMapDocument)
	c.ToolBox.SetGridX(0)
	c.ToolBox.SetGridY(0)
	c.ToolBox.SetMaxWidth(200)
	c.panelMapDocument.AddWidget(c.ToolBox)

	c.mapWidget_ = simplemap.NewMapWidget(c.panelMapDocument)
	c.mapWidget_.SetGridX(1)
	c.mapWidget_.SetGridY(0)
	c.mapWidget_.OnViewChanged = c.ViewChanged
	c.mapWidget_.SetToolSelector(c.ToolBox)
	c.mapWidget_.OnScaleChanged = func(scale float64) {
		if c.OnScaleChanged != nil {
			c.OnScaleChanged(scale)
		}
	}
	c.mapWidget_.OnActionOpenMap = func(resId string) {
		if c.OnActionOpenMap != nil {
			c.OnActionOpenMap(resId)
		}
	}
	c.mapWidget_.OnActionWriteItem = func(item string, value string) {
		if c.OnActionWriteItem != nil {
			c.OnActionWriteItem(item, value)
		}
	}
	c.panelMapDocument.AddWidget(c.mapWidget_)

	c.SetMapDataSource(c)
	c.mapWidget_.Load("", []byte(""))
	c.mapWidget_.ZoomDefault()

	c.panelRight = c.panelMapDocument.AddPanelOnGrid(2, 0)
	c.panelRight.SetMaxWidth(300)
	c.propertiesEditor = simplemap.NewPropertiesEditor(c.panelRight, c.client)
	c.propertiesEditor.SetGridX(0)
	c.propertiesEditor.SetGridY(0)
	c.panelRight.AddWidget(c.propertiesEditor)

	c.mapWidget_.SetOnSelectionChanged(c.OnSelectionChanged)

	///////////////////

	/*c.mapWidget = simplemap.NewMapDocumentWidget(c, c.client)
	c.mapWidget.SetMapDataSource(c)
	c.mapWidget.SetEdit(false)
	c.mapWidget.Load("", []byte(""))
	c.mapWidget.ZoomToDefault()*/

	c.OnScaleChanged = func(scale float64) {
		c.lblScale.SetText("Scale: " + fmt.Sprint(int(scale*100)) + "%")
	}

	c.OnActionOpenMap = func(resId string) {
		c.openMap(resId)
	}

	c.OnActionWriteItem = func(item string, value string) {
		c.client.Write(item, value, nil)
	}

	c.timer = c.Window().NewTimer(250, func() {
		c.Tick()
	})
	c.timer.StartTimer()

	c.loadMaps()
	c.updateButtons()
	c.SetEdit(false)
}

func (c *PanelMaps) Dispose() {
	if c.timer != nil {
		c.timer.StopTimer()
	}
	c.timer = nil
	c.client = nil

	c.SetEdit(false)
	c.mapWidget_.CloseView()
	c.ToolBox = nil
	c.propertiesEditor = nil
	c.mapWidget_ = nil

	c.Panel.Dispose()
}

func (c *PanelMaps) SaveMap() {
	if c.currentResId != "" {
		c.saving = true
		c.updateButtons()
		imgThumbnail := c.GetThumbnail(48, 48)

		var thumbnailBytes bytes.Buffer
		w := bufio.NewWriter(&thumbnailBytes)
		err := png.Encode(w, imgThumbnail)
		if err == nil {
			_ = w.Flush()
		}

		c.client.ResSet(c.currentResId, thumbnailBytes.Bytes(), c.Save(), func(err error) {
			c.splitter.SetLeftCollapsed(false)
			c.saving = false
			c.UpdateLayout()
			c.SetEdit(false)
			c.ZoomToDefault()
			c.updateButtons()
		})
	}
}

func (c *PanelMaps) openMap(resId string) {
	c.client.ResGet(resId, func(item *common_interfaces.ResourcesItem, err error) {
		if err == nil {
			c.SetEdit(false)
			err = c.Load(resId, item.Content)
			if err != nil {
				c.SetCurrentRes("", "")
				uicontrols.ShowErrorMessage(c, err.Error(), "Error")
				return
			}
			c.ZoomToDefault()
			c.SetCurrentRes(resId, item.Info.Type)
		} else {
			c.SetCurrentRes("", "")
			uicontrols.ShowErrorMessage(c, err.Error(), "Error")
		}
	})
}

func (c *PanelMaps) loadSelected() {
	selectedItem := c.lvItems.SelectedItem()
	if selectedItem != nil {
		resId := selectedItem.TempData
		c.openMap(resId)
	} else {
		c.Load("", nil)
		c.SetCurrentRes("", "")
	}
	c.updateButtons()
}

func (c *PanelMaps) SetCurrentRes(resId string, resType string) {
	c.currentResId = resId
	c.currentResType = resType
}

func (c *PanelMaps) UpdateStyle() {
	c.Panel.UpdateStyle()

	activeColor := c.AccentColor()
	inactiveColor := c.InactiveColor()

	c.btnAdd.SetImage(uiresources.ResImageAdjusted("icons/material/content/drawable-hdpi/ic_add_black_48dp.png", activeColor))
	c.btnRename.SetImage(uiresources.ResImageAdjusted("icons/material/content/drawable-hdpi/ic_create_black_48dp.png", activeColor))
	c.btnRemove.SetImage(uiresources.ResImageAdjusted("icons/material/content/drawable-hdpi/ic_clear_black_48dp.png", activeColor))
	c.btnEdit.SetImage(uiresources.ResImageAdjusted("icons/material/content/drawable-hdpi/ic_create_black_48dp.png", activeColor))
	c.btnSave.SetImage(uiresources.ResImageAdjusted("icons/material/content/drawable-hdpi/outline_save_alt_black_48dp.png", activeColor))
	c.btnReject.SetImage(uiresources.ResImageAdjusted("icons/material/communication/drawable-hdpi/outline_cancel_presentation_black_48dp.png", activeColor))

	c.btnItemUpMax.SetImage(uiresources.ResImageAdjusted("icons/custom/arrow_up_stop.png", activeColor))
	c.btnItemUp.SetImage(uiresources.ResImageAdjusted("icons/custom/arrow_up.png", activeColor))
	c.btnItemDown.SetImage(uiresources.ResImageAdjusted("icons/custom/arrow_down.png", activeColor))
	c.btnItemDownMax.SetImage(uiresources.ResImageAdjusted("icons/custom/arrow_down_stop.png", activeColor))

	c.btnZoomIn.SetImage(uiresources.ResImageAdjusted("icons/material/action/drawable-hdpi/ic_zoom_in_black_48dp.png", activeColor))
	c.btnZoomOut.SetImage(uiresources.ResImageAdjusted("icons/material/action/drawable-hdpi/ic_zoom_out_black_48dp.png", activeColor))
	c.btnZoom100.SetImage(uiresources.ResImageAdjusted("icons/custom/zoom_100.png", activeColor))
	c.btnZoomInContainer.SetImage(uiresources.ResImageAdjusted("icons/material/action/drawable-hdpi/ic_open_with_black_48dp.png", activeColor))

	c.btnAdd.SetImageDisabled(uiresources.ResImageAdjusted("icons/material/content/drawable-hdpi/ic_add_black_48dp.png", inactiveColor))
	c.btnRename.SetImageDisabled(uiresources.ResImageAdjusted("icons/material/content/drawable-hdpi/ic_create_black_48dp.png", inactiveColor))
	c.btnRemove.SetImageDisabled(uiresources.ResImageAdjusted("icons/material/content/drawable-hdpi/ic_clear_black_48dp.png", inactiveColor))
	c.btnEdit.SetImageDisabled(uiresources.ResImageAdjusted("icons/material/content/drawable-hdpi/ic_create_black_48dp.png", inactiveColor))
	c.btnSave.SetImageDisabled(uiresources.ResImageAdjusted("icons/material/content/drawable-hdpi/outline_save_alt_black_48dp.png", inactiveColor))
	c.btnReject.SetImageDisabled(uiresources.ResImageAdjusted("icons/material/communication/drawable-hdpi/outline_cancel_presentation_black_48dp.png", inactiveColor))

	c.btnItemUpMax.SetImageDisabled(uiresources.ResImageAdjusted("icons/custom/arrow_up_stop.png", inactiveColor))
	c.btnItemUp.SetImageDisabled(uiresources.ResImageAdjusted("icons/custom/arrow_up.png", inactiveColor))
	c.btnItemDown.SetImageDisabled(uiresources.ResImageAdjusted("icons/custom/arrow_down.png", inactiveColor))
	c.btnItemDownMax.SetImageDisabled(uiresources.ResImageAdjusted("icons/custom/arrow_down_stop.png", inactiveColor))

	c.btnZoomIn.SetImageDisabled(uiresources.ResImageAdjusted("icons/material/action/drawable-hdpi/ic_zoom_in_black_48dp.png", inactiveColor))
	c.btnZoomOut.SetImageDisabled(uiresources.ResImageAdjusted("icons/material/action/drawable-hdpi/ic_zoom_out_black_48dp.png", inactiveColor))
	c.btnZoom100.SetImageDisabled(uiresources.ResImageAdjusted("icons/custom/zoom_100.png", inactiveColor))
	c.btnZoomInContainer.SetImageDisabled(uiresources.ResImageAdjusted("icons/material/action/drawable-hdpi/ic_open_with_black_48dp.png", inactiveColor))
}

func (c *PanelMaps) GetDataItemValue(path string, control simplemap.IMapControl) {
	val := c.client.GetItemValue(path)
	control.UpdateValue(val)
}

func (c *PanelMaps) updateButtons() {
	if c.saving {
		c.btnEdit.SetVisible(false)
		c.btnReject.SetVisible(false)
		c.btnSave.SetVisible(false)
		return
	}

	if c.IsEditing() {
		c.btnEdit.SetVisible(false)
		c.btnReject.SetVisible(true)
		c.btnSave.SetVisible(true)

		c.btnItemUpMax.SetVisible(true)
		c.btnItemUp.SetVisible(true)
		c.btnItemDown.SetVisible(true)
		c.btnItemDownMax.SetVisible(true)
	} else {
		c.btnEdit.SetVisible(true)
		c.btnReject.SetVisible(false)
		c.btnSave.SetVisible(false)

		c.btnItemUpMax.SetVisible(false)
		c.btnItemUp.SetVisible(false)
		c.btnItemDown.SetVisible(false)
		c.btnItemDownMax.SetVisible(false)
	}
}

func (c *PanelMaps) LoadContent(itemUrl string, control simplemap.IMapControl) {
	c.client.ResGet(itemUrl, func(item *common_interfaces.ResourcesItem, err error) {
		if err == nil {
			control.LoadContent(item.Content, err)
		} else {
			control.LoadContent(nil, err)
		}
	})
}

func (c *PanelMaps) GetWidgets(filter string, offset int, maxCount int, toolbox simplemap.IMapToolbox) {
	c.client.ResList("simple_map", filter, offset, maxCount, func(infos common_interfaces.ResourcesInfo, err error) {
		toolbox.SetItems(infos)
	})
}

func (c *PanelMaps) loadMaps() {
	c.client.ResList("simple_map", "", 0, 1000000, func(infos common_interfaces.ResourcesInfo, err error) {
		c.lvItems.RemoveItems()
		for _, info := range infos.Items {
			item := c.lvItems.AddItem(info.Name)
			item.TempData = info.Id
		}
	})
}

func (c *PanelMaps) IsLoaded() bool {
	if c.mapWidget_ == nil {
		return false
	}
	if c.mapWidget_.View() == nil {
		return false
	}
	return true
}

func (c *PanelMaps) SetEdit(edit bool) {
	c.mapWidget_.SetEdit(edit)

	if edit {
		c.ToolBox.SetVisible(true)
		c.panelRight.SetVisible(true)
		c.OnSelectionChanged()
		c.ToolBox.Reset()
	} else {
		c.ToolBox.SetVisible(false)
		c.panelRight.SetVisible(false)
	}
}

func (c *PanelMaps) IsEditing() bool {
	return c.mapWidget_.IsEditing()
}

func (c *PanelMaps) OnSelectionChanged() {

	if len(c.mapWidget_.SelectedItems()) == 0 {
		c.propertiesEditor.SetPropertiesContainer(c.mapWidget_.View())
	}

	if len(c.mapWidget_.SelectedItems()) == 1 {
		c.propertiesEditor.SetPropertiesContainer(c.mapWidget_.SelectedItems()[0].(uiproperties.IPropertiesContainer))
	}

	if len(c.mapWidget_.SelectedItems()) > 1 {
		c.propertiesEditor.SetPropertiesContainer(nil)
	}
}

func (c *PanelMaps) Tick() {
	c.mapWidget_.Tick()
}

func (c *PanelMaps) ViewChanged() {
	if c.OnDocumentChanged != nil {
		c.OnDocumentChanged()
	}
}

func (c *PanelMaps) HasChanges() bool {
	return c.mapWidget_.HasChanges()
}

func (c *PanelMaps) SetOnMouseDrop(OnMouseDrop func(droppedValue interface{}, control simplemap.IMapControl, x int32, y int32)) {
	c.OnMouseDrop = OnMouseDrop
	c.mapWidget_.OnMouseDrop = OnMouseDrop
}

func (c *PanelMaps) SetMapDataSource(mapDataSource simplemap.IMapDataSource) {
	c.ToolBox.SetMapDataSource(mapDataSource)
	c.mapWidget_.SetMapDataSource(mapDataSource)
}

func (c *PanelMaps) AddControl(control simplemap.IMapControl) {
	c.mapWidget_.AddControl(control)
}

func (c *PanelMaps) Save() []byte {
	return c.mapWidget_.Save()
}

func (c *PanelMaps) GetThumbnail(width, height int) image.Image {
	return c.mapWidget_.GetThumbnail(width, height)
}

func (c *PanelMaps) Load(resId string, value []byte) error {
	err := c.mapWidget_.Load(resId, value)
	if err != nil {
		return err
	}
	c.ZoomToDefault()
	return nil
}

func (c *PanelMaps) MoveItemUpMax() {
	c.UpdateLayout()
	c.mapWidget_.MoveItemUpMax()
}

func (c *PanelMaps) MoveItemUp() {
	c.UpdateLayout()
	c.mapWidget_.MoveItemUp()
}

func (c *PanelMaps) MoveItemDown() {
	c.UpdateLayout()
	c.mapWidget_.MoveItemDown()
}

func (c *PanelMaps) MoveItemDownMax() {
	c.UpdateLayout()
	c.mapWidget_.MoveItemDownMax()
}

func (c *PanelMaps) ZoomToDefault() {
	c.UpdateLayout()
	c.mapWidget_.ZoomDefault()
}

func (c *PanelMaps) ZoomIn() {
	c.UpdateLayout()
	c.mapWidget_.ZoomIn()
}

func (c *PanelMaps) ZoomOut() {
	c.UpdateLayout()
	c.mapWidget_.ZoomOut()
}

func (c *PanelMaps) Zoom100() {
	c.UpdateLayout()
	c.mapWidget_.Zoom100()
}

func (c *PanelMaps) ZoomInContainer() {
	c.UpdateLayout()
	c.mapWidget_.ZoomInContainer()
}
