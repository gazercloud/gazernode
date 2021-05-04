package timechart

import (
	"fmt"
	"github.com/gazercloud/gazerui/canvas"
	"github.com/gazercloud/gazerui/ui"
	"golang.org/x/image/colornames"
	"image/color"
	"math"
)

type Series struct {
	verticalScale *VerticalScale
	hScale        *HorizontalScale
	values        []*Value
	vlines        []vline
	color         color.Color

	hoverTime int64
	hoverX    int
	hoverY    int

	// Debug
	drawBegin int64

	id           string
	name         string
	bottomHeader *AreaTopHeader
	dataProvider IDataProvider
	area         *Area
}

func NewSeries(id string, area *Area) *Series {
	var c Series
	c.verticalScale = NewVerticalScale(&c)
	c.area = area
	c.color = colornames.Blue
	c.bottomHeader = NewAreaTopHeader(7)
	c.id = id
	return &c
}

func (c *Series) Dispose() {
	c.values = nil
	c.vlines = nil
	c.hScale = nil
	c.verticalScale = nil
	c.dataProvider = nil
}

func (c *Series) SetDataProvider(provider IDataProvider) {
	c.dataProvider = provider
}

func (c *Series) SetColor(color color.Color) {
	c.color = color
}

func (c *Series) Color() color.Color {
	return c.color
}

func (c *Series) SetName(name string) {
	c.name = name
}

func (c *Series) Id() string {
	return c.id
}

