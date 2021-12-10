package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
)

func (c *Session) cmdCloud(p []string) error {
	if len(p) == 0 {
		resp, err := c.client.CloudState()

		color.Set(color.FgCyan)
		fmt.Println("Cloud Connection State")
		color.Set(color.FgGreen)

		if err != nil {
			return err
		}
		if resp.Connected {
			color.Set(color.FgGreen)
		} else {
			color.Set(color.FgRed)
		}
		fmt.Println("Connected:", resp.Connected)
		if resp.LoggedIn {
			color.Set(color.FgGreen)
		} else {
			color.Set(color.FgRed)
		}
		fmt.Println("LoggedIn:", resp.LoggedIn)
		if resp.UserName != "" {
			color.Set(color.FgGreen)
		} else {
			color.Set(color.FgRed)
		}
		fmt.Println("UserName:", resp.UserName)
		if resp.SessionKey != "" {
			color.Set(color.FgGreen)
		} else {
			color.Set(color.FgRed)
		}
		fmt.Println("SessionKey:", resp.SessionKey)
		if resp.NodeId != "" {
			color.Set(color.FgGreen)
		} else {
			color.Set(color.FgRed)
		}
		fmt.Println("NodeId:", resp.NodeId)
		if resp.CurrentRepeater != "" {
			color.Set(color.FgGreen)
		} else {
			color.Set(color.FgRed)
		}
		fmt.Println("CurrentRepeater:", resp.CurrentRepeater)
		color.Unset()
		return nil
	}

	if p[0] == "logout" {
		err := c.client.CloudLogout()
		return err
	}

	if p[0] == "login" {
		if len(p) != 3 {
			return errors.New("not enough parameters")
		}
		err := c.client.CloudLogin(p[1], p[2])
		return err
	}

	err := errors.New("unknown command")

	return err
}
