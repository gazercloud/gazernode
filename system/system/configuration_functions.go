package system

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/system/cloud"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"github.com/gazercloud/gazernode/utilities/paths"
	"io/ioutil"
)

type Config struct {
	Units      []units_common.UnitInfo
	Channels   []cloud.ChannelFullInfo
	Items      []common_interfaces.ItemConfiguration
	NextItemId uint64
}

func (c *System) SaveConfig() error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	var conf Config
	conf.Units = c.unitsSystem.ListOfUnits()
	conf.Channels = c.cloud.ChannelsFullInfo()

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

	err = ioutil.WriteFile(paths.ProgramDataFolder()+"/gazer/config.json", configBytes, 0666)
	if err != nil {
		return err
	}

	return nil
}

func (c *System) LoadConfig() error {
	configString, err := ioutil.ReadFile(paths.ProgramDataFolder() + "/gazer/config.json")
	if err != nil {
		return err
	}

	var conf Config
	err = json.Unmarshal([]byte(configString), &conf)
	if err != nil {
		return err
	}

	c.nextItemId = conf.NextItemId

	for _, itemConf := range conf.Items {
		var item common_interfaces.Item
		item.Id = itemConf.Id
		item.Name = itemConf.Name

		c.items = append(c.items, &item)
		c.itemsByName[item.Name] = &item
		c.itemsById[item.Id] = &item
	}

	for _, sens := range conf.Units {
		c.unitsSystem.AddUnit(sens.Type, sens.Id, sens.Name, sens.Config)
	}

	for _, ch := range conf.Channels {
		c.cloud.AddChannel(ch.Id, ch.Password, ch.Name)
		c.cloud.AddItems([]string{ch.Id}, ch.Items)
	}
	return nil
}
