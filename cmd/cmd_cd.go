package cmd

import (
	"errors"
	"strings"
)

func (c *Session) cmdCd(p []string) error {
	if len(p) != 1 {
		return errors.New("wrong path")
	}

	newCurrentUnitId := c.currentUnitId
	//newCurrentUnitName := c.currentUnitName
	newCurrentItem := c.currentItem

	if p[0] == ".." {
		if c.currentPathIsItem() {
			newCurrentItem = ""
		} else {
			newCurrentUnitId = ""
			//newCurrentUnitName = ""
		}
	} else {
		if p[0] == "/" {
			newCurrentUnitId = ""
			//newCurrentUnitName = ""
		} else {

			if c.currentPathIsItem() {
				return errors.New("item have no items")
			}

			if c.currentPathIsUnit() {
				// to item
				unitValues, err := c.client.GetUnitValues(c.currentUnitId)
				if err != nil {
					return err
				}
				found := false
				foundItemName := ""
				for _, item := range unitValues {
					shortName := strings.ReplaceAll(item.Name, c.currentUnitId+"/", "")
					if shortName == p[0] {
						found = true
						foundItemName = item.Name
					}
				}

				if !found {
					return errors.New("no item found")
				}

				newCurrentItem = foundItemName
			} else {
				// to unit
				items, err := c.client.GetUnitStateAll()
				if err != nil {
					return err
				}
				found := false
				foundIndex := -1
				for i, item := range items.Items {
					if item.UnitId == p[0] {
						found = true
						foundIndex = i
						break
					}
				}

				if !found {
					return errors.New("no unit found")
				}

				newCurrentUnitId = items.Items[foundIndex].UnitId
				//newCurrentUnitName = items.Items[foundIndex].UnitName
			}
		}
	}

	c.currentUnitId = newCurrentUnitId
	//c.currentUnitName = newCurrentUnitName
	c.currentItem = newCurrentItem
	return nil
}
