package forms

import (
	"crypto/tls"
	"encoding/json"
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/local_user_storage"
	"github.com/gazercloud/gazernode/product/productinfo"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiforms"
	"github.com/gazercloud/gazerui/uistyles"
	"github.com/go-gl/glfw/v3.3/glfw"
	"io/ioutil"
	"net/http"
)

type MainForm struct {
	uiforms.Form
	tabNodes    *uicontrols.TabControl
	nodeWidgets []*PanelNode
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

var MainFormInstance *MainForm

func (c *MainForm) OnInit() {
	MainFormInstance = c
	c.SetTitle("Gazer " + productinfo.Version())
	c.SetIcon(productinfo.Icon())

	c.nodeWidgets = make([]*PanelNode, 0)

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
		c.AddNode()
	}
	c.tabNodes.OnNeedClose = func(index int) {
		c.RemoveNode(index)
	}
	c.loadNodes()

	c.SetTheme(c.GetTheme())

	c.tabNodes.SetCurrentPage(0)

	go updateAdFromSite()

	//MainFormInstance.SetTheme(MainFormInstance.GetTheme())
}

func (c *MainForm) Dispose() {
	c.Form.Dispose()
}

func (c *MainForm) AddNode() {
	dialog := NewNodeConnectionDialog(c.Panel(), nil)
	dialog.OnAccept = func() {
		// Add to preferences
		var conn local_user_storage.NodeConnection
		conn.Address = dialog.Connection.Address
		conn.UserName = dialog.Connection.UserName
		conn.SessionToken = ""
		local_user_storage.Instance().AddConnection(conn)
		connIndex := local_user_storage.Instance().ConnectionCount() - 1

		// Add connection tab
		cl := client.New(c, dialog.Connection.Address, dialog.Connection.UserName, dialog.Connection.Password)
		c.addNodeTab(cl, connIndex)
	}
	dialog.ShowDialog()
}

func (c *MainForm) RemoveNode(index int) {
	c.nodeWidgets = append(c.nodeWidgets[:index], c.nodeWidgets[index+1:]...)
	c.tabNodes.RemovePage(index)
	local_user_storage.Instance().RemoveConnection(index)
}

func (c *MainForm) addNodeTab(cl *client.Client, index int) {
	page := c.tabNodes.AddPage()
	page.SetPanelPadding(0)
	page.SetText("  " + cl.Address() + "  ")
	panelNode := NewPanelNode(page, cl, index)
	page.AddWidgetOnGrid(panelNode, 0, 0)
	c.nodeWidgets = append(c.nodeWidgets, panelNode)
}

func (c *MainForm) loadNodes() {
	connections := local_user_storage.Instance().Connections()
	for i, conn := range connections {
		cl := client.NewWithSessionToken(c, conn.Address, conn.UserName, conn.SessionToken)
		c.addNodeTab(cl, i)
	}

	if len(connections) == 0 {
		c.AddNode()
	}
}

func (c *MainForm) GetTheme() int {
	theme := local_user_storage.Instance().Theme()
	if theme == "light" {
		return uistyles.StyleLight
	}
	if theme == "dark_blue" {
		return uistyles.StyleDarkBlue
	}
	return uistyles.StyleDarkBlue
}

func (c *MainForm) ShowFullScreenValue(show bool, itemId string) {
	//c.panelNode.ShowFullScreenValue(show, itemId)
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

	local_user_storage.Instance().SetTheme(themeStr)
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
