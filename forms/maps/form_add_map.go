package maps

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
)

type FormAddMap struct {
	uicontrols.Dialog
	client      *client.Client
	tp          string
	txtUnitName *uicontrols.TextBox
	btnOK       *uicontrols.Button
}

func NewFormAddMap(parent uiinterfaces.Widget, client *client.Client, tp string) *FormAddMap {
	var c FormAddMap
	c.client = client
	c.tp = tp
	c.InitControl(parent, &c)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	pLeft := pContent.AddPanelOnGrid(0, 0)
	pLeft.SetPanelPadding(0)
	//pLeft.SetBorderRight(1, c.ForeColor())
	pLeft.SetMinWidth(100)
	pRight := pContent.AddPanelOnGrid(1, 0)
	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)

	img := pLeft.AddImageBoxOnGrid(0, 0, uiresources.ResImgCol(uiresources.R_icons_material4_png_maps_layers_materialiconsoutlined_48dp_1x_outline_layers_black_48dp_png, c.ForeColor()))
	img.SetScaling(uicontrols.ImageBoxScaleAdjustImageKeepAspectRatio)
	img.SetMinHeight(64)
	img.SetMinWidth(64)
	pLeft.AddVSpacerOnGrid(0, 1)

	pRight.AddTextBlockOnGrid(0, 0, "Map name:")
	c.txtUnitName = pRight.AddTextBoxOnGrid(1, 0)

	pRight.AddVSpacerOnGrid(0, 10)

	pButtons.AddHSpacerOnGrid(0, 0)
	c.btnOK = pButtons.AddButtonOnGrid(1, 0, "OK", nil)
	c.TryAccept = func() bool {
		if c.txtUnitName.Text() == "" || len(c.txtUnitName.Text()) > 50 {
			uicontrols.ShowErrorMessage(&c, "wrong name", "Error")
			return false
		}

		c.btnOK.SetEnabled(false)
		c.client.ResAdd(c.txtUnitName.Text(), c.tp, []byte("{}"), func(id string, err error) {
			if err == nil {
				c.TryAccept = nil
				c.Accept()
			} else {
				c.btnOK.SetEnabled(true)
				uicontrols.ShowErrorMessage(&c, err.Error(), "error")
			}
		})
		return false
	}

	c.btnOK.SetMinWidth(70)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Cancel", func(event *uievents.Event) {
		c.Reject()
	})
	btnCancel.SetMinWidth(70)

	c.SetAcceptButton(c.btnOK)
	c.SetRejectButton(btnCancel)

	c.OnShow = func() {
		c.txtUnitName.Focus()
	}

	return &c
}

func (c *FormAddMap) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("Add map")
	c.Resize(400, 200)
}
