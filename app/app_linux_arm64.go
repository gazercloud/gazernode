package app

import (
	"fmt"
	"github.com/gazercloud/gazernode/application"
	"github.com/gazercloud/gazernode/utilities/logger"
	"github.com/gazercloud/gazernode/utilities/paths"
)

func RunDesktop() {
	logger.Init(paths.HomeFolder() + "/gazer/log_ui")
	start(application.ServerDataPathArgument)
	logger.Println("Started as console application")
	logger.Println("Press ENTER to stop")
	_, _ = fmt.Scanln()
	stop()
}
