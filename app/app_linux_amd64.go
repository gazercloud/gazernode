package app

import (
	"fmt"
	"github.com/gazercloud/gazernode/application"
	"github.com/gazercloud/gazernode/cmd"
	"github.com/gazercloud/gazernode/utilities/logger"
	"github.com/gazercloud/gazernode/utilities/paths"
)

func RunDesktop() {
	logger.Init(paths.HomeFolder() + "/gazer/log_ui")
	if len(os.Args) == 1 {
		cmd.Console()
		return
	}

	if len(os.Args) == 1 {
		cmd.Console()
		return
	}

	start(application.ServerDataPathArgument)
	logger.Println("Started as console application")
	logger.Println("Press ENTER to stop")
	_, _ = fmt.Scanln()
	stop()

	/*logger.Init(paths.HomeFolder() + "/gazer/log_ui")

	if *runServerFlagPtr {
		start(application.ServerDataPathArgument)
	}

	ui.InitUISystem()

	{
		var form forms.MainForm
		uiforms.StartMainForm(&form)
		form.Dispose()
	}

	if *runServerFlagPtr {
		stop()
	}*/
}
