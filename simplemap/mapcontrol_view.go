package simplemap

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gazercloud/gazerui/canvas"
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uiproperties"
	"github.com/nfnt/resize"
	"github.com/yuin/gopher-lua"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"image/draw"
	"strconv"
	"strings"
	"time"
)

type MapControlView struct {
	MapControl

	code_               *uiproperties.Property
	initcode_           *uiproperties.Property
	adaptiveSizeDesign_ *uiproperties.Property
	adaptiveSize_       *uiproperties.Property
	mapWidth_           *uiproperties.Property
	mapHeight_          *uiproperties.Property
	backColor_          *uiproperties.Property

	backgroundImage_ *uiproperties.Property
	borderColor_     *uiproperties.Property
	borderWidth_     *uiproperties.Property

	loaded_  bool
	loading_ bool
	editing_ bool
	adding_  bool

	layers_ []*MapControlViewLayer

	currentLayer_      *MapControlViewLayer
	lastMouseDownItem_ IMapControl

	modeGroupSelecting_    bool
	needToUpdateLayout_    bool
	lastTimeLayoutUpdated_ int64
	needToSetDefaultSize_  bool
	selectionChanged_      bool

	designTimeProperties []*uiproperties.Property

	OnLoadedInEditor   func(list *uiproperties.PropertiesChangesList)
	OnSelectionChanged func()

	img          *image.RGBA
	imgScaled    image.Image
	imgScaledKey string

	intScale float64

	err error
}

func NewMapControlView(mapWidget *MapWidget, parent IMapControl) *MapControlView {
	var c MapControlView
	c.initMapControl(&c, mapWidget, parent)

	c.adaptiveSizeDesign_ = AddPropertyToControl(&c, "adaptiveSizeDesign", "Adaptive Size", uiproperties.PropertyTypeBool, "View", "")
	c.adaptiveSizeDesign_.SetVisible(false)
	c.designTimeProperties = append(c.designTimeProperties, c.adaptiveSizeDesign_)
	c.adaptiveSize_ = AddPropertyToControl(&c, "adaptiveSize", "Adaptive Size", uiproperties.PropertyTypeBool, "View", "")
	c.adaptiveSize_.SetVisible(false)
	//c.viewOwnProperties = append(c.viewOwnProperties, c.adaptiveSize_)
	c.mapWidth_ = AddPropertyToControl(&c, "mapWidth", "Width", uiproperties.PropertyTypeInt32, "Map", "")
	c.designTimeProperties = append(c.designTimeProperties, c.mapWidth_)
	c.mapHeight_ = AddPropertyToControl(&c, "mapHeight", "Height", uiproperties.PropertyTypeInt32, "Map", "")
	c.designTimeProperties = append(c.designTimeProperties, c.mapHeight_)
	c.code_ = AddPropertyToControl(&c, "code", "Code", uiproperties.PropertyTypeMultiline, "Map", "")
	c.code_.SetVisible(false)
	c.designTimeProperties = append(c.designTimeProperties, c.code_)
	c.initcode_ = AddPropertyToControl(&c, "initcode", "Init Code", uiproperties.PropertyTypeMultiline, "Map", "")
	c.initcode_.SetVisible(false)
	c.designTimeProperties = append(c.designTimeProperties, c.initcode_)

	c.backColor_ = AddPropertyToControl(&c, "backColor", "Color", uiproperties.PropertyTypeColor, "Background", "")
	c.backColor_.SetOwnValue(color.RGBA{})
	c.backColor_.DefaultValue = c.backColor_.ValueOwn()
	c.designTimeProperties = append(c.designTimeProperties, c.backColor_)
	c.backgroundImage_ = AddPropertyToControl(&c, "backgroundImage", "Image", uiproperties.PropertyTypeString, "Background", "file")
	c.designTimeProperties = append(c.designTimeProperties, c.backgroundImage_)
	c.backgroundImage_.OnChanged = func(property *uiproperties.Property, oldValue interface{}, newValue interface{}) {
		c.loadBackgroundImage()
	}
	c.backgroundImage_.DefaultValue = make([]byte, 0)

	c.borderColor_ = AddPropertyToControl(&c, "borderColor", "Color", uiproperties.PropertyTypeColor, "Border", "")
	c.borderColor_.SetOwnValue(colornames.Gray)
	c.borderColor_.DefaultValue = c.borderColor_.ValueOwn()
	c.designTimeProperties = append(c.designTimeProperties, c.borderColor_)
	c.borderWidth_ = AddPropertyToControl(&c, "borderWidth", "Width", uiproperties.PropertyTypeInt32, "Border", "")
	c.borderWidth_.SetOwnValue(4)
	c.borderWidth_.DefaultValue = c.borderWidth_.ValueOwn()
	c.designTimeProperties = append(c.designTimeProperties, c.borderWidth_)

	c.dataSource_ = AddPropertyToControl(&c, "data_source", "Path", uiproperties.PropertyTypeString, "DataSource", "datasource")

	c.adaptiveSize_.OnChanged = func(property *uiproperties.Property, oldValue interface{}, newValue interface{}) {
		c.updateLayout(true)
		//c.SetWidth(c.Width())
		//c.SetHeight(c.Height())
		//c.updateLayout(false)
	}

	for _, prop := range c.designTimeProperties {
		prop.SetVisible(false)
	}

	c.mapWidth_.SetOwnValue(int32(500))
	c.mapHeight_.SetOwnValue(int32(500))
	c.SetWidth(int32(500))
	c.SetHeight(int32(500))

	c.mapWidth_.OnChanged = c.OnMapWidthChanged
	c.mapHeight_.OnChanged = c.OnMapHeightChanged

	l := NewMapControlViewLayer()
	c.layers_ = append(c.layers_, l)
	c.currentLayer_ = l
	c.loaded_ = false
	c.editing_ = false

	return &c
}

