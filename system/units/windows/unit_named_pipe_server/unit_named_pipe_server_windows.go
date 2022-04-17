package unit_system_named_pipe_server

import (
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"github.com/gazercloud/gazernode/utilities/logger"
	"github.com/gazercloud/gazernode/utilities/uom"
	"github.com/google/uuid"
	"net"
	"strings"
	"sync"
	"time"
)

type UnitSystemNamedPipeServer struct {
	units_common.Unit

	// Runtime
	listener          net.Listener
	receivedVariables map[string]string
	connectedClients  map[string]*UnitSystemNamedPipeServerConnectedClient
	mtx               sync.Mutex

	// Config
	config Config
}

type Config struct {
	PipeName string `json:"pipe_name"`
}

const (
	ItemNameStatus = "Status"
)

type UnitSystemNamedPipeServerConnectedClient struct {
	id         string
	connection net.Conn
}

var Image []byte

func init() {
	Image = resources.R_files_sensors_unit_computer_memory_png
}

func New() common_interfaces.IUnit {
	var c UnitSystemNamedPipeServer
	c.receivedVariables = make(map[string]string)
	c.connectedClients = make(map[string]*UnitSystemNamedPipeServerConnectedClient)
	return &c
}

func (c *UnitSystemNamedPipeServer) InternalUnitStart() error {
	var err error
	c.SetMainItem(ItemNameStatus)

	err = json.Unmarshal([]byte(c.GetConfig()), &c.config)
	if err != nil {
		err = errors.New("config error")
		c.SetString(ItemNameStatus, err.Error(), "error")
		return err
	}
	if len(c.config.PipeName) < 1 {
		err = errors.New("empty pipe name")
		c.SetString(ItemNameStatus, err.Error(), "error")
		return err
	}

	go c.Tick()
	return nil
}

func (c *UnitSystemNamedPipeServer) InternalUnitStop() {
}

func (c *UnitSystemNamedPipeServer) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	meta.Add("pipe_name", "Pipe Name", "q_gazer_pipe", "string", "", "", "")
	return meta.Marshal()
}

func (c *UnitSystemNamedPipeServer) Tick() {
	logger.Println("UnitSystemNamedPipeServer Tick begin")
	c.Started = true
	c.SetStringForAll("", uom.STARTED)

	go c.serve() // separated thread

	for !c.Stopping {
		for i := 0; i < 10; i++ {
			if c.Stopping {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

	if c.listener != nil {
		c.listener.Close()
	}

	c.removeAllClients()

	c.Started = false
	c.SetStringForAll("", uom.STOPPED)

	logger.Println("UnitSystemNamedPipeServer Tick end")
}

func (c *UnitSystemNamedPipeServer) serve() {
	logger.Println("UnitSystemNamedPipeServer serve begin")
	var err error
	c.listener, err = winio.ListenPipe("\\\\.\\pipe\\"+c.config.PipeName, nil)
	if err != nil {
		logger.Println("UnitSystemNamedPipeServer ListenPipe error:", err)
		return
	}
	for !c.Stopping {
		conn, err := c.listener.Accept()
		if err != nil {
			logger.Println("UnitSystemNamedPipeServer accepting error", err)
			break
		}
		c.mtx.Lock()
		client := UnitSystemNamedPipeServerConnectedClient{
			id:         uuid.NewString(),
			connection: conn,
		}
		c.connectedClients[client.id] = &client
		go c.serveClient(&client)
		c.mtx.Unlock()
	}
	c.listener.Close()
	logger.Println("UnitSystemNamedPipeServer serve end")
}

func (c *UnitSystemNamedPipeServer) removeClient(id string) {
	c.mtx.Lock()
	if client, ok := c.connectedClients[id]; ok {
		client.connection.Close()
		delete(c.connectedClients, id)
	}
	c.mtx.Unlock()
}

func (c *UnitSystemNamedPipeServer) removeAllClients() {
	c.mtx.Lock()
	for _, client := range c.connectedClients {
		client.connection.Close()
	}
	c.connectedClients = make(map[string]*UnitSystemNamedPipeServerConnectedClient)
	c.mtx.Unlock()
}

func (c *UnitSystemNamedPipeServer) serveClient(client *UnitSystemNamedPipeServerConnectedClient) {
	inputBuffer := make([]byte, 0)
	currentOffset := 0
	for {
		buffer := make([]byte, 1024)
		n, err := client.connection.Read(buffer[currentOffset:])
		if err != nil {
			break
		}

		if n > 0 {
			c.mtx.Lock()
			inputBuffer = append(inputBuffer, buffer[:n]...)

			found := true
			for found {
				found = false
				currentLine := make([]byte, 0)
				for index, b := range inputBuffer {
					if b == 10 || b == 13 {
						// parse currentLine
						if len(currentLine) > 0 {
							parts := strings.Split(string(currentLine), "=")
							if len(parts) > 1 {
								if len(parts[0]) > 0 {
									key := parts[0]
									value := parts[1]

									finalValue := value
									c.receivedVariables[key] = finalValue
									c.SetString(key, finalValue, "")
								}
							}

						}
						inputBuffer = inputBuffer[index+1:]
						found = true
						break
					} else {
						if b >= 32 && b < 128 {
							currentLine = append(currentLine, b)
						}
					}
				}
			}
			c.mtx.Unlock()
		}
	}

	c.removeClient(client.id)
}
