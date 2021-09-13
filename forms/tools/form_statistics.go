package tools

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"
)

type FormStatistics struct {
	uicontrols.Dialog
	client  *client.Client
	lvItems *uicontrols.ListView
	timer   *uievents.FormTimer

	itemLocalRcv  *uicontrols.ListViewItem
	itemLocalSend *uicontrols.ListViewItem
	itemCloudRcv  *uicontrols.ListViewItem
	itemCloudSend *uicontrols.ListViewItem
	itemApiCalls  *uicontrols.ListViewItem

	lastReceivedBytes int
	lastSentBytes     int
	lastStatDT        time.Time

	lastStatSrvCloudSend    int
	lastStatSrvCloudRvc     int
	lastStatSrvApiCallCount int
	lastStatSrvDT           time.Time
}

func NewFormStatistics(parent uiinterfaces.Widget, client *client.Client) *FormStatistics {
	var c FormStatistics
	c.client = client
	c.InitControl(parent, &c)
	return &c
}

func (c *FormStatistics) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("Statistics")
	c.Resize(500, 400)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)

	c.lvItems = pContent.AddListViewOnGrid(0, 1)
	c.lvItems.AddColumn("Name", 260)
	c.lvItems.AddColumn("Value", 100)

	pButtons.AddHSpacerOnGrid(0, 0)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Close", func(event *uievents.Event) {
		c.Reject()
	})
	btnCancel.SetMinWidth(70)

	c.SetRejectButton(btnCancel)

	c.itemLocalRcv = c.lvItems.AddItem("LocalService -> UI app")
	c.itemLocalSend = c.lvItems.AddItem("UI app -> LocalService")
	c.itemCloudRcv = c.lvItems.AddItem("Cloud(Internet) -> LocalService")
	c.itemCloudSend = c.lvItems.AddItem("LocalService -> Cloud(Internet)")
	c.itemApiCalls = c.lvItems.AddItem("API calls")

	c.timer = c.Window().NewTimer(1000, c.timerUpdate)
	c.timer.StartTimer()

	runtime.GC()
	debug.FreeOSMemory()
	runtime.GC()
	debug.FreeOSMemory()
	runtime.GC()
	debug.FreeOSMemory()
}

func (c *FormStatistics) timerUpdate() {
	t := time.Now().UTC()
	duration := t.Sub(c.lastStatDT).Seconds()
	speedSend := float64(0)
	speedReceive := float64(0)
	if duration > 0 {
		sent := client.StatSent()
		received := client.StatReceived()
		speedReceive = float64(received-c.lastReceivedBytes) / duration
		speedSend = float64(sent-c.lastSentBytes) / duration
		c.lastStatDT = time.Now().UTC()
		c.lastSentBytes = sent
		c.lastReceivedBytes = received
	}

	c.itemLocalSend.SetValue(1, strconv.FormatFloat(speedSend/1024, 'f', 1, 64)+" KB/sec")
	c.itemLocalRcv.SetValue(1, strconv.FormatFloat(speedReceive/1024, 'f', 1, 64)+" KB/sec")

	c.client.GetStatistics(func(statistics common_interfaces.Statistics, err error) {
		if err == nil {

			t := time.Now().UTC()
			duration := t.Sub(c.lastStatSrvDT).Seconds()
			speedCloudSend := float64(0)
			speedCloudReceive := float64(0)
			speedApiCalls := float64(0)
			if duration > 0 {
				sent := statistics.CloudSentBytes
				received := statistics.CloudReceivedBytes
				apiCalls := statistics.ApiCalls
				speedCloudSend = float64(sent-c.lastStatSrvCloudSend) / duration
				speedCloudReceive = float64(received-c.lastStatSrvCloudRvc) / duration
				speedApiCalls = float64(apiCalls-c.lastStatSrvApiCallCount) / duration
				c.lastStatSrvDT = time.Now().UTC()
				c.lastStatSrvCloudSend = sent
				c.lastStatSrvCloudRvc = received
				c.lastStatSrvApiCallCount = apiCalls
			}

			c.itemCloudSend.SetValue(1, strconv.FormatFloat(speedCloudSend/1024.0, 'f', 1, 64)+" KB/sec")
			c.itemCloudRcv.SetValue(1, strconv.FormatFloat(speedCloudReceive/1024.0, 'f', 1, 64)+" KB/sec")
			c.itemApiCalls.SetValue(1, strconv.FormatFloat(speedApiCalls, 'f', 1, 64)+" count/sec")
		} else {
		}
	})

}
