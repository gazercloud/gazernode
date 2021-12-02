package cmd

import "fmt"

func (c *Session) cmdDisconnect(p []string) error {
	fmt.Println("Disconnect")
	return nil
}
