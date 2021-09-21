package unit_http_json_items_server

/*
type Item struct {
	Name    string `json:"item_name"`
	UrlPath string `json:"url_path"`
}

type Config struct {
	Port  float64 `json:"port"`
	Items []Item  `json:"items"`
}

type UnitHttpJsonServer struct {
	units_common.Unit
	config Config

	srv *http.Server
	r   *mux.Router
}

func New() common_interfaces.IUnit {
	var c UnitHttpJsonServer
	return &c
}

const (
	ItemNameStatus = "Status"
)

var Image []byte

func init() {
	Image = resources.R_files_sensors_sensor_network_json_items_png
}

func (c *UnitHttpJsonServer) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	meta.Add("port", "Port", "8888", "num", "0", "99999", "")
	t1 := meta.Add("items", "Items", "", "table", "", "", "")
	t1.Add("item_name", "Item Name", "Unit/Item", "string", "", "", "")
	t1.Add("url_path", "URL Path", "item", "string", "", "", "")
	return meta.Marshal()
}

func (c *UnitHttpJsonServer) InternalUnitStart() error {
	var err error

	err = json.Unmarshal([]byte(c.GetConfig()), &c.config)
	if err != nil {
		err = errors.New("config error")
		c.SetString(ItemNameStatus, err.Error(), "error")
		return err
	}

	if c.config.Port < 0 || c.config.Port > 65535 {
		err = errors.New("wrong port")
		c.SetString(ItemNameStatus, err.Error(), "error")
		return err
	}

	c.SetMainItem(ItemNameStatus)

	c.SetString(ItemNameStatus, "", "")

	c.r = mux.NewRouter()

	// Static files
	c.r.NotFoundHandler = http.HandlerFunc(c.processFile)

	c.srv = &http.Server{Addr: ":" + fmt.Sprint(int(c.config.Port))}
	c.srv.Handler = c.r

	go c.Tick()

	return nil
}

func (c *UnitHttpJsonServer) InternalUnitStop() {
	c.Stopping = true
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := c.srv.Shutdown(ctx); err != nil {
		logger.Println(err)
	}
	c.srv = nil
}

func (c *UnitHttpJsonServer) Tick() {
	c.Started = true

	c.SetString(ItemNameStatus, "", "started!")
	err := c.srv.ListenAndServe()
	if err != nil {
		c.SetString(ItemNameStatus, err.Error(), "error")
	}

	dtLastTime := time.Now().UTC()

	for !c.Stopping {
		for {
			if c.Stopping || time.Now().Sub(dtLastTime) > time.Duration(1000)*time.Millisecond {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if c.Stopping {
			c.SetString(ItemNameStatus, "stopped", "")
			break
		}
	}
	c.Started = false
}

func (c *UnitHttpJsonServer) processFile(w http.ResponseWriter, r *http.Request) {
	var err error

	URL := *r.URL
	if URL.Path == "/" || URL.Path == "" {
		URL.Path = "/index.html"
	}

	found := false
	for _, item := range c.config.Items {
		if "/"+item.UrlPath == URL.Path {
			found = true
			var v common_interfaces.ItemValue
			v, err = c.GetItem(item.Name)
			if err != nil {
				w.WriteHeader(502)
				_, _ = w.Write([]byte(err.Error()))
			} else {
				bs, err := json.MarshalIndent(v, "", " ")

				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				_, err = w.Write(bs)
				if err != nil {
					logger.Println("HttpServer sendError w.Write error:", err)
				}
			}
		}
	}

	if !found {
		w.WriteHeader(404)
		_, _ = w.Write([]byte("no item found"))
	}
}
*/
