package simplemap

import (
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uiproperties"
	"golang.org/x/image/colornames"
	"math"
)

type MapControlLine struct {
	MapControl

	x1_ *uiproperties.Property
	y1_ *uiproperties.Property
	x2_ *uiproperties.Property
	y2_ *uiproperties.Property

	point1 *MapControlPoint
	point2 *MapControlPoint

	posOfPoint1AtLastMouseDown_ Point32
	posOfPoint2AtLastMouseDown_ Point32

	color_     *uiproperties.Property
	lineWidth_ *uiproperties.Property

	edges_ *uiproperties.Property

	loadingCoordinates bool
}

func NewMapControlLine(mapWidget *MapWidget, parent IMapControl) *MapControlLine {
	var c MapControlLine
	c.initMapControl(&c, mapWidget, parent)

	//c.points_ = make([]*MapControlPoint, 0)

	c.type_.SetOwnValue("line")

	c.x1_ = AddPropertyToControl(&c, "x1", "X1", uiproperties.PropertyTypeInt32, "Location", "")
	c.y1_ = AddPropertyToControl(&c, "y1", "Y1", uiproperties.PropertyTypeInt32, "Location", "")
	c.x2_ = AddPropertyToControl(&c, "x2", "X2", uiproperties.PropertyTypeInt32, "Location", "")
	c.y2_ = AddPropertyToControl(&c, "y2", "Y2", uiproperties.PropertyTypeInt32, "Location", "")

	c.x1_.SetOwnValue(int32(10))
	c.y1_.SetOwnValue(int32(10))
	c.x2_.SetOwnValue(int32(100))
	c.y2_.SetOwnValue(int32(100))

	c.x1_.OnChanged = c.OnCoordinatesChanged
	c.y1_.OnChanged = c.OnCoordinatesChanged
	c.x2_.OnChanged = c.OnCoordinatesChanged
	c.y2_.OnChanged = c.OnCoordinatesChanged

	c.UpdateXYWH()

	c.color_ = AddPropertyToControl(&c, "color", "Color", uiproperties.PropertyTypeColor, "Line", "")
	c.color_.SetOwnValue(colornames.Gray)

	c.lineWidth_ = AddPropertyToControl(&c, "line_width", "Width", uiproperties.PropertyTypeInt32, "Line", "")
	c.lineWidth_.SetOwnValue(int32(4))

	c.edges_ = AddPropertyToControl(&c, "edges", "Edges", uiproperties.PropertyTypeString, "Line", "edges")
	c.edges_.SetOwnValue("square")
	c.edges_.DefaultValue = c.edges_.ValueOwn()

	return &c
}

func (c *MapControlLine) drawControl(ctx ui.DrawContext) {
	x1 := int(c.scaleValue(c.X1() - c.X()))
	y1 := int(c.scaleValue(c.Y1() - c.Y()))
	x2 := int(c.scaleValue(c.X2() - c.X()))
	y2 := int(c.scaleValue(c.Y2() - c.Y()))

	cc := ctx.GG()
	cc.Push()

	cc.SetLineCapSquare()

	if c.edges_.String() == "square" {
		cc.SetLineCapSquare()
	}
	if c.edges_.String() == "round" {
		cc.SetLineCapRound()
	}
	cc.Translate(float64(ctx.State().TranslateX), float64(ctx.State().TranslateY))

	if c.selected_ {
		cc.SetColor(colornames.Crimson)
	} else {
		cc.SetColor(c.color_.Color())
	}

	cc.SetColor(c.color_.Color())
	cc.SetLineWidth(float64(c.scaleValue(c.lineWidth_.Int32())))
	cc.DrawLine(float64(x1), float64(y1), float64(x2), float64(y2))
	cc.Stroke()
	cc.Pop()
}

func (c *MapControlLine) X1() int32 {
	return c.x1_.ValueOwn().(int32)
}

func (c *MapControlLine) Y1() int32 {
	return c.y1_.ValueOwn().(int32)
}

func (c *MapControlLine) X2() int32 {
	return c.x2_.ValueOwn().(int32)
}

func (c *MapControlLine) Y2() int32 {
	return c.y2_.ValueOwn().(int32)
}

