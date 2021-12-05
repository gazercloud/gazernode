package system

import (
	"github.com/gazercloud/gazernode/common_interfaces"
)

func (c *System) ResAdd(name string, tp string, content []byte) (string, error) {
	return c.resources.Add(name, tp, content)
}

func (c *System) ResSet(id string, thumbnail []byte, content []byte) error {
	return c.resources.Set(id, thumbnail, content)
}

func (c *System) ResGet(id string) (*common_interfaces.ResourcesItem, error) {
	return c.resources.Get(id)
}

func (c *System) ResList(tp string, filter string, offset int, maxCount int) common_interfaces.ResourcesInfo {
	return c.resources.List(tp, filter, offset, maxCount)
}

func (c *System) ResRemove(id string) error {
	return c.resources.Remove(id)
}

func (c *System) ResRename(id string, name string) error {
	return c.resources.Rename(id, name)
}
