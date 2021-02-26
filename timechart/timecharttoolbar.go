package timechart

import (
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uistyles"
)

type ToolBar struct {
	uicontrols.Panel
	chart *TimeChart

	btnUnitedVerticalScale *uicontrols.Button
	unitedVerticalScale    bool

	btnShowAreaHeader *uicontrols.Button
	showAreaHeader    bool

	OnChanged func()
}

func NewToolBar(parent uiinterfaces.Widget, chart *TimeChart) *ToolBar {
	var c ToolBar
	c.chart = chart
	c.InitControl(parent, &c)
	c.btnUnitedVerticalScale = c.AddButtonOnGrid(0, 0, "United Scale", func(event *uievents.Event) {
		c.unitedVerticalScale = !c.unitedVerticalScale
		c.updateButtons()
		c.chart.SetUnitedVerticalScale(c.unitedVerticalScale)
		if c.OnChanged != nil {
			c.OnChanged()
		}
	})
	c.btnUnitedVerticalScale.SetTooltip("United Vertical Scale")
	c.btnUnitedVerticalScale.SetMinWidth(50)
	c.btnShowAreaHeader = c.AddButtonOnGrid(1, 0, "Show qualities", func(event *uievents.Event) {
		c.showAreaHeader = !c.showAreaHeader
		c.updateButtons()
		c.chart.SetShowQualities(c.showAreaHeader)
		if c.OnChanged != nil {
			c.OnChanged()
		}
	})
	c.btnShowAreaHeader.SetTooltip("Show Qualities Frame")
	c.btnShowAreaHeader.SetMinWidth(50)
	c.AddHSpacerOnGrid(2, 0)
	return &c
}

func (c *ToolBar) Dispose() {
	c.chart = nil
	c.btnShowAreaHeader = nil
	c.btnUnitedVerticalScale = nil
	c.Panel.Dispose()
}

func (c *ToolBar) updateButtons() {
	if c.unitedVerticalScale {
		c.btnUnitedVerticalScale.SetForeColor(uistyles.DefaultBackColor)
		c.btnUnitedVerticalScale.SetBackColor(c.ForeColor())
	} else {
		c.btnUnitedVerticalScale.SetForeColor(c.ForeColor())
		c.btnUnitedVerticalScale.SetBackColor(uistyles.DefaultBackColor)
	}
	if c.showAreaHeader {
		c.btnShowAreaHeader.SetForeColor(uistyles.DefaultBackColor)
		c.btnShowAreaHeader.SetBackColor(c.ForeColor())
	} else {
		c.btnShowAreaHeader.SetForeColor(c.ForeColor())
		c.btnShowAreaHeader.SetBackColor(uistyles.DefaultBackColor)
	}
}

func (c *ToolBar) UpdateStyle() {
	c.Panel.UpdateStyle()
	c.updateButtons()
}
