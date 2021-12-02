package cmd

import "fmt"

func (c *Session) cmdUnits(p []string) error {
	resp, err := c.client.GetUnitStateAll()
	if err != nil {
		return err
	}
	for _, item := range resp.Items {
		fmt.Println("["+item.UnitId+"]", item.UnitName, "<"+item.TypeName+">", "= "+item.Value+" "+item.UOM)
	}
	return nil
}
