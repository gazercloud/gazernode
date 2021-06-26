module github.com/gazercloud/gazernode

go 1.16

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fogleman/gg v1.3.0
	github.com/gazercloud/gazerui v1.0.9
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20201108214237-06ea97f0c265
	github.com/go-pg/pg v8.0.7+incompatible // indirect
	github.com/go-ping/ping v0.0.0-20210216210419-25d1413fb7bb
	github.com/google/uuid v1.2.0
	github.com/gorilla/mux v1.8.0
	github.com/jackc/pgx v3.6.2+incompatible // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/josephspurrier/goversioninfo v1.2.0
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
	github.com/kardianos/service v1.2.0
	github.com/kbinani/win v0.3.0
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/pkg/errors v0.9.1 // indirect
	github.com/shirou/gopsutil v3.21.1+incompatible
	github.com/srwiley/oksvg v0.0.0-20210209000435-a757b9cbd472
	github.com/srwiley/rasterx v0.0.0-20200120212402-85cb7272f5e9
	github.com/tarm/serial v0.0.0-20180830185346-98f6abe2eb07
	github.com/yuin/gopher-lua v0.0.0-20200816102855-ee81675732da
	go.bug.st/serial v1.1.2
	golang.org/x/image v0.0.0-20210220032944-ac19c3e999fb
	golang.org/x/sys v0.0.0-20210223212115-eede4237b368
	golang.org/x/text v0.3.5
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
	mellium.im/sasl v0.2.1 // indirect
)

replace github.com/gazercloud/gazerui => ../gazerui
replace github.com/go-gl/glfw/v3.3/glfw => ../glfw/v3.3/glfw
