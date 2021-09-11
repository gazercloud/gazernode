package forms

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/cloud_accounts"
	"github.com/gazercloud/gazernode/crunner"
	"github.com/gazercloud/gazernode/home_client"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
	"time"
)

type NodeConnectionInfo struct {
	Transport string
	Address   string
	UserName  string
	Password  string
}

type NodeConnectionDialog struct {
	uicontrols.Dialog
	client *client.Client
	Id     string
	tp     string
	first  bool

	// Panel Transport
	panelButtonsTransport  *uicontrols.Panel
	btnTransportLocal      *uicontrols.Button
	btnTransportCloud      *uicontrols.Button
	currentTransportButton *uicontrols.Button

	// Local Transport
	panelLocalTransport *uicontrols.Panel
	txtAddress          *uicontrols.TextBox
	txtUserName         *uicontrols.TextBox
	txtPassword         *uicontrols.TextBox

	// Cloud Transport
	panelCloudTransport *uicontrols.Panel
	cmbCloudAccounts    *uicontrols.ComboBox
	lvNodes             *uicontrols.ListView
	selectedAccount     *cloud_accounts.CloudAccount
	selectedNode        string
	selectedSessionKey  string

	// Ok & Cancel Buttons
	btnOK  *uicontrols.Button
	runner *crunner.CRunner

	Connection NodeConnectionInfo
}

func OpenSessionInDialog(parent uiinterfaces.Widget, OnConnected func(cl *client.Client)) {
	dialog := NewNodeConnectionDialog(parent, nil, false)
	dialog.runner = crunner.New(parent.Window())
	dialog.ShowDialog()
	dialog.OnAccept = func() {
		if OnConnected != nil {
			OnConnected(dialog.client)
		}
	}
}

