package public_channel

/*
type Routers struct {
	mtx                sync.Mutex
	currentRouterIndex int
	routers            []string
}

var routersInstance *Routers

func init() {
	routersInstance = &Routers{}
	routersInstance.routers = make([]string, 0)
	routersInstance.routers = append(routersInstance.routers, "r001.gazer.cloud")
	routersInstance.routers = append(routersInstance.routers, "r002.gazer.cloud")
	routersInstance.currentRouterIndex = 0
}

func RoutersCount() int {
	count := 0
	routersInstance.mtx.Lock()
	count = len(routersInstance.routers)
	routersInstance.mtx.Unlock()
	return count
}

func CurrentRouter() string {
	routersInstance.mtx.Lock()
	router := routersInstance.routers[routersInstance.currentRouterIndex]
	routersInstance.mtx.Unlock()
	return router
}

func SetNextRouter() {
	routersInstance.mtx.Lock()
	routersInstance.currentRouterIndex++
	if routersInstance.currentRouterIndex >= len(routersInstance.routers) {
		routersInstance.currentRouterIndex = 0
	}
	logger.Println("Routers: switched to ", routersInstance.routers[routersInstance.currentRouterIndex])
	routersInstance.mtx.Unlock()
}
*/