func (c *MapControlView) loadBackgroundImage() {
	imgData, err := base64.StdEncoding.DecodeString(c.backgroundImage_.String())
	if err == nil {
		img, _, err := image.Decode(bytes.NewBuffer(imgData))
		if err == nil {
			c.img = image.NewRGBA(img.Bounds())
			draw.Draw(c.img, c.img.Bounds(), img, image.Point{}, draw.Src)
		} else {
			c.img = nil
		}
	} else {
		c.img = nil
	}
}

func (c *MapControlView) drawBackImage(ctx ui.DrawContext) {
	if c.width_.Int32() > 0 && c.height_.Int32() > 0 {
		if c.img != nil {
			wX, wY := c.MapWidget.RectClientAreaOnWindow()

			//currentWidth := c.scaleValueIndependentContent(c.Width())
			//currentHeight := c.scaleValueIndependentContent(c.Height())
			currentWidth := int32(float64(c.mapWidth()) * c.scale_)
			currentHeight := int32(float64(c.mapHeight()) * c.scale_)

			clipXo := int32(ctx.State().TranslateX - wX)
			clipYo := int32(ctx.State().TranslateY - wY)
			clipX := int32(ctx.State().TranslateX - wX)
			clipY := int32(ctx.State().TranslateY - wY)
			clipW := currentWidth
			clipH := currentHeight

			if clipX < 0 {
				clipW += clipX
				clipX = 0
			}
			if clipY < 0 {
				clipH += clipY
				clipY = 0
			}
			if clipX+clipW > int32(c.MapWidget.Width()) {
				clipW = clipW - (clipX + clipW - int32(c.MapWidget.Width()))
			}
			if clipY+clipH > int32(c.MapWidget.Height()) {
				clipH = clipH - (clipY + clipH - int32(c.MapWidget.Height()))
			}
			if clipW < 0 {
				clipW = 0
			}
			if clipH < 0 {
				clipH = 0
			}

			percentsX := 0.0
			if clipXo < 0 {
				percentsX = -(float64(clipXo) / float64(currentWidth))
			}
			percentsY := 0.0
			if clipYo < 0 {
				percentsY = -(float64(clipYo) / float64(currentHeight))
			}
			percentsW := float64(clipW) / float64(currentWidth)
			percentsH := float64(clipH) / float64(currentHeight)

			srcPosX := int(float64(c.img.Bounds().Max.X) * percentsX)
			srcPosY := int(float64(c.img.Bounds().Max.Y) * percentsY)
			srcPosW := int(float64(c.img.Bounds().Max.X) * percentsW)
			srcPosH := int(float64(c.img.Bounds().Max.Y) * percentsH)

			if clipW > 0 && clipH > 0 {
				imgToDraw := image.Image(c.imgScaled)
				imgScaledKey := fmt.Sprint(srcPosX, "-", srcPosY, "-", srcPosX+srcPosW, "-", srcPosY+srcPosH, "-", uint(clipW), "-", uint(clipH))
				if c.imgScaledKey != imgScaledKey {
					qq := c.img.SubImage(image.Rect(srcPosX, srcPosY, srcPosX+srcPosW, srcPosY+srcPosH))
					imgToDraw = resize.Resize(uint(clipW), uint(clipH), qq, resize.Bilinear)
					c.imgScaled = imgToDraw
					c.imgScaledKey = imgScaledKey
				}
				ctx.DrawImage(int(percentsX*float64(currentWidth)), int(percentsY*float64(currentHeight)), 0, 0, imgToDraw)
			}

		}
	}

}

func (c *MapControlView) SetAdding() {
	c.adding_ = true
}

func (c *MapControlView) OnMapWidthChanged(property *uiproperties.Property, oldValue interface{}, newValue interface{}) {
	if c.isRootControl_ {
		c.SetWidth(newValue.(int32))
	}
}

func (c *MapControlView) OnMapHeightChanged(property *uiproperties.Property, oldValue interface{}, newValue interface{}) {
	if c.isRootControl_ {
		c.SetHeight(newValue.(int32))
	}
}

func (c *MapControlView) SetIsRootControl(isRootControl bool) {
	c.isRootControl_ = isRootControl

	if c.isRootControl_ {
		for _, prop := range c.GetProperties() {
			prop.SetVisible(false)
		}
		for _, prop := range c.designTimeProperties {
			if prop.Name != "adaptiveSizeDesign" {
				prop.SetVisible(true)
			}
		}
	}
}

func (c *MapControlView) Dispose() {
	for _, layer := range c.layers_ {
		for _, item := range layer.items_ {
			item.Dispose()
		}
	}

	c.currentLayer_ = nil
	c.MapWidget = nil
	c.mapDataSource = nil
	c.parent_ = nil
	c.OnSelectionChanged = nil
	c.OnLoadedInEditor = nil
	c.layers_ = make([]*MapControlViewLayer, 0)
}

