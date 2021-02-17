package simplemap

type MapControlViewLayer struct {
	name_    string
	items_   []IMapControl
	view_    *MapControlView
	visible_ bool
}

func NewMapControlViewLayer() *MapControlViewLayer {
	var c MapControlViewLayer
	c.visible_ = true
	c.items_ = make([]IMapControl, 0)
	c.name_ = "layoutName"
	return &c
}
