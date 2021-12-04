package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/gazercloud/gazernode/gazer_client"
)

func (c *Session) cmdConnect(p []string) error {
	if len(p) != 3 {
		return errors.New("wrong parameters. Usage: connect http://addr:port user password")
	}

	c.client = gazer_client.New(p[0])
	resp, err := c.client.SessionOpen(p[1], p[2])
	if err != nil {
		return err
	}
	color.Set(color.FgGreen)
	fmt.Println("Session:", resp.SessionToken)
	color.Unset()
	c.save()
	return nil
}