func (c *MapControlView) SetEditing(editing bool) {
	c.editing_ = editing
	if !c.editing_ {
		c.setSelectedForItems(false)
	}
}

func (c *MapControlView) isEditing() bool {
	return c.editing_
}

func (c *MapControlView) setMapDataSource(mapDataSource IMapDataSource) {
	c.mapDataSource = mapDataSource
	for _, layer := range c.layers_ {
		for _, item := range layer.items_ {
			item.setMapDataSource(mapDataSource)
		}
	}
}

func (c *MapControlView) updateMapDataSource() {
	c.setMapDataSource(c.mapDataSource)
}

func (c *MapControlView) drawControl(ctx ui.DrawContext) {
	if !c.loaded_ {

		if !c.loading_ {
			if c.mapDataSource != nil {
				c.mapDataSource.LoadContent(c.type_.String(), c.imapControl)
			}
		}

		var w float64
		var h float64
		if c.imapControl.adaptiveSize() {
			w = float64(c.scaleValue(c.Width()))
			h = float64(c.scaleValue(c.Height()))
		} else {
			w = float64(c.scaleValueIndependentContent(c.Width()))
			h = float64(c.scaleValueIndependentContent(c.Height()))
		}
		cc := ctx.GG()

		cc.Push()
		cc.Translate(float64(ctx.State().TranslateX), float64(ctx.State().TranslateY))
		if w < 20 || h < 20 {
			cc.SetColor(uiproperties.ParseHexColor("#333"))
			cc.DrawRectangle(0, 0, w, h)
			cc.Fill()
		} else {
			cc.SetColor(uiproperties.ParseHexColor("#222"))
			cc.DrawRectangle(0, 0, w, h)
			cc.Fill()
			cc.SetColor(uiproperties.ParseHexColor("#555"))
			cc.SetLineWidth(10)
			cc.DrawRectangle(0, 0, w, h)
			cc.Stroke()
		}
		cc.Pop()

		if c.err == nil {
			ctx.SetColor(colornames.White)
			ctx.SetFontSize(12)
			ctx.SetTextAlign(canvas.HAlignCenter, canvas.VAlignCenter)
			ctx.DrawText(0, 0, int(w), int(h), "LOADING ...")
		} else {
			ctx.SetColor(colornames.Red)
			ctx.SetFontSize(14)
			ctx.SetTextAlign(canvas.HAlignCenter, canvas.VAlignCenter)
			text := "ERR: " + c.err.Error()
			text = strings.ReplaceAll(text, " ", "\r\n")
			ctx.DrawText(0, 0, int(w), int(h), text)
		}

		return
	}

	// Draw background
	ctx.SetColor(c.backColor_.Color())
	currentWidth := int(float64(c.mapWidth()) * c.scale_)
	currentHeight := int(float64(c.mapHeight()) * c.scale_)
	ctx.FillRect(0, 0, currentWidth, currentHeight)
	c.drawBackImage(ctx)

	// Draw border
	{
		var w int
		var h int

		if c.imapControl.adaptiveSize() {
			w = int(c.scaleValue(c.Width()))
			h = int(c.scaleValue(c.Height()))
		} else {
			w = int(c.scaleValueIndependentContent(c.Width()))
			h = int(c.scaleValueIndependentContent(c.Height()))
		}

		w = int(c.scaleValueIndependentContent(int32(c.intScale * float64(c.mapWidth()))))
		h = int(c.scaleValueIndependentContent(int32(c.intScale * float64(c.mapHeight()))))

		ctx.SetColor(c.borderColor_.Color())
		ctx.SetStrokeWidth(int(c.scaleValue(c.borderWidth_.Int32())))
		ctx.DrawRect(0, 0, w, h)
	}

	clipXmin, clipYmin, clipXmax, clipYmax := ctx.ClippedRegion()

	// Draw layers
	for _, layer := range c.layers_ {
		if !layer.visible_ {
			continue
		}

		for itemIndex := len(layer.items_) - 1; itemIndex >= 0; itemIndex-- {
			// Translate
			item := layer.items_[itemIndex]
			posX := int(c.scaleValue(item.X()))
			posY := int(c.scaleValue(item.Y()))
			w := c.scaleValue(item.Width())
			h := c.scaleValue(item.Height())
			if (w <= 0 || h <= 0) && item.TypeName() != "line" {
				continue
			}

			if posX+ctx.TranslatedX() > clipXmax {
				continue
			}

			if posY+ctx.TranslatedY() > clipYmax {
				continue
			}

			if posX+ctx.TranslatedX()+int(w) < clipXmin {
				continue
			}

			if posY+ctx.TranslatedY()+int(h) < clipYmin {
				continue
			}

			// Draw item
			ctx.Save()
			ctx.Translate(posX, posY)
			item.draw(ctx)

			if c.MapWidget.IsEditing() && c.isRootControl() && item.TypeName() != "line" {

				cc := ctx.GG()
				cc.Push()
				cc.SetLineCapSquare()
				cc.Translate(float64(ctx.State().TranslateX), float64(ctx.State().TranslateY))

				if item.selected() {
					cc.SetColor(colornames.Crimson)
				} else {
					cc.SetColor(colornames.Yellow)
				}

				cc.SetLineWidth(1.5)
				cc.SetDash(5, 5)
				leftX := c.scaleValue(0)
				topY := c.scaleValue(0)
				x := float64(leftX)
				y := float64(topY)
				w := float64(int(c.scaleValue(item.Width())))
				h := float64(int(c.scaleValue(item.Height())))
				cc.DrawLine(x, y, x+w, y)
				cc.DrawLine(x+w, y, x+w, y+h)
				cc.DrawLine(x+w, y+h, x, y+h)
				cc.DrawLine(x, y+h, x, y)
				cc.Stroke()
				cc.Pop()

				//ctx.SetColor(colornames.Lightgray)
				//ctx.SetStrokeWidth(int(c.scaleValue(1)))
				//ctx.DrawRect(int(leftX), int(topY), int(c.scaleValue(c.Width())), int(c.scaleValue(c.Height())))
			}

			ctx.Load()
		}
	}

	if c.modeGroupSelecting_ {
		x := int(c.scaleValue(c.mouseDownPointInches_.x))
		y := int(c.scaleValue(c.mouseDownPointInches_.y))
		w := int(c.scaleValue(c.mouseMovePointInches_.x - c.mouseDownPointInches_.x))
		h := int(c.scaleValue(c.mouseMovePointInches_.y - c.mouseDownPointInches_.y))
		ctx.SetColor(colornames.Purple)
		ctx.SetStrokeWidth(1)
		ctx.DrawRect(x, y, w, h)
	}
}

