package app

import (
	"allece.com/system/core/grid/stats"
	"allece.com/system/core/paths"
	"allece.com/system/gazer/gazer/httpserver"
	"allece.com/system/gazer/gazer/system"
	"allece.com/system/gazer/gazer_common/hostid"
	"allece.com/system/gazer/gazer_common/logger"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
)

var runServerFlagPtr = flag.Bool("server", false, "Run server")

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
		httpServer.Stop()
	}
}

func RunDesktop() {
	logger.Init(paths.ProgramDataFolder() + "/gazer/log_ui")
	if *runServerFlagPtr {
		start()
	}

	fmt.Println("Started ...")
	fmt.Scanln()

	runtime.GC()
	runtime.GC()
	debug.FreeOSMemory()

	stats.Dump()

	if *runServerFlagPtr {
		stop()
	}
}

func RunAsService() error {
	logger.Init(paths.ProgramDataFolder() + "/gazer/log_service")
	logger.Println("")
	logger.Println("------------------------------")
	logger.Println("Started as Service")
	logger.Println("------------------------------")
	start()
	return nil
}

func StopService() {
	logger.Println("")
	logger.Println("------------------------------")
	logger.Println("Service stopped")
	logger.Println("------------------------------")
	logger.Println("")
	stop()
}
