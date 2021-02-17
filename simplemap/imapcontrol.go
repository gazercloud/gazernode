package simplemap

import (
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uiproperties"
)

type IMapControl interface {
	drawControl(ctx ui.DrawContext)
	TypeName() string
	Name() string
	adaptiveSize() bool

	isView() bool
	isEditing() bool
	isRootControl() bool

	updateLayout(direct bool)
	refreshScale()
	changeNotify(list *uiproperties.PropertiesChangesList)

	SetAdding()

	SetX(x int32)
	X() int32

	SetY(y int32)
	Y() int32

	SetWidth(width int32)
	Width() int32

	SetHeight(height int32)
	Height() int32

	draw(ctx ui.DrawContext)

	scale() float64
	setScale(scale float64)

	selectedExclusive() bool
	setSelectedExclusive(selectedExclusive bool)

	isPointInside(x, y int32) bool

	selected() bool
	setSelected(selected bool)

	isRectIntersect(x1, y1, x2, y2 int32) bool

	rememberLastPosition()

	mouseDown(x int32, y int32, leftButton, centerButton, rightButton, shift, control, alt bool) bool
	mouseMove(lastMouseDownPos, pos Point32, leftButton, centerButton, rightButton, shift, control, alt bool) bool
	mouseUp(x, y int32, leftButton, centerButton, rightButton, shift, control, alt bool) bool

	anchorBottom() bool
	anchorLeft() bool
	anchorRight() bool
	anchorTop() bool

	original_x() int32
	original_y() int32
	original_width() int32
	original_height() int32

	SetOriginalX(x int32)
	SetOriginalY(y int32)
	SetOriginalWidth(width int32)
	SetOriginalHeight(height int32)

	Subclass() string
	//addProperty(name string, prop *ui.Property)
	setNeedToSetDefaultSize(needToSetDefaultSize bool)

	saveBase() *MapItem
	load(value *MapItem)
	Tick()
	setMapDataSource(mapDataSource IMapDataSource)
	SetDataSource(dataSource string)

	//SetPropertyValue(name string, value interface{})
	PropertyValue(name string) interface{}

	UpdateValue(val common_interfaces.ItemValue)
	LoadContent(contentBytes []byte, err error)

	UpdateSizePoints()
	SetType(typeName string)

	Dispose()

	FullDataSource() string

	OnMouseDown(x, y int)
	FindControlUnderPoint(x, y int) IMapControl
	HasAction() bool

	GetFullPathToMapControl() []string
}
