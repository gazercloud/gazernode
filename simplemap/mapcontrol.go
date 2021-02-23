package simplemap

import (
	"fmt"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uiproperties"
	"golang.org/x/image/colornames"
	"image/color"
	"math"
)

type MapControlPoint struct {
	x         int32
	y         int32
	radius    int32
	highlight bool
}

type RectF struct {
	x      float64
	y      float64
	width  float64
	height float64
}

type Rect32 struct {
	x      int32
	y      int32
	width  int32
	height int32
}

func (c *RectF) contains(x, y float64) bool {
	if x > c.x && x < c.x+c.width && y > c.y && y < c.y+c.height {
		return true
	}
	return false
}

func (c *Rect32) contains(x, y int32) bool {
	if x > c.x && x < c.x+c.width && y > c.y && y < c.y+c.height {
		return true
	}
	return false
}

type MapControlVertex int

const (
	MapControlVertexNoVertex    MapControlVertex = 0
	MapControlVertexLeftTop     MapControlVertex = 1
	MapControlVertexRightTop    MapControlVertex = 2
	MapControlVertexRightBottom MapControlVertex = 3
	MapControlVertexLeftBottom  MapControlVertex = 4
	MapControlVertex1           MapControlVertex = 5
	MapControlVertex2           MapControlVertex = 6
)

func NewMapControlPoint() *MapControlPoint {
	var c MapControlPoint
	c.radius = 5
	return &c
}

func (c *MapControlPoint) draw(control *MapControl, ctx ui.DrawContext) {
	col := colornames.Green
	if c.highlight {
		col = colornames.Red
	}
	var x int32
	var y int32

	if control.imapControl.adaptiveSize() {
		// Normal scaling
		x = control.scaleValue(c.x)
		y = control.scaleValue(c.y)
	} else {
		// Special for scaled view
		x = control.scaleValueIndependentContent(c.x)
		y = control.scaleValueIndependentContent(c.y)
	}

	radius := int32(5)

	ctx.SetColor(col)
	ctx.FillRect(int(x-radius), int(y-radius), int(radius*2), int(radius*2))

	if c.highlight {
		ctx.SetColor(col)
	} else {
		ctx.SetColor(colornames.Yellow)
	}
	ctx.SetStrokeWidth(1)
	ctx.DrawRect(int(x-radius), int(y-radius), int(radius*2), int(radius*2))

}

type MapControl struct {
	uiproperties.PropertiesContainer

	selected_ bool

	scale_                   float64
	scaleIndependentContent_ float64
	selectedExclusive_       bool
	isRootControl_           bool

	imapControl IMapControl

	// Size points
	points_           []*MapControlPoint
	pointTopLeft_     *MapControlPoint
	pointTopRight_    *MapControlPoint
	pointBottomLeft_  *MapControlPoint
	pointBottomRight_ *MapControlPoint

	x_          *uiproperties.Property
	y_          *uiproperties.Property
	width_      *uiproperties.Property
	height_     *uiproperties.Property
	name_       *uiproperties.Property
	anchors_    *uiproperties.Property
	type_       *uiproperties.Property
	dataSource_ *uiproperties.Property

	layer_ *MapControlViewLayer

	mouseDownPointInches_ Point32
	mouseMovePointInches_ Point32
	mouseUpPointInches_   Point32

	posOfControlAtLastMouseDown_  Point32
	sizeOfControlAtLastMouseDown_ Point32

	modeResizeItem_   bool
	resizeItemVertex_ MapControlVertex

	resizeItemOriginalPos_  Point32
	resizeItemOriginalSize_ Point32

	parent_          IMapControl
	original_x_      int32
	original_y_      int32
	original_width_  int32
	original_height_ int32

	mapDataSource IMapDataSource
	MapWidget     *MapWidget

	actionElapsed bool

	// CurrentValue
	value common_interfaces.ItemValue
}

func AddPropertyToControl(c IMapControl, name string, displayName string, propertyType uiproperties.PropertyType, groupName string, subType string) *uiproperties.Property {
	if c == nil {
		panic("No control for property " + name)
	}
	p := uiproperties.NewProperty(name, propertyType)
	p.DisplayName = displayName
	p.GroupName = groupName
	p.SubType = subType
	p.Init(name, c)
	c.(uiproperties.IPropertiesContainer).AddProperty(name, p)
	return p
}

