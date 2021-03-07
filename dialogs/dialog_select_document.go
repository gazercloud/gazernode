package dialogs

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/widgets/widget_resources"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type DialogSelectDocument struct {
	uicontrols.Dialog
	client               *client.Client
	widgetSelectDataItem *widget_resources.WidgetResources
	selectedItem         string
}

func NewDialogSelectDocument(parent uiinterfaces.Widget, client *client.Client) *DialogSelectDocument {
	var c DialogSelectDocument
	c.client = client
	c.InitControl(parent, &c)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	c.widgetSelectDataItem = widget_resources.NewWidgetResources(pContent, c.client)
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

func (c *DialogSelectDocument) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("Select resource")
	c.Resize(600, 600)
}

func LookupResource(parent uiinterfaces.Widget, client *client.Client, selected func(key string)) {
	dialog := NewDialogSelectDocument(parent, client)
	dialog.ShowDialog()
	dialog.OnAccept = func() {
		if selected != nil {
			selected(dialog.selectedItem)
		}
	}
}
