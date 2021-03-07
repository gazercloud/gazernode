package actions

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type ActionEditorWidget interface {
	LoadAction(value string)
	SaveAction() string
}

type ActionEditor struct {
	uicontrols.Dialog
	client        *client.Client
	resValue      string
	currentAction *Action

	cmbActionType *uicontrols.ComboBox
	txtComment    *uicontrols.TextBox
	pEditor       *uicontrols.Panel
	wEditor       ActionEditorWidget
	pContent      *uicontrols.Panel
	btnOK         *uicontrols.Button
}

func NewActionEditor(parent uiinterfaces.Widget, value string, client *client.Client) *ActionEditor {
	var c ActionEditor
	c.client = client
	c.resValue = value
	c.InitControl(parent, &c)

	var a Action
	_ = json.Unmarshal([]byte(value), &a)

	c.currentAction = &a

	c.pContent = c.ContentPanel().AddPanelOnGrid(0, 0)
	c.pEditor = c.ContentPanel().AddPanelOnGrid(0, 1)
	//c.ContentPanel().AddVSpacerOnGrid(0, 2)
	pButtons := c.ContentPanel().AddPanelOnGrid(0, 3)

	c.pContent.AddTextBlockOnGrid(0, 0, "Type:")
	c.cmbActionType = c.pContent.AddComboBoxOnGrid(1, 0)
	c.cmbActionType.AddItem("write_item", "write_item")
	c.cmbActionType.AddItem("open_document", "open_document")
	c.cmbActionType.SetCurrentItemKey(a.Type)
	c.cmbActionType.OnCurrentIndexChanged = func(event *uicontrols.ComboBoxEvent) {
		c.buildResult()
		c.LoadEditor(c.currentAction.Type, c.currentAction.Content)
	}

	c.pContent.AddTextBlockOnGrid(0, 1, "Comment:")
	c.txtComment = c.pContent.AddTextBoxOnGrid(1, 1)
	c.txtComment.SetText(a.Comment)

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

	c.LoadEditor(a.Type, a.Content)

	return &c
}

func (c *ActionEditor) buildResult() {
	var a Action
	a.Type, _ = c.cmbActionType.CurrentItemKey().(string)
	a.Comment = c.txtComment.Text()
	if c.wEditor != nil {
		a.Content = c.wEditor.SaveAction()
	}
	c.currentAction = &a
	bs, _ := json.MarshalIndent(a, "", " ")
	c.resValue = string(bs)
}

func (c *ActionEditor) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("Edit action")
	c.Resize(600, 600)
}

func (c *ActionEditor) ActionText() string {
	return c.resValue
}

func (c *ActionEditor) LoadEditor(actionType string, actionContent string) {
	if c.wEditor != nil {
		c.pEditor.RemoveAllWidgets()
		c.wEditor = nil
	}

	switch actionType {
	case "open_document":
		c.wEditor = NewOpenMap(c.pContent, c.client)
	case "write_item":
		c.wEditor = NewWriteItem(c.pContent, c.client)
	}

	if c.wEditor != nil {
		c.pEditor.AddWidgetOnGrid(c.wEditor.(uiinterfaces.Widget), 0, 0)
		c.wEditor.LoadAction(actionContent)
	}
}
