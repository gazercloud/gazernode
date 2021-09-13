package tools

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
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
	c.Resize(500, 400)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)

	// {"addr": "localhost","frame_size": 64,"period": 1000,"timeout": 1000}

	var btn *uicontrols.Button

	btn = pContent.AddButtonOnGrid(0, 0, "OS Resources\r\nMonitoring", func(event *uievents.Event) {
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
	btn.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_hardware_computer_materialiconsoutlined_48dp_1x_outline_computer_black_48dp_png, c.ForeColor()))
	btn.SetImageSize(64, 64)
	btn.SetMouseCursor(ui.MouseCursorPointer)

	btn = pContent.AddButtonOnGrid(1, 0, "Memory\r\nMonitoring", func(event *uievents.Event) {
		c.client.AddUnit("windows_memory", "Windows Memory", "", func(s string, err error) {
			c.Accept()
		})
	})
	btn.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_hardware_memory_materialiconsoutlined_48dp_1x_outline_memory_black_48dp_png, c.ForeColor()))
	btn.SetImageSize(64, 64)
	btn.SetMouseCursor(ui.MouseCursorPointer)

	btn = pContent.AddButtonOnGrid(0, 1, "Networks\r\nMonitoring", func(event *uievents.Event) {
		c.client.AddUnit("windows_network", "Windows Networks", "", func(s string, err error) {
			c.Accept()
		})
	})
	btn.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_notification_network_check_materialiconsoutlined_48dp_1x_outline_network_check_black_48dp_png, c.ForeColor()))
	btn.SetImageSize(64, 64)
	btn.SetMouseCursor(ui.MouseCursorPointer)

	btn = pContent.AddButtonOnGrid(1, 1, "Storage\r\nMonitoring", func(event *uievents.Event) {
		c.client.AddUnit("windows_storage", "Windows Storage", "", func(s string, err error) {
			c.Accept()
		})
	})
	btn.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_device_storage_materialiconsoutlined_48dp_1x_outline_storage_black_48dp_png, c.ForeColor()))
	btn.SetImageSize(64, 64)
	btn.SetMouseCursor(ui.MouseCursorPointer)

	pContent.AddVSpacerOnGrid(0, 2)

	pButtons.AddHSpacerOnGrid(0, 0)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Close", func(event *uievents.Event) {
		c.Reject()
	})
	btnCancel.SetMinWidth(70)

	c.SetRejectButton(btnCancel)
}
