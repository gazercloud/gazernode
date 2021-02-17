package forms

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/product/productinfo"
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type FormAbout struct {
	uicontrols.Dialog
}

func NewFormAbout(parent uiinterfaces.Widget) *FormAbout {
	var c FormAbout
	c.InitControl(parent, &c)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)

	pLeft := pContent.AddPanelOnGrid(0, 0)
	pLeft.SetPanelPadding(0)
	pLeft.SetMinWidth(100)
	img := pLeft.AddImageBoxOnGrid(0, 0, productinfo.Icon64())
	img.SetScaling(uicontrols.ImageBoxScaleAdjustImageKeepAspectRatio)
	//img.SetMinHeight(16)
	//img.SetMinWidth(16)
	img.SetFixedSize(64, 64)
	//pLeft.AddVSpacerOnGrid(0, 1)

	eMailAddress := "admin@gazer.cloud"

	pRight := pContent.AddPanelOnGrid(1, 0)
	pRight.AddTextBlockOnGrid(0, 0, "Gazer version "+productinfo.Version())
	pRight.AddTextBlockOnGrid(0, 1, "Copyright (c) Poluianov Ivan, 2020")
	txtEMail := pRight.AddTextBlockOnGrid(0, 2, "eMail: "+eMailAddress)
	txtEMail.OnClick = func(ev *uievents.Event) {
		client.OpenBrowser("mailto:" + eMailAddress)
	}
	txtEMail.SetMouseCursor(ui.MouseCursorPointer)
	pRight.AddHSpacerOnGrid(1, 0)

	pContent.AddVSpacerOnGrid(0, 1)

	pButtons := c.ContentPanel().AddPanelOnGrid(0, 2)
	btnPoweredBy := pButtons.AddButtonOnGrid(0, 0, "Powered by open-source software", func(event *uievents.Event) {
		formPoweredBy := NewFormPoweredBy(&c)
		formPoweredBy.ShowDialog()
	})
	btnPoweredBy.SetMinWidth(70)

	pButtons.AddHSpacerOnGrid(1, 0)

	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Close", func(event *uievents.Event) {
		c.Reject()
	})
	btnCancel.SetMinWidth(70)

	c.SetRejectButton(btnCancel)

	return &c
}

func (c *FormAbout) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("About")
	c.Resize(500, 300)
}
