package timechart

import (
	"fmt"
	"github.com/gazercloud/gazerui/canvas"
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiresources"
	"github.com/nfnt/resize"
	"golang.org/x/image/colornames"
	"time"
)

type HorizontalScale struct {
	timeChart *TimeChart
	Height    int

	displayMin_ int64
	displayMax_ int64

	defaultDisplayMin_ int64
	defaultDisplayMax_ int64

	width        int
	hoverTime    int64
	allowedSteps []int64
}

func NewHorizontalScale(timeChart *TimeChart) *HorizontalScale {
	var c HorizontalScale
	c.timeChart = timeChart
	c.Height = 48
	c.displayMin_ = 0
	c.displayMax_ = 30

	c.allowedSteps = make([]int64, 0)

	c.allowedSteps = append(c.allowedSteps, 1)      // 1 nSec
	c.allowedSteps = append(c.allowedSteps, 5)      // 5 nSec
	c.allowedSteps = append(c.allowedSteps, 10)     // 10 nSec
	c.allowedSteps = append(c.allowedSteps, 50)     // 50 nSec
	c.allowedSteps = append(c.allowedSteps, 100)    // 100 nSec
	c.allowedSteps = append(c.allowedSteps, 500)    // 500 nSec
	c.allowedSteps = append(c.allowedSteps, 1000)   // 1 mSec
	c.allowedSteps = append(c.allowedSteps, 5000)   // 5 mSec
	c.allowedSteps = append(c.allowedSteps, 10000)  // 10 mSec
	c.allowedSteps = append(c.allowedSteps, 50000)  // 50 mSec
	c.allowedSteps = append(c.allowedSteps, 100000) // 100 mSec
	c.allowedSteps = append(c.allowedSteps, 500000) // 500 mSec

	c.allowedSteps = append(c.allowedSteps, 1*1000000)  // 1 Sec
	c.allowedSteps = append(c.allowedSteps, 2*1000000)  // 2 Sec
	c.allowedSteps = append(c.allowedSteps, 5*1000000)  // 5 Sec
	c.allowedSteps = append(c.allowedSteps, 10*1000000) // 10 Sec
	c.allowedSteps = append(c.allowedSteps, 15*1000000) // 15 Sec
	c.allowedSteps = append(c.allowedSteps, 30*1000000) // 30 Sec

	c.allowedSteps = append(c.allowedSteps, 1*60*1000000)  // 1 Min
	c.allowedSteps = append(c.allowedSteps, 2*60*1000000)  // 2 Min
	c.allowedSteps = append(c.allowedSteps, 5*60*1000000)  // 5 Min
	c.allowedSteps = append(c.allowedSteps, 10*60*1000000) // 10 Min
	c.allowedSteps = append(c.allowedSteps, 15*60*1000000) // 15 Min
	c.allowedSteps = append(c.allowedSteps, 30*60*1000000) // 30 Min

	c.allowedSteps = append(c.allowedSteps, 1*60*60*1000000)  // 1 Hour
	c.allowedSteps = append(c.allowedSteps, 3*60*60*1000000)  // 3 Hour
	c.allowedSteps = append(c.allowedSteps, 6*60*60*1000000)  // 6 Hour
	c.allowedSteps = append(c.allowedSteps, 12*60*60*1000000) // 12 Hour

	c.allowedSteps = append(c.allowedSteps, 1*24*3600*1000000)    // 1 Day
	c.allowedSteps = append(c.allowedSteps, 2*24*3600*1000000)    // 2 Day
	c.allowedSteps = append(c.allowedSteps, 7*24*3600*1000000)    // 7 Day
	c.allowedSteps = append(c.allowedSteps, 15*24*3600*1000000)   // 15 Day
	c.allowedSteps = append(c.allowedSteps, 1*30*24*3600*1000000) // 1 Month
	c.allowedSteps = append(c.allowedSteps, 2*30*24*3600*1000000) // 2 Month
	c.allowedSteps = append(c.allowedSteps, 3*30*24*3600*1000000) // 3 Month
	c.allowedSteps = append(c.allowedSteps, 6*30*24*3600*1000000) // 3 Month
	c.allowedSteps = append(c.allowedSteps, 365*24*3600*1000000)  // Year

	return &c
}

