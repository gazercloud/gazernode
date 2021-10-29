package timechart

import (
	"fmt"
	"github.com/gazercloud/gazerui/canvas"
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiproperties"
	"golang.org/x/image/colornames"
	"image/color"
	"math"
	"time"
)

type InteractiveSelectionModifier int

const (
	InteractiveSelectionModifierLeftMouseButton  InteractiveSelectionModifier = 1
	InteractiveSelectionModifierRightMouseButton InteractiveSelectionModifier = 2
	InteractiveSelectionModifierShiftHorizontal  InteractiveSelectionModifier = 4
	InteractiveSelectionModifierControl          InteractiveSelectionModifier = 8
	InteractiveSelectionModifierAlt              InteractiveSelectionModifier = 16
)

type SelectionAxes int

const (
	SelectionAxesX  SelectionAxes = 0
	SelectionAxesY  SelectionAxes = 1
	SelectionAxesXY SelectionAxes = 2
)

type ResizingDirection int

const (
	ResizingDirectionTop    ResizingDirection = 0
	ResizingDirectionRight  ResizingDirection = 1
	ResizingDirectionBottom ResizingDirection = 2
	ResizingDirectionLeft   ResizingDirection = 3
)

type AllowedSelection struct {
	interactiveModifiers InteractiveSelectionModifier
	axes                 SelectionAxes
}

type Selection struct {
	interactiveModifiers InteractiveSelectionModifier
	axes                 SelectionAxes
	minX                 int64
	maxX                 int64
}

type ResizingSelection struct {
	selection *Selection
	direction ResizingDirection
}

type ChartPoint struct {
	series *Series
	Index  int
	X      int64
	Y      float64
}

type TimeChart struct {
	uicontrols.Control
	areas           []*Area
	horizontalScale *HorizontalScale

	editing_ bool

	selections_                   []*Selection
	allowedSelections             []*AllowedSelection
	resizing                      *ResizingSelection
	movingSelection               *Selection
	leftXScaleMovingHandPositionX int64

	xyZoomModifiers    InteractiveSelectionModifier // Способ зума области (Прямоугольник)
	xZoomModifiers     InteractiveSelectionModifier // Способ зума по оси X
	yZoomModifiers     InteractiveSelectionModifier // Способ зума по оси Y
	chartMoveModifiers InteractiveSelectionModifier // Способ зума по оси Y

	chartMoving              bool
	chartMovingHandPositionX int
	chartMovingHandPositionY int
	chartMovingDisplayXMin   int64
	chartMovingDisplayYMin   float64

	lastMouseDownPointX int

	selectionProcessing *Selection

	//showAreaHeader      bool
	//unitedVerticalScale bool

	displayFullRange   bool
	movingPoint        *ChartPoint
	mousePressedValueX int64

	mouseLeftButtonPressed  bool
	mouseRightButtonPressed bool

	lastMouseX int
	lastMouseY int

	OnMouseDropOnArea func(droppedValue interface{}, area *Area)
	OnZoomed          func()
	OnMoved           func()

	DrawTime      int64
	mouseIsInside bool

	color0 *uiproperties.Property
	color1 *uiproperties.Property
	color2 *uiproperties.Property
	color3 *uiproperties.Property
	color4 *uiproperties.Property

	OnChartContextMenuNeed func(timeChart *TimeChart, area *Area, areaIndex int) uiinterfaces.Menu

	//dataProvider IDataProvider

	//pr DataProvider
}

func NewTimeChart(parent uiinterfaces.Widget) *TimeChart {
	var c TimeChart
	c.InitControl(parent, &c)

	c.color0 = uiproperties.NewProperty("color0", uiproperties.PropertyTypeColor)
	c.AddProperty("color0", c.color0)
	c.color1 = uiproperties.NewProperty("color1", uiproperties.PropertyTypeColor)
	c.AddProperty("color1", c.color1)
	c.color2 = uiproperties.NewProperty("color2", uiproperties.PropertyTypeColor)
	c.AddProperty("color2", c.color2)
	c.color3 = uiproperties.NewProperty("color3", uiproperties.PropertyTypeColor)
	c.AddProperty("color3", c.color3)
	c.color4 = uiproperties.NewProperty("color4", uiproperties.PropertyTypeColor)
	c.AddProperty("color4", c.color4)
	c.UpdateStyle()

	c.areas = make([]*Area, 0)
	c.horizontalScale = NewHorizontalScale(&c)
	//c.dataProvider = &c.pr
	//c.pr.Init()

	c.SetXExpandable(true)
	c.SetYExpandable(true)

	c.addInteractiveSelectionType(true, false, true, false, false, SelectionAxesX)
	c.addInteractiveZoom(true, false, false, false, false, SelectionAxesX)
	c.chartMoveModifiers = InteractiveSelectionModifierRightMouseButton
	c.ZoomShowEntire()

	c.OnContextMenuNeed = func(x, y int) uiinterfaces.Menu {
		if c.OnChartContextMenuNeed != nil {
			area := c.AreaByPoint(x, y)
			areaIndex := -1
			for index, a := range c.areas {
				if a == area {
					areaIndex = index
					break
				}
			}
			return c.OnChartContextMenuNeed(&c, area, areaIndex)
		}
		return nil
	}

	return &c
}

