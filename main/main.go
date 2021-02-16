// Copyright 2018-2021 GazerCloud Authors. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.
// The list of authors can be found in the AUTHORS file

package main

import (
	"github.com/gazercloud/gazernode/app"
	"github.com/gazercloud/gazernode/application"
)

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
