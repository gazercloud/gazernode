package cmd

import (
	"github.com/fatih/color"
	"github.com/gazercloud/gazernode/gazer_client"
)

func (c *Session) cmdDisconnect(p []string) error {
	color.Set(color.FgGreen)
	c.client = gazer_client.New("")
	color.Unset()
	c.save()
	return nil
}
