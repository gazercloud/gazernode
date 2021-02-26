package timechart

import (
	"github.com/gazercloud/gazerui/ui"
	"golang.org/x/image/colornames"
	"image/color"
)

type AreaTopHeader struct {
	height int
}

func (c *AreaTopHeader) Height() int {
	return c.height
}

func NewAreaTopHeader(height int) *AreaTopHeader {
	var c AreaTopHeader
	c.height = height
	return &c
}

func (c *AreaTopHeader) draw(ctx ui.DrawContext, vLines []vline, chart *TimeChart) {
	lastIndexOfVLineWithValue := -1

	type Range struct {
		indexFrom  int
		indexTo    int
		colorIndex int
	}

	ranges := make([]Range, 0)

	var currentRange Range
	currentRange.indexFrom = 0
	colorIndex := 0

	for index, vLine := range vLines {

		if vLine.hasValues {
			lastIndexOfVLineWithValue = index
			colorIndex = 1
			if vLine.hasBadQuality {
				colorIndex = 2
			}
		} else {
			colorIndex = currentRange.colorIndex
		}

		if colorIndex != currentRange.colorIndex {
			currentRange.indexTo = index - 1
			if currentRange.indexTo < 0 {
				currentRange.indexTo = 0
			}
			ranges = append(ranges, currentRange)
			currentRange.indexFrom = index
			currentRange.colorIndex = colorIndex
		}
	}

	if lastIndexOfVLineWithValue < 0 {
		lastIndexOfVLineWithValue = 0
	}

	currentRange.indexTo = lastIndexOfVLineWithValue
	ranges = append(ranges, currentRange)

	for _, r := range ranges {
		height := c.height - 2
		col := color.Color(colornames.Green)
		if r.colorIndex == 0 {
			col = chart.borderColor()
			height -= 3
		}
		if r.colorIndex == 1 {
			col = colornames.Green
			height -= 3
		}
		if r.colorIndex == 2 {
			col = colornames.Crimson
		}

		ctx.SetColor(col)
		ctx.FillRect(vLines[r.indexFrom].X, 1, vLines[r.indexTo].X-vLines[r.indexFrom].X+1, height)
	}
}
