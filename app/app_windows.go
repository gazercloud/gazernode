package app

import (
	"github.com/gazercloud/gazernode/forms"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/utilities/paths"
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uiforms"
)

func RunDesktop() {
	logger.Init(paths.HomeFolder() + "/gazer/log_ui")

	if *runServerFlagPtr {
		start()
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