func (c *HorizontalScale) Dispose() {
	c.timeChart = nil
}

func (c *HorizontalScale) Draw(ctx ui.DrawContext, xOffset int, yOffset int, width int) {
	c.width = width

	dateTextWidth := 100

	textWidth, _, _ := canvas.MeasureText("Roboto", 10, false, false, "2006-01-02", false)
	countOfValues := width / textWidth
	countOfValues /= 2

	displayDatesBlocks := true
	countOfDays := (c.displayMax_ - c.displayMin_) / (24 * 3600 * 1000000)
	maxCountOfDaysForDisplay := int64(width / dateTextWidth)
	if countOfDays > maxCountOfDaysForDisplay {
		displayDatesBlocks = false
	}

	beautifulScale := c.getBeautifulScale(c.displayMin_, c.displayMax_, countOfValues, 0)
	for _, v := range beautifulScale {

		dateStr := time.Unix(0, v*1000).Format("2006-01-02")
		timeStr := time.Unix(0, v*1000).Format("15:04:05")
		us := time.Unix(0, v*1000).Nanosecond() / 1000
		usStr := ""

		if len(beautifulScale) > 1 {
			if beautifulScale[1]-beautifulScale[0] >= 60*1000000 {
				timeStr = time.Unix(0, v*1000).Format("15:04")
			}

			if beautifulScale[1]-beautifulScale[0] < 1000000 {
				if beautifulScale[1]-beautifulScale[0] < 1000 {
					usStr = fmt.Sprint(us, " us")
				} else {
					usStr = fmt.Sprint(us/1000, " ms")
				}
			}
		}

		labelWidth := 100

		off := yOffset

		if usStr != "" {
			ctx.SetColor(colornames.Gray)
			ctx.SetFontSize(10)
			ctx.DrawText(xOffset+c.valueToPixel(v)-labelWidth/2, off, labelWidth, 50, usStr)
			off += 12
		}
		ctx.SetColor(colornames.Lightblue)
		ctx.SetTextAlign(canvas.HAlignCenter, canvas.VAlignTop)
		ctx.DrawText(xOffset+c.valueToPixel(v)-labelWidth/2, off, labelWidth, 50, timeStr)
		off += 16
		if !displayDatesBlocks {
			ctx.SetColor(colornames.Gray)
			ctx.SetTextAlign(canvas.HAlignCenter, canvas.VAlignTop)
			ctx.DrawText(xOffset+c.valueToPixel(v)-labelWidth/2, off, labelWidth, 50, dateStr)
		}
		//ctx.DrawText(xOffset+c.valueToPixel(v), yOffset+15, timeStr, "Roboto", 10, colornames.Gray)
		ctx.SetColor(colornames.Gray)
		ctx.SetStrokeWidth(1)
		ctx.DrawLine(xOffset+c.valueToPixel(v), yOffset-5, xOffset+c.valueToPixel(v), yOffset)
	}

	if displayDatesBlocks {
		beautifulScaleForDates := c.getBeautifulScaleForDates(c.displayMin_, c.displayMax_)
		for _, v := range beautifulScaleForDates {
			date := time.Unix(0, v*1000)
			off := yOffset + 20

			currentColor := colornames.Gray

			isToday := false
			dateNow := time.Now()
			if date.Year() == dateNow.Year() && date.Month() == dateNow.Month() && date.Day() == dateNow.Day() {
				isToday = true
			}

			if isToday {
				currentColor = colornames.Green
			}

			dateStr := date.Format("2006-01-02")
			xPos1 := xOffset + c.valueToPixel(v)
			xPos2 := xOffset + c.valueToPixel(v+24*3600*1000000)

			xPos1Visible := xPos1
			if xPos1Visible < 0 {
				xPos1Visible = 0
			}
			xPos2Visible := xPos2
			if xPos2Visible > c.width {
				xPos2Visible = c.width
			}

			ctx.SetColor(currentColor)
			ctx.DrawLine(xPos1+2, off+5, xPos1+2, off+20)
			ctx.DrawLine(xPos2-2, off+5, xPos2-2, off+20)
			ctx.DrawLine(xPos1+5, off+13, xPos2-5, off+13)

			dateTextHeight := 50
			dateTextPosX := xPos1Visible + (xPos2Visible-xPos1Visible)/2 - (dateTextWidth / 2)
			dateTextPosY := off

			if dateTextPosX+dateTextWidth > c.width {
				dateTextPosX = c.width - dateTextWidth
			}

			if dateTextPosX < xPos1 {
				dateTextPosX = xPos1
			}

			if dateTextPosX < 0 {
				dateTextPosX = 0
			}
			if dateTextPosX+dateTextWidth > xPos2 {
				dateTextPosX = xPos2 - dateTextWidth
			}

			ctx.SetColor(c.timeChart.BackColor())
			ctx.FillRect(dateTextPosX, dateTextPosY, dateTextWidth, dateTextHeight)
			ctx.SetColor(currentColor)
			ctx.DrawText(dateTextPosX, dateTextPosY, dateTextWidth, dateTextHeight, dateStr)
		}
	}

	///////////////////////////////////////////////////////////////////////////////
	// Hover
	if c.timeChart.MouseIsInside() && c.timeChart.lastMouseX > xOffset && len(c.timeChart.areas) > 0 {
		xHover := c.valueToPixel(c.hoverTime) + xOffset
		ctx.SetColor(c.timeChart.InactiveColor())
		ctx.SetStrokeWidth(1)
		ctx.DrawLine(xHover, 0, xHover, yOffset)

		hoverDateStr := time.Unix(0, c.hoverTime*1000).Format("2006:01-02")
		hoverTimeStr := time.Unix(0, c.hoverTime*1000).Format("15:04:05")
		hoverDateTime := hoverTimeStr + "\r\n" + hoverDateStr

		ctx.SetColor(c.timeChart.borderColor())
		ctx.FillRect(xHover-50, yOffset+2, 100, 26)
		ctx.SetColor(c.timeChart.InactiveColor())
		ctx.SetStrokeWidth(1)
		ctx.DrawRect(xHover-50, yOffset+2, 100, 26)

		ctx.SetColor(c.timeChart.textColor())
		ctx.SetFontSize(10)
		ctx.SetTextAlign(canvas.HAlignCenter, canvas.VAlignCenter)
		ctx.DrawText(xHover-50, yOffset+2, 100, 26, hoverDateTime)
	}
	///////////////////////////////////////////////////////////////////////////////

	if c.timeChart.Editing() {
		foreColor, backColor := c.timeChart.borderColor(), c.timeChart.BackColor()

		btnX, btnY, btnW, btnH := 0, c.timeChart.Height()-48, 48, 48
		if c.timeChart.lastMouseX > btnX && c.timeChart.lastMouseX < btnX+btnW {
			if c.timeChart.lastMouseY > btnY && c.timeChart.lastMouseY < btnY+btnH {
				foreColor, backColor = backColor, foreColor
			}
		}

		ctx.SetColor(backColor)
		ctx.FillRect(btnX, btnY, btnW, btnH)
		ctx.SetColor(foreColor)
		ctx.DrawRect(btnX, btnY, btnW, btnH)
		img := resize.Resize(uint(24), uint(24), uiresources.ResImgCol(uiresources.R_icons_material4_png_action_settings_materialicons_48dp_1x_baseline_settings_black_48dp_png, foreColor), resize.Bicubic)
		ctx.DrawImage(btnX+12, btnY+12, btnW, btnH, img)
	}

	/*c.drawDates(ctx, xOffset, yOffset, c.displayMin_, c.displayMax_)
	c.drawTimes(ctx, xOffset, yOffset, c.displayMin_, c.displayMax_)
	c.drawSecs(ctx, xOffset, yOffset, c.displayMin_, c.displayMax_)*/
}

