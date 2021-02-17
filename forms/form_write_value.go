package forms

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
)

type FormWriteValue struct {
	uicontrols.Dialog
	client   *client.Client
	itemName string
}

func NewFormWriteValue(parent uiinterfaces.Widget, client *client.Client, itemName string) *FormWriteValue {
	var c FormWriteValue
	c.client = client
	c.itemName = itemName
	c.InitControl(parent, &c)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	pLeft := pContent.AddPanelOnGrid(0, 0)
	pLeft.SetPanelPadding(0)
	//pLeft.SetBorderRight(1, c.ForeColor())
	pLeft.SetMinWidth(100)
	pRight := pContent.AddPanelOnGrid(1, 0)
	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)

	img := pLeft.AddImageBoxOnGrid(0, 0, uiresources.ResImageAdjusted("icons/material/image/drawable-hdpi/ic_blur_on_black_48dp.png", c.ForeColor()))
	img.SetScaling(uicontrols.ImageBoxScaleAdjustImageKeepAspectRatio)
	img.SetMinHeight(64)
	img.SetMinWidth(64)
	pLeft.AddVSpacerOnGrid(0, 1)

	pRight.AddTextBlockOnGrid(0, 0, "Value:")
	txtValue := pRight.AddTextBoxOnGrid(1, 0)
	pRight.AddVSpacerOnGrid(0, 1)

	pButtons.AddHSpacerOnGrid(0, 0)
	btnOK := pButtons.AddButtonOnGrid(1, 0, "OK", nil)
	btnOK.SetMinWidth(70)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Close", nil)
	btnCancel.SetMinWidth(70)

	c.TryAccept = func() bool {
		c.client.Write(c.itemName, txtValue.Text(), func(err error) {
			if err == nil {
				c.TryAccept = nil
				c.Accept()
			}
		})
		return false
	}

	c.SetAcceptButton(btnOK)
	c.SetRejectButton(btnCancel)

	{
		c.client.GetItemsValues([]string{c.itemName}, func(items []common_interfaces.ItemGetUnitItems, err error) {
			if err == nil {
				for _, item := range items {
					if item.Name == c.itemName {
						txtValue.SetText(item.Value.Value)
						txtValue.SelectAllText()
						txtValue.Focus()
					}
				}
			}
		})
	}

	return &c
}

func (c *FormWriteValue) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("Write value")
	c.Resize(800, 400)
}
