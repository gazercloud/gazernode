package simplemap

import (
	"errors"
	"fmt"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazerui/go-gl/glfw/v3.3/glfw"
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiproperties"
	"github.com/nfnt/resize"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"math"
	"time"
)

type Point struct {
	x int
	y int
}

type Point32 struct {
	x int32
	y int32
}

type PointF struct {
	x float64
	y float64
}

type SnapshotItem struct {
	Name string
	Ds   []byte
}

type MapWidget struct {
	uicontrols.Control
	view_ *MapControlView

	showGrid_ bool
	//editMode_    bool
	modeZoomMap_ bool
	timerZoom_   *uievents.FormTimer

	errors_ []string

	mousePressPoint_         Point
	mouseMovePoint_          Point
	mouseReleasePoint_       Point
	mousePressPointInches_   Point32
	mouseMovePointInches_    Point32
	mouseReleasePointInches_ Point32

	viewOffset_     Point32
	viewOffsetLast_ Point32

	isLeftButtonDown_   bool
	isCenterButtonDown_ bool
	isRightButtonDown_  bool

	modeMovingMap_ bool

	slowZoomTargetScale float64
	slowZoomTargetX     int32
	slowZoomTargetY     int32

	allowAddSnapshots_      bool
	indexOfCurrentSnapshot_ int
	snapshots_              []SnapshotItem
	hasChanges_             bool
	toolSelector            IMapToolSelector

	propertiesChangesStack *uiproperties.PropertiesChangesStack

	mapDataSource IMapDataSource

	OnMouseDrop        func(droppedValue interface{}, control IMapControl, x int32, y int32)
	OnSelectionChanged func()
	OnViewChanged      func()
	OnScaleChanged     func(scale float64)

	thumbnailMaking_ bool

	OnActionOpenMap   func(resId string)
	OnActionWriteItem func(item string, value string)
}

func NewMapWidget(parent uiinterfaces.Widget) *MapWidget {
	var c MapWidget

	c.InitControl(parent, &c)
	c.showGrid_ = true
	c.propertiesChangesStack = uiproperties.NewPropertiesChangesStack()

	c.view_ = NewMapControlView(&c, nil)
	c.view_.SetIsRootControl(true)
	c.view_.OnLoadedInEditor = c.OnLoadedInEditor

	c.allowAddSnapshots_ = true
	c.indexOfCurrentSnapshot_ = 0
	c.SetXExpandable(true)
	c.SetYExpandable(true)

	c.timerZoom_ = c.Window().NewTimer(20, func() {
		/*lastScale := c.view_.scale()
		diffScale := c.slowZoomTargetScale - lastScale
		if math.Abs(diffScale) > 0.01 {
			diffScale = diffScale / 2
			c.view_.setScale(lastScale + diffScale)
			//zoomDirect(targetZoomX_, targetZoomY_, view_->scale() + qAbs(targetScale_ - view_->scale()) / 2)
			ddX := c.slowZoomTargetX - c.viewOffset_.x
			ddY := c.slowZoomTargetY - c.viewOffset_.y
			c.viewOffset_.x = c.viewOffset_.x + ddX / 2
			c.viewOffset_.y = c.viewOffset_.y + ddY / 2
			//c.viewOffset_.y = c.viewOffset_.y + c.slowZoomTargetY
			//notifyControlChanged();
			c.Update("MapWidget")
		} else {
			if math.Abs(diffScale) > 0 {
				c.view_.setScale(c.slowZoomTargetScale)
			}
		}*/
	})
	c.timerZoom_.StartTimer()

	c.ZoomDefault()

	return &c
}

func (c *MapWidget) ControlType() string {
	return "MapWidget"
}

func (c *MapWidget) CloseView() {
	if c.view_ != nil {
		c.view_.Dispose()
		c.view_ = nil
	}
	c.OwnWindow = nil
	c.toolSelector = nil
	c.Dispose()
}

func (c *MapWidget) UpSelectedItem() {
	c.view_.UpSelectedItem()
	c.Update("MapWidget")
}

func (c *MapWidget) DownSelectedItem() {
	c.view_.DownSelectedItem()
	c.Update("MapWidget")
}

func (c *MapWidget) SetCurrentLayer(layer *MapControlViewLayer) {
	c.view_.setCurrentLayer(layer)
	c.Update("MapWidget")
}

