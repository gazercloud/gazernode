package widget_chart

import (
	"encoding/json"
	"fmt"
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/history"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/timechart"
	"github.com/gazercloud/gazernode/widgets/widget_time_filter"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiproperties"
	"sort"
	"strconv"
	"time"
)

type WidgetCharts struct {
	uicontrols.Panel
	client           *client.Client
	timer            *uievents.FormTimer
	timeChartToolbar *timechart.ToolBar
	timeChart        *timechart.TimeChart
	items            []*DocumentChartItem
	timeFilter       *widget_time_filter.TimeFilterWidget
	isActive_        bool
}

func NewWidgetCharts(parent uiinterfaces.Widget, client *client.Client) *WidgetCharts {
	var c WidgetCharts
	c.client = client
	c.InitControl(parent, &c)
	c.SetPanelPadding(0)
	c.items = make([]*DocumentChartItem, 0)
	return &c
}

func (c *WidgetCharts) OnInit() {
	c.timer = c.Window().NewTimer(500, c.timerUpdate)
	c.timer.StartTimer()

	c.timeFilter = widget_time_filter.NewTimeFilterWidget(c)
	c.timeFilter.OnEdited = c.timeFilterChanged
	c.timeFilter.SetGridX(0)
	c.timeFilter.SetGridY(0)
	c.AddWidgetOnGrid(c.timeFilter, 0, 0)

	c.timeChart = timechart.NewTimeChart(c)
	c.AddWidgetOnGrid(c.timeChart, 0, 2)

	c.timeChart.OnMouseDropOnArea = func(droppedValue interface{}, area *timechart.Area) {
		c.AddSeries(droppedValue.(string), area)

	}

	c.timeChartToolbar = timechart.NewToolBar(c, c.timeChart)
	c.AddWidgetOnGrid(c.timeChartToolbar, 0, 1)
	c.timeChartToolbar.SetVisible(false)
}

func (c *WidgetCharts) SetOnChartContextMenuNeed(OnChartContextMenuNeed func(timeChart *timechart.TimeChart, area *timechart.Area, areaIndex int) uiinterfaces.Menu) {
	c.timeChart.OnChartContextMenuNeed = OnChartContextMenuNeed
}

func (c *WidgetCharts) IsActive() bool {
	return c.isActive_
}

func (c *WidgetCharts) SetIsActive(isActive bool) {
	c.isActive_ = isActive
}

func (c *WidgetCharts) Dispose() {
	if c.timer != nil {
		c.timer.StopTimer()
	}
	c.timer = nil
	c.timeChart = nil

	c.client = nil
	c.items = nil
	c.timeFilter = nil
	c.Panel.Dispose()
}

func (c *WidgetCharts) SetEdit(editing bool) {
	c.timeChart.SetEditing(editing)
}

func (c *WidgetCharts) timeFilterChanged() {
	for _, item := range c.items {
		item.serMain.Clear()
		item.serMain.SetName(item.name)
		item.needToReload = true
	}

	c.timeChart.SetHorizRange(c.timeFilter.TimeFrom(), c.timeFilter.TimeTo())
	c.timeChart.ZoomShowEntire()
	c.timeChart.Update("DocumentChart")
}

// Main update cycle
func (c *WidgetCharts) timerUpdateValuesHandler() {
	/*for _, item := range c.items {
		item.checkItemIDS()
	}*/
}

type ChartSettingsSeries struct {
	Item  string `json:"item"`
	Color string `json:"color"`
}

type ChartSettingsArea struct {
	UnitedScale   bool                   `json:"united_scale"`
	ShowQualities bool                   `json:"show_qualities"`
	Series        []*ChartSettingsSeries `json:"series"`
}

type ChartSettings struct {
	Areas []*ChartSettingsArea `json:"areas"`
}