func (c *MapControl) Dispose() {
	c.MapWidget = nil
	c.mapDataSource = nil
	c.parent_ = nil
}

func (c *MapControl) Name() string {
	return c.name_.String()
}

func (c *MapControl) setMapDataSource(mapDataSource IMapDataSource) {
	c.mapDataSource = mapDataSource
}

func (c *MapControl) initMapControl(control IMapControl, mapWidget *MapWidget, parent IMapControl) {
	c.scale_ = 1
	c.parent_ = parent
	c.MapWidget = mapWidget
	c.scaleIndependentContent_ = 1
	c.InitPropertiesContainer()
	c.PropertiesContainer.OnPropertyChanged = c.OnPropertyChanged
	c.imapControl = control

	c.name_ = AddPropertyToControl(control, "name", "Name", uiproperties.PropertyTypeString, "Common", "")
	//c.name_.SetVisible(false)
	c.type_ = AddPropertyToControl(control, "type", "Type", uiproperties.PropertyTypeString, "Common", "")
	c.type_.SetVisible(false)
	c.x_ = AddPropertyToControl(control, "x", "X", uiproperties.PropertyTypeInt32, "Position", "")
	c.y_ = AddPropertyToControl(control, "y", "Y", uiproperties.PropertyTypeInt32, "Position", "")
	c.width_ = AddPropertyToControl(control, "width", "Width", uiproperties.PropertyTypeInt32, "Position", "")
	c.height_ = AddPropertyToControl(control, "height", "Height", uiproperties.PropertyTypeInt32, "Position", "")
	c.anchors_ = AddPropertyToControl(control, "anchors", "Anchors", uiproperties.PropertyTypeString, "Position", "")
	c.anchors_.SetVisible(false)

	c.x_.OnChanged = func(property *uiproperties.Property, oldValue interface{}, newValue interface{}) {
		control.UpdateSizePoints()
	}
	c.y_.OnChanged = func(property *uiproperties.Property, oldValue interface{}, newValue interface{}) {
		control.UpdateSizePoints()
	}
	c.width_.OnChanged = func(property *uiproperties.Property, oldValue interface{}, newValue interface{}) {
		control.UpdateSizePoints()
	}
	c.height_.OnChanged = func(property *uiproperties.Property, oldValue interface{}, newValue interface{}) {
		control.UpdateSizePoints()
	}
}

func (c *MapControl) draw(ctx ui.DrawContext) {
	c.imapControl.drawControl(ctx)
	if c.selectedExclusive_ {
		c.drawControlPoints(ctx)
	}

	if c.selected_ {
		c.drawSelection(ctx)
	}

	if c.actionElapsed {
		ctx.SetColor(color.RGBA{
			R: 255,
			G: 255,
			B: 0,
			A: 100,
		})
		ctx.FillRect(0, 0, int(c.scaleValue(c.Width())), int(c.scaleValue(c.Height())))
		c.actionElapsed = false
	}
}

func (c *MapControl) drawControlPoints(ctx ui.DrawContext) {
	for _, point := range c.points_ {
		point.draw(c, ctx)
	}
}

func (c *MapControl) GetFullPathToMapControl() []string {
	result := make([]string, 0)
	if c.parent_ != nil {
		result = append(result, c.parent_.GetFullPathToMapControl()...)
	}
	result = append(result, c.type_.String())
	return result
}

func (c *MapControl) drawSelection(ctx ui.DrawContext) {
	return
	widthOfLine := 1

	selectedRectDistance := 5
	var x int
	var y int
	var w int
	var h int

	if c.imapControl.adaptiveSize() {
		// Normal scaling
		w = int(c.scaleValue(c.Width()))
		h = int(c.scaleValue(c.Height()))
	} else {
		// Special for scaled view
		w = int(c.scaleValueIndependentContent(c.Width()))
		h = int(c.scaleValueIndependentContent(c.Height()))
	}
	ctx.SetColor(colornames.Blue)
	ctx.SetStrokeWidth(int(widthOfLine))
	ctx.DrawRect(x-selectedRectDistance, y-selectedRectDistance, w+selectedRectDistance*2, h+selectedRectDistance*2)
}

