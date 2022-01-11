package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
)

func (c *Session) cmdLs(p []string) error {
	color.Set(color.FgGreen)
	defer color.Unset()

	if c.currentPathIsItem() {
		resp, err := c.client.GetItemsValues([]string{c.currentItem})
		if err != nil {
			return err
		}

		color.Set(color.FgMagenta)
		fmt.Println("Item", c.currentItem)
		color.Set(color.FgGreen)

		for _, item := range resp {
			fmt.Println("Id", item.Id)
			fmt.Println("Name", item.Name)
			fmt.Println("Value", item.Value)
			fmt.Println("UOM:", item.Value)
			fmt.Println("DT:", item.Value)
		}

		return err
	}

	if c.currentPathIsUnit() {
		resp, err := c.client.GetUnitValues(c.currentUnitId)

		color.Set(color.FgMagenta)
		fmt.Println("Items of Unit", c.currentUnitId)
		color.Set(color.FgGreen)

		for _, item := range resp {
			shortName := strings.ReplaceAll(item.Name, c.currentUnitId+"/", "")
			if strings.HasPrefix(shortName, ".service") {
				continue
			}
			fmt.Println(shortName, "= "+item.Value.Value+" "+item.Value.UOM)
		}
		return err
	}

	resp, err := c.client.GetUnitStateAll()
	if err != nil {
		return err
	}
	color.Set(color.FgMagenta)
	fmt.Println("Units:")
	color.Set(color.FgGreen)
	for _, item := range resp.Items {
		fmt.Println(item.UnitId, "= "+item.Value+" "+item.UOM)
	}

	return nil
}
