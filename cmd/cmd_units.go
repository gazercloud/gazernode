package cmd

import (
	"fmt"
	"github.com/fatih/color"
)

func (c *Session) cmdUnits(p []string) error {
	resp, err := c.client.GetUnitStateAll()
	if err != nil {
		return err
	}
	color.Set(color.FgGreen)
	for _, item := range resp.Items {
		fmt.Println("["+item.UnitId+"]", item.UnitName, "<"+item.TypeName+">", "= "+item.Value+" "+item.UOM)
	}
	color.Unset()
	return nil
}