func (c *MapControl) NeedTranslateOnDraw() bool {
	return true
}

func (c *MapControl) scaleValue(value int32) int32 {
	return int32(float64(value) * c.scale_)
}

func (c *MapControl) unscaleValue(value int32) int32 {
	if c.scale_ > -0.0000001 && c.scale_ < 0.0000001 {
		return 0
	}
	return int32(float64(value) / c.scale_)
}

func (c *MapControl) scaleValueFloat(value int32) float64 {
	return float64(value) * c.scale_
}

func (c *MapControl) scaleValueIndependentContent(value int32) int32 {
	return int32(float64(value) * c.scaleIndependentContent_)
}

func (c *MapControl) scaleIndependentContent() float64 {
	return c.scaleIndependentContent_
}

func (c *MapControl) setScaleIndependentContent(scale float64) {
	c.scaleIndependentContent_ = scale
}

func (c *MapControl) X() int32 {
	return c.x_.Int32()
}

func (c *MapControl) SetX(x int32) {
	if c.x_.Int32() != x {
		c.x_.SetOwnValue(x)
		c.imapControl.UpdateSizePoints()
	}
}

func (c *MapControl) Y() int32 {
	return c.y_.Int32()
}

func (c *MapControl) SetY(y int32) {
	if c.y_.Int32() != y {
		c.y_.SetOwnValue(y)
		c.imapControl.UpdateSizePoints()
	}
}

func (c *MapControl) Width() int32 {
	return c.width_.Int32()
}

func (c *MapControl) SetWidth(width int32) {
	if c.width_.Int32() != width {
		c.width_.SetOwnValue(width)
	}
	c.imapControl.UpdateSizePoints()
}

func (c *MapControl) Height() int32 {
	return c.height_.Int32()
}

func (c *MapControl) SetHeight(height int32) {

	if height == 0 {
		fmt.Println(height)
	}

	if c.height_.Int32() != height {
		c.height_.SetOwnValue(height)
	}
	c.imapControl.UpdateSizePoints()
}

func (c *MapControl) name() string {
	return c.name_.String()
}

func (c *MapControl) setName(name string) {
	c.name_.SetOwnValue(name)
}

func (c *MapControl) dataSource() string {
	return c.dataSource_.String()
}

func (c *MapControl) SetDataSource(dataSource string) {
	c.dataSource_.SetOwnValue(dataSource)
	c.NotifyChangedToContainer(c.dataSource_)
	//dataItemSource_.reset();
}

func (c *MapControl) anchors() int32 {
	return c.anchors_.Int32()
}

func (c *MapControl) setAnchors(anchors int) {
	c.anchors_.SetOwnValue(anchors)
}

func (c *MapControl) UpdateSizePoints() {
	if c.isRootControl_ {
		return
	}

	if c.pointTopLeft_ == nil {
		c.pointTopLeft_ = NewMapControlPoint()
		c.points_ = append(c.points_, c.pointTopLeft_)
	}

	if c.pointTopRight_ == nil {
		c.pointTopRight_ = NewMapControlPoint()
		c.points_ = append(c.points_, c.pointTopRight_)
	}

	if c.pointBottomLeft_ == nil {
		c.pointBottomLeft_ = NewMapControlPoint()
		c.points_ = append(c.points_, c.pointBottomLeft_)
	}

	if c.pointBottomRight_ == nil {
		c.pointBottomRight_ = NewMapControlPoint()
		c.points_ = append(c.points_, c.pointBottomRight_)
	}

	c.pointTopLeft_.x = 0
	c.pointTopLeft_.y = 0

	c.pointTopRight_.x = c.Width()
	c.pointTopRight_.y = 0

	c.pointBottomLeft_.x = 0
	c.pointBottomLeft_.y = c.Height()

	c.pointBottomRight_.x = c.Width()
	c.pointBottomRight_.y = c.Height()
}

func (c *MapControl) setScale(scale float64) {
	c.scale_ = scale
}

func (c *MapControl) scale() float64 {
	return c.scale_
}

