package app

import (
	"flag"
	"github.com/gazercloud/gazernode/application"
	"github.com/gazercloud/gazernode/system/httpserver"
	"github.com/gazercloud/gazernode/system/settings"
	"github.com/gazercloud/gazernode/system/system"
	"github.com/gazercloud/gazernode/utilities/hostid"
	"github.com/gazercloud/gazernode/utilities/logger"
)

var httpServer *httpserver.HttpServer
var sys *system.System
var runServerFlagPtr = flag.Bool("server", false, "Run server")

func start(dataPath string) {
	hostid.InitHostId()

	ss := settings.NewSettings()
	ss.SetServerDataPath(dataPath)

	sys = system.NewSystem(ss)
	httpServer = httpserver.NewHttpServer(sys)
	sys.SetRequester(httpServer)
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

func RunAsService() error {
	logger.Init(application.ServerDataPathArgument + "/log_service")
	logger.Println("Started as Service")
	start(application.ServerDataPathArgument)
	return nil
}

func StopService() {
	logger.Println("Service stopped")
	stop()
}
