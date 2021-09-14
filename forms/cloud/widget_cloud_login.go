package cloud

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"golang.org/x/image/colornames"
)

type WidgetCloudLogin struct {
	uicontrols.Panel
	client *client.Client

	userName    string
	txtEMail    *uicontrols.TextBox
	txtPassword *uicontrols.TextBox
	lblStatus   *uicontrols.TextBlock
	btnLogin    *uicontrols.Button

	OnNeedToLoadState func()
}

func NewWidgetCloudLogin(parent uiinterfaces.Widget, client *client.Client) *WidgetCloudLogin {
	var c WidgetCloudLogin
	c.client = client
	c.InitControl(parent, &c)
	c.SetPanelPadding(0)
	return &c
}

func (c *WidgetCloudLogin) OnInit() {
	pContent := c.AddPanelOnGrid(0, 1)
	pContent.SetPanelPadding(20)

	pLoginForm := pContent.AddPanelOnGrid(0, 0)
	pLoginForm.SetPanelPadding(20)
	pLoginForm.SetBorders(1, colornames.Orange)
	pLoginForm.SetMinWidth(300)
	pLoginForm.SetMaxWidth(500)
	pLoginForm.AddTextBlockOnGrid(0, 0, "E-Mail:")
	c.txtEMail = pLoginForm.AddTextBoxOnGrid(1, 0)
	pLoginForm.AddTextBlockOnGrid(0, 1, "Password:")
	c.txtPassword = pLoginForm.AddTextBoxOnGrid(1, 1)
	c.txtPassword.SetIsPassword(true)
	c.btnLogin = pLoginForm.AddButtonOnGrid(1, 2, "Login", func(event *uievents.Event) {
		c.client.CloudLogin(c.txtEMail.Text(), c.txtPassword.Text(), func(err error) {
			if c.OnNeedToLoadState != nil {
				c.OnNeedToLoadState()
			}
		})
	})
	c.lblStatus = pLoginForm.AddTextBlockOnGrid(1, 3, "-")

	c.AddVSpacerOnGrid(0, 2)

	c.UpdateStyle()
}

func (c *WidgetCloudLogin) Dispose() {
	c.client = nil

	c.Panel.Dispose()
}

func (c *WidgetCloudLogin) SetState(response nodeinterface.CloudStateResponse) {
	if c.userName != response.UserName {
		c.userName = response.UserName
		c.txtEMail.SetText(c.userName)
	}

	if response.LoginStatus == "processing" {
		c.btnLogin.SetEnabled(false)
	} else {
		c.btnLogin.SetEnabled(true)
	}

	c.lblStatus.SetText(response.LoginStatus)
}
