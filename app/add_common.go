package app

import (
	"flag"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/settings"
	"github.com/gazercloud/gazernode/system/httpserver"
	"github.com/gazercloud/gazernode/system/system"
	"github.com/gazercloud/gazernode/utilities/hostid"
)

var httpServer *httpserver.HttpServer
var sys *system.System
var runServerFlagPtr = flag.Bool("server", false, "Run server")

func start() {
	hostid.InitHostId()

	sys = system.NewSystem()
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
	logger.Init(settings.ServerDataPath() + "/log_service")
	logger.Println("Started as Service")
	start()
	return nil
}

func StopService() {
	logger.Println("Service stopped")
	stop()
}
