package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"sort"
	"strings"
)

func (c *Session) cmdItems(p []string) error {
	color.Set(color.FgCyan)
	defer color.Unset()
	fmt.Println("Items:")
	items, err := c.client.GetAllItems()
	if err != nil {
		return err
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})
	for _, item := range items {
		if strings.Contains(item.Name, "/.service/") {
			continue
		}
		if item.Value.UOM != "error" {
			color.Set(color.FgGreen)
		} else {
			color.Set(color.FgRed)
		}
		fmt.Println(item.Name + " = " + item.Value.Value + " " + item.Value.UOM)
	}
	return nil
}
