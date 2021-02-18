package forms

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/settings"
	"github.com/gazercloud/gazerui/canvas"
	"github.com/gazercloud/gazerui/go-gl/glfw/v3.3/glfw"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"math"
	"strconv"
	"strings"
	"time"
)

type PanelFullScreenValue struct {
	uicontrols.Panel
	client   *client.Client
	lblName  *uicontrols.TextBlock
	lblValue *uicontrols.TextBlock
	lblUOM   *uicontrols.TextBlock
	lblDT    *uicontrols.TextBlock
	btnClose *uicontrols.Button
	timer    *uievents.FormTimer
	itemId   string
}

func NewPanelFullScreenValue(parent uiinterfaces.Widget, client *client.Client, itemId string) *PanelFullScreenValue {
	var c PanelFullScreenValue
	c.client = client
	c.itemId = itemId
	c.InitControl(parent, &c)
	return &c
}

func (c *PanelFullScreenValue) OnInit() {
	/*pClose := c.AddPanelOnGrid(0, 0)
	pClose.AddHSpacerOnGrid(0, 0)
	c.btnClose = pClose.AddButtonOnGrid(1, 0, "Close (ESC)", func(event *uievents.Event) {
		MainFormInstance.ShowFullScreenValue(false, "")
	})
	c.btnClose.SetMaxWidth(200)
	c.btnClose.SetImage(uiresources.ResImageAdjusted("icons/material/navigation/drawable-hdpi/ic_fullscreen_exit_black_48dp.png", c.ForeColor()))
	c.btnClose.SetBorders(0, colornames.Aqua)*/

	c.lblName = c.AddTextBlockOnGrid(0, 1, "---")
	c.lblName.SetFontSize(36)
	c.lblName.TextHAlign = canvas.HAlignCenter

	c.lblValue = c.AddTextBlockOnGrid(0, 2, "---")
	c.lblValue.SetFontSize(144)
	c.lblValue.SetYExpandable(true)
	c.lblValue.TextHAlign = canvas.HAlignCenter

	c.lblUOM = c.AddTextBlockOnGrid(0, 3, "---")
	c.lblUOM.SetFontSize(72)
	c.lblUOM.TextHAlign = canvas.HAlignCenter

	c.lblDT = c.AddTextBlockOnGrid(0, 4, "---")
	c.lblDT.SetFontSize(16)
	c.lblDT.TextHAlign = canvas.HAlignCenter

	c.timer = c.Window().NewTimer(500, c.timerUpdate)
	c.timer.StartTimer()

	c.SetOnKeyDown(func(event *uievents.KeyDownEvent) bool {
		if event.Key == glfw.KeyEscape || event.Key == glfw.KeyEnter || event.Key == glfw.KeyKPEnter {
			MainFormInstance.ShowFullScreenValue(false, "")
			return true
		}
		return false
	})
}

func (c *PanelFullScreenValue) Dispose() {
	if c.timer != nil {
		c.timer.StopTimer()
	}
	c.timer = nil

	c.client = nil
	c.lblName = nil
	c.lblValue = nil
	c.lblUOM = nil
	c.lblDT = nil
	c.btnClose = nil
	c.Panel.Dispose()
}

func (c *PanelFullScreenValue) SetItemId(itemId string) {
	c.itemId = itemId
	c.Clear()
}

func (c *PanelFullScreenValue) Clear() {
	c.lblName.SetText("waiting for data [" + c.itemId + "]")
	c.lblValue.SetText("")
	c.lblUOM.SetText("")
	c.lblDT.SetText("")
	c.UpdateLayout()
}

func (c *PanelFullScreenValue) timerUpdate() {
	c.client.GetItemsValues([]string{c.itemId}, func(items []common_interfaces.ItemGetUnitItems, err error) {
		for _, di := range items {
			if di.Name == c.itemId {
				value := di.Value.Value

				{
					if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
						p := message.NewPrinter(language.English)
						value = strings.ReplaceAll(p.Sprint(intValue), ",", " ")
					}
				}

				val := []rune(value)
				fontSize := 14.0

				if len(val) > 1024 {
					val = val[:1024]
					val = append(val, []rune("...")...)
				}
				c.lblName.SetText(c.itemId)
				c.lblValue.SetFontSize(fontSize)
				c.lblValue.SetText(string(val))
				c.lblUOM.SetText(di.Value.UOM)
				c.lblDT.SetText(time.Unix(0, di.Value.DT*1000).Format("2006-01-02 15-04-05"))

				kW := float64(c.lblValue.MinWidth()) / float64(c.Width())
				kH := float64(c.lblValue.MinHeight()) / float64(c.Height())
				targetFontSizeW := 0.0
				targetFontSizeH := 0.0
				if kW != 0 {
					targetFontSizeW = fontSize * (1 / kW)
					targetFontSizeW -= targetFontSizeW / 5
					if targetFontSizeW > 288 {
						targetFontSizeW = 288
					}
				}
				if kH != 0 {
					targetFontSizeH = fontSize * (1 / kH)
					targetFontSizeH -= targetFontSizeH / 5
					if targetFontSizeH > 288 {
						targetFontSizeH = 288
					}
				}

				targetFontSize := math.Min(targetFontSizeW, targetFontSizeH)

				c.lblValue.SetFontSize(targetFontSize)

				if di.Value.UOM == "error" {
					c.lblValue.SetForeColor(settings.BadColor)
					c.lblValue.SetBorderTop(3, settings.BadColor)
					c.lblValue.SetBorderBottom(3, settings.BadColor)
				} else {
					c.lblValue.SetForeColor(settings.GoodColor)
					c.lblValue.SetBorderTop(3, settings.GoodColor)
					c.lblValue.SetBorderBottom(3, settings.GoodColor)
					//c.lblValue.SetForeColor(c.ForeColor())
				}
			}
		}
	})
}
