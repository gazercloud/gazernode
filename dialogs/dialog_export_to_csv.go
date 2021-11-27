package dialogs

import (
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/utilities/filedialog"
	"github.com/gazercloud/gazernode/utilities/paths"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"io/ioutil"
	"runtime"
	"strings"
)

type ExportToCSVTask struct {
	ItemName string
	Values   []*common_interfaces.ItemValue
}

type DialogExportToCSV struct {
	uicontrols.Dialog
	txtExample *uicontrols.TextBox
	task       ExportToCSVTask
	btnOK      *uicontrols.Button

	fieldSeparator string
	lineSeparator  string
}

func NewDialogExportToCSV(parent uiinterfaces.Widget, task ExportToCSVTask) *DialogExportToCSV {
	var c DialogExportToCSV
	c.task = task
	c.InitControl(parent, &c)
	c.Resize(750, 450)
	c.SetTitle("Export data")

	c.fieldSeparator = ";"
	c.lineSeparator = "\r\n"

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	pContent.AddTextBlockOnGrid(0, 0, "Example:")
	c.txtExample = pContent.AddTextBoxOnGrid(0, 1)
	c.txtExample.SetMultiline(true)

	pSettings := pContent.AddPanelOnGrid(0, 2)

	pFieldsSep := pSettings.AddPanelOnGrid(0, 0)
	pFieldsSep.AddTextBlockOnGrid(0, 0, "Fields Separator:")
	txtFieldsSep := pFieldsSep.AddTextBoxOnGrid(1, 0)
	txtFieldsSep.OnTextChanged = func(txtBox *uicontrols.TextBox, oldValue string, newValue string) {
		c.fieldSeparator = newValue
		c.updateExample()
	}
	txtFieldsSep.SetText(c.fieldSeparator)

	pSelectFile := pSettings.AddPanelOnGrid(0, 1)
	pSelectFile.AddTextBlockOnGrid(0, 0, "File Name:")
	txtFileName := pSelectFile.AddTextBoxOnGrid(1, 0)
	txtFileName.OnTextChanged = func(txtBox *uicontrols.TextBox, oldValue string, newValue string) {
		c.updateExample()
	}
	pSelectFile.AddButtonOnGrid(2, 0, "Select ...", func(event *uievents.Event) {
		fileName := strings.ReplaceAll(c.task.ItemName, "/", "_") + ".csv"
		dialog := filedialog.NewFileDialog(&c, true, fileName)
		dialog.ShowDialog()
		dialog.OnAccept = func() {
			txtFileName.SetText(dialog.FileName)
		}
	})

	defaultFileName := paths.DocumentsFolder()
	if runtime.GOOS == "windows" {
		defaultFileName += "\\"
	} else {
		defaultFileName += "/"
	}
	defaultFileName += strings.ReplaceAll(c.task.ItemName, "/", "_") + ".csv"

	txtFileName.SetText(defaultFileName)

	pButtons := pContent.AddPanelOnGrid(0, 3)
	pButtons.SetPanelPadding(0)

	pButtons.AddHSpacerOnGrid(0, 0)
	c.btnOK = pButtons.AddButtonOnGrid(1, 0, "Export", func(event *uievents.Event) {
		bs := []byte(c.export(0))
		err := ioutil.WriteFile(txtFileName.Text(), bs, 0644)
		if err != nil {
			uicontrols.ShowErrorMessage(&c, err.Error(), "error")
			return
		}
		c.Accept()
	})

	c.btnOK.SetMinWidth(70)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Cancel", func(event *uievents.Event) {
		c.Reject()
	})
	btnCancel.SetMinWidth(70)

	c.SetAcceptButton(c.btnOK)
	c.SetRejectButton(btnCancel)

	c.updateExample()

	return &c
}

func (c *DialogExportToCSV) updateExample() {
	str := c.export(5)
	c.txtExample.SetText(str)
}

func (c *DialogExportToCSV) export(maxLines int) string {
	var builder strings.Builder
	/*for i, item := range c.task.Values {
		dt := time.UnixMicro(item.DT)
		line := ""
		line += dt.Format("2006-01-02 15:04:05.999")
		line += c.fieldSeparator
		line += item.Value
		line += c.fieldSeparator
		line += item.UOM
		line += c.lineSeparator
		builder.WriteString(line)
		if maxLines > 0 && i >= maxLines {
			builder.WriteString("...")
			break
		}
	}*/

	return builder.String()
}