func (c *MapControlView) mapWidth() int32 {
	return c.mapWidth_.Int32()
}

func (c *MapControlView) setMapWidth(mapWidth int32) {
	c.mapWidth_.SetOwnValue(mapWidth)
}

func (c *MapControlView) mapHeight() int32 {
	return c.mapHeight_.Int32()
}

func (c *MapControlView) setMapHeight(mapHeight int32) {
	c.mapHeight_.SetOwnValue(mapHeight)
}

func (c *MapControlView) refreshScale() {
	c.setScale(c.scale())
}

func (c *MapControlView) setScale(scale float64) {
	c.setScaleIndependentContent(scale)

	if !c.imapControl.adaptiveSize() {
		intWidth := c.internalWidth()
		intHeight := c.internalHeight()

		iK := float64(0)
		iKo := float64(0)
		if c.Height() != 0 {
			iK = float64(c.Width()) / float64(c.Height())
		}
		if intHeight != 0 {
			iKo = float64(intWidth) / float64(intHeight)
		}

		if iK < iKo {
			if intWidth != 0 {
				c.intScale = float64(c.Width()) / float64(intWidth)
			}
		} else {
			if intHeight != 0 {
				c.intScale = float64(c.Height()) / float64(intHeight)
			}
		}
		scale = c.intScale * scale
	}

	c.scale_ = scale

	//c.SetWidth(int32(c.intScale * float64(c.mapWidth())))
	//c.SetHeight(int32(c.intScale * float64(c.mapWidth())))

	for _, layer := range c.layers_ {
		for _, i := range layer.items_ {
			i.setScale(c.scale_)
		}
	}
}

func (c *MapControlView) currentLayer() *MapControlViewLayer {
	return c.currentLayer_
}

func (c *MapControlView) setCurrentLayer(layer *MapControlViewLayer) {
	c.setSelectedForItems(false)

	for _, l := range c.layers_ {
		if l == layer {
			c.currentLayer_ = l
		}
	}
}

func (c *MapControlView) setSelectedForItem(control IMapControl, selected bool) {
	lastState := control.selected()
	control.setSelected(selected)
	if lastState != selected {
		c.selectionChanged_ = true
	}
}

func (c *MapControlView) setSelectedForItems(selected bool) {
	if c.currentLayer() == nil {
		return
	}

	for _, item := range c.currentLayer_.items_ {
		c.setSelectedForItem(item, selected)
	}

	c.checkExclusiveSelectedItem()
}

func (c *MapControlView) checkExclusiveSelectedItem() {
	if c.currentLayer() == nil {
		return
	}

	countOfSelectedItems := 0
	var exclusiveSelectedItem IMapControl

	for _, item := range c.currentLayer_.items_ {
		if item.selected() {
			exclusiveSelectedItem = item
			countOfSelectedItems++
			if countOfSelectedItems > 1 {
				break
			}
		}
	}

	if countOfSelectedItems == 1 {
		for _, item := range c.currentLayer_.items_ {
			item.setSelectedExclusive(false)
		}
		exclusiveSelectedItem.setSelectedExclusive(true)
	} else {
		for _, item := range c.currentLayer_.items_ {
			item.setSelectedExclusive(false)
		}
	}

	return

	if c.selectionChanged_ {
		c.selectionChanged_ = false
		if c.OnSelectionChanged != nil {
			c.OnSelectionChanged()
		}
	}

}

func (c *MapControlView) internalWidth() int32 {
	return c.mapWidth_.Int32()
}

func (c *MapControlView) internalHeight() int32 {
	return c.mapHeight_.Int32()
}