func (c *MapControl) setSelected(selected bool) {
	if c.layer_ != nil && selected {
		if !c.layer_.visible_ {
			return
		}
		if !selected {
			c.selectedExclusive_ = false
		}
	}
	c.selected_ = selected
}

func (c *MapControl) selectedExclusive() bool {
	return c.selectedExclusive_
}

func (c *MapControl) setSelectedExclusive(selectedExclusive bool) {
	c.selectedExclusive_ = selectedExclusive
}

func (c *MapControl) rememberLastPosition() {
	c.posOfControlAtLastMouseDown_ = Point32{c.X(), c.Y()}
	c.sizeOfControlAtLastMouseDown_ = Point32{c.Width(), c.Height()}
}

func (c *MapControl) pointUnderPoint(x0, y0 int32) *MapControlPoint {
	for _, point := range c.points_ {
		rectOfPoint := Rect32{c.X() + point.x - point.radius, c.Y() + point.y - point.radius, point.radius * 2, point.radius * 2}
		if rectOfPoint.contains(x0, y0) {
			return point
		}
	}
	return nil
}

func (c *MapControl) isPointInside(x1, y1 int32) bool {
	left := c.X()
	right := c.X() + c.Width()
	top := c.Y()
	bottom := c.Y() + c.Height()
	if x1 >= left && x1 <= right && y1 >= top && y1 <= bottom {
		return true
	}
	return false
}

