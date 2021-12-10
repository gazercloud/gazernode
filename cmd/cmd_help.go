package cmd

import (
	"fmt"
	"github.com/fatih/color"
)

func cmdDisplayHelpLine(cmd string, text string) {
	fmt.Println(cmd, " - ", text)
}

func (c Session) cmdHelp(p []string) error {
	color.Set(color.FgYellow)
	fmt.Println("Help")
	fmt.Println()

	color.Set(color.FgMagenta, color.Bold)

	cmdDisplayHelpLine("connect", "Connect to the node")
	cmdDisplayHelpLine("disconnect", "Connect to the node")

	color.Unset()

	return nil
}