func (c *TimeChart) ControlType() string {
	return "TimeChart"
}

func (c *TimeChart) SetWidth(w int) {
	c.Control.SetWidth(w)
}

func (c *TimeChart) SetHeight(h int) {
	c.Control.SetHeight(h)
	for _, a := range c.areas {
		a.SetHeight((h - c.horizontalScale.Height) / len(c.areas))
	}
}

func (c *TimeChart) SetEditing(editing bool) {
	if c.editing_ != editing {
		c.editing_ = editing
		c.Update("TimeChart")
	}
}

func (c *TimeChart) Editing() bool {
	return c.editing_
}

func (c *TimeChart) SelectedTimeRange() (int64, int64) {
	if len(c.selections_) > 0 {
		return c.selections_[0].minX, c.selections_[0].maxX
	}
	return 0, 0
}

func (c *TimeChart) borderColor() color.Color {
	col, _, _, _ := c.BorderColors()
	return col
}

func (c *TimeChart) textColor() color.Color {
	col := c.ForeColor()
	return col
}

func (c *TimeChart) Dispose() {
	c.Control.Dispose()
	for _, area := range c.areas {
		area.Dispose()
	}

	if c.horizontalScale != nil {
		c.horizontalScale.Dispose()
		c.horizontalScale = nil
	}

	c.selections_ = nil
	c.selectionProcessing = nil
	c.allowedSelections = nil
	c.resizing = nil
	c.movingPoint = nil

	c.OnMouseDropOnArea = nil
	c.OnZoomed = nil
	c.OnMoved = nil
}

func (c *TimeChart) AddArea() *Area {
	area := NewArea(c)
	c.Window().UpdateLayout()
	//area.dataProvider = c.dataProvider
	c.areas = append(c.areas, area)
	return area
}

func (c *TimeChart) RemoveAllAreas() {
	for _, area := range c.areas {
		area.Dispose()
	}

	c.areas = make([]*Area, 0)
}

func (c *TimeChart) RemoveAreaByIndex(index int) {
	if index < 0 || index > len(c.areas) {
		return
	}
	c.areas[index].Dispose()
	c.areas = append(c.areas[:index], c.areas[index+1:]...)
	c.Update("TimeChart")
	c.SetHeight(c.Height())
}

func (c *TimeChart) IsChartMoving() bool {
	return c.chartMoving
}

func (c *TimeChart) AddSeries(area *Area, id string) *Series {
	ser := NewSeries(id, area)

	//ser.dataProvider = c.dataProvider

	/*	areaIndex := 0
		for index, a := range c.areas {
			if a == area {
				areaIndex = index
				break
			}
		}*/

	colorIndex := len(area.series)

	var col color.Color = colornames.Black
	switch colorIndex % 5 {
	case 0:
		//col = uiproperties.ParseHexColor("#ff8c00")
		//col = color.RGBA{0, 100, 214, 255}
		col = c.color0.Color()
	case 1:
		//col = colornames.Darkgreen
		col = c.color1.Color()
	case 2:
		//col = colornames.Red
		col = c.color2.Color()
	case 3:
		//col = colornames.Yellow
		col = c.color3.Color()
	case 4:
		//col = colornames.Magenta
		col = c.color4.Color()
	}
	ser.SetColor(col)

	area.series = append(area.series, ser)

	return ser
}

func (c *TimeChart) Areas() []*Area {
	return c.areas
}

func (c *TimeChart) SetShowQualities(showQualities bool) {
	for _, a := range c.areas {
		a.SetShowQualities(showQualities)
	}
	c.Update("TimeChart")
}

