package cmd

import (
	"errors"
	"fmt"
)

func (c *Session) cmdLs(p []string) error {
	if c.currentPathIsItem() {
		err := errors.New("current path is data item")
		return err
	}

	if c.currentPathIsUnit() {
		resp, err := c.client.GetUnitValues(c.currentUnitName)
		for _, item := range resp {
			fmt.Println("["+fmt.Sprint(item.Id)+"]", item.Name, "= "+item.Value.Value+" "+item.Value.UOM)
		}
		return err
	}

	resp, err := c.client.GetUnitStateAll()
	if err != nil {
		return err
	}
	for _, item := range resp.Items {
		fmt.Println("["+item.UnitId+"]", item.UnitName, "<"+item.TypeName+">", "= "+item.Value+" "+item.UOM)
	}

	return nil
}