func (c *MapWidget) NotifyLayersChanged() {
	c.NotifyChanges(nil)
	c.Update("MapWidget")
}

func (c *MapWidget) OnLoadedInEditor(list *uiproperties.PropertiesChangesList) {
	c.propertiesChangesStack.AddList(list)
}

func (c *MapWidget) SetToolSelector(toolSelector IMapToolSelector) {
	c.toolSelector = toolSelector
}

func (c *MapWidget) View() *MapControlView {
	return c.view_
}

func (c *MapWidget) SetOnSelectionChanged(OnSelectionChanged func()) {
	c.OnSelectionChanged = OnSelectionChanged
	if c.view_ != nil {
		c.view_.OnSelectionChanged = OnSelectionChanged
	}
}

func (c *MapWidget) SetMapDataSource(mapDataSource IMapDataSource) {
	c.mapDataSource = mapDataSource
	if c.view_ != nil {
		c.view_.setMapDataSource(mapDataSource)
	}
}

func (c *MapWidget) Draw(canvas ui.DrawContext) {
	t1 := time.Now()

	// Background
	canvas.SetColor(colornames.Black)
	canvas.FillRect(0, 0, c.Width(), c.Height())

	if c.view_ == nil {
		return
	}

	canvas.Save()
	canvas.Translate(int(c.viewOffset_.x), int(c.viewOffset_.y))

	canvas.GG().Push()
	cX, cY := c.RectClientAreaOnWindow()
	canvas.GG().DrawRectangle(float64(cX), float64(cY), float64(c.Width()), float64(c.Height()))
	canvas.GG().Clip()

	t2 := time.Now()

	{
		// Map
		c.drawView(canvas, float64(c.Width()), float64(c.Height()))

		if c.showGrid_ && c.view_.editing_ {
			c.drawGrid(canvas)
		}
	}

	canvas.GG().Pop()

	t3 := time.Now()

	canvas.Load()

	//c.drawGroupSelecting(canvas)

	c.drawErrors(canvas)

	c.drawZoomByMouse(canvas)

	t4 := time.Now()
	logger.Println("map 1", t2.Sub(t1))
	logger.Println("map 2", t3.Sub(t2))
	logger.Println("map 3", t4.Sub(t3))

	if c.view_.currentLayer() != nil {
		canvas.SetColor(colornames.Blue)
		canvas.SetFontFamily("Roboto")
		canvas.SetFontSize(12)
		canvas.DrawText(10, c.Height()-50, 100, 100, c.view_.currentLayer().name_)
	}
}

func (c *MapWidget) drawGrid(canvas ui.DrawContext) {
	if c.view_ == nil {
		return
	}

	width := int32(c.view_.mapWidth())
	height := int32(c.view_.mapHeight())

	gridColor := color.RGBA{R: 100, G: 100, B: 100, A: 100}

	for x := int32(0); x <= width; x += 10 {
		y1 := int(c.view_.scaleValue(int32(0)))
		y2 := int(c.view_.scaleValue(height))
		scaledX := int(c.view_.scaleValue(x))
		canvas.SetColor(gridColor)
		canvas.SetStrokeWidth(1)
		canvas.DrawLine(scaledX, y1, scaledX, y2)
	}

	for y := int32(0); y <= height; y += 10 {
		x1 := int32(0)
		x2 := width
		canvas.SetColor(gridColor)
		canvas.SetStrokeWidth(1)
		canvas.DrawLine(int(c.view_.scaleValue(x1)), int(c.view_.scaleValue(y)), int(c.view_.scaleValue(x2)), int(c.view_.scaleValue(y)))
	}

}

func (c *MapWidget) MouseDrop(ev *uievents.MouseDropEvent) {
	if c.OnMouseDrop != nil {

		c.mousePressPoint_ = Point{ev.X, ev.Y}
		point := c.translateToInches(c.mousePressPoint_)

		c.OnMouseDrop(ev.DroppingObject, c.view_.ControlUnderPoint(int32(point.x), int32(point.y)), point.x, point.y)
	}
}

func (c *MapWidget) drawView(canvas ui.DrawContext, width, height float64) {
	if c.view_ == nil {
		return
	}
	c.view_.refreshScale()
	c.view_.draw(canvas)
}