func (c *MapControlView) mouseDown(x int32, y int32, leftButton, centerButton, rightButton, shift, control, alt bool) bool {

	c.mouseDownPointInches_ = Point32{int32(x), int32(y)}
	c.mouseMovePointInches_ = Point32{int32(x), int32(y)}

	if c.editing_ {
		if c.currentLayer() == nil {
			return false
		}

		if !c.currentLayer_.visible_ {
			return false
		}

		for _, item := range c.currentLayer_.items_ {
			if item.selectedExclusive() {
				// Delegate the event to the child item
				resOfExclusive := item.mouseDown(x, y, leftButton, centerButton, rightButton, shift, control, alt)
				if resOfExclusive {
					c.lastMouseDownItem_ = item
					return true
				}
				break
			}
		}

		// Select/deselect logic
		// Find item under mouse pointer
		var itemInPoint IMapControl
		if len(c.currentLayer_.items_) > 0 {
			for _, item := range c.currentLayer_.items_ {
				if item.isPointInside(x, y) {
					itemInPoint = item
					break
				}
			}
		}

		c.lastMouseDownItem_ = itemInPoint

		// No item found
		if itemInPoint == nil {
			if shift {
				c.modeGroupSelecting_ = true
				return true
			}
			c.setSelectedForItems(false)
			return false
		}

		topItem := itemInPoint
		topItem.rememberLastPosition()

		// Clear selection if needed
		/*needToClearSelection := !control
		if topItem.selected() {
			needToClearSelection = false
		}

		if needToClearSelection {
			for _, item := range c.currentLayer_.items_ {
				c.setSelectedForItem(item, false)
			}
		}

		if !control {
			c.setSelectedForItem(topItem, true)
		} else {
			c.setSelectedForItem(topItem, !topItem.selected())
		}*/

		for _, item := range c.currentLayer_.items_ {
			if item == topItem {
				c.setSelectedForItem(item, true)
			} else {
				c.setSelectedForItem(item, false)
			}
		}

		c.checkExclusiveSelectedItem()

		return true
	} else {
		if !c.imapControl.isRootControl() {
			return c.MapControl.mouseDown(x, y, leftButton, centerButton, rightButton, shift, control, alt)
		}

		/*var itemInPoint IMapControl
		for _, layer := range c.layers_ {
			for _, item := range layer.items_ {
				if item.isPointInside(x, y) {
					itemInPoint = item
					break
				}
			}
		}
		if itemInPoint != nil {
			fmt.Println("Clicked item: ", itemInPoint.TypeName(), itemInPoint.X(), "event", x, y)
			itemInPoint.OnMouseDown(int(x-itemInPoint.X()), int(y-itemInPoint.Y()))
		}*/

	}

	return false
}

func (c *MapControlView) OnMouseDown(x, y int) {
	x = int(float64(x) / c.intScale)
	y = int(float64(y) / c.intScale)
	fmt.Println("OnMouseDown ", c.Name(), x, y)
	var itemInPoint IMapControl
	for _, layer := range c.layers_ {
		for _, item := range layer.items_ {
			if item.isPointInside(int32(x), int32(y)) {
				itemInPoint = item
				break
			}
		}
	}
	if itemInPoint != nil {
		targetX := 1 * float64(x-int(itemInPoint.X()))
		targetY := 1 * float64(y-int(itemInPoint.Y()))
		itemInPoint.OnMouseDown(int(targetX), int(targetY))
	}
}

func (c *MapControlView) FindControlUnderPoint(x, y int) IMapControl {
	x = int(float64(x) / c.intScale)
	y = int(float64(y) / c.intScale)
	for _, layer := range c.layers_ {
		for _, item := range layer.items_ {
			if item.isPointInside(int32(x), int32(y)) {
				res := item.FindControlUnderPoint(x-int(item.X()), y-int(item.Y()))
				if res != nil {
					return res
				}
			}
		}
	}
	return c
}

func (c *MapControlView) mouseMove(lastMouseDownPos, pos Point32, leftButton, centerButton, rightButton, shift, control, alt bool) bool {
	c.mouseMovePointInches_ = pos

	if c.editing_ {
		if c.currentLayer() == nil {
			return false
		}

		if !c.currentLayer_.visible_ {
			return false
		}

		if c.modeGroupSelecting_ {
			x1 := c.mouseDownPointInches_.x
			y1 := c.mouseDownPointInches_.y
			x2 := c.mouseMovePointInches_.x
			y2 := c.mouseMovePointInches_.y

			if x1 > x2 {
				x1, x2 = x2, x1
			}

			if y1 > y2 {
				y1, y2 = y2, y1
			}

			for _, item := range c.currentLayer_.items_ {
				intersect := item.isRectIntersect(x1, y1, x2, y2)
				c.setSelectedForItem(item, intersect)
			}

			c.checkExclusiveSelectedItem()
			return true
		}

		if c.lastMouseDownItem_ != nil || (!leftButton && !centerButton && !rightButton) {
			for _, item := range c.currentLayer_.items_ {
				item.mouseMove(lastMouseDownPos, pos, leftButton, centerButton, rightButton, shift, control, alt)
			}
		}
	} else {
		return c.MapControl.mouseMove(lastMouseDownPos, pos, leftButton, centerButton, rightButton, shift, control, alt)
	}

	return false
}

