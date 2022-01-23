package system

import (
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/system/cloud"
	"github.com/gazercloud/gazernode/system/history"
	"github.com/gazercloud/gazernode/system/resources"
	"github.com/gazercloud/gazernode/system/settings"
	"github.com/gazercloud/gazernode/system/units/units_system"
	"sync"
	"time"
)

type System struct {
	nodeName string
	ss       *settings.Settings

	items       []*common_interfaces.Item
	itemsByName map[string]*common_interfaces.Item
	itemsById   map[uint64]*common_interfaces.Item
	nextItemId  uint64

	requester common_interfaces.Requester

	unitsSystem *units_system.UnitsSystem

	cloudConnection *cloud.Connection

	history   *history.History
	resources *resources.Resources

	users      []*common_interfaces.User
	userByName map[string]*common_interfaces.User
	sessions   map[string]*UserSession

	itemWatchers map[string]*ItemWatcher

	apiCallsCount int

	stopping bool
	stopped  bool

	maintenanceLastValuesDT time.Time

	mtx sync.Mutex
}

func NewSystem(ss *settings.Settings) *System {
	var c System
	c.ss = ss
	c.itemWatchers = make(map[string]*ItemWatcher)
	c.items = make([]*common_interfaces.Item, 0)
	c.itemsByName = make(map[string]*common_interfaces.Item)
	c.itemsById = make(map[uint64]*common_interfaces.Item)

	c.cloudConnection = cloud.NewConnection(c.ss.ServerDataPath())

	c.unitsSystem = units_system.New(&c)
	c.history = history.NewHistory(c.ss)
	c.resources = resources.NewResources(c.ss)

	c.users = make([]*common_interfaces.User, 0)
	c.userByName = make(map[string]*common_interfaces.User)
	c.sessions = make(map[string]*UserSession)

	return &c
}

func (c *System) Settings() *settings.Settings {
	return c.ss
}

func (c *System) SetRequester(requester common_interfaces.Requester) {
	c.requester = requester
	c.cloudConnection.SetRequester(c.requester)
}

func (c *System) Start() {
	c.stopping = false
	c.stopped = false

	c.LoadConfig()
	c.loadSessions()

	items := c.ReadLastValues()
	for _, item := range items {
		if realItem, ok := c.itemsByName[item.Name]; ok {
			realItem.Value = item.Value
		}
	}
	c.cloudConnection.Start()
	c.history.Start()
	c.unitsSystem.Start()

	go c.thMaintenance()
}

func (c *System) Stop() {
	c.stopping = true
	c.unitsSystem.Stop()
	c.history.Stop()
	c.cloudConnection.Stop()
	c.SaveConfig()
	c.saveSessions()

	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
		if c.stopped {
			break
		}
	}

	c.WriteLastValues(c.items)
}

func (c *System) RegApiCall() {
	c.mtx.Lock()
	c.apiCallsCount++
	c.mtx.Unlock()
}

func (c *System) StatGazerNode() (res common_interfaces.StatGazerNode) {
	return
}

func (c *System) StatGazerCloud() (res common_interfaces.StatGazerCloud) {
	res = c.cloudConnection.Stat()
	return
}

func (c *System) thMaintenance() {
	for !c.stopping {
		for i := 0; i < 10; i++ {
			time.Sleep(100 * time.Millisecond)
			if c.stopping {
				break
			}
		}
		if c.stopping {
			break
		}

		c.maintenanceLastValues()
	}
	c.stopped = true
}

func (c *System) maintenanceLastValues() {
	if time.Now().Sub(c.maintenanceLastValuesDT) > 10*time.Second {
		c.maintenanceLastValuesDT = time.Now()
		c.WriteLastValues(c.items)
		c.RemoveOldLastValuesFiles()
	}
}