func (c *TimeChart) SetUnitedVerticalScale(unitedVerticalScale bool) {
	for _, a := range c.areas {
		a.SetUnitedScale(unitedVerticalScale)
	}
	c.Update("TimeChart")
}

func (c *TimeChart) USecPerPixel() float64 {
	w := float64(c.horizontalScale.width)
	r := float64(c.horizontalScale.displayMax_ - c.horizontalScale.displayMin_)
	return r / w
}

func (c *TimeChart) HorMin() int64 {
	return c.horizontalScale.displayMin_
}

func (c *TimeChart) HorMax() int64 {
	return c.horizontalScale.displayMax_
}

func (c *TimeChart) SetHorizRange(min, max int64) {
	c.horizontalScale.SetDefaultDisplayRange(min, max)
	c.horizontalScale.SetDisplayRange(min, max)
	if c.OnZoomed != nil {
		c.OnZoomed()
	}
}

func (c *TimeChart) ResetHorizontalRange(ctx *canvas.CanvasDirect) {
	hMin, hMax, err := c.calcHorizontalLimits()
	if err == nil {
		diff := hMax - hMin
		if diff > 0 {
			diff = diff / 10
		} else {
			diff = 1000000 // 1 sec for one point chart data
		}

		c.horizontalScale.SetDefaultDisplayRange(hMin-diff, hMax+diff)
	}
}

func (c *TimeChart) SetDefaultDisplayRange(defaultDisplayMin, defaultDisplayMax int64) {
	c.horizontalScale.SetDefaultDisplayRange(defaultDisplayMin, defaultDisplayMax)
	c.Update("TimeChart")
}

func (c *TimeChart) Draw(ctx ui.DrawContext) {

	//ctx.SetColor(colornames.Black)
	//ctx.FillRect(0, 0, c.Width(), c.Height())

	//c.SetInnerSizeDirect(c.Control.InnerWidth(), len(c.areas)*AREA_HIGHT)
	ctx.Save()
	//ctx.Clip(ctx.TranslatedX()+c.ScrollOffsetX(), ctx.TranslatedY()+c.ScrollOffsetY(), c.Width(), c.Height())

	timeBegin := time.Now().UnixNano()

	if c.displayFullRange {
		c.ZoomShowEntire()
	}

	leftScalesWidth := c.calcLeftScalesWidth()
	c.horizontalScale.Draw(ctx, leftScalesWidth, c.Height()-c.horizontalScale.Height, c.Width()-leftScalesWidth)
	yOffset := 0
	for _, area := range c.areas {
		//ctx.Clip(ctx.TranslatedX(), ctx.TranslatedY(), c.Width(), area.height)
		area.Draw(ctx, leftScalesWidth, c.horizontalScale, c.Width(), c.ForeColor(), yOffset)
		yOffset += area.height
	}
	//ctx.DrawLine(leftScalesWidth, 0, leftScalesWidth, yOffset, 1, colornames.Gray)
	ctx.SetColor(c.ForeColor())
	ctx.SetStrokeWidth(1)
	ctx.DrawLine(0, c.Height()-c.horizontalScale.Height, c.Width(), c.Height()-c.horizontalScale.Height)

	if c.selectionProcessing != nil {
		if c.selectionProcessing.axes == SelectionAxesX {
			x1 := c.horizontalScale.valueToPixel(c.selectionProcessing.minX) + c.calcLeftScalesWidth()
			x2 := c.horizontalScale.valueToPixel(c.selectionProcessing.maxX) + c.calcLeftScalesWidth()
			leftLength := x2 - x1
			//ctx.DrawRect(x1, 0, leftLength, c.Height(), colornames.Blue, 1)
			ctx.SetColor(color.RGBA{0, 0, 255, 50})
			ctx.FillRect(x1, 0, leftLength, c.Height())
		}
	}

	for _, selection := range c.selections_ {
		x1 := c.horizontalScale.valueToPixel(selection.minX) + c.calcLeftScalesWidth()
		x2 := c.horizontalScale.valueToPixel(selection.maxX) + c.calcLeftScalesWidth()
		leftLength := x2 - x1
		ctx.SetColor(color.RGBA{0, 255, 0, 50})
		ctx.FillRect(x1, 0, leftLength, c.Height())
		//ctx.DrawRect(x1, 0, leftLength, c.Height(), colornames.Green, 1)
	}

	timeEnd := time.Now().UnixNano()
	//fmt.Println("Series draw time: ", (timeEnd-timeBegin)/1000000)
	c.DrawTime = (timeEnd - timeBegin) / 1000000

	ctx.Load()
}

