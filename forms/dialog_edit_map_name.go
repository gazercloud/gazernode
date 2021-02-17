package forms

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiinterfaces"
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
	c.txtText = pContent.AddTextBoxOnGrid(0, 0)
	c.txtText.SetText(text)

	pContent.AddVSpacerOnGrid(0, 5)

	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)
	pButtons.AddHSpacerOnGrid(0, 0)
	c.btnOK = pButtons.AddButtonOnGrid(1, 0, "OK", nil)
	c.TryAccept = func() bool {
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

	return &c
}
