package widget_time_filter

import (
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
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
	/*txtBlockLast := pButtons.AddTextBlockOnGrid(0, 0, "Last:")
	txtBlockLast.TextHAlign = canvas.HAlignRight

	c.AddButton(pButtons.AddButtonOnGrid(1, 0, "1 min", c.rbTimeChecked))
	c.AddButton(pButtons.AddButtonOnGrid(2, 0, "5 min", c.rbTimeChecked))
	c.AddButton(pButtons.AddButtonOnGrid(3, 0, "10 min", c.rbTimeChecked))
	c.AddButton(pButtons.AddButtonOnGrid(4, 0, "30 min", c.rbTimeChecked))
	c.AddButton(pButtons.AddButtonOnGrid(5, 0, "60 min", c.rbTimeChecked))*/

	c.pButtons.AddTextBlockOnGrid(0, 0, "Time filter:")
	c.AddButton(1, "1m", "last1min", "Last 1 minute")
	c.AddButton(2, "5m", "last5min", "Last 5 minutes")
	c.AddButton(3, "10m", "last10min", "Last 10 minutes")
	c.AddButton(4, "30m", "last30min", "Last 30 minutes")
	c.AddButton(5, "60m", "last60min", "Last Hour")
	c.AddButton(6, "current hour", "current_hour", "Current Hour")
	c.AddButton(7, "today", "current_day", "Today")
	c.AddButton(8, "previous hour", "previous_hour", "Previous Hour")
	c.AddButton(9, "yesterday", "previous_day", "Previous Day")
	c.AddButton(10, "custom", "custom", "Custom")

	c.cmbIntervals = c.pButtons.AddComboBoxOnGrid(10, 0)
	c.cmbIntervals.SetMinWidth(170)
	c.cmbIntervals.SetXExpandable(false)
	//c.cmbIntervals.SetMaxWidth(200)

	c.cmbIntervals.AddItem("Last 1 min", "last1min")
	c.cmbIntervals.AddItem("Last 5 min", "last5min")
	c.cmbIntervals.AddItem("Last 10 min", "last10min")
	c.cmbIntervals.AddItem("Last 30 min", "last30min")
	c.cmbIntervals.AddItem("Last 60 min", "last60min")

	c.cmbIntervals.AddItem("Current Hour", "current_hour")
	c.cmbIntervals.AddItem("Current Day", "current_day")

	c.cmbIntervals.AddItem("Previous Hour", "previous_hour")
	c.cmbIntervals.AddItem("Previous Day", "previous_day")
	c.cmbIntervals.AddItem("Custom", "custom")
	c.cmbIntervals.SetVisible(false)

	c.cmbIntervals.OnCurrentIndexChanged = func(event *uicontrols.ComboBoxEvent) {
		c.updateCustomDateTime()
	}

	/*txtBlockCurrent := pButtons.AddTextBlockOnGrid(0, 1, "Current:")
	txtBlockCurrent.TextHAlign = canvas.HAlignRight

	c.AddButton(pButtons.AddButtonOnGrid(1, 1, "Hour", c.rbTimeChecked))
	c.AddButton(pButtons.AddButtonOnGrid(2, 1, "Day", c.rbTimeChecked))
	c.AddButton(pButtons.AddButtonOnGrid(3, 1, "Week", c.rbTimeChecked))
	c.AddButton(pButtons.AddButtonOnGrid(4, 1, "Month", c.rbTimeChecked))
	c.AddButton(pButtons.AddButtonOnGrid(5, 1, "Year", c.rbTimeChecked))

	txtBlockPrev := pButtons.AddTextBlockOnGrid(0, 2, "Prev:")
	txtBlockPrev.TextHAlign = canvas.HAlignRight
	c.rbTimePrev1Hour = pButtons.AddButtonOnGrid(1, 2, "Hour", c.rbTimeChecked)
	c.AddButton(c.rbTimePrev1Hour)
	c.rbTimePrev1Day = pButtons.AddButtonOnGrid(2, 2, "Day", c.rbTimeChecked)
	c.AddButton(c.rbTimePrev1Day)
	c.rbTimePrev1Week = pButtons.AddButtonOnGrid(3, 2, "Week", c.rbTimeChecked)
	c.AddButton(c.rbTimePrev1Week)
	c.rbTimePrev1Month = pButtons.AddButtonOnGrid(4, 2, "Month", c.rbTimeChecked)
	c.AddButton(c.rbTimePrev1Month)
	c.rbTimePrev1Year = pButtons.AddButtonOnGrid(5, 2, "Year", c.rbTimeChecked)
	c.AddButton(c.rbTimePrev1Year)


	c.rbTimeCustom = pCustom.AddButtonOnGrid(0, 0, "Custom", c.rbTimeChecked)
	c.AddButton(c.rbTimeCustom)*/

	//pCustom := c.AddPanelOnGrid(2, 0)
	c.dtPickerFrom = c.AddDateTimePickerOnGrid(11, 0)
	dtFrom := time.Now().Add(-1 * time.Hour)
	dtFrom = dtFrom.Add(time.Duration(-dtFrom.Nanosecond()))
	c.dtPickerFrom.SetDateTime(dtFrom)
	//c.dtPickerFrom.DateTimeChanged = c.rbTimeChecked

	dtTo := time.Now()
	dtTo = dtTo.Add(time.Duration(-dtTo.Nanosecond()))
	c.dtPickerTo = c.AddDateTimePickerOnGrid(12, 0)
	c.dtPickerTo.SetDateTime(dtTo)
	//c.dtPickerTo.DateTimeChanged = c.rbTimeChecked

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
	b.SetMinWidth(40)
	b.SetUserData("key", key)
	b.SetTooltip(tooltip)
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