func (c *MapWidget) drawViewBorder(canvas ui.DrawContext) {
	if c.view_ == nil {
		return
	}

	width := c.view_.mapWidth()
	height := c.view_.mapHeight()
	canvas.SetColor(colornames.Brown)
	canvas.SetStrokeWidth(1)
	canvas.DrawRect(0, 0, int(c.view_.scaleValue(width)), int(c.view_.scaleValue(height)))
}

func (c *MapWidget) drawErrors(canvas ui.DrawContext) {
	offsetY := 30
	for _, txt := range c.errors_ {
		canvas.SetColor(colornames.Crimson)
		canvas.DrawText(10, offsetY, 100, 100, txt)
		offsetY += 15
	}
}

func (c *MapWidget) drawZoomByMouse(canvas ui.DrawContext) {
	if c.view_ == nil {
		return
	}
	if !c.modeZoomMap_ {
		return
	}

	x1 := c.mousePressPoint_.x
	y1 := c.mousePressPoint_.y

	width := c.mouseMovePoint_.x - c.mousePressPoint_.x
	height := c.mouseMovePoint_.y - c.mousePressPoint_.y

	if width < 0 {
		x1 = c.mouseMovePoint_.x
		width = c.mousePressPoint_.x - c.mouseMovePoint_.x
	}

	if height < 0 {
		y1 = c.mouseMovePoint_.y
		height = c.mousePressPoint_.y - c.mouseMovePoint_.y
	}
	canvas.SetColor(colornames.Aliceblue)
	canvas.SetStrokeWidth(1)
	canvas.DrawRect(x1, y1, width, height)
}

func (c *MapWidget) ZoomDefault() {

	if c.view_ == nil {
		return
	}

	centerX := int32(c.Width() / 2)
	centerY := int32(c.Height() / 2)
	halfOfWidth := float64(centerX - c.view_.mapWidth()/2)
	halfOfHeight := float64(centerY - c.view_.mapHeight()/2)
	c.viewOffset_ = Point32{int32(halfOfWidth / c.view_.scale()), int32(halfOfHeight / c.view_.scale())}

	//c.zoomDirect(int32(halfOfWidth/c.view_.scale()), int32(halfOfHeight/c.view_.scale()), 1)
	c.zoomSlow(int32(halfOfWidth/c.view_.scale()), int32(halfOfHeight/c.view_.scale()), 1)
	c.view_.setSelectedForItems(false)
	//c.notifyControlChanged();
	c.Update("MapWidget")
}

func (c *MapWidget) MoveItemUpMax() {
	if c.view_ != nil {
		c.view_.UpSelectedItemMax()
	}
}

func (c *MapWidget) MoveItemUp() {
	if c.view_ != nil {
		c.view_.UpSelectedItem()
	}
}

func (c *MapWidget) MoveItemDown() {
	if c.view_ != nil {
		c.view_.DownSelectedItem()
	}
}

func (c *MapWidget) MoveItemDownMax() {
	if c.view_ != nil {
		c.view_.DownSelectedItemMax()
	}
}

func (c *MapWidget) ZoomIn() {
	if c.view_ == nil {
		return
	}
	centerOfScreen := c.translateToInches(Point{c.Width() / 2, c.Height() / 2})
	c.zoomSlow(centerOfScreen.x, centerOfScreen.y, c.view_.scale_*1.1)
	c.Update("MapWidget")
}

func (c *MapWidget) ZoomOut() {
	if c.view_ == nil {
		return
	}
	centerOfScreen := c.translateToInches(Point{c.Width() / 2, c.Height() / 2})
	c.zoomSlow(centerOfScreen.x, centerOfScreen.y, c.view_.scale_*0.9)
	c.Update("MapWidget")
}

func (c *MapWidget) Zoom100() {
	c.ZoomDefault()
}

func (c *MapWidget) ZoomInContainer() {
	intWidth := c.view_.Width()
	intHeight := c.view_.Height()

	iK := float64(0)
	iKo := float64(0)
	if c.Height() != 0 {
		iK = float64(c.Width()) / float64(c.Height())
	}
	if intHeight != 0 {
		iKo = float64(intWidth) / float64(intHeight)
	}
	var intScale float64

	if iK < iKo {
		if intWidth != 0 {
			intScale = float64(c.Width()) / float64(intWidth)
		}
	} else {
		if intHeight != 0 {
			intScale = float64(c.Height()) / float64(intHeight)
		}
	}
	c.zoomSlow(0, 0, intScale)

	/*centerX := int32(c.Width() / 2)
	centerY := int32(c.Height() / 2)
	halfOfWidth := float64(centerX - c.view_.mapWidth()/2)
	halfOfHeight := float64(centerY - c.view_.mapHeight()/2)*/
	fullWidthOfView := float64(c.view_.Width()) * c.view_.scale()
	fullHeightOfView := float64(c.view_.Height()) * c.view_.scale()

	c.viewOffset_ = Point32{int32(c.Width()/2 - int(fullWidthOfView/2)), int32(c.Height()/2 - int(fullHeightOfView/2))}

	c.Update("MapWidget")
}

