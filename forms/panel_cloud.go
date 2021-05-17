package forms

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"golang.org/x/image/colornames"
)

type PanelCloud struct {
	uicontrols.Panel
	client *client.Client
	timer  *uievents.FormTimer
}

func NewPanelCloud(parent uiinterfaces.Widget, client *client.Client) *PanelCloud {
	var c PanelCloud
	c.client = client
	c.InitControl(parent, &c)

	return &c
}

func (c *PanelCloud) OnInit() {
	pHeader := c.AddPanelOnGrid(0, 0)
	txtHeader := pHeader.AddTextBlockOnGrid(0, 0, "Cloud")
	txtHeader.SetFontSize(24)

	pContent := c.AddPanelOnGrid(0, 1)
	pContent.SetPanelPadding(20)

	pLoginForm := pContent.AddPanelOnGrid(0, 0)
	pLoginForm.SetPanelPadding(20)
	pLoginForm.SetBorders(1, colornames.Orange)
	pLoginForm.SetMinWidth(300)
	pLoginForm.SetMaxWidth(300)
	pLoginForm.AddTextBlockOnGrid(0, 0, "E-Mail:")
	pLoginForm.AddTextBoxOnGrid(1, 0)
	pLoginForm.AddTextBlockOnGrid(0, 1, "Password:")
	pLoginForm.AddTextBoxOnGrid(1, 1)
	pLoginForm.AddButtonOnGrid(1, 2, "Login", nil)

	pContent.AddHSpacerOnGrid(1, 0)

	pVSpacer := c.AddPanelOnGrid(0, 2)
	pVSpacer.SetPanelPadding(0)
	pVSpacer.AddVSpacerOnGrid(0, 0)

	c.timer = c.Window().NewTimer(500, c.timerUpdate)
	c.timer.StartTimer()

	c.UpdateStyle()
}

func (c *PanelCloud) Dispose() {
	if c.timer != nil {
		c.timer.StopTimer()
	}
	c.timer = nil

	c.client = nil

	c.Panel.Dispose()
}

func (c *PanelCloud) FullRefresh() {
}

func (c *PanelCloud) UpdateStyle() {
	c.Panel.UpdateStyle()

	//activeColor := c.AccentColor()
	//inactiveColor := c.InactiveColor()
}

func (c *PanelCloud) timerUpdate() {
	if !c.IsVisible() {
		return
	}
}
