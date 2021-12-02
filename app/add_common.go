package app

import (
	"flag"
	"github.com/gazercloud/gazernode/application"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/settings"
	"github.com/gazercloud/gazernode/system/httpserver"
	"github.com/gazercloud/gazernode/system/system"
	"github.com/gazercloud/gazernode/utilities/hostid"
)

var httpServer *httpserver.HttpServer
var sys *system.System
var runServerFlagPtr = flag.Bool("server", false, "Run server")
var cloudLoginFlagPtr = flag.Bool("login", false, "Cloud Login")
var cloudUserNameFlagPtr = flag.String("username", "", "Cloud User Name")
var cloudPasswordFlagPtr = flag.String("password", "", "Cloud Password")
var cloudNodeFlagPtr = flag.String("node", "", "Node Id")
var cloudLogoutFlagPtr = flag.Bool("logout", false, "Cloud Logout")

func cmd111(arguments []string) bool {
	logger.Println("Arguments:")
	logger.Println(arguments)

	if *cloudLoginFlagPtr {
		if *cloudUserNameFlagPtr != "" && *cloudPasswordFlagPtr != "" && *cloudNodeFlagPtr != "" {
			logger.Println("LOGIN: ", *cloudUserNameFlagPtr, *cloudPasswordFlagPtr, *cloudNodeFlagPtr)
		} else {
			logger.Println("Error: no username/password/node_id provided")
		}
	}

	return false
}

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
