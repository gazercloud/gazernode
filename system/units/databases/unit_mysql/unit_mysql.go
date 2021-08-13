package unit_mysql

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"github.com/jackc/pgx"

	"time"
)

type UnitMySQL struct {
	units_common.Unit
	address  string
	user     string
	password string
	database string
	query    string
	asTable  bool
	periodMs int
}

func New() common_interfaces.IUnit {
	var c UnitMySQL
	return &c
}

const (
	ItemNameResult = "Result"
)

var Image []byte

func init() {
	Image = resources.R_files_sensors_unit_database_mysql_png
}

func (c *UnitMySQL) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	meta.Add("addr", "address", "localhost:5432", "string", "", "", "")
	meta.Add("user", "user", "", "string", "", "", "")
	meta.Add("password", "password", "", "string", "", "", "")
	meta.Add("database", "database", "", "string", "", "", "")
	meta.Add("query", "Query", "", "text", "", "", "")
	meta.Add("as_table", "As table", "", "bool", "", "", "")
	meta.Add("period", "Period, ms", "1000", "num", "0", "999999", "")
	return meta.Marshal()
}

func (c *UnitMySQL) InternalUnitStart() error {
	var err error
	c.SetString(ItemNameResult, "", "")
	c.SetMainItem(ItemNameResult)

	type Config struct {
		Addr     string  `json:"addr"`
		User     string  `json:"user"`
		Password string  `json:"password"`
		Database string  `json:"database"`
		Query    string  `json:"query"`
		Period   float64 `json:"period"`
		AsTable  bool    `json:"as_table"`
	}

	var config Config
	err = json.Unmarshal([]byte(c.GetConfig()), &config)
	if err != nil {
		err = errors.New("config error")
		c.SetString(ItemNameResult, err.Error(), "error")
		return err
	}

	c.query = config.Query
	if c.query == "" {
		err = errors.New("wrong file")
		c.SetString(ItemNameResult, err.Error(), "error")
		return err
	}

	c.periodMs = int(config.Period)
	if c.periodMs < 100 {
		err = errors.New("wrong period")
		c.SetString(ItemNameResult, err.Error(), "error")
		return err
	}

	c.address = config.Addr
	c.user = config.User
	c.password = config.Password
	c.database = config.Database
	c.asTable = config.AsTable

	go c.Tick()
	return nil
}

func (c *UnitMySQL) InternalUnitStop() {
}

func (c *UnitMySQL) Tick() {
	var db *pgx.Conn
	var err error

	c.Started = true
	dtOperationTime := time.Now().UTC()
	for !c.Stopping {
		for {
			if c.Stopping || time.Now().Sub(dtOperationTime) > time.Duration(c.periodMs)*time.Millisecond {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if c.Stopping {
			break
		}
		dtOperationTime = time.Now().UTC()

		if db == nil {
			db, err = pgx.Connect(pgx.ConnConfig{
				Host:                 c.address,
				Port:                 5432,
				Database:             c.database,
				User:                 c.user,
				Password:             c.password,
				TLSConfig:            nil,
				UseFallbackTLS:       false,
				FallbackTLSConfig:    nil,
				Logger:               nil,
				LogLevel:             0,
				Dial:                 nil,
				RuntimeParams:        nil,
				OnNotice:             nil,
				CustomConnInfo:       nil,
				CustomCancel:         nil,
				PreferSimpleProtocol: false,
				TargetSessionAttrs:   "",
			})

			if err != nil {
				db = nil
				continue
			}
		}

		var res *pgx.Rows
		res, err = db.Query(c.query)
		if err == nil {
			if c.asTable {
				var values []interface{}

				for res.Next() {
					values, err = res.Values()
					if err == nil {
						if len(values) > 0 {
							name := ""
							value := ""
							uom := ""
							if len(values) > 0 {
								name = fmt.Sprint(values[0])
								if len(values) > 1 {
									value = fmt.Sprint(values[1])
								}
								if len(values) > 2 {
									uom = fmt.Sprint(values[2])
								}
								c.SetString(name, value, uom)
								c.SetError("")
							}
						} else {
							err = errors.New("no values returned")
						}
					}
				}
			} else {
				var values []interface{}
				if res.Next() {
					values, err = res.Values()
				} else {
					err = errors.New("no values returned")
				}
				if err == nil {
					if len(values) > 0 {
						c.SetString(ItemNameResult, fmt.Sprint(values[0]), "")
						c.SetError("")
					} else {
						err = errors.New("no values returned")
					}
				}
			}
			res.Close()
		}

		if err != nil {
			c.SetString(ItemNameResult, err.Error(), "error")
			c.SetError(err.Error())
		}
	}

	if db != nil {
		_ = db.Close()
		db = nil
	}

	c.SetString(ItemNameResult, "", "stopped")
	c.Started = false
}
