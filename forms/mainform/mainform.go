package mainform

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/crunner"
	"github.com/gazercloud/gazernode/forms/tools"
	"github.com/gazercloud/gazernode/product/productinfo"
	"github.com/gazercloud/gazernode/utilities/paths"
	"github.com/gazercloud/gazernode/workspace"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiforms"
	"github.com/gazercloud/gazerui/uistyles"
	"github.com/go-gl/glfw/v3.3/glfw"
	"io/ioutil"
	"net/http"
	"os"
)

type MainForm struct {
	uiforms.Form
	tabNodes          *uicontrols.TabControl
	nodeWidgets       []*PanelNode
	currentNodeWidget *PanelNode

	loadingConnections            []workspace.NodeConnection
	currentConnectionLoadingIndex int
	runner                        *crunner.CRunner

	loadNodesDelay int
}

type AdFromSite struct {
	Text string `json:"text"`
	Url  string `json:"url"`
}

var adFromSite AdFromSite

func updateAdFromSite() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://gazer.cloud/download/ad.json")
	if err == nil {
		content, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(content, &adFromSite)
		resp.Body.Close()
	}
}

func (c *MainForm) OnInit() {
	tools.MainFormInstance = c
	c.loadNodesDelay = 0
	c.SetTitle("Gazer " + productinfo.Version())
	c.SetIcon(productinfo.Icon())

	c.nodeWidgets = make([]*PanelNode, 0)
	c.runner = crunner.New(c)

	winWidth := 1300
	winHeight := 700

	mon := glfw.GetPrimaryMonitor()
	if mon != nil {
		_, _, w, h := mon.GetWorkarea()
		winWidth = w - w/10
		winHeight = h - h/10
	}

	if winWidth < 1000 {
		winWidth = 1000
	}
	if winWidth > 1920 {
		winWidth = 1920
	}
	if winHeight < 500 {
		winHeight = 500
	}
	if winHeight > 1080 {
		winHeight = 1080
	}

	c.Resize(winWidth, winHeight)

	c.Panel().SetPanelPadding(0)

	c.tabNodes = c.Panel().AddTabControlOnGrid(0, 0)
	c.tabNodes.SetShowAddButton(true)
	c.tabNodes.OnAddButtonPressed = func() {
		c.AddNode(false)
	}
	c.tabNodes.OnNeedClose = func(index int) {
		uicontrols.ShowQuestionMessageOKCancel(c.Panel(), "Remove connection to node?", "Confirmation", func() {
			c.RemoveNode(index)
		}, nil)
	}

	c.tabNodes.OnPageSelected = func(index int) {
		if index >= 0 && index < len(c.nodeWidgets) {
			c.currentNodeWidget = c.nodeWidgets[index]
		}
	}

	c.SetTheme(c.GetTheme())

	c.tabNodes.SetCurrentPage(0)

	go updateAdFromSite()

	//MainFormInstance.SetTheme(MainFormInstance.GetTheme())
	//c.loadNodes()

	c.runner.Call(func(thParameters interface{}) (interface{}, error) {
		return nil, nil
	}, func(result interface{}, err error) {
		c.Maximize()

	})

	c.MakeTimerAndStart(100, func(timer *uievents.FormTimer) {
		if c.loadNodesDelay < 1 {
			c.loadNodesDelay++
			return
		}
		c.loadNodes()
		timer.StopTimer()
	})

}

func (c *MainForm) Dispose() {
	c.Form.Dispose()
}

func (c *MainForm) checkFirstOpening() {
	if len(c.nodeWidgets) == 0 {
		first := true
		_, err := os.Stat(paths.HomeGazerFolder() + "/was_launched")
		if err == nil {
			first = false
		}
		c.AddNode(first)
	}
}

func (c *MainForm) AddNode(first bool) {
	hasNodes := true
	if len(c.nodeWidgets) == 0 {
		hasNodes = false
	}

	OpenSessionInDialog(c.Panel(), hasNodes, first, func(cl *client.Client) {
		var conn workspace.NodeConnection
		conn.Transport = cl.Transport()
		conn.Address = cl.Address()
		conn.UserName = cl.UserName()
		conn.SessionToken = cl.SessionToken()
		workspace.Instance().AddConnection(conn)

		connIndex := workspace.Instance().ConnectionCount() - 1
		c.addNodeTab(cl, connIndex)
		_ = ioutil.WriteFile(paths.HomeGazerFolder()+"/was_launched", []byte(""), 0600)
	})

	/*dialog := NewNodeConnectionDialog(c.Panel(), nil, first)
	dialog.OnAccept = func() {

		// Add to preferences
		var conn local_user_storage.NodeConnection
		conn.Transport = dialog.Connection.Transport
		conn.Address = dialog.Connection.Address
		conn.UserName = dialog.Connection.UserName
		conn.SessionToken = ""
		local_user_storage.Instance().AddConnection(conn)
		connIndex := local_user_storage.Instance().ConnectionCount() - 1

		// Add connection tab
		cl := client.New(c, dialog.Connection.Address, dialog.Connection.UserName, dialog.Connection.Password, dialog.Connection.Transport)
		c.addNodeTab(cl, connIndex)
	}
	dialog.ShowDialog()*/
}

