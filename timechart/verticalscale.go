package timechart

import (
	"fmt"
	"github.com/gazercloud/gazerui/canvas"
	"github.com/gazercloud/gazerui/ui"
	"image/color"
	"math"
	"strings"
)

type VerticalScale struct {
	Width  int
	Height int

	displayMin float64
	displayMax float64

	series *Series
}

func NewVerticalScale(series *Series) *VerticalScale {
	var c VerticalScale
	c.Width = 100
	c.series = series

	c.displayMin = 0
	c.displayMax = 60
	return &c
}

func (c *VerticalScale) SetDisplayRange(displayMin, displayMax float64) {
	c.displayMin = displayMin
	c.displayMax = displayMax
}

func (c *VerticalScale) signsAfterPoint(str string) int {
	pointPosition := strings.Index(str, ".")
	if pointPosition < 0 {
		pointPosition = len(str)
	}
	return len(str) - pointPosition
}

func (c *VerticalScale) prepareValuesText(values []ScaleValue) []ScaleValue {
	result := make([]ScaleValue, 0)
	mapsCountByStringLen := make(map[int]int)
	for i, v := range values {
		txt := fmt.Sprintf("%.4f", v.Value)
		signsAfterPoint := c.signsAfterPoint(txt)
		if _, ok := mapsCountByStringLen[signsAfterPoint]; ok {
			mapsCountByStringLen[signsAfterPoint]++
		} else {
			mapsCountByStringLen[signsAfterPoint] = 1
		}
		values[i].ValueText = txt
	}

	countOfTarget := 0
	maxTxtLenForDisplay := 10
	for i := 0; i < 10; i++ {
		if _, ok := mapsCountByStringLen[i]; ok {
			countOfTarget += mapsCountByStringLen[i]
		}
		if countOfTarget > len(values)/3 {
			maxTxtLenForDisplay = i
			break
		}
	}

	for _, v := range values {
		if c.signsAfterPoint(v.ValueText) <= maxTxtLenForDisplay {
			var scaleValue ScaleValue
			scaleValue.Value = v.Value
			scaleValue.ValueText = v.ValueText
			result = append(result, scaleValue)
		}
	}
	return result
}

func (c *VerticalScale) prepareValuesTextSimple(values []ScaleValue) []ScaleValue {
	result := make([]ScaleValue, 0)
	for _, v := range values {
		var scaleValue ScaleValue
		scaleValue.Value = v.Value
		scaleValue.ValueText = fmt.Sprintf("%.6g", v.Value)
		result = append(result, scaleValue)
	}
	return result
}

func (c *VerticalScale) Draw(ctx ui.DrawContext, xOffset int, yOffset int, col color.Color) {
	fontSize := float64(12)
	count := float64(c.Height) / (fontSize + fontSize/2)

	beautifulScale := c.getBeautifulScale(c.displayMin, c.displayMax, int(count))
	beautifulScaleOptimazed := c.prepareValuesTextSimple(beautifulScale)
	for _, v := range beautifulScaleOptimazed {
		_, hText, _ := canvas.MeasureText("Roboto", fontSize, false, false, fmt.Sprint(v), false)
		ctx.SetColor(col)
		ctx.SetFontSize(12)
		ctx.SetStrokeWidth(1)
		ctx.SetTextAlign(canvas.HAlignRight, canvas.VAlignCenter)
		//strVal := fmt.Sprintf("%.4g", v)
		ctx.DrawText(xOffset+3, yOffset+c.getPointOnY(v.Value)-hText/2, c.Width-10, hText, v.ValueText)
		ctx.DrawLine(xOffset+c.Width-3, yOffset+c.getPointOnY(v.Value), xOffset+c.Width, yOffset+c.getPointOnY(v.Value))
	}

	ctx.SetStrokeWidth(1)
	ctx.SetColor(col)
	ctx.DrawLine(xOffset+c.Width, yOffset, xOffset+c.Width, yOffset+c.Height)

	c.drawCurrentValue(ctx, xOffset, yOffset)
}