func (c *WidgetCharts) Save() []byte {
	var settings ChartSettings
	settings.Areas = make([]*ChartSettingsArea, 0)
	for _, a := range c.timeChart.Areas() {
		var area ChartSettingsArea
		area.Series = make([]*ChartSettingsSeries, 0)
		area.UnitedScale = a.UnitedScale()
		area.ShowQualities = a.ShowQualities()
		for _, s := range a.Series() {
			var ser ChartSettingsSeries
			ser.Item = s.Id()
			r, g, b, aa := s.Color().RGBA()
			ser.Color = fmt.Sprintf("#%02X%02X%02X%02X", r/256, g/256, b/256, aa/256)
			area.Series = append(area.Series, &ser)
		}
		settings.Areas = append(settings.Areas, &area)
	}
	bs, _ := json.MarshalIndent(settings, "", " ")
	return bs
}

func (c *WidgetCharts) Load(data []byte) {
	c.timeChart.RemoveAllAreas()
	var settings ChartSettings
	err := json.Unmarshal(data, &settings)
	if err != nil {
		return
	}
	var currentArea *timechart.Area
	for _, a := range settings.Areas {
		currentArea = c.timeChart.AddArea()
		currentArea.SetUnitedScale(a.UnitedScale)
		currentArea.SetShowQualities(a.ShowQualities)
		for _, s := range a.Series {
			ser := c.AddSeries(s.Item, currentArea)
			col := uiproperties.ParseHexColor(s.Color)
			if len(s.Color) > 0 {
				ser.SetColor(col)
			}
		}
	}
}

func (c *WidgetCharts) SetDataItems(items []string) {
	c.timeChart.RemoveAllAreas()
	area := c.timeChart.AddArea()

	for _, i := range items {
		item := NewDocumentChartItem(c.OwnWindow, c.timeChart, area, i, c.client)
		c.items = append(c.items, item)
	}
}

func (c *WidgetCharts) SetShowQualities(showQualities bool) {
	if c.timeChart != nil {
		c.timeChart.SetShowQualities(showQualities)
	}
}

func (c *WidgetCharts) AddSeries(name string, area *timechart.Area) *timechart.Series {
	//c.timeChart.RemoveAllAreas()
	if area == nil {
		area = c.timeChart.AddArea()
		area.SetShowQualities(true)
	}

	item := NewDocumentChartItem(c.OwnWindow, c.timeChart, area, name, c.client)
	c.items = append(c.items, item)
	return item.serMain
}

func (c *WidgetCharts) timerUpdate() {
	c.timeChart.SetDefaultDisplayRange(c.timeFilter.TimeFrom(), c.timeFilter.TimeTo())
	for _, item := range c.items {
		item.Clean()
	}
}

type DocumentChartItem struct {
	client     *client.Client
	area       *timechart.Area
	serMain    *timechart.Series
	loaded     bool
	dataItemId int64
	treeItemId int64
	name       string

	timeFrom int64
	timeTo   int64

	groupTimeRange int64
	needToReload   bool
	allowUpdate    bool

	idsLoaded     bool
	idsNeedToLoad bool
	lastGetDT     time.Time

	values map[int64]*DocumentChartValues

	Window uiinterfaces.Window
}

func NewDocumentChartItem(window uiinterfaces.Window, chart *timechart.TimeChart, area *timechart.Area, name string, client *client.Client) *DocumentChartItem {
	var item DocumentChartItem
	item.client = client
	item.Window = window
	item.area = area
	item.serMain = chart.AddSeries(area, name)
	item.serMain.SetDataProvider(&item)
	item.needToReload = true
	item.idsNeedToLoad = true
	item.idsLoaded = false
	item.name = name
	item.values = make(map[int64]*DocumentChartValues)
	return &item
}

func (c *DocumentChartItem) Dispose() {

	for _, v := range c.values {
		v.Dispose()
	}

	if c.serMain != nil {
		c.serMain.SetDataProvider(nil)
	}
	c.area = nil
	c.serMain = nil
	c.values = nil
	c.Window = nil
}

