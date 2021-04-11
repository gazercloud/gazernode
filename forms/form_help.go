package forms

import (
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
)

type FormHelp struct {
	uicontrols.Dialog
	resValue      string
	txtUnitName   *uicontrols.TextBox
	btnOK         *uicontrols.Button
	OnTextChanged func(txtMultiline *FormHelp, oldText string, newText string)
}

func NewFormHelp(parent uiinterfaces.Widget, title string, value string) *FormHelp {
	var c FormHelp
	c.resValue = value
	c.InitControl(parent, &c)
	c.SetTitle("Help - " + title)
	c.Resize(450, 450)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	pLeft := pContent.AddPanelOnGrid(0, 0)
	pRight := pContent.AddPanelOnGrid(1, 0)

	img := pLeft.AddImageBoxOnGrid(0, 0, uiresources.ResImgCol(uiresources.R_icons_material4_png_action_help_outline_materialiconsoutlined_48dp_1x_outline_help_outline_black_48dp_png, c.ForeColor()))
	img.SetScaling(uicontrols.ImageBoxScaleAdjustImageKeepAspectRatio)
	img.SetMinHeight(64)
	img.SetMinWidth(64)
	pLeft.AddVSpacerOnGrid(0, 1)

	c.txtUnitName = pRight.AddTextBoxOnGrid(1, 0)
	c.txtUnitName.SetText(value)
	c.txtUnitName.SetMultiline(true)
	c.txtUnitName.SetReadOnly(true)
	c.txtUnitName.OnTextChanged = func(txtBox *uicontrols.TextBox, oldValue string, newValue string) {
		if c.OnTextChanged != nil {
			c.OnTextChanged(&c, oldValue, newValue)
		}
	}

	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)
	pButtons.AddHSpacerOnGrid(0, 0)
	btnClose := pButtons.AddButtonOnGrid(2, 0, "Close", func(event *uievents.Event) {
		c.Reject()
	})
	btnClose.SetMinWidth(70)
	c.SetRejectButton(btnClose)

	return &c
}
