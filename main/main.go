// Copyright 2018-2021 GazerCloud Authors. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.
// The list of authors can be found in the AUTHORS file

package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gazercloud/gazernode/app"
	"github.com/gazercloud/gazernode/application"
	"net"
	"time"
)

func ttt() {
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: time.Second * 3}, "tcp", "51.38.98.192:1077", &tls.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.Close()
}

func main() {
	application.Name = "Gazer"
	application.ServiceName = "Gazer"
	application.ServiceDisplayName = "Gazer"
	application.ServiceDescription = "Gazer Service"
	application.ServiceRunFunc = app.RunAsService
	application.ServiceStopFunc = app.StopService

	if !application.TryService() {
		app.RunDesktop()
	}
}