func (c *HorizontalScale) MouseDown(event *uievents.MouseDownEvent) {
	if c.timeChart.Editing() {
		btnX, btnY, btnW, btnH := 0, c.timeChart.Height()-48, 48, 48
		if c.timeChart.lastMouseX > btnX && c.timeChart.lastMouseX < btnX+btnW {
			if c.timeChart.lastMouseY > btnY && c.timeChart.lastMouseY < btnY+btnH {
				dialog := NewSettingsDialog(c.timeChart)
				dialog.ShowDialog()
			}
		}
	}
}

func (c *HorizontalScale) drawDates1(ctx *canvas.CanvasDirect, xOffset int, yOffset int, min int64, max int64) {
	yOffset += 30
	usInDay := int64(86400) * int64(1000000)
	t := time.Unix(0, min*1000)
	pixelsInItem := c.valueToPixel(min+usInDay) - c.valueToPixel(min)
	if pixelsInItem < 10 {
		return
	}
	beginDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).UnixNano() / 1000
	for date := beginDate; date < max; date += usInDay {
		dateOffsetInPixels := c.valueToPixel(date) + xOffset
		ctx.DrawLine(dateOffsetInPixels, yOffset, dateOffsetInPixels, yOffset+15, 1, colornames.Yellow)
		dateStr := time.Unix(0, date*1000).Format("2006-01-02")
		ctx.DrawTextMultiline(xOffset+c.valueToPixel(date)+10, yOffset, pixelsInItem, 15, canvas.HAlignLeft, canvas.VAlignTop, dateStr, colornames.Gray, "Roboto", 10, false)
	}
}

