package app

import (
	"flag"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/settings"
	"github.com/gazercloud/gazernode/system/httpserver"
	"github.com/gazercloud/gazernode/system/repeater_bin_client"
	"github.com/gazercloud/gazernode/system/system"
	"github.com/gazercloud/gazernode/utilities/hostid"
)

var httpServer *httpserver.HttpServer
var sys *system.System
var runServerFlagPtr = flag.Bool("server", false, "Run server")

var repeaterConnection *repeater_bin_client.RepeaterBinClient
var chProcessingData chan repeater_bin_client.BinFrameTask

func start() {
	hostid.InitHostId()

	chProcessingData = make(chan repeater_bin_client.BinFrameTask)
	repeaterConnection = repeater_bin_client.New("db05.gazer.cloud:1077", "user", "password", chProcessingData)
	repeaterConnection.Start()

	sys = system.NewSystem()
	httpServer = httpserver.NewHttpServer(sys, chProcessingData)
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