func (c *TimeChart) calcLeftScalesWidth() int {
	leftScalesWidth := 0
	for _, area := range c.areas {
		areaLeftScalesWidth := area.calcLeftScalesWidth()
		if areaLeftScalesWidth > leftScalesWidth {
			leftScalesWidth = areaLeftScalesWidth
		}
	}
	return leftScalesWidth
}

func (c *TimeChart) calcHorizontalLimits() (int64, int64, error) {
	min := int64(100000000000000000)
	max := int64(-100000000000000000)

	if len(c.areas) < 1 {
		return 0, 0, fmt.Errorf("Area is empty")
	}

	for _, ser := range c.areas {
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

func (c *TimeChart) AreaByPoint(x, y int) *Area {
	yOffset := 0
	for _, area := range c.areas {
		if y > yOffset && y < yOffset+area.height {
			return area
		}
		yOffset += area.height
	}
	return nil
}

func (c *TimeChart) AreaYOffset(area *Area) int {
	yOffset := 0
	for _, a := range c.areas {
		if area == a {
			return yOffset
		}
		yOffset += area.height
	}
	return 0
}

func (c *TimeChart) MouseDrop(ev *uievents.MouseDropEvent) {
	if c.OnMouseDropOnArea != nil {
		fmt.Println("!!!!!!!! drop", ev.Y)
		c.OnMouseDropOnArea(ev.DroppingObject, c.AreaByPoint(ev.X, ev.Y))
	}
}

func (c *TimeChart) MouseWheel(event *uievents.MouseWheelEvent) {
	if event.Modifiers.Shift {
		// Shift
		timeRange := c.horizontalScale.displayMax_ - c.horizontalScale.displayMin_
		delta := timeRange / int64(-event.Delta*20)
		c.horizontalScale.SetDisplayRange(c.horizontalScale.displayMin_+delta, c.horizontalScale.displayMax_+delta)
		c.Update("TimeChart")
	} else {
		// Zoom
		c.displayFullRange = false
		timeRange := c.horizontalScale.displayMax_ - c.horizontalScale.displayMin_
		hoverValue := c.horizontalScale.pixelToValue(c.lastMouseX - c.calcLeftScalesWidth())
		koefX := float64(hoverValue-c.horizontalScale.displayMin_) / float64(timeRange)
		if koefX < 0 {
			koefX = 0
		}

		if koefX > 1 {
			koefX = 1
		}

		if event.Delta != 0 {
			delta := timeRange / int64(-event.Delta*20)
			c.horizontalScale.SetDisplayRange(c.horizontalScale.displayMin_-int64(float64(delta)*koefX), c.horizontalScale.displayMax_+int64(float64(delta)*(1-koefX)))
		}
		c.Update("TimeChart")
	}
}

func (c *TimeChart) buildInteractiveSelectionType(leftMouseButton, rightMouseButton, shiftKey, controlKey, altKey bool) InteractiveSelectionModifier {
	var interactiveSelectionModifier InteractiveSelectionModifier
	interactiveSelectionModifier = 0

	if leftMouseButton {
		interactiveSelectionModifier |= InteractiveSelectionModifierLeftMouseButton
	}
	if rightMouseButton {
		interactiveSelectionModifier |= InteractiveSelectionModifierRightMouseButton
	}
	if shiftKey {
		interactiveSelectionModifier |= InteractiveSelectionModifierShiftHorizontal
	}
	if controlKey {
		interactiveSelectionModifier |= InteractiveSelectionModifierControl
	}
	if altKey {
		interactiveSelectionModifier |= InteractiveSelectionModifierAlt
	}

	return interactiveSelectionModifier
}

func (c *TimeChart) getResizingSelection(x, y int) *ResizingSelection {
	var result *ResizingSelection
	for i := 0; i < len(c.selections_); i++ {
		pSelection := c.selections_[i]
		if pSelection.axes == SelectionAxesX {
			minXPixels := c.horizontalScale.valueToPixel(pSelection.minX) + c.calcLeftScalesWidth()
			if math.Abs(float64(minXPixels-x)) < 5 {
				result = &ResizingSelection{}
				result.selection = pSelection
				result.direction = ResizingDirectionLeft
				break
			}
			maxXPixels := c.horizontalScale.valueToPixel(pSelection.maxX) + c.calcLeftScalesWidth()
			if math.Abs(float64(maxXPixels-x)) < 5 {
				result = &ResizingSelection{}
				result.selection = pSelection
				result.direction = ResizingDirectionRight
				break
			}
		}
	}
	return result
}

func (c *TimeChart) getMovingSelection(x, y int) *Selection {
	valueX := c.horizontalScale.pixelToValue(x - c.calcLeftScalesWidth())

	for i := 0; i < len(c.selections_); i++ {
		pSelection := c.selections_[i]
		if pSelection.axes == SelectionAxesX {
			if valueX > pSelection.minX && valueX < pSelection.maxX {
				return pSelection
			}
		}
		if pSelection.axes == SelectionAxesXY {
			if valueX > pSelection.minX && valueX < pSelection.maxX {
				return pSelection
			}
		}
	}

	return nil
}

func (c *TimeChart) ZoomShowEntire() {
	c.displayFullRange = true
	/*areaCount := len(c.areas)
	if areaCount > 0 {
		for _, area := range c.areas_ {
			//area.setHeight(defaultAreaHeight());
			//area->zoomShowEntire();
		}
	}*/

	c.horizontalScale.ResetToDefaultRange()
}

func (c *TimeChart) MouseDown(event *uievents.MouseDownEvent) {
	processedByChild := false

	c.lastMouseDownPointX = event.X
	c.SetMouseCursor(ui.MouseCursorArrow)

	area := c.AreaByPoint(event.X, event.Y)
	if area != nil {
		processedByChild = area.MouseDown(event.X, event.Y+c.AreaYOffset(area))
	}

	if processedByChild {
		c.Update("TimeChart")
		return
	}

	mouseLeftButtonPressed := event.Button&uievents.MouseButtonLeft != 0
	mouseRightButtonPressed := event.Button&uievents.MouseButtonRight != 0
	c.mousePressedValueX = c.horizontalScale.pixelToValue(event.X - c.calcLeftScalesWidth())

	shiftKey := event.Modifiers.Shift
	controlKey := event.Modifiers.Control
	altKey := event.Modifiers.Alt

	modifiers := c.buildInteractiveSelectionType(mouseLeftButtonPressed, mouseRightButtonPressed, shiftKey, controlKey, altKey)
	var pAllowedSelection *AllowedSelection

	for i := 0; i < len(c.allowedSelections); i++ {
		if c.allowedSelections[i].interactiveModifiers == modifiers {
			pAllowedSelection = c.allowedSelections[i]
			break
		}
	}

	modeDetected := false

	pResizingSelection := c.getResizingSelection(event.X, event.Y)
	if !modeDetected && pResizingSelection != nil {
		c.resizing = pResizingSelection
		modeDetected = true
	}

	if !modeDetected {
		c.movingSelection = c.getMovingSelection(event.X, event.Y)
		if c.movingSelection != nil {
			c.leftXScaleMovingHandPositionX = c.horizontalScale.pixelToValue(event.X-c.calcLeftScalesWidth()) - c.movingSelection.minX
			modeDetected = true
		}
	}

	if !modeDetected && modifiers == c.chartMoveModifiers {
		c.chartMoving = true
		c.chartMovingHandPositionX = event.X
		c.chartMovingHandPositionY = event.Y
		c.chartMovingDisplayXMin = c.horizontalScale.displayMin_
		c.displayFullRange = false
		modeDetected = true
		c.SetMouseCursor(ui.MouseCursorResizeHor)
		c.Update("TimeChart")
	}

	if !modeDetected {
		if pAllowedSelection != nil {
			c.selections_ = make([]*Selection, 0)
		}

		if pAllowedSelection != nil {

			c.selectionProcessing = &Selection{}
			c.selectionProcessing.interactiveModifiers = pAllowedSelection.interactiveModifiers
			c.selectionProcessing.axes = pAllowedSelection.axes
			c.selectionProcessing.minX = c.horizontalScale.pixelToValue(event.X - c.calcLeftScalesWidth())
			c.selectionProcessing.maxX = c.horizontalScale.pixelToValue(event.X - c.calcLeftScalesWidth())
			modeDetected = true
		}
	}

	c.horizontalScale.MouseDown(event)
}

func MinInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func MaxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func (c *TimeChart) MouseUp(event *uievents.MouseUpEvent) {

	c.SetMouseCursor(ui.MouseCursorArrow)

	processedByChild := false

	area := c.AreaByPoint(event.X, event.Y)
	if area != nil {
		processedByChild = area.MouseUp(event.X, event.Y+c.AreaYOffset(area))
	}

	if processedByChild {
		c.Update("TimeChart")
		return
	}

	if c.resizing != nil {
		c.displayFullRange = false
		c.resizing = nil
	}

	if c.movingSelection != nil {
		c.movingSelection = nil
	}

	if c.chartMoving {
		c.displayFullRange = false
		c.chartMoving = false

		if c.OnMoved != nil {
			c.OnMoved()
		}
	}

	if c.movingPoint != nil {
		c.movingPoint = nil
	}

	if c.selectionProcessing != nil {
		finishValueX := c.horizontalScale.pixelToValue(event.X - c.calcLeftScalesWidth())
		c.selectionProcessing.minX = MinInt64(c.mousePressedValueX, finishValueX)
		c.selectionProcessing.maxX = MaxInt64(c.mousePressedValueX, finishValueX)

		isEmpty := true
		if c.selectionProcessing.axes == SelectionAxesX {
			if c.selectionProcessing.maxX-c.selectionProcessing.minX > 0 {
				isEmpty = false
			}
		}

		if !isEmpty {
			needToAddSelection := true
			if c.selectionProcessing.interactiveModifiers == c.xZoomModifiers {
				needToAddSelection = false
				if math.Abs(float64(c.lastMouseDownPointX-event.X)) > 10 {
					if c.mousePressedValueX > c.horizontalScale.pixelToValue(event.X-c.calcLeftScalesWidth()) {
						c.ZoomShowEntire()
					} else {
						// Zoom X
						c.horizontalScale.SetDisplayRange(c.selectionProcessing.minX, c.selectionProcessing.maxX)
						c.displayFullRange = false
					}
					if c.OnZoomed != nil {
						c.OnZoomed()
					}
				}
			}

			if needToAddSelection {
				c.selections_ = append(c.selections_, c.selectionProcessing)
			}
		}
		c.selectionProcessing = nil
	}

	c.mouseLeftButtonPressed = false
	c.mouseRightButtonPressed = false

	c.Update("TimeChart")
}

func (c *TimeChart) MouseValidateDrop(event *uievents.MouseValidateDropEvent) {
	event.AllowDrop = true
	area := c.AreaByPoint(event.X, event.Y)
	fmt.Println("validate drop", event.Y)

	for _, ar := range c.areas {
		ar.SetHighlighting(false)
	}

	if area != nil {
		if c.OwnWindow.CurrentDraggingObject() != nil {
			area.SetHighlighting(true)
		}
	}
}

func (c *TimeChart) MouseMove(event *uievents.MouseMoveEvent) {
	processedByChild := false
	c.SetMouseCursor(ui.MouseCursorNotDefined)

	//fmt.Println("mouse move ", event.Y)

	area := c.AreaByPoint(event.X, event.Y)

	for _, ar := range c.areas {
		ar.SetHighlighting(false)
	}

	if area != nil {
		processedByChild = area.MouseMove(event.X, event.Y+c.AreaYOffset(area))

		if c.OwnWindow.CurrentDraggingObject() != nil {
			area.SetHighlighting(true)
		}
	}

	if processedByChild {
		c.Update("TimeChart")
		return
	}

	{
		time := c.horizontalScale.pixelToValue(event.X - c.calcLeftScalesWidth())
		// Hover value
		for _, area := range c.areas {
			area.setHoverTimePoint(time)
			area.setHoverPoint(event.X-c.calcLeftScalesWidth(), event.Y-c.AreaYOffset(area))
		}
		c.horizontalScale.setHoverTimePoint(time)
	}

	c.lastMouseX = event.X
	c.lastMouseY = event.Y

	pResizingSelection := c.getResizingSelection(event.X, event.Y)
	if pResizingSelection != nil {
		c.SetMouseCursor(ui.MouseCursorResizeHor)
	}

	if c.chartMoving {
		c.SetMouseCursor(ui.MouseCursorResizeHor)
		x := event.X
		deltaX := c.chartMovingHandPositionX - x
		width := c.horizontalScale.displayMax_ - c.horizontalScale.displayMin_
		deltaXValue := c.horizontalScale.pixelToValueRange(deltaX)
		fmt.Println("chart moving delta", deltaX, (c.chartMovingDisplayXMin+deltaXValue)/1000000, (c.chartMovingDisplayXMin+deltaXValue+width)/1000000)
		c.horizontalScale.SetDisplayRange(c.chartMovingDisplayXMin+deltaXValue, c.chartMovingDisplayXMin+deltaXValue+width)
	}

	if c.resizing != nil {
		c.SetMouseCursor(ui.MouseCursorResizeHor)
		if c.resizing.direction == ResizingDirectionLeft {
			c.resizing.selection.minX = c.horizontalScale.pixelToValue(event.X - c.calcLeftScalesWidth())
			if c.resizing.selection.minX > c.resizing.selection.maxX {
				c.resizing.selection.maxX = c.resizing.selection.minX
				c.resizing.direction = ResizingDirectionRight
			}
		}

		if c.resizing.direction == ResizingDirectionRight {
			c.resizing.selection.maxX = c.horizontalScale.pixelToValue(event.X - c.calcLeftScalesWidth())
			if c.resizing.selection.maxX < c.resizing.selection.minX {
				c.resizing.selection.minX = c.resizing.selection.maxX
				c.resizing.direction = ResizingDirectionLeft
			}
		}
	}

	if c.movingSelection != nil {
		c.SetMouseCursor(ui.MouseCursorResizeHor)
		finishValueX := c.horizontalScale.pixelToValue(event.X - c.calcLeftScalesWidth())
		//longlong finishValueY = timeAxis_.pixelToValue(ev->pos().y() - leftScaleMargin());
		width := c.movingSelection.maxX - c.movingSelection.minX
		newMinX := finishValueX - c.leftXScaleMovingHandPositionX
		newMaxX := newMinX + width

		if newMinX < c.horizontalScale.displayMin_ {
			newMinX = c.horizontalScale.displayMin_
			newMaxX = newMinX + width
		}
		if newMaxX > c.horizontalScale.displayMax_ {
			newMaxX = c.horizontalScale.displayMax_
			newMinX = newMaxX - width
		}

		c.movingSelection.minX = newMinX
		c.movingSelection.maxX = newMaxX
	}

	if c.selectionProcessing != nil {
		finishValueX := c.horizontalScale.pixelToValue(event.X - c.calcLeftScalesWidth())
		c.selectionProcessing.minX = MinInt64(c.mousePressedValueX, finishValueX)
		c.selectionProcessing.maxX = MaxInt64(c.mousePressedValueX, finishValueX)
		c.SetMouseCursor(ui.MouseCursorResizeHor)
	}

	c.Update("TimeChart")
}

func (c *TimeChart) addInteractiveSelectionType(leftMouseButton, rightMouseButton, shiftKey, controlKey, altKey bool, axes SelectionAxes) {
	pAllowedSelectionSelect := &AllowedSelection{}
	pAllowedSelectionSelect.axes = axes
	pAllowedSelectionSelect.interactiveModifiers = c.buildInteractiveSelectionType(leftMouseButton, rightMouseButton, shiftKey, controlKey, altKey)
	c.allowedSelections = append(c.allowedSelections, pAllowedSelectionSelect)
}

func (c *TimeChart) addInteractiveZoom(leftMouseButton, rightMouseButton, shiftKey, controlKey, altKey bool, axes SelectionAxes) {
	c.addInteractiveSelectionType(leftMouseButton, rightMouseButton, shiftKey, controlKey, altKey, axes)
	modifiers := c.buildInteractiveSelectionType(leftMouseButton, rightMouseButton, shiftKey, controlKey, altKey)
	switch axes {
	case SelectionAxesX:
		c.xZoomModifiers = modifiers
		break
	case SelectionAxesY:
		c.yZoomModifiers = modifiers
		break
	case SelectionAxesXY:
		c.xyZoomModifiers = modifiers
		break
	}
}

func (c *TimeChart) MouseEnter() {
	c.mouseIsInside = true
}

func (c *TimeChart) MouseLeave() {
	c.mouseIsInside = false
}

func (c *TimeChart) MouseIsInside() bool {
	return c.mouseIsInside
}

func (c *TimeChart) MinHeight() int {
	return c.horizontalScale.Height + 200 // HorScale + 200 for chart
}

/*func (c *TimeChart) MinWidth() int {
	return 300
}*/
