package timechart

type DataProviderTimeLevel struct {
	key            string
	groupTimeRange int64
	values         []*Value
}

type DataProviderItem struct {
	timeLevels map[int64]*DataProviderTimeLevel
}

type DataProvider struct {
	items map[string]*DataProviderItem
}

func (c *DataProvider) Init() {
	c.items = make(map[string]*DataProviderItem)
}

func (c *DataProvider) GetData(key string, minTime, maxTime int64, groupTimeRange int64) []*Value {

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

	if item, ok := c.items[key]; ok {
		return item.GetData(key, minTime, maxTime, groupTimeRange)
	}
	item := &DataProviderItem{}
	item.Init(key)
	c.items[key] = item
	return item.GetData(key, minTime, maxTime, groupTimeRange)
}

func (c *DataProviderItem) Init(key string) {
	c.timeLevels = make(map[int64]*DataProviderTimeLevel)
}

func (c *DataProviderItem) GetData(key string, minTime, maxTime int64, groupTimeRange int64) []*Value {
	if timeLevel, ok := c.timeLevels[groupTimeRange]; ok {
		return timeLevel.GetData(minTime, maxTime)
	}

	timeLevel := &DataProviderTimeLevel{}
	timeLevel.Init(key, groupTimeRange)
	c.timeLevels[groupTimeRange] = timeLevel
	return timeLevel.GetData(minTime, maxTime)
}

func (c *DataProviderTimeLevel) Init(key string, groupTimeRange int64) {
	c.key = key
	c.groupTimeRange = groupTimeRange
}

func (c *DataProviderTimeLevel) GetData(minTime, maxTime int64) []*Value {
	result := make([]*Value, 0)
	for _, v := range c.values {
		if v.DatetimeFirst == minTime {
			result = append(result, v)
		}
	}
	return result
}
