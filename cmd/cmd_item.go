package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"sort"
	"strings"
)

func (c *Session) cmdItem(p []string) error {
	if len(p) < 1 {
		return errors.New("wrong parameters")
	}

	cmd := p[0]
	p = p[1:]

	switch cmd {
	case "list":
		{
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
	case "write":
		{
			itemName := ""
			itemValue := ""

			if len(p) == 0 {
				return errors.New("not enough parameters")
			}

			if len(p) == 1 {
				if c.currentItem == "" {
					return errors.New("no current item")
				}
				itemName = c.currentItem
				itemValue = p[0]
			}

			if len(p) == 2 {
				itemName = p[0]
				itemValue = p[1]
			}

			if len(p) > 2 {
				return errors.New("too many parameters")
			}

			err := c.client.Write(itemName, itemValue)
			return err
		}
	case "remove":
		{
			if len(p) != 1 {
				return errors.New("wrong parameters")
			}
			err := c.client.DataItemRemove([]string{p[0]})
			if err == nil {
				color.Set(color.FgGreen)
				fmt.Println("item", p[0], "removed")
				color.Unset()
			}
			return err
		}
	}
	return errors.New("wrong parameters")
}
