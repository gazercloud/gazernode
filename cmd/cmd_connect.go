package cmd

import (
	"errors"
	"fmt"
	"github.com/gazercloud/gazernode/gazer_client"
)

func (c *Session) cmdConnect(p []string) error {
	if len(p) != 3 {
		return errors.New("wrong parameters")
	}

	c.client = gazer_client.New(p[0])
	resp, err := c.client.SessionOpen(p[1], p[2])
	if err != nil {
		return err
	}
	fmt.Println("Session:", resp.SessionToken)
	return nil
}
