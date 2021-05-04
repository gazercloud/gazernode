package system

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/settings"
	"github.com/gazercloud/gazernode/system/public_channel"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"io/ioutil"
)

type Config struct {
	Users      []common_interfaces.User              `json:"users"`
	Units      []units_common.UnitInfo               `json:"units"`
	Channels   []public_channel.ChannelFullInfo      `json:"channels"`
	Items      []common_interfaces.ItemConfiguration `json:"items"`
	NextItemId uint64                                `json:"next_item_id"`
}

func (c *System) SaveConfig() error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	var conf Config
	conf.Units = c.unitsSystem.Units()
	conf.Channels = c.cloud.ChannelsFullInfo()
	conf.Users = make([]common_interfaces.User, 0)
	for _, u := range c.users {
		conf.Users = append(conf.Users, *u)
	}

	conf.Items = make([]common_interfaces.ItemConfiguration, 0)
	for _, item := range c.items {
		var itemConf common_interfaces.ItemConfiguration
		itemConf.Id = item.Id
		itemConf.Name = item.Name
		conf.Items = append(conf.Items, itemConf)
	}

	conf.NextItemId = c.nextItemId

	configBytes, err := json.MarshalIndent(&conf, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(settings.ServerDataPath()+"/config.json", configBytes, 0666)
	if err != nil {
		return err
	}

	return nil
}

func (c *System) LoadConfig() error {
	configString, err := ioutil.ReadFile(settings.ServerDataPath() + "/config.json")
	if err == nil {
		var conf Config
		err = json.Unmarshal([]byte(configString), &conf)
		if err != nil {
			return err
		}

		c.users = make([]*common_interfaces.User, 0)
		for _, u := range conf.Users {
			us := u
			c.users = append(c.users, &us)
			c.userByName[us.Name] = &us
		}

		c.nextItemId = conf.NextItemId

		realMaxItemId := uint64(0)
		for _, itemConf := range conf.Items {
			var item common_interfaces.Item
			item.Id = itemConf.Id
			item.Name = itemConf.Name

			c.items = append(c.items, &item)
			c.itemsByName[item.Name] = &item
			c.itemsById[item.Id] = &item

			if item.Id > realMaxItemId {
				realMaxItemId = item.Id
			}
		}

		if c.nextItemId < realMaxItemId+1 {
			c.nextItemId = realMaxItemId + 1
		}

		for _, sens := range conf.Units {
			c.unitsSystem.AddUnit(sens.Type, sens.Id, sens.Name, sens.Config)
		}

		for _, ch := range conf.Channels {
			c.cloud.AddChannel(ch.Id, ch.Password, ch.Name)
			c.cloud.AddItems([]string{ch.Id}, ch.Items)
		}
	}

	if len(c.users) == 0 {
		logger.Println("System loadUsers adding default user")
		var u common_interfaces.User
		u.Name = DefaultUserName
		u.PasswordHash = c.hashPassword(DefaultUserPassword)
		c.users = append(c.users, &u)
		c.userByName[u.Name] = &u
	}

	return err
}
