package app

import (
	"fmt"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/utilities/paths"
	"github.com/stianeikeland/go-rpio/v4"
	"time"
)

func ttt() {
	err := rpio.Open()

	if err != nil {
		return
	}

	pin10 := rpio.Pin(10)
	pin10.PullUp()
	pin10.Input()

	for {
		st := pin10.Read()
		time.Sleep(200 * time.Millisecond)
		logger.Println("r", st)
	}

	rpio.Close()
}

func RunDesktop() {
	logger.Init(paths.HomeFolder() + "/gazer/log_ui")
	ttt()
	return
	start()
	logger.Println("Started as console application")
	logger.Println("Press ENTER to stop")
	_, _ = fmt.Scanln()
	stop()
}
