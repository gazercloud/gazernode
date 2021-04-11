package forms

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type FormWizard struct {
	uicontrols.Dialog
	client *client.Client
}

func NewFormWizard(parent uiinterfaces.Widget, client *client.Client) *FormWizard {
	var c FormWizard
	c.client = client
	c.InitControl(parent, &c)
	return &c
}

func (c *FormWizard) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("Setup Wizard")
	c.Resize(900, 500)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)

	pContent.AddButtonOnGrid(0, 0, "Ping to gazer.cloud", func(event *uievents.Event) {
		c.client.AddUnit("network_ping", "Ping to gazer.cloud", `{"addr": "gazer.cloud","frame_size": 64,"period": 1000,"timeout": 1000}`, func(s string, err error) {
			c.Accept()
		})
	})
	// {"addr": "localhost","frame_size": 64,"period": 1000,"timeout": 1000}
	pContent.AddButtonOnGrid(1, 0, "Operation system", func(event *uievents.Event) {
		c.client.AddUnit("windows_memory", "Windows Memory", "", func(s string, err error) {
			c.client.AddUnit("windows_network", "Windows Network", "", func(s string, err error) {
				c.client.AddUnit("windows_storage", "Windows Storage", "", func(s string, err error) {
					c.client.AddUnit("windows_process", "Explorer", `{"period": 1000,"process_name": "explorer.exe"}`, func(s string, err error) {
						c.Accept()
					})
				})
			})
		})
	})
	// {"period": 1000,"process_name": "explorer.exe"}
	pContent.AddButtonOnGrid(0, 1, "Network", func(event *uievents.Event) {
	})
	pContent.AddButtonOnGrid(1, 1, "Process", func(event *uievents.Event) {
	})

	pContent.AddVSpacerOnGrid(0, 2)

	pButtons.AddHSpacerOnGrid(0, 0)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Close", func(event *uievents.Event) {
		c.Reject()
	})
	btnCancel.SetMinWidth(70)

	c.SetRejectButton(btnCancel)
}
