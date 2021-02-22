package forms

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type NodeConnectionInfo struct {
	Address  string
	UserName string
	Password string
}

type NodeConnectionDialog struct {
	uicontrols.Dialog
	client      *client.Client
	Id          string
	tp          string
	txtAddress  *uicontrols.TextBox
	txtUserName *uicontrols.TextBox
	txtPassword *uicontrols.TextBox
	btnOK       *uicontrols.Button

	Connection NodeConnectionInfo
}

func NewNodeConnectionDialog(parent uiinterfaces.Widget) *NodeConnectionDialog {
	var c NodeConnectionDialog
	c.InitControl(parent, &c)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	pLeft := pContent.AddPanelOnGrid(0, 0)
	pLeft.SetPanelPadding(0)
	//pLeft.SetBorderRight(1, c.ForeColor())
	pLeft.SetMinWidth(100)
	pRight := pContent.AddPanelOnGrid(1, 0)
	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)

	pLeft.AddVSpacerOnGrid(0, 1)

	pRight.AddTextBlockOnGrid(0, 0, "Address:")
	c.txtAddress = pRight.AddTextBoxOnGrid(1, 0)
	c.txtAddress.SetText("localhost:8084")
	pRight.AddTextBlockOnGrid(0, 1, "User name:")
	c.txtUserName = pRight.AddTextBoxOnGrid(1, 1)
	pRight.AddTextBlockOnGrid(0, 2, "Password:")
	c.txtPassword = pRight.AddTextBoxOnGrid(1, 2)

	pRight.AddVSpacerOnGrid(0, 10)

	pButtons.AddHSpacerOnGrid(0, 0)
	c.btnOK = pButtons.AddButtonOnGrid(1, 0, "OK", nil)
	c.TryAccept = func() bool {
		c.btnOK.SetEnabled(false)
		c.Connection.Address = c.txtAddress.Text()
		c.Connection.UserName = c.txtUserName.Text()
		c.Connection.Password = c.txtPassword.Text()
		return true
	}

	c.btnOK.SetMinWidth(70)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Cancel", func(event *uievents.Event) {
		c.Reject()
	})
	btnCancel.SetMinWidth(70)

	c.SetAcceptButton(c.btnOK)
	c.SetRejectButton(btnCancel)

	return &c
}

func (c *NodeConnectionDialog) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("Connection to node")
	c.Resize(400, 200)
}
