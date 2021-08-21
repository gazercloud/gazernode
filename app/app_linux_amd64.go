package app

import (
	"github.com/gazercloud/gazernode/application"
	"github.com/gazercloud/gazernode/forms"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/utilities/paths"
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uiforms"
)

func RunDesktop() {

	/*logger.Init(paths.HomeFolder() + "/gazer/log_ui")
	start(application.ServerDataPathArgument)
	logger.Println("Started as console application")
	logger.Println("Press ENTER to stop")
	_, _ = fmt.Scanln()
	stop()*/

	logger.Init(paths.HomeFolder() + "/gazer/log_ui")

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
	}
}