func (c *MapControlLine) SetX1(x1 int32) {
	c.x1_.SetOwnValue(x1)
}

func (c *MapControlLine) SetY1(y1 int32) {
	c.y1_.SetOwnValue(y1)
}

func (c *MapControlLine) SetX2(x2 int32) {
	c.x2_.SetOwnValue(x2)
}

func (c *MapControlLine) SetY2(y2 int32) {
	c.y2_.SetOwnValue(y2)
}

func (c *MapControlLine) UpdateXYWH() {
	x := int32(math.Min(float64(c.X1()), float64(c.X2())))
	y := int32(math.Min(float64(c.Y1()), float64(c.Y2())))
	width := int32(math.Abs(float64(c.X1() - c.X2())))
	height := int32(math.Abs(float64(c.Y1() - c.Y2())))
	c.SetX(x)
	c.SetY(y)
	c.SetWidth(width)
	c.SetHeight(height)
}

func (c *MapControlLine) OnCoordinatesChanged(property *uiproperties.Property, oldValue interface{}, newValue interface{}) {
	if c.loadingCoordinates {
		return
	}
	c.UpdateXYWH()
}

func (c *MapControlLine) TypeName() string {
	return "line"
}

func (c *MapControlLine) adaptiveSize() bool {
	return true
}

func (c *MapControlLine) UpdateSizePoints() {
	if c.isRootControl_ {
		return
	}

	if c.point1 == nil {
		c.point1 = NewMapControlPoint()
		c.points_ = append(c.points_, c.point1)
	}

	if c.point2 == nil {
		c.point2 = NewMapControlPoint()
		c.points_ = append(c.points_, c.point2)
	}

	c.point1.x = c.X1() - c.X()
	c.point1.y = c.Y1() - c.Y()

	c.point2.x = c.X2() - c.X()
	c.point2.y = c.Y2() - c.Y()
}

func (c *MapControlLine) mouseDown(x0, y0 int32, leftButton, centerButton, rightButton, shift, control, alt bool) bool {
	//x0 = x0 + c.X()
	//y0 = y0 + c.Y()
	c.posOfControlAtLastMouseDown_ = Point32{c.X(), c.Y()}
	c.posOfPoint1AtLastMouseDown_ = Point32{c.X1(), c.Y1()}
	c.posOfPoint2AtLastMouseDown_ = Point32{c.X2(), c.Y2()}

	if c.selected_ && c.selectedExclusive_ {
		pointUnderMousePoint := c.pointUnderPoint(x0, y0)
		if pointUnderMousePoint != nil {
			if c.point1 != nil && c.point2 != nil {
				c.rememberLastPosition()
				c.modeResizeItem_ = true
				c.mouseDownPointInches_ = Point32{x0, y0}
				if pointUnderMousePoint == c.point1 {
					c.resizeItemVertex_ = MapControlVertex1
				}
				if pointUnderMousePoint == c.point2 {
					c.resizeItemVertex_ = MapControlVertex2
				}
				return true
			}
		}
	}

	if c.isPointInside(x0, y0) {
		return true
	}

	return false
}

