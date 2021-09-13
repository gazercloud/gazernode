package tools

import (
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type FormAbout struct {
	uicontrols.Dialog
}

func NewFormAbout(parent uiinterfaces.Widget) *FormAbout {
	var c FormAbout
	c.InitControl(parent, &c)

	return &c
}

func (c *FormAbout) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("About")
	c.Resize(500, 300)
}