func (c *MapControl) mouseDown(x0, y0 int32, leftButton, centerButton, rightButton, shift, control, alt bool) bool {
	c.posOfControlAtLastMouseDown_ = Point32{c.X(), c.Y()}

	if c.selected_ && c.selectedExclusive_ {
		pointUnderMousePoint := c.pointUnderPoint(x0, y0)
		if pointUnderMousePoint != nil {
			if c.pointTopLeft_ != nil && c.pointTopRight_ != nil && c.pointBottomLeft_ != nil && c.pointBottomRight_ != nil {
				c.rememberLastPosition()
				c.modeResizeItem_ = true
				c.mouseDownPointInches_ = Point32{x0, y0}
				c.resizeItemOriginalPos_ = Point32{c.X(), c.Y()}
				c.resizeItemOriginalSize_ = Point32{c.Width(), c.Height()}
				if pointUnderMousePoint == c.pointTopLeft_ {
					c.resizeItemVertex_ = MapControlVertexLeftTop
				}
				if pointUnderMousePoint == c.pointTopRight_ {
					c.resizeItemVertex_ = MapControlVertexRightTop
				}
				if pointUnderMousePoint == c.pointBottomLeft_ {
					c.resizeItemVertex_ = MapControlVertexLeftBottom
				}
				if pointUnderMousePoint == c.pointBottomRight_ {
					c.resizeItemVertex_ = MapControlVertexRightBottom
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

func (c *MapControl) TypeName() string {
	return "control"
}

func (c *MapControl) drawControl(ctx ui.DrawContext) {
}

func (c *MapControl) OnMouseDown(x, y int) {
}

func (c *MapControl) FindControlUnderPoint(x, y int) IMapControl {
	return c.imapControl
}

func (c *MapControl) HasAction() bool {
	return false
}

func (c *MapControl) align(value int32) int32 {
	typPointsAlign := int32(10)
	newValue := math.Round(float64(value) / float64(typPointsAlign))
	return int32(newValue) * typPointsAlign
}

func (c *MapControl) saveSizesAsOriginal() {
	if c.parent_ != nil {
		if c.imapControl.isView() && c.imapControl.isRootControl() && c.imapControl.isEditing() {
			c.original_x_ = c.X()
			c.original_y_ = c.Y()
			c.original_width_ = c.Width()
			c.original_height_ = c.Height()
		}

	}
}

func (c *MapControl) SetOriginalX(x int32) {
	c.original_x_ = x
}

func (c *MapControl) SetOriginalY(y int32) {
	c.original_y_ = y
}

func (c *MapControl) SetOriginalWidth(width int32) {
	c.original_width_ = width
}

func (c *MapControl) SetOriginalHeight(height int32) {
	c.original_height_ = height
}

func (c *MapControl) mouseMove(lastMouseDownPos, pos Point32, leftButton, centerButton, rightButton, shift, control, alt bool) bool {

	needToUpdatePropertiesContainer := false

	c.mouseMovePointInches_ = pos
	if c.selected_ && leftButton {
		delta := Point32{pos.x - lastMouseDownPos.x, pos.y - lastMouseDownPos.y}

		newX := c.posOfControlAtLastMouseDown_.x + delta.x
		newY := c.posOfControlAtLastMouseDown_.y + delta.y
		newX = c.align(newX)
		newY = c.align(newY)

		if newX != c.X() || newY != c.Y() {
			c.SetX(newX)
			c.SetY(newY)
			c.saveSizesAsOriginal()
			needToUpdatePropertiesContainer = true
		}
	}

	for _, point := range c.points_ {
		point.highlight = false
	}

	if c.modeResizeItem_ {

		X0 := c.resizeItemOriginalPos_.x
		Y0 := c.resizeItemOriginalPos_.y
		W0 := c.resizeItemOriginalSize_.x
		H0 := c.resizeItemOriginalSize_.y
		W1 := W0
		H1 := H0
		X1 := X0
		Y1 := Y0

		deltaX := c.mouseMovePointInches_.x - c.mouseDownPointInches_.x
		deltaY := c.mouseMovePointInches_.y - c.mouseDownPointInches_.y

		deltaX = c.align(deltaX)
		deltaY = c.align(deltaY)

		if c.resizeItemVertex_ == MapControlVertexRightBottom || c.resizeItemVertex_ == MapControlVertexRightTop {
			W1 = c.align(W0 + deltaX)
		} else {
			X1 = c.align(X0 + deltaX)
			W1 = (X0 + W0) - X1
		}

		if c.resizeItemVertex_ == MapControlVertexRightTop || c.resizeItemVertex_ == MapControlVertexLeftTop {
			Y1 = c.align(Y0 + deltaY)
			H1 = (Y0 + H0) - Y1
		} else {
			H1 = c.align(H0 + deltaY)
		}

		if W1 < 0 {
			W1 = int32(math.Abs(float64(W1)))
			X1 = X1 - W1
		}

		if H1 < 0 {
			H1 = int32(math.Abs(float64(H1)))
			Y1 = Y1 - H1
		}

		if X1 != c.X() || W1 != c.Width() || Y1 != c.Y() || H1 != c.Height() {
			fmt.Println(X1, c.X(), Y1, c.Y(), W1, c.Width(), H1, c.Height())
			c.SetX(X1)
			c.SetWidth(W1)
			c.SetY(Y1)
			c.SetHeight(H1)

			c.imapControl.updateLayout(true)

			c.saveSizesAsOriginal()

			c.imapControl.refreshScale()

			// Highlight of resize point (depends of resize type)
			switch c.resizeItemVertex_ {
			case MapControlVertexLeftTop:
				c.pointTopLeft_.highlight = true
			case MapControlVertexRightTop:
				c.pointTopRight_.highlight = true
			case MapControlVertexRightBottom:
				c.pointBottomRight_.highlight = true
			case MapControlVertexLeftBottom:
				c.pointBottomLeft_.highlight = true
			}

			needToUpdatePropertiesContainer = true
		}

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

func (c *MapControl) mouseUp(x0, y0 int32, leftButton, centerButton, rightButton, shift, control, alt bool) bool {
	if c.selected_ {
		X_changed := c.posOfControlAtLastMouseDown_.x != c.X()
		Y_changed := c.posOfControlAtLastMouseDown_.y != c.Y()
		W_changed := c.sizeOfControlAtLastMouseDown_.x != c.Width()
		H_changed := c.sizeOfControlAtLastMouseDown_.y != c.Height()

		if X_changed || Y_changed || W_changed || H_changed {
			changes := uiproperties.NewPropertiesChangesList()
			changes.AddItem(c, c.x_.Name, c.X())
			changes.AddItem(c, c.y_.Name, c.Y())
			changes.AddItem(c, c.width_.Name, c.Width())
			changes.AddItem(c, c.height_.Name, c.Height())
			c.imapControl.changeNotify(changes)
		}
	}

	c.modeResizeItem_ = false
	c.rememberLastPosition()
	c.posOfControlAtLastMouseDown_ = Point32{c.X(), c.Y()}
	return false
}

func (c *MapControl) selected() bool {
	return c.selected_
}

func (c *MapControl) isRectIntersect(x1, y1, x2, y2 int32) bool {
	left := c.X()
	right := c.X() + c.Width()
	top := c.Y()
	bottom := c.Y() + c.Height()

	if left >= x1 && left <= x2 && top >= y1 && top <= y2 {
		return true
	}

	if left >= x1 && left <= x2 && bottom >= y1 && bottom <= y2 {
		return true
	}
	if right >= x1 && right <= x2 && top >= y1 && top <= y2 {
		return true
	}
	if right >= x1 && right <= x2 && bottom >= y1 && bottom <= y2 {
		return true
	}
	return false
}

func (c *MapControl) isEditing() bool {
	return false
}

func (c *MapControl) isRootControl() bool {
	return c.isRootControl_
}

func (c *MapControl) isView() bool {
	return false
}

func (c *MapControl) refreshScale() {
	c.setScale(c.scale())
}

// left bottom right top
// left = 1100 = 0x0C = 12
// left top bottom = 1101 = 13

func (c *MapControl) anchorBottom() bool {
	return c.anchors_.Int32()&0x04 != 0
}

func (c *MapControl) anchorLeft() bool {
	return c.anchors_.Int32()&0x08 != 0
}

func (c *MapControl) anchorRight() bool {
	return c.anchors_.Int32()&0x02 != 0
}

func (c *MapControl) anchorTop() bool {
	return c.anchors_.Int32()&0x01 != 0
}

func (c *MapControl) changeNotify(list *uiproperties.PropertiesChangesList) {
	topLevelOrSecondLevel := false
	if c.isRootControl_ {
		topLevelOrSecondLevel = true
	}
	if c.parent_ != nil && c.parent_.isRootControl() {
		topLevelOrSecondLevel = true
	}
	if topLevelOrSecondLevel {
		if c.MapWidget != nil {
			c.MapWidget.NotifyChanges(list)
		}
	}
}

func (c *MapControl) original_x() int32 {
	return c.original_x_
}

func (c *MapControl) original_y() int32 {
	return c.original_y_
}

func (c *MapControl) original_width() int32 {
	return c.original_width_
}

func (c *MapControl) original_height() int32 {
	return c.original_height_
}

func (c *MapControl) updateLayout(direct bool) {
}

func (c *MapControl) Subclass() string {
	return "default"
}

func (c *MapControl) adaptiveSize() bool {
	return false
}

func (c *MapControl) setNeedToSetDefaultSize(needToSetDefaultSize bool) {

}

func (c *MapControl) saveBase() *MapItem {
	result := NewMapItem()

	for _, prop := range c.GetProperties() {
		result.Props = append(result.Props, prop.SaveToStruct())
	}

	return result
}

func (c *MapControl) load(m *MapItem) {
	if m != nil {
		for _, prop := range c.GetProperties() {
			for _, p := range m.Props {
				if p.Name == prop.Name {
					prop.SetOwnValue(p.Value)
				}
			}
		}
	}

	c.imapControl.UpdateSizePoints()
}

func (c *MapControl) OnPropertyChanged(prop *uiproperties.Property) {
	c.imapControl.UpdateSizePoints()

	changes := uiproperties.NewPropertiesChangesList()
	changes.AddItem(c, prop.Name, prop.Value())
	c.imapControl.changeNotify(changes)
}

func (c *MapControl) UpdateValue(value common_interfaces.ItemValue) {
	c.value = value
}

func (c *MapControl) SetAdding() {
}

func (c *MapControl) LoadContent(contentBytes []byte, err error) {
}

func (c *MapControl) Tick() {
}

func (c *MapControl) SetType(typeName string) {
	c.type_.SetOwnValue(typeName)
}

func (c *MapControl) FullDataSource() string {
	result := c.dataSource()
	if c.parent_ != nil {
		parentDataSource := c.parent_.FullDataSource()
		if len(parentDataSource) > 0 {
			result = parentDataSource + "/" + result
		}
	}
	return result
}
