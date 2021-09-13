package widget_time_filter

import (
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
	"github.com/gazercloud/gazerui/uistyles"
	"strings"
	"time"
)

type TimeFilterWidget struct {
	uicontrols.Panel

	lblLive *uicontrols.TextBlock

	cmbIntervals *uicontrols.ComboBox

	/*rbTimeLast1Min  *uicontrols.Button
	rbTimeLast5Min  *uicontrols.Button
	rbTimeLast10Min *uicontrols.Button
	rbTimeLast30Min *uicontrols.Button
	rbTimeLast60Min *uicontrols.Button

	rbTimeCurrent1Hour  *uicontrols.Button
	rbTimeCurrent1Day   *uicontrols.Button
	rbTimeCurrent1Week  *uicontrols.Button
	rbTimeCurrent1Month *uicontrols.Button
	rbTimeCurrent1Year  *uicontrols.Button

	rbTimePrev1Hour  *uicontrols.Button
	rbTimePrev1Day   *uicontrols.Button
	rbTimePrev1Week  *uicontrols.Button
	rbTimePrev1Month *uicontrols.Button
	rbTimePrev1Year  *uicontrols.Button

	rbTimeCustom *uicontrols.Button*/

	pButtons *uicontrols.Panel
	buttons  []*uicontrols.Button

	lblCustom *uicontrols.TextBlock

	dtPickerFrom *uicontrols.DateTimePicker
	dtPickerTo   *uicontrols.DateTimePicker

	OnEdited func()
}

func NewTimeFilterWidget(parent uiinterfaces.Widget) *TimeFilterWidget {
	var c TimeFilterWidget
	c.InitControl(parent, &c)
	c.SetPanelPadding(0)

	c.pButtons = c.AddPanelOnGrid(0, 0)
	c.pButtons.SetPanelPadding(0)

	//c.pButtons.AddTextBlockOnGrid(0, 0, "Time filter:")
	c.AddButton(1, "5m", "last5min", "Last 5 minutes")
	c.AddButton(2, "60m", "last60min", "Last Hour")
	c.AddButton(3, "24h", "last24hours", "Last 24 hours")
	c.AddButton(4, "7d", "last7days", "Last Week")
	c.AddButton(5, "C", "custom", "Custom")

	c.cmbIntervals = c.pButtons.AddComboBoxOnGrid(15, 0)
	c.cmbIntervals.SetMinWidth(170)
	c.cmbIntervals.SetXExpandable(false)
	//c.cmbIntervals.SetMaxWidth(200)

	c.cmbIntervals.AddItem("Last 1 min", "last1min")
	c.cmbIntervals.AddItem("Last 5 min", "last5min")
	c.cmbIntervals.AddItem("Last 10 min", "last10min")
	c.cmbIntervals.AddItem("Last 30 min", "last30min")
	c.cmbIntervals.AddItem("Last 60 min", "last60min")
	c.cmbIntervals.AddItem("Last 24 hours", "last24hours")
	c.cmbIntervals.AddItem("Last 7 days", "last7days")

	c.cmbIntervals.AddItem("Current Hour", "current_hour")
	c.cmbIntervals.AddItem("Current Day", "current_day")
	c.cmbIntervals.AddItem("Current Week", "current_week")
	c.cmbIntervals.AddItem("Current Month", "current_month")
	c.cmbIntervals.AddItem("Current Year", "current_year")

	c.cmbIntervals.AddItem("Previous Hour", "previous_hour")
	c.cmbIntervals.AddItem("Previous Day", "previous_day")
	c.cmbIntervals.AddItem("Custom", "custom")
	//c.cmbIntervals.SetVisible(false)

	c.cmbIntervals.OnCurrentIndexChanged = func(event *uicontrols.ComboBoxEvent) {
		c.updateCustomDateTime()
	}

	c.dtPickerFrom = c.AddDateTimePickerOnGrid(11, 0)
	dtFrom := time.Now().Add(-1 * time.Hour)
	dtFrom = dtFrom.Add(time.Duration(-dtFrom.Nanosecond()))
	c.dtPickerFrom.SetDateTime(dtFrom)

	dtTo := time.Now()
	dtTo = dtTo.Add(time.Duration(-dtTo.Nanosecond()))
	c.dtPickerTo = c.AddDateTimePickerOnGrid(12, 0)
	c.dtPickerTo.SetDateTime(dtTo)

	c.AddHSpacerOnGrid(13, 0)

	c.cmbIntervals.SetCurrentItemIndex(1)
	c.updateButtonsColors()

	c.updateEnables()

	return &c
}