func (c *Series) Draw(ctx ui.DrawContext, scaleXOffset int, xOffset int, height int, hScale *HorizontalScale, bottomHeaderHeight int, index int) {

	if c.dataProvider == nil {
		return
	}

	if hScale.width < 1 {
		return
	}

	w := float64(hScale.width)
	r := float64(hScale.displayMax_ - hScale.displayMin_)
	timePerPixel := int64(r / w)

	c.values = c.dataProvider.GetData("", hScale.displayMin_, hScale.displayMax_, timePerPixel)
	c.verticalScale.Height = height - bottomHeaderHeight

	if !c.area.unitedVerticalScale || index == 0 {
		c.verticalScale.Draw(ctx, scaleXOffset, 0, c.color)
	}

	c.hScale = hScale

	min, max, err := c.calcVerticalLimits()
	if c.area.unitedVerticalScale {
		min, max, err = c.area.calcVerticalLimits()
	}

	if err == nil {
		var valueMargin float64
		diff := max - min
		if diff > 0 {
			valueMargin = diff / 10
		} else {
			valueMargin = 1 // For one point series data
		}

		c.verticalScale.SetDisplayRange(min-valueMargin, max+valueMargin)
	}

	ctx.Save()
	ctx.Translate(0, 0)
	c.vlines = make([]vline, hScale.width)

	var vlineFirst vline
	var vlineLast vline

	firstPoint := -1
	lastPoint := -1

	displayedMinY_ := float64(100000000000000)
	displayedMaxY_ := -displayedMinY_

	//lastPointX := 0
	//lastPointY := 0
	hasLastPoint := false
	countPointsWithValue := 0

	for x := 0; x < len(c.vlines); x++ {
		c.vlines[x].X = x
	}

	for index, value := range c.values {
		valueX := value.DatetimeFirst
		x := hScale.valueToPixel(valueX)

		if x >= 0 && x < len(c.vlines) {
			if firstPoint == -1 {
				firstPoint = index
			}
			lastPoint = index

			c.fillVLine(&c.vlines[x], value, x)

			// MinMax Vertical scale

			if value.MinValue < displayedMinY_ {
				displayedMinY_ = value.MinValue
			}
			if value.MaxValue > displayedMaxY_ {
				displayedMaxY_ = value.MaxValue
			}
		}

		//ctx.FillRect(x + xOffset - 2, c.verticalScale.getPointOnY(value.value) - 2, 5, 5, colornames.Blue)
		if hasLastPoint {
			//ctx.DrawLine(lastPointX, lastPointY, x+xOffset, c.verticalScale.getPointOnY(value.value), 1, colornames.Red)
		}

		hasLastPoint = true
		//lastPointX = x + xOffset
		//lastPointY = c.verticalScale.getPointOnY(value.value)
	}

	// Prepare first & last points
	if firstPoint > 0 {
		c.fillVLine(&vlineFirst, c.values[firstPoint-1], hScale.valueToPixel(c.values[firstPoint-1].DatetimeFirst))
	}

	if lastPoint < len(c.values)-1 {
		c.fillVLine(&vlineLast, c.values[lastPoint+1], hScale.valueToPixel(c.values[lastPoint+1].DatetimeFirst))
	}

	var lastPointX int
	var lastPointY int
	var lastPointYMin int
	var lastPointYMax int
	var previousHasEnd bool

	if firstPoint > 0 {
		c.drawVLine(ctx, &vlineFirst, &lastPointX, &lastPointY, &lastPointYMin, &lastPointYMax, &previousHasEnd, xOffset)
	}

	var firstVLine vline

	for i := 0; i < len(c.vlines); i++ {
		c.drawVLine(ctx, &c.vlines[i], &lastPointX, &lastPointY, &lastPointYMin, &lastPointYMax, &previousHasEnd, xOffset)
		if c.vlines[i].hasValues {
			countPointsWithValue++
			if countPointsWithValue == 1 {
				firstVLine = c.vlines[i]
			}
		}
	}

	if lastPoint < len(c.values)-1 {
		c.drawVLine(ctx, &vlineLast, &lastPointX, &lastPointY, &lastPointYMin, &lastPointYMax, &previousHasEnd, xOffset)
	}

	if countPointsWithValue == 1 {
		ctx.SetColor(c.color)
		ctx.FillRect(firstVLine.X+xOffset-1, firstVLine.minYp-1, 3, 3)
	}

	//ctx.SetColor(colornames.Red)
	//ctx.SetFontSize(14)
	//ctx.SetFontFamily("")

	textHeight := 20
	ctx.SetColor(c.color)
	ctx.SetFontSize(14)
	ctx.SetTextAlign(canvas.HAlignLeft, canvas.VAlignTop)
	ctx.DrawText(xOffset+10, textHeight*index+10, 300, 50, c.id)

	if len(c.area.timeChart.selections_) > 0 {
		// Selection statistics
		statMinValue := math.MaxFloat64
		statMaxValue := -math.MaxFloat64
		statAvgValue := float64(0)
		statAvgCount := 0

		for i := 0; i < len(c.area.timeChart.selections_); i++ {
			pSelection := c.area.timeChart.selections_[i]
			for _, aa := range c.values {
				if aa.DatetimeFirst >= pSelection.minX && aa.DatetimeLast <= pSelection.maxX {
					if aa.MaxValue > statMaxValue {
						statMaxValue = aa.MaxValue
					}
					if aa.MinValue < statMinValue {
						statMinValue = aa.MinValue
					}
					statAvgValue += aa.AvgValue
					statAvgCount++
				}
			}
		}

		statAvgValue = statAvgValue / float64(statAvgCount)

		statYOffset := 0

		statYOffset += 20
		ctx.DrawText(xOffset+10, textHeight*index+10+statYOffset, 300, 50, "Min: "+fmt.Sprint(statMinValue))
		statYOffset += 20
		ctx.DrawText(xOffset+10, textHeight*index+10+statYOffset, 300, 50, "Max: "+fmt.Sprint(statMaxValue))
		statYOffset += 20
		ctx.DrawText(xOffset+10, textHeight*index+10+statYOffset, 300, 50, "Avg: "+fmt.Sprint(statAvgValue))
		statYOffset += 20
	}

	loadingDiapasons := c.dataProvider.GetLoadingDiapasons()
	for _, d := range loadingDiapasons {
		ctx.SetColor(color.RGBA{
			R: 200,
			G: 200,
			B: 0,
			A: 50,
		})
		x1 := hScale.valueToPixel(d.MinTime) + c.verticalScaleWidth()
		x2 := hScale.valueToPixel(d.MaxTime) + c.verticalScaleWidth()
		//ctx.FillRect(x1, 0, x2-x1, 5)
		ii := 0
		for x := x1; x < x2; x += 10 {
			ii++
			if (ii % 2) == 0 {
				ctx.SetColor(color.RGBA{
					R: 200,
					G: 200,
					B: 0,
					A: 50,
				})
			} else {
				ctx.SetColor(color.RGBA{
					R: 100,
					G: 0,
					B: 100,
					A: 50,
				})
			}

			ctx.FillRect(x, 0, x+10, 5)
		}
	}

	ctx.Load()
}

func (c *Series) DrawBottomHeader(ctx ui.DrawContext, xOffset int, yOffset int, seriesIndex int, namesLineHeight int) {
	if c.vlines == nil {
		return
	}
	ctx.Save()
	ctx.Translate(xOffset, yOffset+seriesIndex*c.bottomHeader.Height()+namesLineHeight)
	c.bottomHeader.draw(ctx, c.vlines, c.area.timeChart)
	ctx.Load()
}

