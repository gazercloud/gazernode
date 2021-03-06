package system

import (
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/history"
	"github.com/gazercloud/gazernode/system/cloud"
	"github.com/gazercloud/gazernode/system/last_values"
	"github.com/gazercloud/gazernode/system/public_channel"
	"github.com/gazercloud/gazernode/system/resources"
	"github.com/gazercloud/gazernode/system/units/units_system"
	"sync"
)

type System struct {
	items       []*common_interfaces.Item
	itemsByName map[string]*common_interfaces.Item
	itemsById   map[uint64]*common_interfaces.Item
	nextItemId  uint64

	requester common_interfaces.Requester

	unitsSystem *units_system.UnitsSystem

	publicChannels  *public_channel.Cloud
	cloudConnection *cloud.Connection

	history   *history.History
	resources *resources.Resources

	users      []*common_interfaces.User
	userByName map[string]*common_interfaces.User
	sessions   map[string]*UserSession

	apiCallsCount int

	stopping bool

	mtx sync.Mutex
}

func NewSystem() *System {
	var c System
	c.items = make([]*common_interfaces.Item, 0)
	c.itemsByName = make(map[string]*common_interfaces.Item)
	c.itemsById = make(map[uint64]*common_interfaces.Item)

	c.cloudConnection = cloud.NewConnection()

	c.publicChannels = public_channel.NewCloud(&c)
	c.unitsSystem = units_system.New(&c)
	c.history = history.NewHistory()
	c.resources = resources.NewResources()

	c.users = make([]*common_interfaces.User, 0)
	c.userByName = make(map[string]*common_interfaces.User)
	c.sessions = make(map[string]*UserSession)

	return &c
}

func (c *System) SetRequester(requester common_interfaces.Requester) {
	c.requester = requester
	c.cloudConnection.SetRequester(c.requester)
}

func (c *System) Start() {
	c.LoadConfig()
	c.loadSessions()

	items := last_values.Read()
	for _, item := range items {
		if realItem, ok := c.itemsByName[item.Name]; ok {
			realItem.Value = item.Value
		}
	}
	c.cloudConnection.Start()
	c.publicChannels.Start()
	c.history.Start()
	c.unitsSystem.Start()

}

func (c *System) Stop() {
	c.stopping = true
	c.unitsSystem.Stop()
	c.publicChannels.Stop()
	c.history.Stop()
	c.cloudConnection.Stop()
	c.SaveConfig()
	c.saveSessions()
	last_values.Write(c.items)
}

func (c *System) RegApiCall() {
	c.mtx.Lock()
	c.apiCallsCount++
	c.mtx.Unlock()
}

func (c *System) Exec(function string, request []byte, host string) ([]byte, error) {
	return c.requester.RequestJson(function, request, host)
}

func (c *System) StatGazerNode() (res common_interfaces.StatGazerNode) {
	return
}

func (c *System) StatGazerCloud() (res common_interfaces.StatGazerCloud) {
	res = c.cloudConnection.Stat()
	return
}