func (c *DocumentChartItem) Clean() {
	t := time.Now().UTC()
	found := true
	for found {
		found = false
		for key, v := range c.values {
			if t.Sub(v.lastGetDT).Milliseconds() > 5000 {
				delete(c.values, key)
				//logger.Println("DocumentChartItem clear ", key)
				found = true
				break
			}
		}
	}
}

func AlignGroupTimeRange(groupTimeRange int64) int64 {
	if groupTimeRange < 60*1000000 {
		groupTimeRange = 1
	}

	if groupTimeRange >= 60*1000000 && groupTimeRange < 60*60*1000000 {
		groupTimeRange = 60 * 1000000 // By minute
	}

	if groupTimeRange >= 60*60*1000000 && groupTimeRange < 24*60*60*1000000 {
		groupTimeRange = 60 * 60 * 1000000 // By Hour
	}

	if groupTimeRange >= 24*60*60*1000000 {
		groupTimeRange = 24 * 60 * 60 * 1000000 // By day
	}
	return groupTimeRange
}

func (c *DocumentChartItem) GetData(key string, minTime, maxTime int64, groupTimeRange int64) []*timechart.Value {

	groupTimeRange = AlignGroupTimeRange(groupTimeRange)

	//groupTimeRange = 1000000

	if values, ok := c.values[groupTimeRange]; ok {
		return values.GetData(minTime, maxTime)
	}

	values := &DocumentChartValues{}
	values.name = c.name
	values.client = c.client
	values.groupTimeRange = groupTimeRange
	values.values = make([]*timechart.Value, 0)
	values.dataItemId = c.dataItemId
	values.level = 0
	values.loadedRanges = make([]*TimeRange, 0)
	c.values[groupTimeRange] = values
	return values.GetData(minTime, maxTime)
}

type TimeRange struct {
	timeFrom int64
	timeTo   int64
}

type LoadingTask struct {
	started  bool
	timeFrom int64
	timeTo   int64
}

type DocumentChartValues struct {
	name           string
	client         *client.Client
	groupTimeRange int64
	dataItemId     int64
	level          int64
	values         []*timechart.Value

	timeFrom  int64
	timeTo    int64
	lastGetDT time.Time

	loadedRanges  []*TimeRange
	loadingRanges []*LoadingTask
}

func (c *DocumentChartValues) Dispose() {
	c.values = nil
	c.loadedRanges = nil
	c.loadingRanges = nil
}

func (c *DocumentChartValues) GetData(minTime, maxTime int64) []*timechart.Value {
	c.checkValues(minTime, maxTime)
	c.lastGetDT = time.Now().UTC()

	result := make([]*timechart.Value, 0, 4000)

	// local Filter by time
	for _, v := range c.values {
		if v.DatetimeFirst >= minTime && v.DatetimeLast <= maxTime {
			result = append(result, v)
		}
	}

	return result
}

