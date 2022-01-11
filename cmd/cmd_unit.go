package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
)

func (c *Session) cmdUnit(p []string) error {
	if len(p) < 1 {
		return errors.New("wrong parameters. Need to specify unit command")
	}

	cmd := p[0]
	p = p[1:]

	switch cmd {
	case "list":
		{
			resp, err := c.client.GetUnitStateAll()
			if err != nil {
				return err
			}
			color.Set(color.FgGreen)
			for _, item := range resp.Items {
				fmt.Println(item.UnitId, "<"+item.Status+">", "= "+item.Value+" "+item.UOM)
			}
			color.Unset()
			return nil
		}
	case "start":
		{
			if len(p) != 1 {
				return errors.New("wrong parameters")
			}

			unitStateAll, err := c.client.GetUnitStateAll()
			if err != nil {
				return err
			}
			for _, u := range unitStateAll.Items {
				if u.UnitId == p[0] {
					err = c.client.StartUnits([]string{u.UnitId})
					if err == nil {
						color.Set(color.FgGreen)
						fmt.Println("unit", u.UnitId, "started")
						color.Unset()
						return nil
					}
					return err
				}
			}

			return errors.New("no unit found")
		}
	case "stop":
		{
			if len(p) != 1 {
				return errors.New("wrong parameters")
			}
			unitStateAll, err := c.client.GetUnitStateAll()
			if err != nil {
				return err
			}
			for _, u := range unitStateAll.Items {
				if u.UnitId == p[0] {
					err = c.client.StopUnits([]string{u.UnitId})
					if err == nil {
						color.Set(color.FgGreen)
						fmt.Println("unit", u.UnitId, "stopped")
						color.Unset()
						return nil
					}
					return err
				}
			}

			return errors.New("no unit found")
		}
	}

	return errors.New("unknown unit command")
}
