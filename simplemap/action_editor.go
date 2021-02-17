package simplemap

import (
	"encoding/json"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type ActionEditor struct {
	uicontrols.Dialog
	resValue         string
	txtActionType    *uicontrols.TextBox
	txtActionContent *uicontrols.TextBox
	btnOK            *uicontrols.Button
}

func NewActionEditor(parent uiinterfaces.Widget, value string) *ActionEditor {
	var c ActionEditor
	c.resValue = value
	c.InitControl(parent, &c)

	var a Action
	_ = json.Unmarshal([]byte(value), &a)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)

	pContent.AddTextBlockOnGrid(0, 0, "Type:")
	c.txtActionType = pContent.AddTextBoxOnGrid(0, 1)
	c.txtActionType.SetText(a.Type)
	c.txtActionType.OnTextChanged = func(txtBox *uicontrols.TextBox, oldValue string, newValue string) {
		c.buildResult()
	}

	pContent.AddTextBlockOnGrid(0, 2, "Content:")
	c.txtActionContent = pContent.AddTextBoxOnGrid(0, 3)
	c.txtActionContent.SetText(a.Content)
	c.txtActionContent.SetMultiline(true)
	c.txtActionContent.OnTextChanged = func(txtBox *uicontrols.TextBox, oldValue string, newValue string) {
		c.buildResult()
	}

	pButtons.AddHSpacerOnGrid(0, 0)
	c.btnOK = pButtons.AddButtonOnGrid(1, 0, "OK", func(event *uievents.Event) {
		c.Accept()
	})
	c.TryAccept = func() bool {
		c.btnOK.SetEnabled(false)
		c.buildResult()
		c.TryAccept = nil
		c.Accept()
		return false
	}

	c.btnOK.SetMinWidth(70)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Cancel", func(event *uievents.Event) {
		c.Reject()
	})
	btnCancel.SetMinWidth(70)

	c.SetRejectButton(btnCancel)

	return &c
}

func (c *ActionEditor) buildResult() {
	var a Action
	a.Type = c.txtActionType.Text()
	a.Content = c.txtActionContent.Text()
	bs, _ := json.MarshalIndent(a, "", " ")
	c.resValue = string(bs)
}

func (c *ActionEditor) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("Edit action")
	c.Resize(400, 400)
}

func (c *ActionEditor) ActionText() string {
	return c.resValue
}