func (c *DocumentChartValues) requestHistory(task *LoadingTask) {
	//logger.Println("DocumentChartValues requestHistory sec:", (task.timeTo-task.timeFrom)/1000000, " q size: ", len(c.loadingRanges))

	c.client.ReadHistory(c.name, task.timeFrom, task.timeTo, func(result *history.ReadResult, err error) {
		logger.Println("client ReadHistory ", c.groupTimeRange)
		if err == nil {
			resultItems := make([]*timechart.Value, 0)
			if result != nil {
				{

					rawValues := make([]*common_interfaces.ItemValue, 0)
					rawValuesGroupIndex := make([]int64, 0)

					for _, item := range result.Items {
						rawValues = append(rawValues, item)
						groupIndex := (item.DT - result.DTBegin) / c.groupTimeRange
						rawValuesGroupIndex = append(rawValuesGroupIndex, groupIndex)
					}

					lastGroupIndex := int64(-1)
					var currentValueRange *timechart.Value

					for index := range rawValuesGroupIndex {
						r := rawValues[index]
						valudValue := false

						if lastGroupIndex != rawValuesGroupIndex[index] {
							if currentValueRange != nil {
								resultItems = append(resultItems, currentValueRange)
								currentValueRange = nil
							}
							lastGroupIndex = rawValuesGroupIndex[index]
						}

						if currentValueRange == nil {
							currentValueRange = &timechart.Value{}
							currentValueRange.DatetimeFirst = r.DT - (r.DT % c.groupTimeRange)
							currentValueRange.DatetimeLast = r.DT - (r.DT % c.groupTimeRange) + c.groupTimeRange - 1
							currentValueRange.Qualities = make([]int64, 0)
							currentValueRange.MinValue = 1000000000000
							currentValueRange.MaxValue = -1000000000000
							currentValueRange.AvgValue = 0
							currentValueRange.FirstValue = 0
							currentValueRange.LastValue = 0
						}

						if r.UOM != "error" {
							valueAsFloat, err := strconv.ParseFloat(r.Value, 64)
							if err == nil {
								valudValue = true

								if valueAsFloat < currentValueRange.MinValue {
									currentValueRange.MinValue = valueAsFloat
								}
								if valueAsFloat > currentValueRange.MaxValue {
									currentValueRange.MaxValue = valueAsFloat
								}
								currentValueRange.AvgValue += valueAsFloat
								if currentValueRange.CountOfValues > 0 {
									currentValueRange.AvgValue /= 2
								}

								if currentValueRange.CountOfValues == 0 {
									currentValueRange.FirstValue = valueAsFloat
								}

								currentValueRange.LastValue = valueAsFloat

								currentValueRange.CountOfValues++
							}
						}

						if r.UOM != "error" && valudValue {
							foundGood := false
							for _, q := range currentValueRange.Qualities {
								if q == 192 {
									foundGood = true
								}
							}
							if !foundGood {
								currentValueRange.Qualities = append(currentValueRange.Qualities, 192)
							}
						} else {
							foundBad := false
							for _, q := range currentValueRange.Qualities {
								if q == 0 {
									foundBad = true
								}
							}
							if !foundBad {
								currentValueRange.Qualities = append(currentValueRange.Qualities, 0)
							}
						}

					}

				}

				// Apply incoming data
				c.insertValues(resultItems, task.timeFrom, task.timeTo)
			}
		} else {
			logger.Println("DocumentChart timerUpdateValuesHandler error: " + err.Error())
		}

		// Remove loading task
		for index, rng := range c.loadingRanges {
			if task.timeFrom == rng.timeFrom && task.timeTo == rng.timeTo {
				c.loadingRanges = append(c.loadingRanges[:index], c.loadingRanges[index+1:]...)
				break
			}
		}
	})
	/*c.client.GetDataItemHistoryRanges(c.dataItemId, 0, task.timeFrom, task.timeTo, c.groupTimeRange, func(dataItemHistory *datastorage.DataItemHistoryRanges, err error) {
	})*/
}

func (c *DocumentChartValues) checkValues(timeFrom, timeTo int64) {

	// Full requested range already loaded
	for _, rng := range c.loadedRanges {
		if timeFrom >= rng.timeFrom && timeTo <= rng.timeTo {
			return
		}
	}

	// Part of data already loaded
	for _, rng := range c.loadedRanges {
		if timeFrom >= rng.timeFrom && timeFrom <= rng.timeTo {
			timeFrom = rng.timeTo + 1
			break
		}
	}

	// Already loading
	for _, rng := range c.loadingRanges {
		if timeFrom == rng.timeFrom && timeTo == rng.timeTo {
			return
		}
	}

	// Make task for loading
	/*if timeTo - timeFrom > 3600000000 {
		for t := timeFrom; t < timeTo; t += 3600000000 {
			tBegin := t
			tEnd := t + 3600000000
			if tEnd > timeTo {
				tEnd = timeTo
			}
			if tEnd <= tBegin {
				continue
			}
			var task LoadingTask
			task.timeFrom = tBegin
			task.timeTo = tEnd
			c.loadingRanges = append(c.loadingRanges, &task)
			c.requestHistory(&task)
		}
	} else {*/
	if len(c.loadingRanges) < 2 {
		var task LoadingTask
		task.timeFrom = timeFrom
		task.timeTo = timeTo
		c.loadingRanges = append(c.loadingRanges, &task)
		c.requestHistory(&task)
	}
	//}
}