func (c *MapControlView) mouseUp(x, y int32, leftButton, centerButton, rightButton, shift, control, alt bool) bool {
	if c.editing_ {
		if c.currentLayer() == nil {
			return false
		}

		if !c.currentLayer_.visible_ {
			return false
		}

		c.lastMouseDownItem_ = nil

		if c.modeGroupSelecting_ {
			c.modeGroupSelecting_ = false
			return false
		}

		for _, item := range c.currentLayer_.items_ {
			item.mouseUp(x, y, leftButton, centerButton, rightButton, shift, control, alt)
		}
	} else {
		return c.MapControl.mouseUp(x, y, leftButton, centerButton, rightButton, shift, control, alt)
	}

	return false
}

func (c *MapControlView) isView() bool {
	return true
}

func (c *MapControlView) adaptiveSize() bool {
	return c.adaptiveSize_.Bool() && c.adaptiveSizeDesign_.Bool()
}

func (c *MapControlView) updateLayout(direct bool) {
	allowUpdate := true
	if !direct {
		if !c.needToUpdateLayout_ {
			return
		}

		c.needToUpdateLayout_ = false

		if c.isRootControl_ {
			now := time.Now().UnixNano()
			period := now - c.lastTimeLayoutUpdated_
			if period < 100000 {
				allowUpdate = false
			}

			if !allowUpdate {
				return
			}
		}
	}

	//updateStyles(currentTheme_);

	if c.imapControl.adaptiveSize() {
		for _, layer := range c.layers_ {
			for _, i := range layer.items_ {
				if i.anchorRight() && i.anchorLeft() {
					i.SetWidth(i.original_width() + (c.Width() - c.mapWidth()))
				}

				if i.anchorRight() && !i.anchorLeft() {
					i.SetX(i.original_x() + (c.Width() - c.mapWidth()))
				}

				if i.anchorBottom() && i.anchorTop() {
					i.SetHeight(i.original_height() + (c.Height() - c.mapHeight()))
				}

				if i.anchorBottom() && !i.anchorTop() {
					i.SetY(i.original_y() + (c.Height() - c.mapHeight()))
				}

				i.updateLayout(direct)
			}
		}
	} else {
		c.restoreOriginalSizesOfChildren()
		for _, layer := range c.layers_ {
			for _, i := range layer.items_ {
				i.updateLayout(direct)
			}
		}
	}
	c.lastTimeLayoutUpdated_ = time.Now().UnixNano()
}

func (c *MapControlView) restoreOriginalSizesOfChildren() {
	for _, layer := range c.layers_ {
		for _, i := range layer.items_ {
			i.SetX(i.original_x())
			i.SetY(i.original_y())
			i.SetWidth(i.original_width())
			i.SetHeight(i.original_height())
		}
	}
}

func (c *MapControlView) TypeName() string {
	return "view"
}

func (c *MapControlView) addItem(control IMapControl, x0, y0 int32) {
	if c.currentLayer_ == nil {
		return
	}

	control.SetX(x0)
	control.SetY(y0)
	c.currentLayer_.items_ = append([]IMapControl{control}, c.currentLayer_.items_...)

	if control.isView() {
		control.setNeedToSetDefaultSize(true)
	}

	//c.childLoaded();
	//c.changeNotify("added item " + control.TypeName())
	c.updateMapDataSource()
	c.refreshScale()
}

func (c *MapControlView) AddLayer(name string) {
	layer := NewMapControlViewLayer()
	layer.name_ = name
	c.layers_ = append(c.layers_, layer)
	c.changeNotify(nil)
}

func (c *MapControlView) RemoveLayer(layer *MapControlViewLayer) {
	foundIndex := -1
	for index, l := range c.layers_ {
		if l == layer {
			foundIndex = index
			break
		}
	}
	if foundIndex > -1 {
		c.layers_ = append(c.layers_[:foundIndex], c.layers_[foundIndex+1:]...)
		c.changeNotify(nil)
	}
}

func (c *MapControlView) UpSelectedItem() {
	if len(c.selectedItems()) == 1 {
		foundIndex := -1
		selectedItem := c.selectedItems()[0]
		for index, item := range c.currentLayer_.items_ {
			if item == selectedItem {
				foundIndex = index
			}
		}

		if foundIndex > 0 {
			c.currentLayer_.items_[foundIndex] = c.currentLayer_.items_[foundIndex-1]
			c.currentLayer_.items_[foundIndex-1] = selectedItem
			c.changeNotify(nil)
		}
	}
}

func (c *MapControlView) UpSelectedItemMax() {
	if len(c.selectedItems()) == 1 {
		foundIndex := -1
		selectedItem := c.selectedItems()[0]
		for index, item := range c.currentLayer_.items_ {
			if item == selectedItem {
				foundIndex = index
			}
		}

		if foundIndex > 0 {
			c.currentLayer_.items_ = append(c.currentLayer_.items_[:foundIndex], c.currentLayer_.items_[foundIndex+1:]...)
			c.currentLayer_.items_ = append([]IMapControl{selectedItem}, c.currentLayer_.items_...)
			c.changeNotify(nil)
		}

	}
}

