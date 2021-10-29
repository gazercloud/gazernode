package unit_http_json_requester

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type UnitHttpJsonRequester struct {
	units_common.Unit
	addr              string
	timeoutMs         int
	periodMs          int
	receivedVariables map[string]string
}

func New() common_interfaces.IUnit {
	var c UnitHttpJsonRequester
	c.receivedVariables = make(map[string]string)
	return &c
}

const (
	ItemNameStatus = "Status"
)

var Image []byte

func init() {
	Image = resources.R_files_sensors_unit_network_json_requester_png
}

func (c *UnitHttpJsonRequester) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	meta.Add("addr", "Address", "localhost:445", "string", "", "", "")
	meta.Add("period", "Period, ms", "1000", "num", "0", "999999", "0")
	meta.Add("timeout", "Timeout, ms", "1000", "num", "0", "999999", "0")
	return meta.Marshal()
}

func (c *UnitHttpJsonRequester) InternalUnitStart() error {
	var err error

	type Config struct {
		Addr    string  `json:"addr"`
		Timeout float64 `json:"timeout"`
		Period  float64 `json:"period"`
	}

	var config Config
	err = json.Unmarshal([]byte(c.GetConfig()), &config)
	if err != nil {
		err = errors.New("config error")
		c.SetString(ItemNameStatus, err.Error(), "error")
		return err
	}

	c.addr = config.Addr
	if c.addr == "" {
		err = errors.New("wrong address")
		c.SetString(ItemNameStatus, err.Error(), "error")
		return err
	}

	c.timeoutMs = int(config.Timeout)
	if c.timeoutMs < 100 {
		err = errors.New("wrong timeout (<100)")
		c.SetString(ItemNameStatus, err.Error(), "error")
		return err
	}
	if c.timeoutMs > 10000 {
		err = errors.New("wrong timeout (>10000)")
		c.SetString(ItemNameStatus, err.Error(), "error")
		return err
	}

	c.periodMs = int(config.Period)
	if c.periodMs < 100 {
		err = errors.New("wrong period")
		c.SetString(ItemNameStatus, err.Error(), "error")
		return err
	}

	if c.periodMs < c.timeoutMs {
		err = errors.New("wrong period (<timeout)")
		c.SetString(ItemNameStatus, err.Error(), "error")
		return err
	}
	if c.periodMs < 100 {
		err = errors.New("wrong period (<100)")
		c.SetString(ItemNameStatus, err.Error(), "error")
		return err
	}
	if c.periodMs > 60000 {
		err = errors.New("wrong period (>60000)")
		c.SetString(ItemNameStatus, err.Error(), "error")
		return err
	}

	c.receivedVariables = make(map[string]string)

	c.SetMainItem(ItemNameStatus)

	c.SetString(ItemNameStatus, "", "")
	go c.Tick()
	return nil
}

func (c *UnitHttpJsonRequester) InternalUnitStop() {
}

func (c *UnitHttpJsonRequester) Tick() {
	var err error
	c.Started = true
	dtLastTime := time.Now().UTC()

	for !c.Stopping {
		for {
			if c.Stopping || time.Now().Sub(dtLastTime) > time.Duration(c.periodMs)*time.Millisecond {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if c.Stopping {
			c.SetString(ItemNameStatus, "stopped", "")
			break
		}
		dtLastTime = time.Now().UTC()

		var resp string
		resp, err = c.HttpCall(c.addr)
		if err == nil {

			var unm map[string]interface{}
			err = json.Unmarshal([]byte(resp), &unm)
			if err == nil {
				for key, value := range unm {
					valueAsString := fmt.Sprint(value)
					c.SetString(key, valueAsString, "")
					c.receivedVariables[key] = valueAsString
				}
				c.SetString(ItemNameStatus, "ok", "")
			} else {
				c.SetString(ItemNameStatus, err.Error(), "error")
			}

		} else {
			c.SetString(ItemNameStatus, err.Error(), "error")
		}
	}

	for vName, _ := range c.receivedVariables {
		c.SetString(vName, "", "stopped")
	}

	c.SetString(ItemNameStatus, "", "stopped")
	c.Started = false
}

func (c *UnitHttpJsonRequester) HttpCall(url string) (responseString string, err error) {
	var client *http.Client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{},
	}
	client = &http.Client{Transport: tr}
	client.Timeout = 1 * time.Second

	var response *http.Response
	response, err = client.Get(url)
	if err == nil {
		content, _ := ioutil.ReadAll(response.Body)
		responseString = strings.TrimSpace(string(content))

		response.Body.Close()
	}

	client.CloseIdleConnections()

	return responseString, err
}
