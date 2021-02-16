package cloud

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/protocols/gazer_cloud/bin_client"
	"github.com/gazercloud/gazernode/protocols/gazer_cloud/cloud_structures/protocol"
	"github.com/gazercloud/gazernode/utilities"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type Channel struct {
	client       *http.Client
	iDataStorage common_interfaces.IDataStorage

	channelId string
	password  string
	name      string

	items         []string
	itemsToRemove []string

	started  bool
	stopping bool

	worker string

	binClient        *bin_client.BinClient
	chProcessingData chan bin_client.BinFrameTask

	wrongWorkers map[string]time.Time

	subscribers int

	mtx sync.Mutex
}

func NewChannel(iDataStorage common_interfaces.IDataStorage, channelId string, password string, name string) *Channel {
	var c Channel
	c.channelId = channelId
	c.iDataStorage = iDataStorage
	c.password = password
	c.name = name
	c.items = make([]string, 0)
	c.itemsToRemove = make([]string, 0)
	c.chProcessingData = make(chan bin_client.BinFrameTask)
	c.wrongWorkers = make(map[string]time.Time)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{},
	}

	c.client = &http.Client{Transport: tr}
	if len(channelId) < 1 {
		c.GazerReg()
	}

	c.started = false
	c.stopping = false

	c.Start()

	return &c
}

func (c *Channel) Start() {
	if c.started {
		return
	}
	c.started = true
	c.stopping = false
	go c.thWorker()
	go c.thIncomingTraffic()
}