func (c *MainForm) RemoveNode(index int) {
	c.nodeWidgets = append(c.nodeWidgets[:index], c.nodeWidgets[index+1:]...)
	c.tabNodes.RemovePage(index)
	workspace.Instance().RemoveConnection(index)
}

func (c *MainForm) addNodeTab(cl *client.Client, index int) {
	page := c.tabNodes.AddPage()
	page.SetPanelPadding(0)
	page.SetText("  " + cl.Address() + "  ")
	panelNode := NewPanelNode(page, cl, index)
	page.AddWidgetOnGrid(panelNode, 0, 0)
	panelNode.OnNodeNameUpdated = func(nodeName string) {
		page.SetText(nodeName)
	}

	panelNode.OnNeedToConnect = func(userName string, nodeId string, sessionKey string) {
		newClient := client.NewWithSessionToken(c, nodeId, userName, sessionKey, string(client.TransportTypeCloudBin))

		var conn workspace.NodeConnection
		conn.Transport = newClient.Transport()
		conn.Address = newClient.Address()
		conn.UserName = newClient.UserName()
		conn.SessionToken = newClient.SessionToken()
		workspace.Instance().AddConnection(conn)

		connIndex := workspace.Instance().ConnectionCount() - 1
		c.addNodeTab(newClient, connIndex)
	}

	c.nodeWidgets = append(c.nodeWidgets, panelNode)
	c.tabNodes.SetCurrentPage(len(c.nodeWidgets) - 1)
}

func (c *MainForm) loadNodes() {
	c.currentConnectionLoadingIndex = 0
	c.loadingConnections = workspace.Instance().Connections()

	if len(c.loadingConnections) > 0 {

		closed := false

		loadingDialog := uicontrols.NewDialog(c.Panel(), "Loading nodes", 500, 500)
		txtProgress := loadingDialog.ContentPanel().AddTextBlockOnGrid(0, 0, "")
		txtProgress.SetXExpandable(true)
		loadingDialog.ContentPanel().AddVSpacerOnGrid(0, 1)
		loadingDialog.ShowDialog()
		loadingDialog.OnReject = func() {
			closed = true
		}

		c.MakeTimerAndStart(100, func(timer *uievents.FormTimer) {
			if c.currentConnectionLoadingIndex == len(c.loadingConnections) {
				if loadingDialog != nil && !closed {
					loadingDialog.Close()
				}
				loadingDialog = nil
				timer.StopTimer()
			} else {
				conn := c.loadingConnections[c.currentConnectionLoadingIndex]
				txtProgress.SetText(txtProgress.Text() + "loading node " + conn.String() + " (" + fmt.Sprint(c.currentConnectionLoadingIndex) + " / " + fmt.Sprint(len(c.loadingConnections)) + ")\r\n")
				cl := client.NewWithSessionToken(c, conn.Address, conn.UserName, conn.SessionToken, conn.Transport)
				c.addNodeTab(cl, c.currentConnectionLoadingIndex)
				c.currentConnectionLoadingIndex++
			}
		})
	} else {
		c.checkFirstOpening()
	}
}

func (c *MainForm) GetTheme() int {
	theme := workspace.Instance().Theme()
	if theme == "light" {
		return uistyles.StyleLight
	}
	if theme == "dark_blue" {
		return uistyles.StyleDarkBlue
	}
	return uistyles.StyleDarkBlue
}

func (c *MainForm) ShowFullScreenValue(show bool, itemId string) {
	if c.currentNodeWidget != nil {
		c.currentNodeWidget.ShowFullScreenValue(show, itemId)
	}
}

func (c *MainForm) SetTheme(theme int) {
	uistyles.CurrentStyle = theme
	c.UpdateStyle()

	for _, nodeWidget := range c.nodeWidgets {
		nodeWidget.StylizeButton()
	}

	themeStr := "dark_blue"
	if theme == uistyles.StyleLight {
		themeStr = "light"
	}
	if theme == uistyles.StyleDarkBlue {
		themeStr = "dark_blue"
	}

	workspace.Instance().SetTheme(themeStr)
}

func (c *MainForm) OnClose() bool {
	return true
}

func Substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}

func (c *MainForm) ShowChartGroup(resId string) {
	c.currentNodeWidget.btnPanelCharts.Press()
	c.currentNodeWidget.panelCharts.SelectChartGroup(resId)
}