func NewNodeConnectionDialog(parent uiinterfaces.Widget, cl *client.Client, first bool) *NodeConnectionDialog {
	var c NodeConnectionDialog
	c.client = cl
	c.first = first
	c.InitControl(parent, &c)

	if c.client == nil || c.client.SessionToken() == "" {
		c.SetTitle("Connect to node")
		c.Resize(450, 500)

		pContent := c.ContentPanel().AddPanelOnGrid(0, 0)

		// Image
		pLeft := pContent.AddPanelOnGrid(0, 0)
		pLeft.SetPanelPadding(0)
		pLeft.SetBorderRight(1, c.ForeColor())
		imgBox := pLeft.AddImageBoxOnGrid(0, 0, uiresources.ResImgCol(uiresources.R_icons_material4_png_action_perm_identity_materialiconsoutlined_48dp_1x_outline_perm_identity_black_48dp_png, c.ForeColor()))
		imgBox.SetMinWidth(64)
		imgBox.SetMinHeight(64)
		imgBox.SetScaling(uicontrols.ImageBoxScaleAdjustImageKeepAspectRatio)
		pLeft.AddVSpacerOnGrid(0, 1)

		// Main
		pRight := pContent.AddPanelOnGrid(1, 0)
		pRight.SetPanelPadding(0)
		c.panelButtonsTransport = pRight.AddPanelOnGrid(0, 0)
		c.panelLocalTransport = pRight.AddPanelOnGrid(0, 1)
		c.panelCloudTransport = pRight.AddPanelOnGrid(0, 2)

		// Transport Buttons
		c.btnTransportLocal = c.panelButtonsTransport.AddButtonOnGrid(0, 0, "DIRECT", func(event *uievents.Event) {
			c.currentTransportButton = c.btnTransportLocal
			c.updateButtons()
		})
		c.btnTransportCloud = c.panelButtonsTransport.AddButtonOnGrid(1, 0, "CLOUD", func(event *uievents.Event) {
			c.currentTransportButton = c.btnTransportCloud
			c.updateButtons()
		})

		// Local Transport

		c.panelLocalTransport.AddTextBlockOnGrid(0, 1, "Address:")
		c.txtAddress = c.panelLocalTransport.AddTextBoxOnGrid(1, 1)
		c.txtAddress.SetText("localhost")
		c.txtAddress.SetTabIndex(1)

		c.panelLocalTransport.AddTextBlockOnGrid(0, 2, "User name:")
		c.txtUserName = c.panelLocalTransport.AddTextBoxOnGrid(1, 2)
		c.txtUserName.SetTabIndex(2)

		c.panelLocalTransport.AddTextBlockOnGrid(0, 3, "Password:")
		c.txtPassword = c.panelLocalTransport.AddTextBoxOnGrid(1, 3)
		c.txtPassword.SetTabIndex(3)
		c.txtPassword.SetIsPassword(true)

		c.panelLocalTransport.AddVSpacerOnGrid(0, 10)

		c.cmbCloudAccounts = c.panelCloudTransport.AddComboBoxOnGrid(0, 0)
		for _, ac := range cloud_accounts.List() {
			c.cmbCloudAccounts.AddItem(ac.Email, ac)
		}
		c.cmbCloudAccounts.OnCurrentIndexChanged = func(event *uicontrols.ComboBoxEvent) {
			c.lvNodes.RemoveItems()
			c.loadNodes()
		}

		c.lvNodes = c.panelCloudTransport.AddListViewOnGrid(0, 1)
		c.lvNodes.AddColumn("NodeId", 100)
		c.lvNodes.AddColumn("NodeName", 200)
		c.lvNodes.OnSelectionChanged = func() {
			if len(c.lvNodes.SelectedItems()) != 1 {
				c.selectedNode = ""
				return
			}

			c.selectedNode = c.lvNodes.SelectedItems()[0].Value(0)
		}

		/*if c.client != nil {
			if c.client.Transport() == string(client.TransportTypeLocal) {
				c.currentTransportButton = c.btnTransportLocal
			}
			if c.client.Transport() == string(client.TransportTypeCloudHttps) || c.client.Transport() == string(client.TransportTypeCloudBin) {
				c.currentTransportButton = c.btnTransportCloud
			}

			c.txtAddress.SetText(c.client.Address())
			c.txtAddress.SetEnabled(false)
		} else {
			c.txtAddress.SelectAllText()
			if first {
				c.txtUserName.SetText("admin")
				c.txtPassword.SetText("admin")
			}
		}*/

		pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)
		pButtons.AddHSpacerOnGrid(0, 0)
		c.btnOK = pButtons.AddButtonOnGrid(1, 0, "OK", nil)
		c.btnOK.SetTabIndex(4)
		c.TryAccept = func() bool {
			c.btnOK.SetEnabled(false)
			return c.makeClient()
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

		c.btnTransportCloud.Press()
		if len(c.cmbCloudAccounts.Items) > 0 {
			c.cmbCloudAccounts.SetCurrentItemIndex(0)
		}

		c.updateButtons()
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

func (c *NodeConnectionDialog) updateButtons() {
	if c.currentTransportButton == c.btnTransportLocal {
		c.btnTransportLocal.SetForeColor(c.AccentColor())
		c.btnTransportCloud.SetForeColor(nil)
		c.panelLocalTransport.SetVisible(true)
		c.panelCloudTransport.SetVisible(false)
	}
	if c.currentTransportButton == c.btnTransportCloud {
		c.btnTransportLocal.SetForeColor(nil)
		c.btnTransportCloud.SetForeColor(c.AccentColor())
		c.panelLocalTransport.SetVisible(false)
		c.panelCloudTransport.SetVisible(true)
	}
}

func (c *NodeConnectionDialog) transport() string {
	if c.currentTransportButton == c.btnTransportLocal {
		return string(client.TransportTypeLocal)
	}
	if c.currentTransportButton == c.btnTransportCloud {
		return string(client.TransportTypeCloudBin)
	}
	return ""
}

func (c *NodeConnectionDialog) loadNodes() {
	var err error
	var ac *cloud_accounts.CloudAccount
	if len(c.cmbCloudAccounts.Items) > 0 {
		if c.cmbCloudAccounts.CurrentItemIndex >= 0 && c.cmbCloudAccounts.CurrentItemIndex < len(c.cmbCloudAccounts.Items) {
			var ok bool
			ac, ok = c.cmbCloudAccounts.CurrentItemKey().(*cloud_accounts.CloudAccount)
			if ok {
				ac = c.cmbCloudAccounts.CurrentItemKey().(*cloud_accounts.CloudAccount)
			}
		}
	}

	c.selectedAccount = ac

	if ac == nil {
		return
	}

	c.lvNodes.RemoveItems()
	cl := home_client.NewWithSession(ac.SessionKey)
	c.selectedSessionKey, err = cl.SessionActivate()
	var regNodesBytes []byte
	regNodesBytes, err = cl.Call("s-registered-nodes", nil)
	if err != nil {
		c.lvNodes.AddItem(err.Error())
		return
	}

	type NodeResponseItem struct {
		Id              string `json:"id"`
		Name            string `json:"name"`
		CurrentRepeater string `json:"current_repeater"`
	}

	type NodesResponse struct {
		Items []NodeResponseItem `json:"items"`
	}

	var nodesResp NodesResponse

	err = json.Unmarshal(regNodesBytes, &nodesResp)
	if err != nil {
		c.lvNodes.AddItem(err.Error())
		return
	}

	for _, node := range nodesResp.Items {
		c.lvNodes.AddItem2(node.Id, node.Name)
	}
}

func (c *NodeConnectionDialog) makeClient() bool {
	if c.client != nil && c.selectedSessionKey != "" {
		return true
	}

	c.runner.Call(func(thParameters interface{}) (interface{}, error) {
		if c.transport() == string(client.TransportTypeCloudBin) {
			c.client = client.NewWithSessionToken(c.Window(), c.selectedNode, c.selectedAccount.Email, c.selectedSessionKey, c.transport())
		}

		if c.transport() == string(client.TransportTypeLocal) {
			c.client = client.New(c.Window(), c.txtAddress.Text(), c.txtUserName.Text(), c.txtPassword.Text(), c.transport())
			c.client.SessionOpen(c.txtUserName.Text(), c.txtPassword.Text(), nil)
			for i := 0; i < 30; i++ {
				time.Sleep(100 * time.Millisecond)
				if c.client.SessionToken() != "" {
					break
				}
			}
			c.selectedSessionKey = c.client.SessionToken()
		}
		return "", nil
	}, func(result interface{}, err error) {

		c.Accept()
	})

	return false
}
