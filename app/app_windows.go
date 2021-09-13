package app

import (
	"github.com/gazercloud/gazernode/application"
	"github.com/gazercloud/gazernode/forms/mainform"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/utilities/paths"
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uiforms"
)

func RunDesktop() {
	logger.Init(paths.HomeFolder() + "/gazer/log_ui")

	if *runServerFlagPtr {
		start(application.ServerDataPathArgument)
	}

	ui.InitUISystem()

	{
		var form mainform.MainForm
		uiforms.StartMainForm(&form)
		form.Dispose()
	}

	if *runServerFlagPtr {
		stop()
	}
}
