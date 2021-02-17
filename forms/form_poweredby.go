package forms

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type FormPoweredBy struct {
	uicontrols.Dialog
	lvItems *uicontrols.ListView
}

func NewFormPoweredBy(parent uiinterfaces.Widget) *FormPoweredBy {
	var c FormPoweredBy
	c.InitControl(parent, &c)
	return &c
}

func (c *FormPoweredBy) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("Powered By")
	c.Resize(900, 500)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)

	pButtonsTop := pContent.AddPanelOnGrid(0, 0)
	pButtonsTop.AddButtonOnGrid(0, 0, "Open the project website", func(event *uievents.Event) {
		for _, item := range c.lvItems.SelectedItems() {
			client.OpenBrowser("https://" + item.Value(0))
		}
	})
	pButtonsTop.AddButtonOnGrid(1, 0, "License text", func(event *uievents.Event) {
		for _, item := range c.lvItems.SelectedItems() {
			client.OpenBrowser(item.Value(2))
		}
	})
	pButtonsTop.AddHSpacerOnGrid(2, 0)

	c.lvItems = pContent.AddListViewOnGrid(0, 1)
	c.lvItems.AddColumn("Name", 260)
	c.lvItems.AddColumn("License", 100)
	c.lvItems.AddColumn("License URL", 470)

	pButtons.AddHSpacerOnGrid(0, 0)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Close", func(event *uievents.Event) {
		c.Reject()
	})
	btnCancel.SetMinWidth(70)

	c.SetRejectButton(btnCancel)

	c.lvItems.AddItem3("github.com/golang/go", "BSD", "https://github.com/golang/go/blob/master/LICENSE")
	c.lvItems.AddItem3("github.com/go-gl/glfw", "BSD", "https://github.com/go-gl/glfw/blob/master/LICENSE")
	c.lvItems.AddItem3("github.com/go-gl/gl", "MIT", "https://github.com/go-gl/gl/blob/master/LICENSE")
	c.lvItems.AddItem3("github.com/golang/freetype", "BSD", "https://github.com/golang/freetype/blob/master/LICENSE")
	c.lvItems.AddItem3("github.com/StackExchange/wmi", "MIT", "https://github.com/StackExchange/wmi/blob/master/LICENSE")
	c.lvItems.AddItem3("github.com/go-ole/go-ole", "MIT", "https://github.com/go-ole/go-ole/blob/master/LICENSE")
	c.lvItems.AddItem3("github.com/gorilla/css", "BSD", "https://github.com/gorilla/css/blob/master/LICENSE")
	c.lvItems.AddItem3("github.com/gorilla/mux", "BSD", "https://github.com/gorilla/mux/blob/master/LICENSE")
	c.lvItems.AddItem3("github.com/kardianos/osext", "BSD", "https://github.com/kardianos/osext/blob/master/LICENSE")
	c.lvItems.AddItem3("github.com/kardianos/service", "zlib", "https://github.com/kardianos/service/blob/master/LICENSE")
	c.lvItems.AddItem3("github.com/mattn/go-sqlite3", "MIT", "https://github.com/mattn/go-sqlite3/blob/master/LICENSE")
	c.lvItems.AddItem3("github.com/nfnt/resize", "ISC", "https://github.com/nfnt/resize/blob/master/LICENSE")
	c.lvItems.AddItem3("github.com/shirou/gopsutil", "BSD", "https://github.com/shirou/gopsutil/blob/master/LICENSE")
	c.lvItems.AddItem3("github.com/sparrc/go-ping", "MIT", "https://github.com/go-ping/ping/blob/master/LICENSE")
	c.lvItems.AddItem3("github.com/tarm/serial", "BSD", "https://github.com/tarm/serial/blob/master/LICENSE")
	//c.lvItems.AddItem3("", "", "")
}
