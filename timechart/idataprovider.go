package timechart

type IDataProvider interface {
	GetData(key string, minTime, maxTime int64, groupTimeRange int64) []*Value
}
