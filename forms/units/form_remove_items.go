package units

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type FormRemoveItems struct {
	uicontrols.Dialog
	client *client.Client
	items  []string
}

func NewFormRemoveItems(parent uiinterfaces.Widget, client *client.Client, items []string) *FormRemoveItems {
	var c FormRemoveItems
	c.client = client
	c.items = items
	c.InitControl(parent, &c)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	pLeft := pContent.AddPanelOnGrid(0, 0)
	pLeft.SetPanelPadding(0)
	//pLeft.SetBorderRight(1, c.ForeColor())
	pLeft.SetMinWidth(100)
	pLeft.AddTextBlockOnGrid(0, 0, "Remove selected items?")
	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)

	pLeft.AddVSpacerOnGrid(0, 1)

	pButtons.AddHSpacerOnGrid(0, 0)
	btnOK := pButtons.AddButtonOnGrid(1, 0, "OK", nil)
	btnOK.SetMinWidth(70)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Cancel", nil)
	btnCancel.SetMinWidth(70)

	c.TryAccept = func() bool {
		c.client.DataItemRemove(c.items, func(err error) {
			if err == nil {
				c.TryAccept = nil
				c.Accept()
			} else {
				uicontrols.ShowErrorMessage(&c, err.Error(), "Error")
			}
		})
		return false
	}

	c.SetAcceptButton(btnOK)
	c.SetRejectButton(btnCancel)

	return &c
}

func (c *FormRemoveItems) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("Remove data items")
	c.Resize(400, 200)
}
