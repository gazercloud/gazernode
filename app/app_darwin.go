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
	ui.InitUISystem()
	logger.Init(paths.HomeFolder() + "/gazer/log_ui")

	if *runServerFlagPtr {
		start(application.ServerDataPathArgument)
	}

	{
		var form mainform.MainForm
		uiforms.StartMainForm(&form)
		form.Dispose()
	}

	if *runServerFlagPtr {
		stop()
	}
}
