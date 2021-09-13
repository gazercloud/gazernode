package app

import (
	"fmt"
	"github.com/gazercloud/gazernode/forms/mainform"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/utilities"
	"github.com/gazercloud/gazernode/utilities/paths"
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uiforms"
)

func RunDesktop() {
	ui.InitUISystem()

	logger.Init(paths.HomeFolder() + "/gazer/log_ui")

	fmt.Println("Is ROOT:", utilities.IsRoot())

	if *runServerFlagPtr {
		start()
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