func (c *MapWidget) zoomDirect1(x, y int32, newScale float64) {
	if newScale < 0.1 {
		newScale = 0.1
	}
	if newScale > 20 {
		newScale = 20
	}
	lastScale := c.view_.scale()
	diffScale := lastScale - newScale
	c.view_.setScale(newScale)
	c.viewOffset_.x = c.viewOffset_.x + int32(float64(x)*diffScale)
	c.viewOffset_.y = c.viewOffset_.y + int32(float64(y)*diffScale)
	if c.OnScaleChanged != nil {
		c.OnScaleChanged(c.view_.scale())
	}
	//notifyControlChanged();
	c.Update("MapWidget")
}

func (c *MapWidget) zoomSlow(x, y int32, newScale float64) {
	c.zoomDirect1(x, y, newScale)
	return
	/*c.slowZoomTargetScale = newScale
	lastScale := c.view_.scale()
	diffScale := lastScale - newScale
	c.slowZoomTargetX = int32(c.viewOffset_.x + int32(float64(x)*diffScale))
	c.slowZoomTargetY = int32(c.viewOffset_.y + int32(float64(y)*diffScale))*/
}

func (c *MapWidget) translateToInches(pixels Point) Point32 {
	var result Point32
	if c.view_ == nil {
		return Point32{0, 0}
	}

	// Init
	result.x = int32(pixels.x)
	result.y = int32(pixels.y)

	// To DPI
	//result.x = int(result.x)
	//result.y = int(result.y)

	// To Scale
	result.x = int32(float64(result.x) / c.view_.scale())
	result.y = int32(float64(result.y) / c.view_.scale())

	// To offset
	result = Point32{result.x - int32(float64(c.viewOffset_.x)/c.view_.scale()), result.y - int32(float64(c.viewOffset_.y)/c.view_.scale())}

	return result
}

func (c *MapWidget) MouseDown(event *uievents.MouseDownEvent) {
	if c.view_ == nil {
		return
	}

	c.isLeftButtonDown_ = event.Button == uievents.MouseButtonLeft
	c.isCenterButtonDown_ = false
	c.isRightButtonDown_ = event.Button == uievents.MouseButtonRight

	c.mousePressPoint_ = Point{event.X, event.Y}
	c.mousePressPointInches_ = c.translateToInches(c.mousePressPoint_)
	if c.applyCurrentTool(c.mousePressPointInches_) {
		return
	}

	c.mouseMovePoint_ = Point{event.X, event.Y}
	c.mouseMovePointInches_ = c.translateToInches(c.mouseMovePoint_)

	processed := false

	fmt.Println("----------------")

	processed = c.view_.mouseDown(c.mousePressPointInches_.x, c.mousePressPointInches_.y, c.isLeftButtonDown_, c.isCenterButtonDown_, c.isRightButtonDown_, event.Modifiers.Shift, event.Modifiers.Control, event.Modifiers.Alt)
	/*if c.isLeftButtonDown_ && c.view_.editing_ {
		if c.view_ != nil {
			processed = c.view_.mouseDown(c.mousePressPointInches_.x, c.mousePressPointInches_.y, c.isLeftButtonDown_, c.isCenterButtonDown_, c.isRightButtonDown_, event.Modifiers.Shift, event.Modifiers.Control, event.Modifiers.Alt)
			if !processed {
				c.view_.setSelectedForItems(false)
			}
		}
	}*/

	if c.isRightButtonDown_ {
		if !processed {
			c.beginMoveMap()
		}
	}

	if !processed && c.isLeftButtonDown_ {
		c.beginZoomMap()
	}

	if !c.IsEditing() && event.Button == uievents.MouseButtonLeft {
		if c.view_ != nil {
			c.view_.OnMouseDown(int(c.mousePressPointInches_.x), int(c.mousePressPointInches_.y))
		}
	}

	c.Update("MapWidget")
}