func (c *Channel) Stop() {
	if !c.started {
		return
	}

	c.stopping = true
	for i := 0; i < 10; i++ {
		if !c.started {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (c *Channel) Stat() *utilities.Statistics {
	var result *utilities.Statistics
	c.mtx.Lock()
	if c.binClient != nil {
		result = c.binClient.Stat()
	}
	c.mtx.Unlock()
	return result
}

func (c *Channel) GazerReg() {
	for i := 0; i < RoutersCount(); i++ {
		router := CurrentRouter()

		link := "https://" + router + "/api/request?fn=reg"

		response, err := c.client.Get(link)
		type Response struct {
			Channel  string `json:"ch"`
			Password string `json:"p"`
			Error    string `json:"e"`
		}

		var resp Response

		if err != nil {
			logger.Println(err)
		} else {
			content, _ := ioutil.ReadAll(response.Body)
			response.Body.Close()
			s := strings.TrimSpace(string(content))
			if json.Unmarshal([]byte(s), &resp) == nil {
				c.channelId = resp.Channel
				c.password = resp.Password
				break
			}
		}
		SetNextRouter()
	}

	c.client.CloseIdleConnections()
}

func (c *Channel) updateWorker() {
	if len(c.worker) > 0 {
		return
	}

	for i := 0; i < RoutersCount(); i++ {
		router := CurrentRouter()

		logger.Println("updateWorker")

		link := "https://" + router + "/api/request?fn=get_workers"

		response, err := c.client.Get(link)

		type ResponseGetWorkersItem struct {
			Host   string `json:"host"`
			Scores int    `json:"scores"`
		}

		type ResponseGetWorkers struct {
			Workers []ResponseGetWorkersItem `json:"workers"`
		}
		var resp ResponseGetWorkers

		if err != nil {
			logger.Println(err)
		} else {
			content, _ := ioutil.ReadAll(response.Body)
			s := strings.TrimSpace(string(content))
			response.Body.Close()

			if json.Unmarshal([]byte(s), &resp) == nil {
				if resp.Workers != nil {
					sort.Slice(resp.Workers, func(i, j int) bool {
						return resp.Workers[i].Scores > resp.Workers[j].Scores
					})

					if len(resp.Workers) > 0 {
						correctWorker := ""

						for _, w := range resp.Workers {
							if _, ok := c.wrongWorkers[w.Host]; !ok {
								correctWorker = w.Host
								break
							}
						}

						if len(correctWorker) > 0 {
							c.worker = correctWorker
							if c.binClient != nil {
								c.binClient.Stop()
							}
							//c.worker = "w002.gazer.cloud"
							c.binClient = bin_client.New(c.worker+":1077", "public", "public", c.chProcessingData)
							c.binClient.Start()
							logger.Println("updateWorker result: ", c.worker)
							break
						} else {
							c.wrongWorkers = make(map[string]time.Time)
						}
					}
				}
			}
		}

		SetNextRouter()
	}

	c.client.CloseIdleConnections()
}

type Item struct {
	Name  string
	Value string
}

type ChannelFullInfo struct {
	Id       string
	Password string
	Name     string
	Items    []string
}

type ChannelInfo struct {
	Id   string
	Name string
}

func (c *Channel) thWorker() {
	for !c.stopping {
		items, err := c.GetValues()
		if err == nil {
			err = c.Write(items, false)
			if err != nil {
				logger.Println("sending to cloud error", err)
			}
		}
		var itemsToRemove []common_interfaces.Item
		itemsToRemove, err = c.GetRemovedValues()
		if len(itemsToRemove) > 0 {
			logger.Println("Removed: ", itemsToRemove)
			if err == nil {
				err = c.Write(itemsToRemove, true)
				if err != nil {
					logger.Println("sending to cloud error", err)
				}
			}
			c.itemsToRemove = make([]string, 0)
		}
		for i := 0; i < 10; i++ {
			if c.stopping {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
	c.started = false
}

func (c *Channel) processData(task bin_client.BinFrameTask) {
	if task.Frame.Channel == "#count_of_subscribers#" {

		var resp protocol.ResponseCountOfSubscribers
		err := json.Unmarshal(task.Frame.Data, &resp)
		if err != nil {
			logger.Println("System #count_of_subscribers# error", err)
			return
		}

		if resp.Channel == c.channelId {
			c.subscribers = resp.Count
		}

		//logger.Println("System #count_of_subscribers# from ", task.Client.GetRemoteAddr(), " = ", resp.)

		/*addedChannels := c.AddTranslateChannelsToClient(task.Client, channels)
		task.Client.SendNeedChannelOK(addedChannels)*/

		return
	}
}

func (c *Channel) thIncomingTraffic() {
	for !c.stopping {
		var frame bin_client.BinFrameTask
		select {
		case frame = <-c.chProcessingData:
			c.processData(frame)
		case <-time.After(50 * time.Millisecond):
		}
	}

	c.started = false
}

func (c *Channel) AddItems(items []string) error {
	c.mtx.Lock()
	err := c.addItems(items)
	c.mtx.Unlock()
	return err
}

func (c *Channel) addItems(items []string) error {
	for _, itemToAdd := range items {
		itemExists := false
		for _, existingItem := range c.items {
			if existingItem == itemToAdd {
				itemExists = true
				break
			}
		}

		if !itemExists {
			c.items = append(c.items, itemToAdd)
		}
	}
	return nil
}

func (c *Channel) RemoveItems(items []string) error {
	c.mtx.Lock()
	c.removeItems(items)
	c.mtx.Unlock()
	return nil
}

func (c *Channel) removeItems(items []string) {
	for _, itemToRemove := range items {
		for existingItemIndex, existingItem := range c.items {
			if existingItem == itemToRemove {
				c.itemsToRemove = append(c.itemsToRemove, itemToRemove)
				c.items = append(c.items[:existingItemIndex], c.items[existingItemIndex+1:]...)
				break
			}
		}
	}
}

func (c *Channel) RemoveAllItems() error {
	c.mtx.Lock()
	c.itemsToRemove = append(c.itemsToRemove, c.items...)
	c.items = make([]string, 0)
	c.mtx.Unlock()
	return nil
}

func (c *Channel) Write(items []common_interfaces.Item, removed bool) error {
	c.updateWorker()
	return c.GazerWrite(c.worker, c.channelId, c.password, items, removed, c.subscribers > 0)
}

func (c *Channel) GazerWrite(host string, channelId string, psw string, items []common_interfaces.Item, removed bool, fullData bool) error {
	if len(c.channelId) < 1 {
		return errors.New("no channel registered")
	}

	if len(c.worker) < 1 {
		return errors.New("no worker")
	}

	type ResponseGetItem struct {
		Name  string `json:"n"`
		Value string `json:"v"`
		UOM   string `json:"u"`
		DT    int64  `json:"t"`
		Flags string `json:"f"`
	}

	type ResponseGet struct {
		Items []ResponseGetItem `json:"is"`
	}

	var req ResponseGet
	req.Items = make([]ResponseGetItem, 0)
	for _, item := range items {
		if fullData || item.Name == ".service/name" {
			var respItem ResponseGetItem
			if removed {
				respItem.Name = item.Name
				respItem.Value = ""
				respItem.UOM = ""
				respItem.Flags = "d"
				respItem.DT = item.Value.DT
			} else {
				respItem.Name = item.Name
				respItem.Value = item.Value.Value
				respItem.UOM = item.Value.UOM
				respItem.Flags = ""
				respItem.DT = item.Value.DT
			}
			req.Items = append(req.Items, respItem)
		}
	}

	dataBytes, _ := json.Marshal(&req)

	compressed := dataBytes
	compressed = make([]byte, 0)
	compressed = append(compressed, []byte("G_JSON__")...)
	compressed = append(compressed, dataBytes...)

	var frame bin_client.BinFrame
	frame.Channel = channelId
	frame.Password = psw
	frame.Data = compressed

	if len(compressed) > 100*1024 {
		return errors.New("too big frame")
	}

	if c.binClient != nil {
		c.binClient.SendData(&frame)
		if c.binClient.LastError() != nil {
			c.wrongWorkers[c.worker] = time.Now().UTC()
			c.worker = ""
		}
	}
	return nil
}

func (c *Channel) GetValues() ([]common_interfaces.Item, error) {
	values := make([]common_interfaces.Item, 0)
	c.mtx.Lock()
	{
		// name of the channel
		var itemName common_interfaces.Item
		itemName.Id = 0
		itemName.Value.DT = time.Now().UTC().UnixNano() / 1000
		itemName.Value.Value = c.name
		itemName.Name = ".service/name"
		values = append(values, itemName)
	}

	for _, itemName := range c.items {
		val, err := c.iDataStorage.GetItem(itemName)
		if err == nil {
			values = append(values, val)
		}
	}
	c.mtx.Unlock()
	return values, nil
}

func (c *Channel) GetRemovedValues() ([]common_interfaces.Item, error) {
	values := make([]common_interfaces.Item, 0)
	c.mtx.Lock()
	for _, itemName := range c.itemsToRemove {
		val, err := c.iDataStorage.GetItem(itemName)
		if err == nil {
			val.Value.Flags = "d"
			values = append(values, val)
		} else {
			val = common_interfaces.Item{}
			val.Id = 42
			val.Name = itemName
			val.Value.UOM = ""
			val.Value.Value = ""
			val.Value.DT = time.Now().UTC().UnixNano() / 1000
			val.Value.Flags = "d"
			values = append(values, val)
		}
	}
	c.mtx.Unlock()
	return values, nil
}

func (c *Channel) RenameItems(oldPrefix string, newPrefix string) {
	c.mtx.Lock()
	itemsToRemove := make([]string, 0)
	itemsToAdd := make([]string, 0)

	for _, item := range c.items {
		if strings.HasPrefix(item, oldPrefix) {
			itemsToRemove = append(itemsToRemove, item)
			newItem := strings.Replace(item, oldPrefix, newPrefix, 1)
			itemsToAdd = append(itemsToAdd, newItem)
		}
	}

	c.removeItems(itemsToRemove)
	c.addItems(itemsToAdd)

	c.mtx.Unlock()
}
