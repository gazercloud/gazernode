package unit_ssl

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"net"
	"strings"
	"time"
)

type UnitSSL struct {
	units_common.Unit
	domain            string
	timeoutMs         int
	periodMs          int
	receivedVariables map[string]string
}

func New() common_interfaces.IUnit {
	var c UnitSSL
	c.receivedVariables = make(map[string]string)
	return &c
}

const (
	ItemNameDaysLeft = "DaysLeft"
	ItemNameStatus   = "Status"
	ItemNameTime     = "Time"
	ItemNameIP       = "IP"
	ItemNameDomain   = "Domain"
)

var Image []byte

func init() {
	Image = resources.R_files_sensors_unit_network_ssl_png
}

func (c *UnitSSL) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	meta.Add("domain", "Domain", "example.com", "string", "", "", "")
	meta.Add("period", "Period, ms", "10000", "num", "0", "999999", "0")
	meta.Add("timeout", "Timeout, ms", "1000", "num", "0", "999999", "0")
	return meta.Marshal()
}

func (c *UnitSSL) InternalUnitStart() error {
	var err error

	type Config struct {
		Domain  string  `json:"domain"`
		Timeout float64 `json:"timeout"`
		Period  float64 `json:"period"`
	}

	c.SetString(ItemNameDaysLeft, "", "")

	var config Config
	err = json.Unmarshal([]byte(c.GetConfig()), &config)
	if err != nil {
		err = errors.New("config error")
		c.SetString(ItemNameStatus, err.Error(), "error")
		return err
	}

	c.domain = config.Domain
	if c.domain == "" {
		err = errors.New("wrong domain")
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

	c.SetMainItem(ItemNameDaysLeft)

	c.SetString(ItemNameStatus, "", "")

	c.SetString(ItemNameTime, "", "")
	c.SetString(ItemNameDomain, "", "")
	c.SetString(ItemNameIP, "", "")

	go c.Tick()
	return nil
}

func (c *UnitSSL) InternalUnitStop() {
}

func (c *UnitSSL) Tick() {
	var err error
	var lastIP string

	c.Started = true
	dtLastTime := time.Now().UTC().Add(-1 * time.Hour)

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

		var resolvedAddr *net.IPAddr
		resolvedAddr, err = net.ResolveIPAddr("", c.domain)
		if err == nil {
			ip := resolvedAddr.IP.String()
			if ip != lastIP {
				lastIP = ip
				c.SetString(ItemNameIP, ip, "-")
			}

			timeBegin := time.Now()

			var conn *tls.Conn
			conn, err = tls.DialWithDialer(&net.Dialer{Timeout: time.Second * 1}, "tcp", c.domain+":443", &tls.Config{})
			if err != nil {
				c.SetString(ItemNameStatus, err.Error(), "error")
				c.SetString(ItemNameDaysLeft, "", "error")
				c.SetString(ItemNameDaysLeft, "", "error")
			} else {
				if conn.ConnectionState().PeerCertificates != nil {
					for _, cert := range conn.ConnectionState().PeerCertificates {

						found := false

						for _, domainName := range cert.DNSNames {
							if strings.Contains(domainName, c.domain) {
								found = true
							}
						}

						if found {
							c.SetFloat64(ItemNameDaysLeft, cert.NotAfter.Sub(time.Now()).Hours()/24, "days", 3)
						}
					}

					c.SetString(ItemNameStatus, "ok", "")
				}
				conn.Close()
			}

			timeEnd := time.Now()
			duration := timeEnd.Sub(timeBegin)
			c.SetInt(ItemNameTime, int(duration.Milliseconds()), "ms")

		} else {
			c.SetString(ItemNameStatus, err.Error(), "")
			c.SetString(ItemNameIP, "", "error")
		}
	}

	c.SetString(ItemNameTime, "", "stopped")
	c.SetString(ItemNameDomain, "", "stopped")
	c.SetString(ItemNameIP, "", "stopped")

	c.SetString(ItemNameStatus, "", "stopped")
	c.SetString(ItemNameDaysLeft, "", "stopped")

	c.Started = false
}