func (c *MapWidget) MouseUp(event *uievents.MouseUpEvent) {
	if c.view_ == nil {
		return
	}

	c.isLeftButtonDown_ = event.Button == uievents.MouseButtonLeft
	c.isCenterButtonDown_ = false
	c.isRightButtonDown_ = event.Button == uievents.MouseButtonRight

	//lastMovePoint := c.mousePressPoint_
	c.mouseReleasePoint_ = Point{event.X, event.Y}
	c.mouseReleasePointInches_ = c.translateToInches(c.mouseReleasePoint_)

	if c.view_ != nil {
		c.view_.mouseUp(c.mousePressPointInches_.x, c.mousePressPointInches_.y, c.isLeftButtonDown_, c.isCenterButtonDown_, c.isRightButtonDown_, event.Modifiers.Shift, event.Modifiers.Control, event.Modifiers.Alt)
	}

	if c.modeMovingMap_ {
		c.modeMovingMap_ = false
	}

	if c.modeZoomMap_ {
		c.modeZoomMap_ = false
		x1 := c.mousePressPoint_.x
		y1 := c.mousePressPoint_.y
		x2 := c.mouseReleasePoint_.x
		y2 := c.mouseReleasePoint_.y

		if math.Abs(float64(x1-x2)) > 10 && math.Abs(float64(y1-y2)) > 10 {

			if x2 < x1 || y2 < y1 {
				c.ZoomDefault()
			} else {
				p1 := c.translateToInches(Point{x1, y1})
				//pc := c.translateToInches(Point{(x2 + x1) / 2, (y2 + y1) / 2})
				p2 := c.translateToInches(Point{x2, y2})
				internalRect := Rect32{p1.x, p1.y, p2.x - p1.x, p2.y - p1.y}
				dyPoints := p2.y - p1.y
				dxPoints := p2.x - p1.x

				intWidth := float64(internalRect.width)
				intHeight := float64(internalRect.height)

				iK := float64(c.Width() / c.Height())
				iKo := intWidth / intHeight
				//var intScale float64

				if iK < iKo {
					if dxPoints > 5 {
						needScale := float64(c.Width()) / float64(dxPoints)
						c.zoomSlow(int32((x2-x1)/2), int32((y2-y1)/2), needScale)
						c.viewOffset_.x = int32(-float64(p1.x) * c.view_.scale())
						c.viewOffset_.y = int32(-float64(p1.y) * c.view_.scale())
					}
				} else {
					if dyPoints > 5 {
						needScale := float64(c.Height()) / float64(dyPoints)
						c.zoomSlow(0, 0, needScale)
						c.viewOffset_.x = int32(-float64(p1.x) * c.view_.scale())
						c.viewOffset_.y = int32(-float64(p1.y) * c.view_.scale())
					}
				}

			}
		}
	}

	//c.notifyControlChanged();
	c.Update("MapWidget")
}

func (c *MapWidget) MouseWheel(event *uievents.MouseWheelEvent) {
	if c.view_ == nil {
		return
	}

	delta := event.Delta
	c.mouseMovePoint_ = Point{event.X, event.Y}
	c.mouseMovePointInches_ = c.translateToInches(c.mouseMovePoint_)

	kScale := 0.0
	if delta < 0 {
		kScale = 0.9
	} else {
		kScale = 1.1
	}

	c.zoomSlow(c.mouseMovePointInches_.x, c.mouseMovePointInches_.y, c.view_.scale()*kScale)
	//c.zoomDirect(c.mouseMovePointInches_.x, c.mouseMovePointInches_.y, c.view_.scale()*kScale)

	c.Update("MapWidget")
}

