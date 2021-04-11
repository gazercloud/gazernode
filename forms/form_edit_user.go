package forms

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
)

type FormEditUser struct {
	uicontrols.Dialog
	client      *client.Client
	txtUnitName *uicontrols.TextBox
	txtPassword *uicontrols.TextBox
}

func NewFormEditUser(parent uiinterfaces.Widget, client *client.Client, name string) *FormEditUser {
	var c FormEditUser
	c.client = client
	c.InitControl(parent, &c)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	pLeft := pContent.AddPanelOnGrid(0, 0)
	pLeft.SetPanelPadding(0)
	//pLeft.SetBorderRight(1, c.ForeColor())
	pLeft.SetMinWidth(100)
	pRight := pContent.AddPanelOnGrid(1, 0)
	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)

	img := pLeft.AddImageBoxOnGrid(0, 0, uiresources.ResImgCol(uiresources.R_icons_material4_png_action_perm_identity_materialiconsoutlined_48dp_1x_outline_perm_identity_black_48dp_png, c.ForeColor()))
	img.SetScaling(uicontrols.ImageBoxScaleAdjustImageKeepAspectRatio)
	img.SetMinHeight(64)
	img.SetMinWidth(64)
	pLeft.AddVSpacerOnGrid(0, 1)

	pRight.AddTextBlockOnGrid(0, 0, "User name:")
	c.txtUnitName = pRight.AddTextBoxOnGrid(1, 0)
	c.txtUnitName.SetText(name)
	c.txtUnitName.SetReadOnly(true)
	pRight.AddTextBlockOnGrid(0, 1, "Password:")
	c.txtPassword = pRight.AddTextBoxOnGrid(1, 1)
	c.txtPassword.SetIsPassword(true)

	pRight.AddVSpacerOnGrid(0, 10)

	pButtons.AddHSpacerOnGrid(0, 0)
	btnOK := pButtons.AddButtonOnGrid(1, 0, "OK", nil)
	btnOK.SetMinWidth(70)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Cancel", nil)
	btnCancel.SetMinWidth(70)

	c.TryAccept = func() bool {
		c.client.UserSetPassword(c.txtUnitName.Text(), c.txtPassword.Text(), func(response nodeinterface.UserSetPasswordResponse, err error) {
			if err == nil {
				c.TryAccept = nil
				c.Accept()
			} else {
				uicontrols.ShowErrorMessage(&c, err.Error(), "error")
			}
		})
		return false
	}

	c.SetAcceptButton(btnOK)
	c.SetRejectButton(btnCancel)

	c.OnShow = func() {
		c.txtPassword.Focus()
	}

	return &c
}

func (c *FormEditUser) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("Edit user")
	c.Resize(400, 200)
}