func (c *MapControlView) DownSelectedItem() {
	if len(c.selectedItems()) == 1 {
		foundIndex := -1
		selectedItem := c.selectedItems()[0]
		for index, item := range c.currentLayer_.items_ {
			if item == selectedItem {
				foundIndex = index
			}
		}

		if foundIndex < len(c.currentLayer_.items_)-1 {
			c.currentLayer_.items_[foundIndex] = c.currentLayer_.items_[foundIndex+1]
			c.currentLayer_.items_[foundIndex+1] = selectedItem
			c.changeNotify(nil)
		}
	}
}

func (c *MapControlView) DownSelectedItemMax() {
	if len(c.selectedItems()) == 1 {
		foundIndex := -1
		selectedItem := c.selectedItems()[0]
		for index, item := range c.currentLayer_.items_ {
			if item == selectedItem {
				foundIndex = index
			}
		}

		if foundIndex > -1 {
			c.currentLayer_.items_ = append(c.currentLayer_.items_[:foundIndex], c.currentLayer_.items_[foundIndex+1:]...)
			c.currentLayer_.items_ = append(c.currentLayer_.items_, selectedItem)
			c.changeNotify(nil)
		}
	}
}

func (c *MapControlView) setNeedToSetDefaultSize(needToSetDefaultSize bool) {
	c.needToSetDefaultSize_ = needToSetDefaultSize
}

func (c *MapControlView) selectedItems() []IMapControl {
	result := make([]IMapControl, 0)
	if c.currentLayer() == nil {
		return result
	}

	for _, item := range c.currentLayer_.items_ {
		if item.selected() {
			result = append(result, item)
		}
	}
	return result
}

func (c *MapControlView) removeItem(control IMapControl) {
	if c.currentLayer_ == nil {
		return
	}

	items := make([]IMapControl, 0)
	for _, item := range c.currentLayer_.items_ {
		if item != control {
			items = append(items, item)
		}
	}

	c.currentLayer_.items_ = items

	//c.changeNotify("removing item")
}

func (c *MapControlView) removeAllLayers() {
	c.layers_ = make([]*MapControlViewLayer, 0)
}

func (c *MapControlView) saveView() []byte {
	result := NewMap()

	result.RootItem = c.saveBase()

	for _, layer := range c.layers_ {
		dsLayer := NewMapLayer()
		dsLayer.Name = layer.name_
		dsLayer.Visible = layer.visible_
		for _, item := range layer.items_ {
			dsLayer.Items = append(dsLayer.Items, item.saveBase())
		}
		result.Layers = append(result.Layers, dsLayer)
	}

	bs, _ := json.MarshalIndent(result, "", " ")

	return bs
}

func (c *MapControlView) loadView(resId string, value []byte) error {
	c.img = nil

	var m Map
	err := json.Unmarshal(value, &m)
	if err != nil {
		return err
	}

	c.removeAllLayers()

	if c.isRootControl() {
		// Load root view
		c.load(m.RootItem)
		c.SetType(resId)

	} else {
		// Load view content inside other view
		// Just own properties
		// Base properties must be loaded already
		for _, prop := range c.designTimeProperties {
			for _, vp := range m.RootItem.Props {
				if prop.Name == vp.Name {
					prop.SetOwnValue(vp.Value)
				}
			}
		}
	}

	if c.adding_ {

		parentView, ok := c.parent_.(*MapControlView)
		if ok {
			maxWidthOfInsertedView := parentView.mapWidth() / 4
			targetWidth := c.mapWidth()
			targetHeight := c.mapHeight()

			if c.mapWidth() > maxWidthOfInsertedView {
				k := float64(c.mapWidth()) / float64(maxWidthOfInsertedView)
				targetWidth = int32(float64(c.mapWidth()) / k)
				targetHeight = int32(float64(c.mapHeight()) / k)
			}

			centerX := c.X() + targetWidth/2
			centerY := c.Y() + targetHeight/2
			c.SetX(centerX - targetWidth/2)
			c.SetY(centerY - targetHeight/2)
			c.SetWidth(targetWidth)
			c.SetHeight(targetHeight)
		}

		if c.adaptiveSizeDesign_.Bool() {
			c.adaptiveSize_.SetOwnValue(true)
		}

		c.adding_ = false
	}

	for _, dsLayer := range m.Layers {
		layer := NewMapControlViewLayer()
		c.layers_ = append(c.layers_, layer)

		layer.name_ = dsLayer.Name
		layer.visible_ = dsLayer.Visible

		for _, dsItem := range dsLayer.Items {
			for _, vp := range dsItem.Props {
				if vp.Name == "type" {
					control := c.makeControlByType(vp.Value)
					control.load(dsItem)
					layer.items_ = append(layer.items_, control)
					break
				}
			}
		}
	}

	if len(c.layers_) == 0 {
		layer := NewMapControlViewLayer()
		layer.name_ = "DefaultLayer"
		c.layers_ = append(c.layers_, layer)
	}

	c.setCurrentLayer(c.layers_[0])

	c.updateMapDataSource()

	c.saveSizesAsOriginalForTopLevelControl()
	c.loadBackgroundImage()

	c.loaded_ = true

	return nil
}

func (c *MapControlView) saveSizesAsOriginalForTopLevelControl() {
	for _, layer := range c.layers_ {
		for _, item := range layer.items_ {
			item.SetOriginalX(item.X())
			item.SetOriginalY(item.Y())
			item.SetOriginalWidth(item.Width())
			item.SetOriginalHeight(item.Height())
		}
	}
}

