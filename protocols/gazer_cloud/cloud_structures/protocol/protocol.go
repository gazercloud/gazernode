package protocol

// get
type RequestGet struct {
	ChannelId string `json:"ch"`
}

// issue
type RequestIssue struct {
	Host      string `json:"h"`
	ChannelId string `json:"ch"`
	Password  string `json:"p"`
}

// where

type RequestWhere struct {
	Channels []string `json:"chs"`
}

type ResponseWhereItem struct {
	ChannelId string `json:"ch"`
	Host      string `json:"h"`
}

type ResponseWhere struct {
	Items []ResponseWhereItem `json:"items"`
}

// get_workers
type ResponseGetWorkersItem struct {
	Host   string `json:"host"`
	Scores int    `json:"scores"`
}

type ResponseGetWorkers struct {
	Workers []ResponseGetWorkersItem `json:"workers"`
}

// reg
type ResponseReg struct {
	ChannelId string `json:"ch"`
	Password  string `json:"p"`
}

// clear
type RequestClear struct {
	ChannelId string `json:"ch"`
	Password  string `json:"p"`
}

// channels
type ResponseChannels struct {
	Channels []string `json:"chs"`
}

// workerStatistics

type ResponseWorkerStatistics struct {
	Channels    int64 `json:"chs"`
	Connections int64 `json:"cs"`
	Memory      int64 `json:"m"`
}

type ResponseCountOfSubscribers struct {
	Channel string `json:"ch"`
	Count   int    `json:"c"`
}