func (c *HorizontalScale) drawDates(ctx *canvas.CanvasDirect, xOffset int, yOffset int, min int64, max int64) {
	yOffset += 30
	beautifulScale := c.getBeautifulScale(c.displayMin_, c.displayMax_, 10, int64(86400)*int64(1000000))

	for _, date := range beautifulScale {
		dateOffsetInPixels := c.valueToPixel(date) + xOffset
		ctx.DrawLine(dateOffsetInPixels, yOffset, dateOffsetInPixels, yOffset+15, 1, colornames.Yellow)
		dateStr := time.Unix(0, date*1000).Format("2006-01-02")
		ctx.DrawTextMultiline(xOffset+c.valueToPixel(date)+10, yOffset, 100, 15, canvas.HAlignLeft, canvas.VAlignTop, dateStr, colornames.Gray, "Roboto", 10, false)
	}
}

func (c *HorizontalScale) drawTimes(ctx *canvas.CanvasDirect, xOffset int, yOffset int, min int64, max int64) {
	yOffset += 15

	beautifulScale := c.getBeautifulScale(c.displayMin_, c.displayMax_, 10, int64(3600)*int64(1000000))

	usInHour := int64(3600) * int64(1000000)
	//t := time.Unix(0, min*1000)
	pixelsInItem := c.valueToPixel(min+usInHour) - c.valueToPixel(min)
	if pixelsInItem < 30 {
		return
	}
	//beginDate := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, time.Local).UnixNano() / 1000
	labelWidth := 100

	for _, date := range beautifulScale {
		dateOffsetInPixels := c.valueToPixel(date) + xOffset
		ctx.DrawLine(dateOffsetInPixels, yOffset, dateOffsetInPixels, yOffset+3, 1, colornames.Yellow)
		timeStr := time.Unix(0, date*1000).Format("15:04")
		ctx.DrawTextMultiline(xOffset+c.valueToPixel(date)-labelWidth/2, yOffset+2, labelWidth, 15, canvas.HAlignCenter, canvas.VAlignTop, timeStr, colornames.Gray, "Roboto", 10, false)
	}
}