func (c *MapControlView) makeControlByType(typeName string) IMapControl {
	var control IMapControl

	if typeName == "line" {
		control = NewMapControlLine(c.MapWidget, c)
	}

	if typeName == "text" {
		mapControlText := NewMapControlText(c.MapWidget, c)
		mapControlText.SetBorderTypeRect()
		mapControlText.SetWidth(200)
		mapControlText.SetHeight(100)
		mapControlText.SetBorderWidth(4)
		control = mapControlText
	}

	if typeName == "circle" {
		mapControlText := NewMapControlText(c.MapWidget, c)
		mapControlText.SetBorderTypeCircle()
		mapControlText.SetWidth(200)
		mapControlText.SetHeight(200)
		mapControlText.SetBorderWidth(4)
		control = mapControlText
	}

	if control == nil {
		control = NewMapControlView(c.MapWidget, c)
		control.SetType(typeName)
		control.SetWidth(200)
		control.SetHeight(200)
	}

	return control
}

func (c *MapControlView) Tick() {
	c.execLua()

	for _, layer := range c.layers_ {
		for _, item := range layer.items_ {
			item.Tick()
		}
	}
}

func (c *MapControlView) LoadContent(contentBytes []byte, err error) {
	if err != nil {
		c.err = err
		return
	}

	fullPath := c.GetFullPathToMapControl()
	for _, pathItem := range fullPath[:len(fullPath)-1] {
		if pathItem == c.type_.String() {
			c.err = errors.New("recursion detected")
			break
		}
	}
	if c.err != nil {
		return
	}

	c.err = c.loadView(c.TypeName(), contentBytes)
	if c.err != nil {
		return
	}

	c.loaded_ = true
	c.err = nil
	c.updateLayout(true)
	c.refreshScale()
}

func (c *MapControlView) SaveOriginalProperties() {
	listOfAllProps := uiproperties.NewPropertiesChangesList()
	for _, prop := range c.GetProperties() {
		listOfAllProps.AddItem(c, prop.Name, prop.Value())
	}

	for _, layer := range c.layers_ {
		for _, item := range layer.items_ {
			for _, prop := range item.(uiproperties.IPropertiesContainer).GetProperties() {
				listOfAllProps.AddItem(item.(uiproperties.IPropertiesContainer), prop.Name, prop.Value())
			}
		}
	}

	if c.OnLoadedInEditor != nil {
		c.OnLoadedInEditor(listOfAllProps)
	}
}

func (c *MapControlView) ControlUnderPoint(x, y int32) IMapControl {
	var itemInPoint IMapControl
	if len(c.currentLayer_.items_) > 0 {
		for _, item := range c.currentLayer_.items_ {
			if item.isPointInside(x, y) {
				itemInPoint = item
				break
			}
		}
	}
	return itemInPoint
}

func (c *MapControlView) execLua() {
	if c.editing_ {
		return
	}

	L := lua.NewState()
	defer L.Close()
	L.SetGlobal("setValue", L.NewFunction(c.luaSetValue))
	L.SetGlobal("getDataItemValueAsString", L.NewFunction(c.luaGetDataItemValueAsString))
	L.SetGlobal("getDataItemValueAsDouble", L.NewFunction(c.luaGetDataItemValueAsDouble))
	/*if err := L.DoString("setValue('text', getDataItemValue())"); err != nil {
	}*/
	if err := L.DoString(c.code_.String()); err != nil {
		fmt.Println("Lua error: " + err.Error())
	}
}

func (c *MapControlView) luaSetValue(L *lua.LState) int {
	propItem := L.ToString(1)
	propName := L.ToString(2)
	propValue := L.ToString(3)

	for _, layer := range c.layers_ {
		for _, item := range layer.items_ {
			if item.Name() == propItem {
				prop := item.(uiproperties.IPropertiesContainer).Property(propName)
				if prop != nil {
					prop.SetOwnValue(propValue)
				}
			}
		}
	}

	fmt.Println("---------------------- LUA ------------------")
	return 0
}

func (c *MapControlView) luaGetDataItemValueAsString(L *lua.LState) int {
	signsStr := L.ToString(1)
	signs, err := strconv.ParseInt(signsStr, 10, 32)
	if err != nil {
		signs = 3
	}

	if signs < 0 {
		signs = 0
	}

	if signs > 6 {
		signs = 6
	}

	val := "<novalue>"

	if c.mapDataSource != nil {
		c.mapDataSource.GetDataItemValue(c.dataSource(), c.imapControl)
		dblValue, _ := strconv.ParseFloat(c.value.Value, 64)
		val = fmt.Sprintf("%."+fmt.Sprint(signs)+"f", dblValue)
	} else {
		val = "<noDataSource>"
	}

	L.Push(lua.LString(val))

	return 1
}

func (c *MapControlView) luaGetDataItemValueAsDouble(L *lua.LState) int {

	val := 0.0

	if c.mapDataSource != nil {
		c.mapDataSource.GetDataItemValue(c.FullDataSource(), c.imapControl)
		val, _ = strconv.ParseFloat(c.value.Value, 64)
	} else {
	}

	L.Push(lua.LNumber(val))

	return 1
}