func (c *MapControlLine) mouseMove(lastMouseDownPos, pos Point32, leftButton, centerButton, rightButton, shift, control, alt bool) bool {
	//pos = Point32{pos.x + c.X(), pos.y + c.Y()}

	needToUpdatePropertiesContainer := false

	c.mouseMovePointInches_ = pos
	if c.selected_ && leftButton && !c.modeResizeItem_ {
		delta := Point32{pos.x - lastMouseDownPos.x, pos.y - lastMouseDownPos.y}

		newX1 := c.posOfPoint1AtLastMouseDown_.x + delta.x
		newY1 := c.posOfPoint1AtLastMouseDown_.y + delta.y
		c.SetX1(c.align(newX1))
		c.SetY1(c.align(newY1))

		newX2 := c.posOfPoint2AtLastMouseDown_.x + delta.x
		newY2 := c.posOfPoint2AtLastMouseDown_.y + delta.y
		c.SetX2(c.align(newX2))
		c.SetY2(c.align(newY2))

		c.saveSizesAsOriginal()

		needToUpdatePropertiesContainer = true
	}

	for _, point := range c.points_ {
		point.highlight = false
	}

	if c.modeResizeItem_ {

		deltaX := c.mouseMovePointInches_.x - c.mouseDownPointInches_.x
		deltaY := c.mouseMovePointInches_.y - c.mouseDownPointInches_.y

		deltaX = c.align(deltaX)
		deltaY = c.align(deltaY)

		if c.resizeItemVertex_ == MapControlVertex1 {
			c.SetX1(c.align(pos.x))
			c.SetY1(c.align(pos.y))
		}

		if c.resizeItemVertex_ == MapControlVertex2 {
			c.SetX2(c.align(pos.x))
			c.SetY2(c.align(pos.y))
		}

		c.imapControl.updateLayout(true)

		c.saveSizesAsOriginal()

		c.imapControl.refreshScale()

		// Highlight of resize point (depends of resize type)
		switch c.resizeItemVertex_ {
		case MapControlVertex1:
			c.point1.highlight = true
		case MapControlVertex2:
			c.point2.highlight = true
		}

		needToUpdatePropertiesContainer = true

	} else {

		// Highlight of resize point (depends of mouse pointer)
		pointUnderMousePoint := c.pointUnderPoint(pos.x, pos.y)
		if pointUnderMousePoint != nil {
			pointUnderMousePoint.highlight = true
		}
	}

	if needToUpdatePropertiesContainer {
		c.NotifyChangedToContainer(nil)
	}

	return false
}

func (c *MapControlLine) mouseUp(x0, y0 int32, leftButton, centerButton, rightButton, shift, control, alt bool) bool {
	x0 = x0 + c.X()
	y0 = y0 + c.Y()

	if c.selected_ {
		X_changed := c.posOfControlAtLastMouseDown_.x != c.X()
		Y_changed := c.posOfControlAtLastMouseDown_.y != c.Y()
		W_changed := c.sizeOfControlAtLastMouseDown_.x != c.Width()
		H_changed := c.sizeOfControlAtLastMouseDown_.y != c.Height()

		if X_changed || Y_changed || W_changed || H_changed {
			//c.imapControl.changeNotify(changeText)
		}
	}

	c.modeResizeItem_ = false
	c.rememberLastPosition()
	c.posOfControlAtLastMouseDown_ = Point32{c.X(), c.Y()}
	return false
}

func (c *MapControlLine) pointUnderPoint(x0, y0 int32) *MapControlPoint {
	x0 -= c.X()
	y0 -= c.Y()
	for _, point := range c.points_ {
		rectOfPoint := Rect32{point.x - point.radius, point.y - point.radius, point.radius * 2, point.radius * 2}
		if rectOfPoint.contains(x0, y0) {
			return point
		}
	}
	return nil
}

func (c *MapControlLine) isPointInside(x0, y0 int32) bool {
	x1 := float64(c.X1())
	x2 := float64(c.X2())
	y1 := float64(c.Y1())
	y2 := float64(c.Y2())

	// In rectangle of line - for performance
	rx1 := x1
	rx2 := x2
	ry1 := y1
	ry2 := y2
	if rx1 > rx2 {
		rx1, rx2 = rx2, rx1
	}
	if y1 > y2 {
		ry1, ry2 = ry2, ry1
	}
	inRectangle := false
	distForRect := float64(c.lineWidth_.Int32())
	if distForRect < 5 {
		distForRect = 5
	}
	if float64(x0) >= rx1-distForRect && float64(x0) <= rx2+distForRect && float64(y0) >= ry1-distForRect && float64(y0) <= ry2+distForRect {
		inRectangle = true
	}
	//inRectangle = true

	// Determine distance to line
	nearLine := false
	if inRectangle {
		d1 := math.Abs((y2-y1)*float64(x0) - (x2-x1)*float64(y0) + x2*y1 - y2*x1)
		d2 := math.Sqrt((y2-y1)*(y2-y1) + (x2-x1)*(x2-x1))
		distance := d1 / d2
		allowedDistance := c.lineWidth_.Int32() / 2
		if allowedDistance < 5 {
			allowedDistance = 5
		}
		if c.scaleValue(int32(distance)) < allowedDistance {
			nearLine = true
		}
	}

	return nearLine
}
