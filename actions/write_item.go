package actions

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/dialogs"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type WriteItem struct {
	uicontrols.Panel
	client *client.Client

	txtDataItem *uicontrols.TextBoxExt
	txtValue    *uicontrols.TextBox
}

func NewWriteItem(parent uiinterfaces.Widget, client *client.Client) *WriteItem {
	var c WriteItem
	c.client = client
	c.InitControl(parent, &c)
	c.SetPanelPadding(0)

	c.AddTextBlockOnGrid(0, 0, "Data Item:")
	c.txtDataItem = c.AddTextBoxExtOnGrid(1, 0, "", func(textBoxExt *uicontrols.TextBoxExt) {
		dialogs.LookupDataItem(&c, c.client, "", "", func(key string) {
			textBoxExt.SetText(key)
		})
	})

	c.AddTextBlockOnGrid(0, 1, "Value:")
	c.txtValue = c.AddTextBoxOnGrid(1, 1)

	c.AddVSpacerOnGrid(0, 2)

	return &c
}

func (c *WriteItem) LoadAction(value string) {
	var a WriteItemAction
	err := json.Unmarshal([]byte(value), &a)
	if err != nil {
		return
	}
	c.txtDataItem.SetText(a.Item)
	c.txtValue.SetText(a.Value)
}

func (c *WriteItem) SaveAction() string {
	var a WriteItemAction
	a.Item = c.txtDataItem.Text()
	a.Value = c.txtValue.Text()
	bs, _ := json.MarshalIndent(a, "", " ")
	return string(bs)
}
