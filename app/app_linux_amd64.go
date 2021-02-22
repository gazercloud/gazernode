package app

import (
	"github.com/gazercloud/gazernode/system/httpserver"
	"github.com/gazercloud/gazernode/system/system"
	"github.com/gazercloud/gazernode/utilities/hostid"
	"os"
	"path/filepath"
)

func CurrentExePath() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

var httpServer *httpserver.HttpServer
var sys *system.System

func start() {
	hostid.InitHostId()

	sys = system.NewSystem()
	httpServer = httpserver.NewHttpServer(sys)
	sys.Start()
	httpServer.Start()
}

func stop() {
	if sys != nil {
		sys.Stop()
	}
	if httpServer != nil {
		_ = httpServer.Stop()
	}
}

func RunDesktop() {
	start()

	ui.InitUISystem()

	{
		var form forms.MainForm
		uiforms.StartMainForm(&form)
		form.Dispose()
	}

	stop()
}

func RunAsService() error {
	//logger.Init(paths.ProgramDataFolder() + "/gazer/log_service")
	//logger.Println("Started as Service")
	start()
	return nil
}

func StopService() {
	//logger.Println("Service stopped")
	stop()
}