func (c *TimeFilterWidget) AddButton(pos int, text string, key string, tooltip string) {
	b := c.pButtons.AddButtonOnGrid(pos, 0, text, func(event *uievents.Event) {
		k := event.Sender.(*uicontrols.Button).UserData("key").(string)
		c.selectIntervalComboItem(k)
	})
	b.SetMinWidth(50)
	b.SetUserData("key", key)
	b.SetTooltip(tooltip)
	//b.SetImageSize(32, 24)
	b.SetMouseCursor(ui.MouseCursorPointer)
	c.buttons = append(c.buttons, b)
}

func (c *TimeFilterWidget) Dispose() {
	c.lblLive = nil

	c.lblCustom = nil

	c.dtPickerFrom = nil
	c.dtPickerTo = nil

	c.buttons = nil

	c.cmbIntervals = nil

	c.Panel.Dispose()
}

func (c *TimeFilterWidget) ControlType() string {
	return "TimeFilterWidget"
}

func (c *TimeFilterWidget) UpdateStyle() {
	c.Panel.UpdateStyle()
	c.updateButtonsColors()
}

func (c *TimeFilterWidget) imageByKey(key string) []byte {
	switch key {
	case "last1min":
		return uiresources.R_icons_custom_time_filter_icon_1min_png
	case "last5min":
		return uiresources.R_icons_custom_time_filter_icon_5min_png
	case "last10min":
		return uiresources.R_icons_custom_time_filter_icon_10min_png
	case "last30min":
		return uiresources.R_icons_custom_time_filter_icon_30min_png
	case "last60min":
		return uiresources.R_icons_custom_time_filter_icon_60min_png
	case "current_hour":
		return uiresources.R_icons_custom_time_filter_icon_cH_png
	case "current_day":
		return uiresources.R_icons_custom_time_filter_icon_cD_png
	case "previous_hour":
		return uiresources.R_icons_custom_time_filter_icon_pH_png
	case "previous_day":
		return uiresources.R_icons_custom_time_filter_icon_pD_png
	case "custom":
		return uiresources.R_icons_custom_time_filter_icon_C_png
	}
	return nil
}

func (c *TimeFilterWidget) updateButtonsColors() {
	key := c.cmbIntervals.Items[c.cmbIntervals.CurrentItemIndex].UserData("key").(string)
	for _, b := range c.buttons {
		if b.UserData("key").(string) == key {
			b.SetForeColor(uistyles.DefaultBackColor)
			b.SetBackColor(c.ForeColor())
		} else {
			b.SetForeColor(nil)
			b.SetBackColor(nil)
		}

		//b.SetImageSize(32, 24)
		//b.SetImage(uiresources.ResImgCol(c.imageByKey(b.UserData("key").(string)), b.ForeColor()))
	}
}

func (c *TimeFilterWidget) selectIntervalComboItem(key string) {
	for i, item := range c.cmbIntervals.Items {
		if item.UserData("key").(string) == key {
			c.cmbIntervals.SetCurrentItemIndex(i)
			break
		}
	}
}

func (c *TimeFilterWidget) updateCustomDateTime() {

	now := time.Now()
	from := now
	to := now
	customTime := false

	key := c.cmbIntervals.Items[c.cmbIntervals.CurrentItemIndex].UserData("key").(string)

	if key == "previous_year" {
		now = now.AddDate(-1, 0, 0)
		from = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()).UTC()
		to = from.AddDate(1, 0, 0).Add(-1 * time.Nanosecond)
		customTime = true
	}
	if key == "previous_month" {
		now = now.AddDate(0, -1, 0)
		from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).UTC()
		to = from.AddDate(0, 1, 0).Add(-1 * time.Nanosecond)
		customTime = true
	}
	if key == "previous_week" {
		now = now.AddDate(0, 0, -int(now.Weekday())-7)
		from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).UTC()
		to = time.Date(now.Year(), now.Month(), now.Day()+6, 23, 59, 59, 999999999, now.Location()).UTC()
		customTime = true
	}
	if key == "previous_day" {
		now = now.AddDate(0, 0, -1)
		from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).UTC()
		to = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location()).UTC()
		customTime = true
	}
	if key == "previous_hour" {
		now = now.Add(-1 * time.Hour)
		from = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location()).UTC()
		to = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 59, 59, 999999999, now.Location()).UTC()
		customTime = true
	}

	if key == "current_year" {
		from = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()).UTC()
		to = from.AddDate(1, 0, 0).Add(-1 * time.Nanosecond)
		customTime = true
	}
	if key == "current_month" {
		from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).UTC()
		to = from.AddDate(0, 1, 0).Add(-1 * time.Nanosecond)
		customTime = true
	}
	if key == "current_week" {
		now = now.AddDate(0, 0, -int(now.Weekday()))
		from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).UTC()
		to = time.Date(now.Year(), now.Month(), now.Day()+6, 23, 59, 59, 999999999, now.Location()).UTC()
		customTime = true
	}
	if key == "current_day" {
		from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).UTC()
		to = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location()).UTC()
		customTime = true
	}

	if key == "current_hour" {
		from = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location()).UTC()
		to = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 59, 59, 999999999, now.Location()).UTC()
		customTime = true
	}

	if customTime {
		c.dtPickerFrom.SetDateTime(from.Local())
		c.dtPickerTo.SetDateTime(to.Local())
	}

	if c.OnEdited != nil {
		c.OnEdited()
	}

	c.updateEnables()
	c.updateButtonsColors()
}