func (c *HorizontalScale) drawSecs(ctx *canvas.CanvasDirect, xOffset int, yOffset int, min int64, max int64) {
	yOffset += 0
	//usInHour := int64(60) * int64(1000000)
	//t := time.Unix(0, min*1000)

	beautifulScale := c.getBeautifulScale(c.displayMin_, c.displayMax_, 20, 0)
	labelWidth := 100

	//beginDate := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.Local).UnixNano() / 1000
	for _, date := range beautifulScale {
		dateOffsetInPixels := c.valueToPixel(date) + xOffset
		ctx.DrawLine(dateOffsetInPixels, yOffset, dateOffsetInPixels, yOffset+2, 1, colornames.Yellow)
		timeStr := time.Unix(0, date*1000).Format("05.999999")
		ctx.DrawTextMultiline(xOffset+c.valueToPixel(date)-labelWidth/2, yOffset+2, labelWidth, 15, canvas.HAlignCenter, canvas.VAlignTop, timeStr, colornames.Gray, "Roboto", 10, false)
	}
}

func (c *HorizontalScale) SetDisplayRange(displayMin, displayMax int64) {
	c.displayMin_ = displayMin
	c.displayMax_ = displayMax
}

func (c *HorizontalScale) SetDefaultDisplayRange(defaultDisplayMin, defaultDisplayMax int64) {
	c.defaultDisplayMin_ = defaultDisplayMin
	c.defaultDisplayMax_ = defaultDisplayMax
}

func (c *HorizontalScale) ResetToDefaultRange() {
	c.displayMin_ = c.defaultDisplayMin_
	c.displayMax_ = c.defaultDisplayMax_
}

func (c *HorizontalScale) getBeautifulScaleForDates(min int64, max int64) []int64 {
	_, timeOffset := time.Now().Zone()

	dates := make([]int64, 0)
	min += int64(timeOffset) * 1000000
	max += int64(timeOffset) * 1000000

	min = min - (min % (24 * 3600 * 1000000))
	max = max + ((24 * 3600 * 1000000) - (max % (24 * 3600 * 1000000)))

	add := int64(24 * 3600 * 1000000)
	for t := min; t < max; t += add {
		dates = append(dates, t-int64(timeOffset)*1000000)
	}

	return dates
}

func (c *HorizontalScale) getBeautifulScale(min int64, max int64, countOfPoints int, minStep int64) []int64 {
	var scale []int64

	if max < min {
		return scale
	}

	if max == min {
		scale = append(scale, min)
		return scale
	}

	_, timeOffset := time.Now().Zone()

	min += int64(timeOffset) * 1000000
	max += int64(timeOffset) * 1000000

	diapason := max - min

	// Raw step - ugly
	step := int64(1)
	if countOfPoints != 0 {
		step = diapason / int64(countOfPoints)
	}
	newMin := min

	for i := 0; i < len(c.allowedSteps); i++ {
		st := c.allowedSteps[i]
		if st < minStep {
			continue
		}
		if step < st {
			step = st // Beautiful step
			break
		}
	}
	newMin = newMin - (newMin % step) // New begin point

	// Make points
	for i := 0; i < countOfPoints; i++ {
		if newMin < max && newMin > min {
			scale = append(scale, newMin-int64(timeOffset)*1000000)
		}
		newMin += step
	}
	return scale
}

func (c *HorizontalScale) valueToPixel(value int64) int {
	chartPixels := c.width
	displayRange := c.displayMax_ - c.displayMin_
	offsetOfValueFromMin := value - c.displayMin_
	onePixelValue := float64(chartPixels) / float64(displayRange)
	return int(onePixelValue * float64(offsetOfValueFromMin))
}

func (c *HorizontalScale) pixelToValue(x int) int64 {
	chartPixels := c.width
	displayRange := c.displayMax_ - c.displayMin_
	onePixelValue := float64(chartPixels) / float64(displayRange)
	return int64(float64(x)/float64(onePixelValue) + float64(c.displayMin_))
}

func (c *HorizontalScale) pixelToValueRange(xRange int) int64 {
	chartPixels := c.width
	displayRange := c.displayMax_ - c.displayMin_
	onePixelValue := float64(chartPixels) / float64(displayRange)
	return int64(float64(xRange) / float64(onePixelValue))
}

func (c *HorizontalScale) setHoverTimePoint(time int64) {
	c.hoverTime = time
}
