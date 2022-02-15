package system

import (
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/system/protocols/nodeinterface"
)

func (c *System) ResAdd(name string, tp string, content []byte) (string, error) {
	return c.resources.Add(name, tp, content)
}

func (c *System) ResSet(id string, suffix string, offset int64, content []byte) error {
	return c.resources.Set(id, suffix, offset, content)
}

func (c *System) ResGet(id string, offset int64, size int64) (nodeinterface.ResourceGetResponse, error) {
	return c.resources.Get(id, offset, size)
}

func (c *System) ResGetThumbnail(id string) (*common_interfaces.ResourcesItem, error) {
	return c.resources.GetThumbnail(id)
}

func (c *System) ResList(tp string, filter string, offset int, maxCount int) common_interfaces.ResourcesInfo {
	return c.resources.List(tp, filter, offset, maxCount)
}

func (c *System) ResRemove(id string) error {
	return c.resources.Remove(id)
}

func (c *System) ResRename(id string, props []nodeinterface.PropItem) error {
	return c.resources.Rename(id, props)
}

func (c *System) ResourcePropSet(resourceId string, props []nodeinterface.PropItem) error {
	return c.resources.PropSet(resourceId, props)
}

func (c *System) ResourcePropGet(resourceId string) ([]nodeinterface.PropItem, error) {
	return c.resources.PropGet(resourceId)
}