func (c *TimeFilterWidget) updateEnables() {
	key := c.cmbIntervals.Items[c.cmbIntervals.CurrentItemIndex].UserData("key").(string)
	if key == "custom" {
		c.dtPickerFrom.SetVisible(true)
		c.dtPickerTo.SetVisible(true)
	} else {
		c.dtPickerFrom.SetVisible(false)
		c.dtPickerTo.SetVisible(false)
	}
}

func (c *TimeFilterWidget) TimeFrom() int64 {
	var timeFrom int64
	stepTime := int64(10000000)
	key := c.cmbIntervals.Items[c.cmbIntervals.CurrentItemIndex].UserData("key").(string)

	if key == "last7days" {
		timeFrom = (time.Now().UTC().UnixNano() - 7*24*3600*1000000000) / 1000
		stepTime = 7 * 5760 * 1000000
	}

	if key == "last24hours" {
		timeFrom = (time.Now().UTC().UnixNano() - 24*3600*1000000000) / 1000
		stepTime = 5760 * 1000000
	}

	if key == "last60min" {
		timeFrom = (time.Now().UTC().UnixNano() - 3600*1000000000) / 1000
		stepTime = 240 * 1000000
	}
	if key == "last30min" {
		timeFrom = (time.Now().UTC().UnixNano() - 1800*1000000000) / 1000
		stepTime = 120 * 1000000
	}
	if key == "last10min" {
		timeFrom = (time.Now().UTC().UnixNano() - 600*1000000000) / 1000
		stepTime = 60 * 1000000
	}
	if key == "last5min" {
		timeFrom = (time.Now().UTC().UnixNano() - 300*1000000000) / 1000
		stepTime = 30 * 1000000
	}
	if key == "last1min" {
		timeFrom = (time.Now().UTC().UnixNano() - 60*1000000000) / 1000
		stepTime = 10 * 1000000
	}

	if !strings.HasPrefix(key, "last") {
		// Static Diapason
		timeFrom = c.dtPickerFrom.DateTime().UTC().UnixNano() / 1000
	} else {
		// Dynamic diapason
		deltaTime := timeFrom % stepTime
		timeFrom += stepTime - deltaTime
	}
	return timeFrom
}

func (c *TimeFilterWidget) TimeTo() int64 {
	var timeTo int64
	stepTime := int64(10000000)
	key := c.cmbIntervals.Items[c.cmbIntervals.CurrentItemIndex].UserData("key").(string)
	if strings.HasPrefix(key, "last") {

		if key == "last7days" {
			stepTime = 7 * 24 * 240 * 1000000
		}
		if key == "last24hours" {
			stepTime = 24 * 240 * 1000000
		}
		if key == "last60min" {
			stepTime = 240 * 1000000
		}
		if key == "last30min" {
			stepTime = 120 * 1000000
		}
		if key == "last10min" {
			stepTime = 60 * 1000000
		}
		if key == "last5min" {
			stepTime = 30 * 1000000
		}
		if key == "last1min" {
			stepTime = 10 * 1000000
		}

		timeTo = time.Now().UTC().UnixNano() / 1000

		deltaTime := timeTo % stepTime
		timeTo += stepTime - deltaTime
	} else {
		timeTo = c.dtPickerTo.DateTime().UTC().UnixNano() / 1000
	}
	return timeTo
}
