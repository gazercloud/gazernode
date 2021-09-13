package maps

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
)

type DialogEditMapName struct {
	uicontrols.Dialog
	id        string
	client    *client.Client
	txtText   *uicontrols.TextBox
	btnOK     *uicontrols.Button
	btnCancel *uicontrols.Button
}

func NewDialogEditMapName(parent uiinterfaces.Widget, client *client.Client, id string, text string) *DialogEditMapName {
	var c DialogEditMapName
	c.id = id
	c.client = client
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
	c.txtText = pRight.AddTextBoxOnGrid(1, 0)
	c.txtText.SetText(text)

	pRight.AddVSpacerOnGrid(0, 10)

	pButtons.AddHSpacerOnGrid(0, 0)
	c.btnOK = pButtons.AddButtonOnGrid(1, 0, "OK", nil)
	c.TryAccept = func() bool {
		if c.txtText.Text() == "" || len(c.txtText.Text()) > 50 {
			uicontrols.ShowErrorMessage(&c, "wrong name", "Error")
			return false
		}

		c.client.ResRename(c.id, c.txtText.Text(), func(err error) {
			c.TryAccept = nil
			c.Accept()
		})
		return false
	}
	c.btnOK.SetMinWidth(70)

	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Cancel", nil)
	btnCancel.SetMinWidth(70)

	c.SetAcceptButton(c.btnOK)
	c.SetRejectButton(btnCancel)

	c.Resize(500, 300)
	c.SetTitle("Edit")

	c.OnShow = func() {
		c.txtText.Focus()
	}

	return &c
}
