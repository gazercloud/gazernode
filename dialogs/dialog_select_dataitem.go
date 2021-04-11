package dialogs

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/widgets/widget_dataitems"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type DialogSelectDataItem struct {
	uicontrols.Dialog
	client               *client.Client
	widgetSelectDataItem *widget_dataitems.WidgetDataItems
	selectedItem         string
}

func NewDialogSelectDataItem(parent uiinterfaces.Widget, client *client.Client, text1 string, text2 string) *DialogSelectDataItem {
	var c DialogSelectDataItem
	c.client = client
	c.InitControl(parent, &c)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	c.widgetSelectDataItem = widget_dataitems.NewWidgetDataItems(pContent, c.client, text1, text2)
	pContent.AddWidgetOnGrid(c.widgetSelectDataItem, 0, 0)
	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)

	pButtons.AddHSpacerOnGrid(0, 0)
	btnOK := pButtons.AddButtonOnGrid(1, 0, "OK", nil)
	btnOK.SetMinWidth(70)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Cancel", nil)
	btnCancel.SetMinWidth(70)

	c.TryAccept = func() bool {
		c.selectedItem = c.widgetSelectDataItem.SelectedItem()
		return true
	}

	c.SetAcceptButton(btnOK)
	c.SetRejectButton(btnCancel)

	return &c
}

func (c *DialogSelectDataItem) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("Select data item")
	c.Resize(600, 600)
}

func LookupDataItem(parent uiinterfaces.Widget, client *client.Client, text1 string, text2 string, selected func(key string)) {
	dialog := NewDialogSelectDataItem(parent, client, text1, text2)
	dialog.ShowDialog()
	dialog.OnAccept = func() {
		if selected != nil {
			selected(dialog.selectedItem)
		}
	}
}
