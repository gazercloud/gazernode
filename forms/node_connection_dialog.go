package forms

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
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

func NewNodeConnectionDialog(parent uiinterfaces.Widget, client *client.Client) *NodeConnectionDialog {
	var c NodeConnectionDialog
	c.client = client
	c.InitControl(parent, &c)

	if c.client == nil || c.client.SessionToken() == "" {
		c.SetTitle("Connection to node")
		c.Resize(400, 200)

		pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
		pLeft := pContent.AddPanelOnGrid(0, 0)
		pLeft.SetPanelPadding(0)
		pLeft.SetBorderRight(1, c.ForeColor())
		imgBox := pLeft.AddImageBoxOnGrid(0, 0, uiresources.ResImgCol(uiresources.R_icons_material4_png_action_perm_identity_materialiconsoutlined_48dp_1x_outline_perm_identity_black_48dp_png, c.ForeColor()))
		imgBox.SetMinWidth(64)
		imgBox.SetMinHeight(64)
		imgBox.SetScaling(uicontrols.ImageBoxScaleAdjustImageKeepAspectRatio)
		pRight := pContent.AddPanelOnGrid(1, 0)
		pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)

		pLeft.AddVSpacerOnGrid(0, 1)

		pRight.AddTextBlockOnGrid(0, 0, "Address:")
		c.txtAddress = pRight.AddTextBoxOnGrid(1, 0)
		c.txtAddress.SetText("localhost")
		c.txtAddress.SetTabIndex(1)

		if c.client != nil {
			c.txtAddress.SetText(c.client.Address())
			c.txtAddress.SetEnabled(false)
		} else {
			c.txtAddress.SelectAllText()
		}

		pRight.AddTextBlockOnGrid(0, 1, "User name:")
		c.txtUserName = pRight.AddTextBoxOnGrid(1, 1)
		c.txtUserName.SetTabIndex(2)
		pRight.AddTextBlockOnGrid(0, 2, "Password:")
		c.txtPassword = pRight.AddTextBoxOnGrid(1, 2)
		c.txtPassword.SetTabIndex(3)

		pRight.AddVSpacerOnGrid(0, 10)
		pButtons.AddHSpacerOnGrid(0, 0)
		c.btnOK = pButtons.AddButtonOnGrid(1, 0, "OK", nil)
		c.btnOK.SetTabIndex(4)
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
		btnCancel.SetTabIndex(5)

		c.SetAcceptButton(c.btnOK)
		c.SetRejectButton(btnCancel)
		if c.client != nil {
			c.txtUserName.Focus()
		} else {
			c.txtAddress.Focus()
		}
	} else {
		c.SetTitle("Connection to node")
		c.Resize(450, 300)

		pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
		pLeft := pContent.AddPanelOnGrid(0, 0)
		pLeft.SetPanelPadding(0)
		pLeft.SetBorderRight(1, c.ForeColor())
		imgBox := pLeft.AddImageBoxOnGrid(0, 0, uiresources.ResImgCol(uiresources.R_icons_material4_png_action_perm_identity_materialiconsoutlined_48dp_1x_outline_perm_identity_black_48dp_png, c.ForeColor()))
		imgBox.SetMinWidth(64)
		imgBox.SetMinHeight(64)
		imgBox.SetScaling(uicontrols.ImageBoxScaleAdjustImageKeepAspectRatio)
		pRight := pContent.AddPanelOnGrid(1, 0)

		pLeft.AddVSpacerOnGrid(0, 1)

		pRight.AddTextBlockOnGrid(0, 0, "Address:")
		pRight.AddTextBlockOnGrid(1, 0, c.client.Address())

		pRight.AddTextBlockOnGrid(0, 1, "User name:")
		pRight.AddTextBlockOnGrid(1, 1, c.client.UserName())

		pRight.AddButtonOnGrid(1, 2, "Logout", func(event *uievents.Event) {
			uicontrols.ShowQuestionMessageOKCancel(pRight, "Do you want to logout?", "Confirmation", func() {
				c.client.SessionRemove(c.client.SessionToken(), func(response nodeinterface.SessionRemoveResponse, err error) {
					c.Reject()
				})
			}, nil)
		})
		pRight.AddHSpacerOnGrid(1, 3)

		pRight.AddVSpacerOnGrid(0, 10)
		//c.txtAddress.Focus()
	}

	return &c
}
