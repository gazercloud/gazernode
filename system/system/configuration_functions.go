package system

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"github.com/gazercloud/gazernode/utilities/logger"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
)

type Config struct {
	Name       string                                `json:"name"`
	Users      []common_interfaces.UserConfiguration `json:"users"`
	Units      []units_common.UnitInfo               `json:"units"`
	Items      []common_interfaces.ItemConfiguration `json:"items"`
	NextItemId uint64                                `json:"next_item_id"`
}

func (c *System) SaveConfig() error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	var conf Config
	conf.Name = c.nodeName
	conf.Units = c.unitsSystem.Units()
	conf.Users = make([]common_interfaces.UserConfiguration, 0)
	for _, u := range c.users {
		var userConfig common_interfaces.UserConfiguration
		userConfig.Name = u.Name
		userConfig.PasswordHash = u.PasswordHash
		userConfig.Properties = make([]*common_interfaces.ItemProperty, 0)
		for _, p := range u.Properties {
			userConfig.Properties = append(userConfig.Properties, &common_interfaces.ItemProperty{
				Name:  p.Name,
				Value: p.Value,
			})
		}
		sort.Slice(userConfig.Properties, func(i, j int) bool {
			return userConfig.Properties[i].Name < userConfig.Properties[j].Name
		})

		conf.Users = append(conf.Users, userConfig)
	}

	conf.Items = make([]common_interfaces.ItemConfiguration, 0)
	for _, item := range c.items {
		var itemConf common_interfaces.ItemConfiguration
		itemConf.Id = item.Id
		itemConf.Name = item.Name
		itemConf.Properties = make([]*common_interfaces.ItemProperty, 0)
		for _, p := range item.Properties {
			itemConf.Properties = append(itemConf.Properties, &common_interfaces.ItemProperty{
				Name:  p.Name,
				Value: p.Value,
			})
		}
		sort.Slice(itemConf.Properties, func(i, j int) bool {
			return itemConf.Properties[i].Name < itemConf.Properties[j].Name
		})
		conf.Items = append(conf.Items, itemConf)
	}

	conf.NextItemId = c.nextItemId

	configBytes, err := json.MarshalIndent(&conf, "", " ")
	if err != nil {
		return err
	}

	_, err = os.Stat(c.ss.ServerDataPath())
	if err != nil {
		err = os.MkdirAll(c.ss.ServerDataPath(), 0777)
		if err != nil {
			logger.Println("System SaveConfig MkdirAll error: ", err)
		}
	}

	err = ioutil.WriteFile(c.ss.ServerDataPath()+"/config.json", configBytes, 0666)
	if err != nil {
		return err
	}

	return nil
}

func (c *System) LoadConfig() error {
	configString, err := ioutil.ReadFile(c.ss.ServerDataPath() + "/config.json")
	if err == nil {
		var conf Config
		err = json.Unmarshal([]byte(configString), &conf)
		if err != nil {
			return err
		}

		c.nodeName = conf.Name

		c.users = make([]*common_interfaces.User, 0)
		for _, u := range conf.Users {
			var user common_interfaces.User
			user.Name = u.Name
			user.PasswordHash = u.PasswordHash
			user.Properties = make(map[string]*common_interfaces.ItemProperty)
			for _, p := range u.Properties {
				user.Properties[p.Name] = &common_interfaces.ItemProperty{
					Name:  p.Name,
					Value: p.Value,
				}
			}

			c.users = append(c.users, &user)
			c.userByName[user.Name] = &user
		}

		c.nextItemId = conf.NextItemId

		realMaxItemId := uint64(0)
		for _, itemConf := range conf.Items {
			var item common_interfaces.Item
			item.Id = itemConf.Id
			item.Name = itemConf.Name
			item.Properties = make(map[string]*common_interfaces.ItemProperty)
			item.TranslateToItems = make(map[uint64]*common_interfaces.Item)
			for _, p := range itemConf.Properties {
				item.Properties[p.Name] = &common_interfaces.ItemProperty{
					Name:  p.Name,
					Value: p.Value,
				}
			}

			c.items = append(c.items, &item)
			c.itemsByName[item.Name] = &item
			c.itemsById[item.Id] = &item

			if item.Id > realMaxItemId {
				realMaxItemId = item.Id
			}
		}

		c.mtx.Lock()
		c.applyItemsProperties()
		c.mtx.Unlock()

		if c.nextItemId < realMaxItemId+1 {
			c.nextItemId = realMaxItemId + 1
		}

		for _, unitConfig := range conf.Units {
			unit, errAddUnit := c.unitsSystem.AddUnit(unitConfig.Type, unitConfig.Id, unitConfig.DisplayName, unitConfig.Config, false)
			if errAddUnit == nil {
				unit.PropSet(unitConfig.Properties)
			}
		}

		/*for _, ch := range conf.Channels {
			c.publicChannels.AddChannel(ch.Id, ch.Password, ch.Name, ch.NeedToStartAfterLoad)
			c.publicChannels.AddItems([]string{ch.Id}, ch.Items)
		}*/
	}

	for _, user := range c.users {
		if user.PasswordHash == "" {
			user.PasswordHash = c.hashPassword(user.Name)
			logger.Println("WARNING! Password for user " + user.Name + " set to " + user.Name)
		}
	}

	if len(c.users) == 0 {
		logger.Println("System loadUsers adding default user")
		passwordBuffer := make([]byte, 8)
		binary.LittleEndian.PutUint64(passwordBuffer, rand.Uint64())
		for i := 0; i < 8; i++ {
			xorKey := byte(rand.Intn(255))
			passwordBuffer[i] = passwordBuffer[i] ^ xorKey
		}
		userPassword := hex.EncodeToString(passwordBuffer)

		_, err = os.Stat(c.ss.ServerDataPath())
		if err != nil {
			err = os.MkdirAll(c.ss.ServerDataPath(), 0777)
			if err != nil {
				logger.Println("System LoadConfig loadUsers MkdirAll error: ", err)
			}
		}

		defaultAdminPasswordFilename := c.ss.ServerDataPath() + "/default_admin_password.txt"

		err = ioutil.WriteFile(defaultAdminPasswordFilename, []byte(userPassword), 0655)
		if err != nil {
			logger.Println("cannot write ", defaultAdminPasswordFilename)
		}

		var u common_interfaces.User
		u.Name = DefaultUserName
		u.Properties = make(map[string]*common_interfaces.ItemProperty)
		u.PasswordHash = c.hashPassword(userPassword)
		c.users = append(c.users, &u)
		c.userByName[u.Name] = &u
	}

	return err
}
