package timechart

type LoadingDiapason struct {
	MinTime, MaxTime int64
}

type LoadedDiapason struct {
	MinTime, MaxTime int64
	TimeRange        int64
}

type IDataProvider interface {
	GetData(key string, minTime, maxTime int64, groupTimeRange int64) ([]*Value, string)
	GetLoadingDiapasons() []LoadingDiapason
	GetLoadedDiapasons() []LoadedDiapason
}
