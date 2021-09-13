package tools

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/protocols/lookup"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"strings"
)

type FormLookup struct {
	uicontrols.Dialog
	client       *client.Client
	txtFilter    *uicontrols.TextBox
	lvItems      *uicontrols.ListView
	lookupResult lookup.Result
	selectedKey  string
}

func NewFormLookup(parent uiinterfaces.Widget, lookupResult lookup.Result) *FormLookup {
	var c FormLookup
	c.lookupResult = lookupResult
	c.InitControl(parent, &c)
	return &c
}

func (c *FormLookup) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("Lookup Dialog")
	c.Resize(500, 400)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)

	pFilter := pContent.AddPanelOnGrid(0, 0)
	pFilter.AddTextBlockOnGrid(0, 0, "Filter")
	c.txtFilter = pFilter.AddTextBoxOnGrid(1, 0)
	c.txtFilter.OnTextChanged = func(txtBox *uicontrols.TextBox, oldValue string, newValue string) {
		c.updateItems()
	}

	c.lvItems = pContent.AddListViewOnGrid(0, 1)

	pButtons.AddHSpacerOnGrid(0, 0)
	btnOK := pButtons.AddButtonOnGrid(1, 0, "OK", nil)
	btnOK.SetMinWidth(70)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Cancel", nil)
	btnCancel.SetMinWidth(70)
	c.SetAcceptButton(btnOK)
	c.SetRejectButton(btnCancel)

	c.TryAccept = func() bool {
		if c.lvItems.SelectedItem() != nil {
			c.TryAccept = nil
			c.selectedKey = c.lvItems.SelectedItem().TempData
			c.Accept()
		}
		return false
	}

	c.updateItems()
}

func (c *FormLookup) updateItems() {
	c.lvItems.RemoveColumns()
	c.lvItems.RemoveItems()

	keyColumnPos := 0

	colWidth := c.Width() - 50
	if len(c.lookupResult.Columns) > 0 {
		colWidth = (c.Width() - 50) / len(c.lookupResult.Columns)
	}

	for index, col := range c.lookupResult.Columns {
		c.lvItems.AddColumn(col.DisplayName, colWidth)
		if col.Name == c.lookupResult.KeyColumn {
			keyColumnPos = index
		}
	}

	for _, row := range c.lookupResult.Rows {
		inFilter := false

		if c.txtFilter.Text() != "" {
			for _, cell := range row.Cells {
				if strings.Contains(strings.ToLower(cell), strings.ToLower(c.txtFilter.Text())) {
					inFilter = true
				}
			}
		} else {
			inFilter = true
		}

		if inFilter {
			item := c.lvItems.AddItem("")
			for index, cell := range row.Cells {
				item.SetValue(index, cell)
				if keyColumnPos == index {
					item.TempData = cell
				}
			}
		}
	}
}

func LookupDialog(parent uiinterfaces.Widget, client *client.Client, entity string, selected func(key string)) {
	client.Lookup(entity, func(result lookup.Result, err error) {
		dialog := NewFormLookup(parent, result)
		dialog.ShowDialog()
		dialog.OnAccept = func() {
			if selected != nil {
				selected(dialog.selectedKey)
			}
		}
	})
}