func (c *VerticalScale) getPointOnY(value float64) int {
	chartPixels := c.Height
	yDelta := c.displayMax - c.displayMin
	onePixelValue := float64(1)
	if math.Abs(yDelta) > 0.0001 {
		onePixelValue = float64(chartPixels) / float64(yDelta)
	}
	return int(float64(chartPixels) - onePixelValue*(value-c.displayMin))
}

func (c *VerticalScale) convertYtoValue(y int) float64 {
	chartPixels := c.Height
	yDelta := c.displayMax - c.displayMin
	onePixelValue := float64(1)
	if math.Abs(yDelta) > 0.0001 {
		onePixelValue = float64(chartPixels) / float64(yDelta)
	}
	return float64(y) / onePixelValue
}

func (c *VerticalScale) getValueByY(y int) float64 {
	chartPixels := c.Height
	yDelta := c.displayMax - c.displayMin
	onePixelValue := float64(1)
	if math.Abs(yDelta) > 0.0001 {
		onePixelValue = float64(chartPixels) / float64(yDelta)
	}
	return float64(c.Height-y)/onePixelValue + c.displayMin
}

type ScaleValue struct {
	Value     float64
	ValueText string
}

func (c *VerticalScale) getBeautifulScale(min float64, max float64, countOfPoints int) []ScaleValue {
	var scale []ScaleValue
	scale = make([]ScaleValue, 0)

	if max < min {
		return scale
	}

	if max == min {
		scale = append(scale, ScaleValue{Value: min})
		return scale
	}

	diapason := max - min

	// Некрасивый шаг
	step := diapason / float64(countOfPoints)

	// Порядок
	log := math.Ceil(math.Log10(float64(step)))
	// Красивый шаг = степень 10-ки
	step10 := math.Pow10(int(log))

	// деление на 2 - это тоже красиво
	for float64(diapason)/float64(step10/2) < float64(countOfPoints) {
		step10 = step10 / 2
	}

	// Определяем новый минимум
	newMin := float64(min) - math.Mod(float64(min), float64(step10))

	// Генерируем точки
	for i := 0; i < countOfPoints; i++ {
		if newMin < max && newMin > min {
			scale = append(scale, ScaleValue{Value: newMin})
		}
		newMin += step10
	}
	return scale
}

func (c *VerticalScale) drawCurrentValue(ctx ui.DrawContext, xOffset int, yOffset int) {
	if c.series == nil {
		return
	}

	if len(c.series.values) < 1 {
		return
	}

	// LastValue
	/*lastValue := c.series.values[len(c.series.values)-1].LastValue
	lastValueAsString := strconv.FormatFloat(lastValue, 'f', 2, 64)
	ctx.SetColor(c.series.color)
	ctx.FillRect(xOffset, 0, c.Width, 16)
	ctx.SetColor(colornames.White)
	ctx.SetFontSize(12)
	ctx.DrawText(xOffset+3, 0, 100, 20, lastValueAsString)*/

	/*ctx.Save()
	ctx.Translate(0, yOffset)
	// Hover Value
	yHover := c.series.hoverY
	hoverText := strconv.FormatFloat(c.getValueByY(yHover-yOffset), 'f', 2, 64)
	_, hText, _ := canvas.MeasureText("Roboto", 12, false, false, hoverText, false)
	r, g, b, _ := c.series.color.RGBA()
	hoverBackColor := color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 192}
	ctx.SetColor(hoverBackColor)
	ctx.FillRect(xOffset, yHover-hText/2-yOffset, c.Width, hText)

	ctx.SetColor(colornames.White)
	ctx.SetFontSize(12)
	ctx.DrawText(xOffset, yHover-hText/2-yOffset, 100, 20, hoverText)
	ctx.Load()*/
}
