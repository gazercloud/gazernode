package simplemap

import (
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type DataSourceFormatEditor struct {
	uicontrols.Dialog
	text        string
	txtFilePath *uicontrols.TextBox
	textInit    string

	panelDetails *uicontrols.Panel
}

func NewDataSourceFormatEditor(parent uiinterfaces.Widget, textInit string) *DataSourceFormatEditor {
	var c DataSourceFormatEditor
	c.InitControl(parent, &c)
	c.SetTitle("Select file")
	c.Resize(600, 400)
	c.textInit = textInit
	c.text = textInit

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	c.txtFilePath = pContent.AddTextBoxOnGrid(0, 0)
	c.txtFilePath.SetText(c.text)
	c.txtFilePath.OnTextChanged = func(txtBox *uicontrols.TextBox, oldValue string, newValue string) {
		c.text = newValue
	}
	pContent.AddTextBlockOnGrid(0, 1, "Format: \r\n{v} = value \r\n{uom} = unit of measure \r\n{d} = date \r\n{t} = time")
	pContent.AddVSpacerOnGrid(0, 5)

	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)
	pButtons.AddHSpacerOnGrid(0, 0)
	btnOK := pButtons.AddButtonOnGrid(1, 0, "OK", nil)
	c.TryAccept = func() bool {
		return true
	}

	btnOK.SetMinWidth(70)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Cancel", func(event *uievents.Event) {
		c.Reject()
	})
	btnCancel.SetMinWidth(70)

	c.SetAcceptButton(btnOK)
	c.SetRejectButton(btnCancel)

	return &c
}

func EditDataSourceFormat(parent uiinterfaces.Widget, initText string, selected func(filePath string)) {
	dialog := NewDataSourceFormatEditor(parent, initText)
	dialog.ShowDialog()
	dialog.OnAccept = func() {
		if selected != nil {
			selected(dialog.text)
		}
	}
}