func (c *Series) drawVLine(ctx ui.DrawContext, vline *vline, lastPointX *int, lastPointYMin *int, lastPointYMax *int, lastPointY *int, previousHasEnd *bool, xOffset int) {
	//return
	if !vline.hasValues {
		return
	}

	lineWidth := 1

	if vline.hasY {
		vline.minYp = c.verticalScale.getPointOnY(vline.minY)
		vline.maxYp = c.verticalScale.getPointOnY(vline.maxY)

		// from perious to begin point
		if *previousHasEnd && vline.hasBegin {
			firstYp := c.verticalScale.getPointOnY(vline.firstY)
			needToDraw := true
			if vline.X-*lastPointX < 2 {
				if (vline.minYp < *lastPointYMin && vline.minYp > *lastPointYMax) || (vline.maxYp < *lastPointYMin && vline.maxYp > *lastPointYMax) {
					needToDraw = false
				}
			}

			if needToDraw {
				ctx.SetColor(c.color)
				ctx.SetStrokeWidth(lineWidth)
				ctx.DrawLine(*lastPointX+xOffset, *lastPointY, vline.X+xOffset, firstYp)
			}
		}
		// min/max line on this X
		if vline.minY != vline.maxY {
			ctx.SetColor(c.color)
			ctx.SetStrokeWidth(lineWidth)
			ctx.DrawLine(vline.X+xOffset, vline.minYp, vline.X+xOffset, vline.maxYp)
		} else {
			//ctx.FillRect(vline.X+xOffset - 1, vline.minYp - 1, 3, 3, c.color)
			//ctx.MixPixel(vline.X+xOffset, vline.minYp, c.color)
		}
	}

	if vline.hasEnd {
		*previousHasEnd = true
		*lastPointX = vline.X
	} else {
		*previousHasEnd = false
	}

	*lastPointY = c.verticalScale.getPointOnY(vline.lastY)
	*lastPointYMin = vline.minYp
	*lastPointYMax = vline.maxYp
}

func (c *Series) AddValue(v Value) {
	//c.values = append(c.values, &v)
}

func (c *Series) Clear() {
	//c.values = make([]*Value, 0)
}

func (c *Series) verticalScaleWidth() int {
	return c.verticalScale.Width
}

func (c *Series) RemoveItemsByTime(timeFrom, timeTo int64) {
	/*indexFrom := 0
	indexTo := 0
	for index, value := range c.values {
		if value.DatetimeFirst >= timeFrom {
			indexFrom = index
			break
		}
	}
	for index, value := range c.values {
		if value.DatetimeFirst > timeTo {
			indexTo = index
			break
		}
	}

	c.values = append(c.values[:indexFrom], c.values[indexTo:]...)*/
}

func (c *Series) fillVLine(vline *vline, value *Value, x int) {
	vline.X = x

	good := !value.hasBadQuality() && value != nil

	if good {

		//valueY := value.AvgValue

		//y := valueY

		if !vline.hasY {
			vline.hasY = true
			if !vline.hasBadQuality {
				vline.hasBegin = true
			}
			vline.hasEnd = true
			vline.hasValues = true
			vline.firstY = value.FirstValue
			vline.lastY = value.LastValue
			vline.minY = value.MinValue
			vline.maxY = value.MaxValue
			vline.minYValue = value.MinValue
			vline.maxYValue = value.MaxValue
		} else {
			vline.hasEnd = true
			vline.hasValues = true
			vline.lastY = value.LastValue
			if value.MinValue < vline.minY {
				vline.minY = value.MinValue
			}
			if value.MaxValue > vline.maxY {
				vline.maxY = value.MaxValue
			}
		}
	} else {
		vline.hasBadQuality = true
		vline.hasEnd = false
		if !vline.hasValues {
			vline.hasBegin = false

		}
		vline.hasValues = true
	}
}

func (c *Series) calcHorizontalLimits() (int64, int64, error) {
	min := int64(100000000000000000)
	max := int64(-100000000000000000)

	if len(c.values) < 1 {
		return 0, 0, fmt.Errorf("Series is empty")
	}

	for _, v := range c.values {
		if v.DatetimeFirst < min {
			min = v.DatetimeFirst
		}
		if v.DatetimeFirst > max {
			max = v.DatetimeFirst
		}
	}

	return min, max, nil
}

func (c *Series) calcVerticalLimits() (float64, float64, error) {
	min := float64(100000000000000000)
	max := float64(-100000000000000000)

	if len(c.values) < 1 {
		return 0, 0, fmt.Errorf("Series is empty")
	}

	for _, v := range c.values {
		if v.hasGoodQuality() {
			if v.MinValue < min {
				min = v.MinValue
			}
			if v.MaxValue > max {
				max = v.MaxValue
			}
		}
	}

	return min, max, nil
}

type vline struct {
	X         int
	hasValues bool
	minYValue float64
	maxYValue float64
	minY      float64
	maxY      float64
	firstY    float64
	lastY     float64

	minYp int
	maxYp int

	hasY          bool
	hasBegin      bool
	hasEnd        bool
	hasBadQuality bool
}

func (c *Series) setHoverTimePoint(time int64) {
	c.hoverTime = time
}

func (c *Series) setHoverPoint(x, y int) {
	c.hoverX = x
	c.hoverY = y
}