func (c *MapWidget) MouseMove(event *uievents.MouseMoveEvent) {
	if c.view_ == nil {
		return
	}

	c.mouseMovePoint_ = Point{event.X, event.Y}
	c.mouseMovePointInches_ = c.translateToInches(c.mouseMovePoint_)

	if c.view_ != nil {
		c.view_.mouseMove(c.mousePressPointInches_, c.mouseMovePointInches_, c.isLeftButtonDown_, c.isCenterButtonDown_, c.isRightButtonDown_, event.Modifiers.Shift, event.Modifiers.Control, event.Modifiers.Alt)
		controlUnderPoint := c.view_.FindControlUnderPoint(int(c.mouseMovePointInches_.x), int(c.mouseMovePointInches_.y))
		if controlUnderPoint != nil {
			if controlUnderPoint.HasAction() && !c.IsEditing() {
				c.Window().SetMouseCursor(ui.MouseCursorPointer)
			} else {
				c.Window().SetMouseCursor(ui.MouseCursorArrow)
			}
		}
	}

	if c.modeMovingMap_ {
		offset := Point32{int32(c.mouseMovePoint_.x - c.mousePressPoint_.x), int32(c.mouseMovePoint_.y - c.mousePressPoint_.y)}
		pointDpi := Point32{offset.x, offset.y}
		c.viewOffset_ = Point32{c.viewOffsetLast_.x + pointDpi.x, c.viewOffsetLast_.y + pointDpi.y}
	} else {
		if c.isLeftButtonDown_ || c.isCenterButtonDown_ || c.isRightButtonDown_ {
			//notifyControlChanged();
		}
	}

	c.Update("MapWidget")
}

func (c *MapWidget) AddControl(control IMapControl) {
	c.view_.addItem(control, control.X(), control.Y())
}

func (c *MapWidget) beginMoveMap() {
	c.modeMovingMap_ = true
	c.viewOffsetLast_ = c.viewOffset_
}

func (c *MapWidget) beginZoomMap() {
	c.modeZoomMap_ = true
	c.viewOffsetLast_ = c.viewOffset_
}

func (c *MapWidget) applyCurrentTool(point Point32) bool {

	currentToolName := c.toolSelector.CurrentTool()

	if c.view_ == nil {
		return false
	}

	if !c.view_.editing_ {
		return false
	}

	if currentToolName == "" {
		return false
	}

	if c.view_ == nil {
		return false
	}

	layer := c.view_.currentLayer()
	if layer == nil {
		return false
	}

	var control = c.view_.makeControlByType(currentToolName)

	c.toolSelector.ResetCurrentTool()
	if control == nil {
		return true
	}

	control.SetAdding()
	c.view_.addItem(control, point.x-control.Width()/2, point.y-control.Height()/2)

	return true
}

func (c *MapWidget) KeyDown(event *uievents.KeyDownEvent) bool {
	if event.Key == glfw.KeyEscape {
		if len(c.errors_) == 0 {
			c.ZoomDefault()
		} else {
			c.errors_ = make([]string, 0)
		}
		c.Update("MapWidget")
		return true
	}
	if event.Key == glfw.KeyZ && event.Modifiers.Control {
		propListForSet := c.propertiesChangesStack.Undo()
		for _, propItem := range propListForSet.Items() {
			prop := propItem.PropContainer.Property(propItem.Name)
			if prop != nil {
				prop.SetOwnValue(propItem.Value)
			}
		}
		c.Update("MapWidget")
		return true
	}
	if event.Key == 'g' {
		c.showGrid_ = !c.showGrid_
		c.Update("MapWidget")
		return true
	}
	if event.Key == 'a' && event.Modifiers.Control {
		c.selectAllItems()
		c.Update("MapWidget")
		return true
	}
	if event.Key == glfw.KeyKPSubtract {
		centerOfScreen := c.translateToInches(Point{c.Width() / 2, c.Height() / 2})
		c.zoomSlow(centerOfScreen.x, centerOfScreen.y, c.view_.scale_*0.8)
		c.Update("MapWidget")
		return true
	}
	if event.Key == glfw.KeyKPAdd {
		centerOfScreen := c.translateToInches(Point{c.Width() / 2, c.Height() / 2})
		c.zoomSlow(centerOfScreen.x, centerOfScreen.y, c.view_.scale_*1.1)
		c.Update("MapWidget")
		return true
	}
	if c.view_.editing_ && event.Key == glfw.KeyDelete {
		c.deleteSelectedItems()
		//c.notifyControlChanged()
		c.Update("MapWidget")
		return true
	}

	c.Update("MapWidget")
	return false
}

func (c *MapWidget) deleteSelectedItems() {
	items := c.SelectedItems()
	for _, item := range items {
		c.view_.removeItem(item)
	}
	c.view_.setSelectedForItems(false)
	c.view_.selectionChanged_ = true
	c.view_.checkExclusiveSelectedItem()
	c.Update("MapWidget")
}

