package timechart

import (
	"errors"
	"fmt"
	"github.com/gazercloud/gazerui/ui"
	"image/color"
)

/*const (
	AREA_HIGHT = 300
)
*/

type Area struct {
	timeChart           *TimeChart
	series              []*Series
	height              int
	hScale              *HorizontalScale
	unitedVerticalScale bool

	hoverX int
	hoverY int

	highlighting  bool
	dataProvider  IDataProvider
	showQualities bool
}

func NewArea(timeChart *TimeChart) *Area {
	var c Area
	c.timeChart = timeChart
	c.series = make([]*Series, 0)
	c.height = 0

	return &c
}

func (c *Area) SetHeight(h int) {
	c.height = h
}

func (c *Area) TimeChart() *TimeChart {
	return c.timeChart
}

func (c *Area) RemoveSeriesByIndex(index int) {
	if index < 0 || index > len(c.series) {
		return
	}
	c.series[index].Dispose()
	c.series = append(c.series[:index], c.series[index+1:]...)
	c.timeChart.Update("TimeChart")
}

func (c *Area) Dispose() {
	for _, s := range c.series {
		s.Dispose()
	}
	c.series = nil

	c.hScale = nil
	c.dataProvider = nil
	c.timeChart = nil
}

func (c *Area) SetShowQualities(showQualities bool) {
	c.showQualities = showQualities
	c.timeChart.Update("TimeChartArea")
}

func (c *Area) ShowQualities() bool {
	return c.showQualities
}

func (c *Area) SetUnitedScale(unitedVerticalScale bool) {
	c.unitedVerticalScale = unitedVerticalScale
	c.timeChart.Update("TimeChartArea")
}

func (c *Area) UnitedScale() bool {
	return c.unitedVerticalScale
}

func (c *Area) SetHighlighting(highlighting bool) {
	c.highlighting = highlighting
}

func (c *Area) Series() []*Series {
	return c.series
}

func (c *Area) Draw(ctx ui.DrawContext, leftScaleWidth int, hScale *HorizontalScale, width int, foreColor color.Color, yOffset int) {

	ctx.Save()
	ctx.Translate(0, yOffset)

	c.hScale = hScale
	namesLineHeight := 0
	ctx.SetStrokeWidth(1)
	ctx.SetColor(foreColor)
	//ctx.DrawLine(0, 0, 100, 0)

	///////////////////////////////////////////////////////////////////////////////////
	// Hover line
	if c.timeChart.MouseIsInside() && c.timeChart.lastMouseY > yOffset && c.timeChart.lastMouseY < (yOffset+c.height) {
		ctx.SetStrokeWidth(1)
		ctx.SetColor(c.timeChart.InactiveColor())
		ctx.DrawLine(0, c.hoverY, width, c.hoverY)
	}
	///////////////////////////////////////////////////////////////////////////////////

	scaleXOffset := 0

	for index, ser := range c.series {
		ser.Draw(ctx, scaleXOffset, leftScaleWidth, c.height, c.hScale, c.bottomHeaderHeight(), index)
		scaleXOffset += ser.verticalScale.Width
	}

	if c.showQualities {
		for index, ser := range c.series {
			ser.DrawBottomHeader(ctx, leftScaleWidth, c.height-c.bottomHeaderHeight(), index, namesLineHeight)
		}
	}

	/*namesLineXOffset := leftScaleWidth
	for index, ser := range c.series {
		ctx.SetColor(ser.color)
		ctx.SetFontSize(12)
		ctx.DrawText(namesLineXOffset, 0, 100, 20, fmt.Sprint(index)+" - "+ser.name)
		namesLineXOffset += 200
		//ctx.DrawLine(namesLineXOffset, namesLineHeight / 2, namesLineXOffset + 10, namesLineHeight / 2, 1, ser.color)
	}*/

	//ctx.SetStrokeWidth(1)
	//ctx.SetColor(colornames.Darkgray)
	//ctx.DrawLine(0, c.topHeaderHeight()+namesLineHeight, width, c.topHeaderHeight()+namesLineHeight)
	ctx.SetStrokeWidth(1)
	ctx.SetColor(c.timeChart.borderColor())
	ctx.DrawLine(0, c.height, width, c.height)
	if c.highlighting {
		ctx.SetColor(color.RGBA{255, 0, 0, 50})
		ctx.FillRect(0, 0, width, c.height)
	}

	ctx.Load()
}

func (c *Area) calcVerticalLimits() (float64, float64, error) {
	min := float64(100000000000000000)
	max := float64(-100000000000000000)

	found := false

	for _, v := range c.series {
		sMin, sMax, sErr := v.calcVerticalLimits()
		if sErr == nil {
			if sMin < min {
				min = sMin
			}
			if sMax > max {
				max = sMax
			}
			found = true
		}
	}

	if !found {
		return 0, 1, errors.New("no data")
	}

	return min, max, nil

}

func (c *Area) calcLeftScalesWidth() int {
	leftScalesWidth := 0
	for index, ser := range c.series {
		if !c.unitedVerticalScale || index == 0 {
			leftScalesWidth += ser.verticalScaleWidth()
		}
	}
	return leftScalesWidth
}

func (c *Area) calcHorizontalLimits() (int64, int64, error) {
	min := int64(100000000000000000)
	max := int64(-100000000000000000)

	if len(c.series) < 1 {
		return 0, 0, fmt.Errorf("Area is empty")
	}

	for _, ser := range c.series {
		sMin, sMax, err := ser.calcHorizontalLimits()
		if err == nil {
			if sMin < min {
				min = sMin
			}
			if sMax > max {
				max = sMax
			}
		}
	}

	return min, max, nil
}

func (c *Area) MouseDown(x, y int) bool {
	return false
}

func (c *Area) MouseUp(x, y int) bool {
	return false
}

func (c *Area) MouseMove(x, y int) bool {
	return false
}

func (c *Area) setHoverTimePoint(time int64) {
	for _, ser := range c.series {
		ser.setHoverTimePoint(time)
	}
}

func (c *Area) setHoverPoint(x, y int) {
	c.hoverX = x
	c.hoverY = y

	for _, ser := range c.series {
		ser.setHoverPoint(x, y)
	}
}

func (c *Area) bottomHeaderHeight() int {
	result := 0
	if c.showQualities {
		for _, s := range c.series {
			result += s.bottomHeader.Height()
		}
	}
	return result
}
