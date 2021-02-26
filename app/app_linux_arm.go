package app

import (
	"fmt"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/utilities/paths"
)

func RunDesktop() {
	logger.Init(paths.HomeFolder() + "/gazer/log_ui")
	start()
	logger.Println("Started as console application")
	logger.Println("Press ENTER to stop")
	_, _ = fmt.Scanln()
	stop()
}