func (c *MapWidget) SelectedItems() []IMapControl {
	if c.view_ == nil {
		return make([]IMapControl, 0)
	}
	return c.view_.selectedItems()
}

func (c *MapWidget) selectAllItems() {
	c.view_.setSelectedForItems(true)
	//c.notifyControlChanged()
	c.Update("MapWidget")
}

func (c *MapWidget) Save() []byte {
	c.hasChanges_ = false
	return c.view_.saveView()
}

func (c *MapWidget) HasChanges() bool {
	return c.hasChanges_
}

func (c *MapWidget) GetThumbnail(width, height int) image.Image {
	c.thumbnailMaking_ = true
	context := ui.NewDrawContextSWSpecial(int(c.view_.mapWidth()), int(c.view_.mapHeight()))
	lastScale := c.view_.scale()
	c.view_.setScale(1)
	c.view_.draw(context)
	c.view_.setScale(lastScale)
	img := resize.Resize(uint(width), uint(height), context.GraphContextImage(), resize.Bilinear)
	c.thumbnailMaking_ = false
	return img
}

func (c *MapWidget) Load(resId string, value []byte) error {

	if len(value) == 0 {
		if c.view_ != nil {
			c.view_.Dispose()
		}
		c.view_ = nil
		return errors.New("no data")
	}

	if c.view_ == nil {
		c.view_ = NewMapControlView(c, nil)
		c.view_.SetIsRootControl(true)
		c.view_.OnLoadedInEditor = c.OnLoadedInEditor
		c.view_.setMapDataSource(c.mapDataSource)
	}

	err := c.view_.loadView(resId, value)

	if err != nil {
		return err
	}
	c.view_.SaveOriginalProperties()
	c.view_.SetEditing(false)
	c.Update("MapWidget")
	c.view_.OnSelectionChanged = c.OnSelectionChanged
	return nil
}

func (c *MapWidget) Tick() {
	if c.view_ != nil {
		c.view_.Tick()
	}
}

func (c *MapWidget) SetEdit(edit bool) {
	if c.view_ != nil {
		c.view_.SetEditing(edit)
		if edit {
			if c.OnSelectionChanged != nil {
				c.OnSelectionChanged()
			}
		}
	}
	c.Update("MapWidget")
}

func (c *MapWidget) IsEditing() bool {
	if c.view_ != nil {
		return c.view_.isEditing()
	}
	return false
}

func (c *MapWidget) NotifyChanges(list *uiproperties.PropertiesChangesList) {
	if list != nil {
		c.propertiesChangesStack.AddList(list)
	}
	c.hasChanges_ = true
	if c.OnViewChanged != nil {
		c.OnViewChanged()
	}
}

func (c *MapWidget) MakeSnapshot(action string) {
	if c.view_ == nil {
		return
	}

	if !c.view_.editing_ {
		return
	}

	if !c.allowAddSnapshots_ {
		return
	}

	if c.indexOfCurrentSnapshot_ < len(c.snapshots_) {
		c.snapshots_ = c.snapshots_[:c.indexOfCurrentSnapshot_]
	}

	var snapshotItem SnapshotItem
	snapshotItem.Name = action
	snapshotItem.Ds = c.view_.saveView()
	c.snapshots_ = append(c.snapshots_, snapshotItem)
	c.indexOfCurrentSnapshot_++

	c.hasChanges_ = true

	if c.OnViewChanged != nil {
		c.OnViewChanged()
	}
}

func (c *MapWidget) RectClientAreaOnWindow() (int, int) {
	if !c.thumbnailMaking_ {
		return c.Control.RectClientAreaOnWindow()
	}
	return 0, 0
}

func (c *MapWidget) Width() int {
	if !c.thumbnailMaking_ {
		return c.Control.Width()
	}
	return int(c.view_.mapWidth())
}

func (c *MapWidget) Height() int {
	if !c.thumbnailMaking_ {
		return c.Control.Height()
	}
	return int(c.view_.mapHeight())
}

func (c *MapWidget) ActionOpenMap(resId string) {
	if c.OnActionOpenMap != nil {
		c.OnActionOpenMap(resId)
	}
}

func (c *MapWidget) ActionWriteItem(item string, value string) {
	if c.OnActionOpenMap != nil {
		c.OnActionWriteItem(item, value)
	}
}
