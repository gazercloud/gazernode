package cloud

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/cloud_accounts"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type WidgetCloud struct {
	uicontrols.Panel
	client *client.Client
	timer  *uievents.FormTimer

	wHeader *WidgetCloudHeader
	wLogin  *WidgetCloudLogin
	wMain   *WidgetCloudMain

	currentState         nodeinterface.CloudStateResponse
	firstTimeStateLoaded bool

	OnNeedToConnect func(userName string, nodeId string, sessionKey string)
}

func NewWidgetCloud(parent uiinterfaces.Widget, client *client.Client) *WidgetCloud {
	var c WidgetCloud
	c.client = client
	c.InitControl(parent, &c)
	//c.SetPanelPadding(0)
	return &c
}

func (c *WidgetCloud) OnInit() {
	c.wHeader = NewWidgetCloudHeader(c, c.client)
	c.AddWidgetOnGrid(c.wHeader, 0, 0)
	c.wHeader.OnNeedToLoadState = func() {
		c.timerUpdate()
	}

	c.wLogin = NewWidgetCloudLogin(c, c.client)
	c.AddWidgetOnGrid(c.wLogin, 0, 1)
	c.wLogin.SetVisible(false)
	c.wLogin.OnNeedToLoadState = func() {
		c.timerUpdate()
	}

	c.wMain = NewWidgetCloudMain(c, c.client)
	c.AddWidgetOnGrid(c.wMain, 0, 1)
	c.wMain.SetVisible(false)

	c.wMain.OnNeedToConnect = func(userName string, nodeId string, sessionKey string) {
		if c.OnNeedToConnect != nil {
			c.OnNeedToConnect(userName, nodeId, sessionKey)
		}
	}

	c.timer = c.Window().NewTimer(1000, c.timerUpdate)
	c.timer.StartTimer()

	c.UpdateStyle()
	c.updateVisibility()
}

func (c *WidgetCloud) Dispose() {
	if c.timer != nil {
		c.timer.StopTimer()
	}
	c.timer = nil
	c.client = nil

	c.Panel.Dispose()
}

func (c *WidgetCloud) timerUpdate() {
	/*if !c.IsVisibleRec() && c.firstTimeStateLoaded {
		return
	}*/

	c.client.CloudState(func(response nodeinterface.CloudStateResponse, err error) {
		if err != nil {
			return
		}
		c.firstTimeStateLoaded = true
		c.currentState = response
		c.updateVisibility()
		c.wLogin.SetState(response)
		c.wMain.SetState(response)
		c.wHeader.SetState(response)

		cloud_accounts.Set(response.UserName, response.SessionKey)
	})
}

func (c *WidgetCloud) updateVisibility() {
	if c.wLogin == nil || c.wMain == nil {
		return
	}

	if c.currentState.LoggedIn {
		c.wLogin.SetVisible(false)
		c.wMain.SetVisible(true)
	} else {
		c.wLogin.SetVisible(true)
		c.wMain.SetVisible(false)
	}
}

func (c *WidgetCloud) IsSomethingWrong() bool {
	return c.wHeader.IsSomethingWrong()
}