func (c *DocumentChartValues) insertValues(readResult []*timechart.Value, timeFrom, timeTo int64) {
	indexOfBeginForDelete := -1
	indexOfBeginForDeleteFound := false
	indexOfEndForDelete := -1
	indexOfEndForDeleteFound := false
	for index, v := range c.values {
		if v.DatetimeFirst > timeFrom {
			indexOfBeginForDelete = index
			indexOfBeginForDeleteFound = true
			break
		}
	}

	for index := len(c.values) - 1; index >= 0; index-- {
		v := c.values[index]
		if v.DatetimeLast < timeTo {
			indexOfEndForDelete = index
			indexOfEndForDeleteFound = true
			break
		}
	}

	if indexOfEndForDelete > indexOfBeginForDelete && indexOfBeginForDeleteFound && indexOfEndForDeleteFound {
		c.values = append(c.values[:indexOfBeginForDelete], c.values[indexOfEndForDelete:]...)
	}

	for _, dataItemValue := range readResult {
		//val := strings.ReplaceAll(dataItemValue.Value, "\n", "")
		//valueAsFloat, err := strconv.ParseFloat(val, 64)
		/*if dataItemValue.Value == "" || dataItemValue.UOM == "error" || dataItemValue.UOM == "stopped" {
			err = errors.New("W")
		}*/

		/*var v timechart.Value
		v.DatetimeFirst = dataItemValue.DT
		v.DatetimeLast = dataItemValue.DT
		v.FirstValue = valueAsFloat
		v.LastValue = valueAsFloat
		v.MinValue = valueAsFloat
		v.MaxValue = valueAsFloat
		v.AvgValue = valueAsFloat*/
		/*if err == nil {
			dataItemValue.Qualities = []int64{192}
		} else {
			dataItemValue.Qualities = []int64{0}
		}*/
		c.values = append(c.values, dataItemValue)
	}

	sort.Slice(c.values, func(i, j int) bool { return c.values[i].DatetimeFirst < c.values[j].DatetimeFirst })

	loadedRange := &TimeRange{}
	loadedRange.timeFrom = timeFrom
	loadedRange.timeTo = timeTo
	c.loadedRanges = append(c.loadedRanges, loadedRange)

	sort.Slice(c.loadedRanges, func(i, j int) bool { return c.loadedRanges[i].timeFrom < c.loadedRanges[j].timeFrom })

	// Crossing time ranges
	for {
		foundCross := false
		lastTimeTo := int64(0)
		crossIndexOfSecond := 0
		for index, rng := range c.loadedRanges {
			if index == 0 {
				lastTimeTo = rng.timeTo
				continue
			}
			if rng.timeFrom <= (lastTimeTo + 1) {
				crossIndexOfSecond = index
				foundCross = true
				break
			}
			lastTimeTo = rng.timeTo
		}

		if foundCross {
			c.loadedRanges[crossIndexOfSecond-1].timeTo = c.loadedRanges[crossIndexOfSecond].timeTo
			c.loadedRanges = append(c.loadedRanges[:crossIndexOfSecond], c.loadedRanges[crossIndexOfSecond+1:]...)
		}

		if !foundCross {
			break
		}
	}

	// Last time range must be less than last timestamp of values
	if len(c.values) > 0 {
		lastDateTime := c.values[len(c.values)-1].DatetimeFirst
		for index, rng := range c.loadedRanges {
			if lastDateTime >= rng.timeFrom && lastDateTime <= rng.timeTo {
				rng.timeTo = lastDateTime
				if rng.timeTo <= rng.timeFrom {
					c.loadedRanges = append(c.loadedRanges[:index], c.loadedRanges[index+1:]...)
				}
				break
			}
		}
	}
}
