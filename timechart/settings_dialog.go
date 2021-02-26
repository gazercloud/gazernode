package timechart

import (
	"fmt"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"image/color"
)

type SettingsDialog struct {
	uicontrols.Dialog
	timeChart        *TimeChart
	lvAreas          *uicontrols.ListView
	lvSeries         *uicontrols.ListView
	panelAreaProps   *uicontrols.Panel
	chkShowQualities *uicontrols.CheckBox
	chkUnitedScale   *uicontrols.CheckBox
	panelSeriesProps *uicontrols.Panel
	cpSeriesColor    *uicontrols.ColorPicker
	txtUnitName      *uicontrols.TextBox
	btnOK            *uicontrols.Button

	loading            bool
	currentAreaIndex   int
	currentSeriesIndex int
}

func NewSettingsDialog(parent *TimeChart) *SettingsDialog {
	var c SettingsDialog
	c.timeChart = parent
	c.InitControl(parent, &c)

	c.currentAreaIndex = -1
	c.currentSeriesIndex = -1

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)

	c.lvAreas = pContent.AddListViewOnGrid(0, 0)
	c.lvAreas.AddColumn("Name", 150)
	c.lvAreas.OnSelectionChanged = func() {
		if c.loading {
			return
		}
		c.currentAreaIndex = c.lvAreas.SelectedItemIndex()
		c.load()
	}

	c.lvAreas.OnItemClicked = func(item *uicontrols.ListViewItem) {
		if c.loading {
			return
		}
		c.currentSeriesIndex = -1
		c.load()
	}

	c.panelAreaProps = pContent.AddPanelOnGrid(1, 0)

	c.lvSeries = c.panelAreaProps.AddListViewOnGrid(0, 0)
	c.lvSeries.AddColumn("Name", 150)
	c.lvSeries.OnSelectionChanged = func() {
		if c.loading {
			return
		}
		c.currentSeriesIndex = c.lvSeries.SelectedItemIndex()
		c.load()
	}

	pArea := c.panelAreaProps.AddPanelOnGrid(0, 1)
	c.chkShowQualities = pArea.AddCheckBoxOnGrid(0, 0, "Show qualities")

	c.chkShowQualities.OnCheckedChanged = func(checkBox *uicontrols.CheckBox, checked bool) {
		if c.currentAreaIndex > -1 {
			area := c.lvAreas.SelectedItem().UserData("area").(*Area)
			area.SetShowQualities(c.chkShowQualities.IsChecked())
		}
	}

	c.chkUnitedScale = pArea.AddCheckBoxOnGrid(0, 1, "United scale")
	c.chkUnitedScale.OnCheckedChanged = func(checkBox *uicontrols.CheckBox, checked bool) {
		if c.currentAreaIndex > -1 {
			area := c.lvAreas.SelectedItem().UserData("area").(*Area)
			area.SetUnitedScale(c.chkUnitedScale.IsChecked())
		}
	}

	c.panelSeriesProps = pContent.AddPanelOnGrid(2, 0)
	c.panelSeriesProps.AddTextBlockOnGrid(0, 0, "Color:")
	c.cpSeriesColor = c.panelSeriesProps.AddColorPickerOnGrid(0, 1)
	c.cpSeriesColor.SetMinHeight(50)
	c.panelSeriesProps.AddVSpacerOnGrid(0, 5)

	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)
	pButtons.AddHSpacerOnGrid(0, 0)
	c.btnOK = pButtons.AddButtonOnGrid(1, 0, "OK", nil)
	c.TryAccept = func() bool {
		c.btnOK.SetEnabled(false)
		c.TryAccept = nil
		c.Accept()
		return false
	}
	c.btnOK.SetMinWidth(70)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Cancel", func(event *uievents.Event) {
		c.Reject()
	})
	btnCancel.SetMinWidth(70)

	c.SetAcceptButton(c.btnOK)
	c.SetRejectButton(btnCancel)

	c.load()

	c.cpSeriesColor.OnColorChanged = func(colorPicker *uicontrols.ColorPicker, color color.Color) {
		if c.currentSeriesIndex > -1 {
			series := c.lvSeries.SelectedItem().UserData("series").(*Series)
			series.SetColor(color)
		}
	}

	return &c
}

func (c *SettingsDialog) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("Chart settings")
	c.Resize(600, 400)
}

func (c *SettingsDialog) load() {
	c.loading = true
	c.lvSeries.RemoveItems()
	c.lvAreas.RemoveItems()
	for i, a := range c.timeChart.areas {
		item := c.lvAreas.AddItem("Area #" + fmt.Sprint(i))
		item.SetUserData("area", a)
	}
	if c.currentAreaIndex > -1 {
		c.panelAreaProps.SetEnabled(true)
		c.lvAreas.SelectItem(c.currentAreaIndex)
		area := c.lvAreas.SelectedItem().UserData("area").(*Area)

		c.chkShowQualities.SetChecked(area.ShowQualities())
		c.chkUnitedScale.SetChecked(area.UnitedScale())

		for _, s := range area.series {
			item := c.lvSeries.AddItem(s.id)
			item.SetUserData("series", s)
		}

		if c.currentSeriesIndex >= c.lvSeries.ItemsCount() {
			c.currentSeriesIndex = -1
		}

		if c.currentSeriesIndex > -1 {
			c.lvSeries.SelectItem(c.currentSeriesIndex)
		}
	} else {
		c.panelAreaProps.SetEnabled(false)
	}

	c.loadSeriesProps()
	c.loading = false
}

func (c *SettingsDialog) loadSeriesProps() {
	if c.lvSeries.SelectedItem() == nil {
		c.currentSeriesIndex = -1
	}

	if c.currentSeriesIndex > -1 {
		series := c.lvSeries.SelectedItem().UserData("series").(*Series)
		c.panelSeriesProps.SetEnabled(true)
		c.cpSeriesColor.SetColor(series.color)
	} else {
		c.cpSeriesColor.SetColor(color.Transparent)
		c.panelSeriesProps.SetEnabled(false)
	}
}
