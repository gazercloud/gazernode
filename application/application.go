package application

import (
	"flag"
	"fmt"
	"github.com/gazercloud/gazernode/utilities/paths"
	"github.com/kardianos/osext"
	"github.com/kardianos/service"
	"log"
	"os"
)

var Name string
var Description string
var ServiceName string
var ServiceDisplayName string
var ServiceDescription string
var ServiceRunFunc func() error
var ServiceStopFunc func()

type Application struct {
	Name    string
	Version string
}

var App Application

func SetAppPath() {
	exePath, _ := osext.ExecutableFolder()
	os.Chdir(exePath)

}

func init() {
	SetAppPath()
}

var ServerDataPathArgument string

func TryService() bool {
	setupFlagPtr := flag.Bool("setup", false, "Install to /usr/local/bin")
	serviceFlagPtr := flag.Bool("service", false, "Run as service")
	installFlagPtr := flag.Bool("install", false, "Install service")
	uninstallFlagPtr := flag.Bool("uninstall", false, "Uninstall service")
	startFlagPtr := flag.Bool("start", false, "Start service")
	stopFlagPtr := flag.Bool("stop", false, "Stop service")
	serverPath := flag.String("path", paths.ProgramDataFolder1()+"/gazer", "Server data path")

	flag.Parse()

	ServerDataPathArgument = *serverPath

	if *setupFlagPtr {
		setupPosix()
		return true
	}

	if *serviceFlagPtr {
		runService()
		return true
	}

	if *installFlagPtr {
		InstallService()
		return true
	}

	if *uninstallFlagPtr {
		UninstallService()
		return true
	}

	if *startFlagPtr {
		StartService()
		return true
	}

	if *stopFlagPtr {
		StopService()
		return true
	}

	return false
}

func NewSvcConfig() *service.Config {
	var SvcConfig = &service.Config{
		Name:        ServiceName,
		DisplayName: ServiceDisplayName,
		Description: ServiceDescription,
	}
	SvcConfig.Arguments = append(SvcConfig.Arguments, "-service")
	return SvcConfig
}

func InstallService() {
	fmt.Println("Service installing")
	prg := &program{}
	s, err := service.New(prg, NewSvcConfig())
	if err != nil {
		log.Fatal(err)
	}
	err = s.Install()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Service installed")
}

func UninstallService() {
	fmt.Println("Service uninstalling")
	prg := &program{}
	s, err := service.New(prg, NewSvcConfig())
	if err != nil {
		log.Fatal(err)
	}
	err = s.Uninstall()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Service uninstalled")
}

func StartService() {
	fmt.Println("Service starting")
	prg := &program{}
	s, err := service.New(prg, NewSvcConfig())
	if err != nil {
		log.Println(err)
	}
	err = s.Start()
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println("Service started")
	}
}

func StopService() {
	fmt.Println("Service stopping")
	prg := &program{}
	s, err := service.New(prg, NewSvcConfig())
	if err != nil {
		log.Println(err)
	}
	err = s.Stop()
	if err != nil {
		log.Println(err)
		return
	} else {
		fmt.Println("Service stopped")
	}
}

func runService() {
	prg := &program{}
	s, err := service.New(prg, NewSvcConfig())
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		_ = logger.Error(err)
	}
}

var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	return ServiceRunFunc()
}

func (p *program) Stop(s service.Service) error {
	ServiceStopFunc()
	return nil
}
