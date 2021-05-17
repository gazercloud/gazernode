package cloud

import (
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/settings"
	"io/ioutil"
	"net"
	"sync"
	"time"
)

type Connection struct {
	mtx  sync.Mutex
	addr string
	conn net.Conn

	started   bool
	stopping  bool
	sessionId string
	userName  string
	password  string
	nodeId    string
	requester common_interfaces.Requester

	disconnectedSent bool
}

func NewConnection() *Connection {
	var c Connection
	c.addr = "db05.gazer.cloud:1077"
	c.userName = "qwe"
	c.password = "qwe"
	c.nodeId = "123"
	return &c
}

func (c *Connection) SetRequester(requester common_interfaces.Requester) {
	c.requester = requester
}

func (c *Connection) Connected() bool {
	return c.conn != nil
}

func (c *Connection) Start() {
	if c.started {
		return
	}
	c.started = true
	c.stopping = false
	c.LoadSession()
	go c.thConn()
}

func (c *Connection) Started() bool {
	return c.started
}

func (c *Connection) Stop() {
	if !c.started {
		return
	}
	c.stopping = true

	if c.conn != nil {
		_ = c.conn.Close()
	}
	for i := 0; i < 100; i++ {
		time.Sleep(10 * time.Millisecond)
		if !c.started {
			break
		}
	}
}

func (c *Connection) LoadSession() error {
	configString, err := ioutil.ReadFile(settings.ServerDataPath() + "/cloud_session.json")
	if err == nil {
		var config SessionConfig
		err = json.Unmarshal(configString, &config)
		if err == nil {
			logger.Println("CloudConnection LoadSession:", config.Key)
			c.sessionId = config.Key
		} else {
			logger.Println("CloudConnection LoadSession unmarshal error:", err)
		}
	} else {
		logger.Println("CloudConnection LoadSession read file error:", err)
	}
	return err
}

func (c *Connection) SaveSession() error {
	var config SessionConfig
	config.Key = c.sessionId
	bs, err := json.MarshalIndent(config, "", " ")
	if err == nil {
		err = ioutil.WriteFile(settings.ServerDataPath()+"/cloud_session.json", bs, 0600)
	}
	return err
}

func (c *Connection) Login(userName string, password string) {
	c.mtx.Lock()
	c.userName = userName
	c.password = password
	c.sessionId = ""
	c.mtx.Unlock()
}

func (c *Connection) Logout() {
	logger.Println("CloudConnection logout")

	// Send LogOut frame to the cloud
	var err error
	var frame BinFrame
	frame.Header.Src = ""
	frame.Header.Dest = ""
	frame.Header.Function = "#logout"
	frame.Header.TransactionId = ""
	frame.Header.SessionId = ""
	frame.Data = nil
	c.SendData(&frame)

	// Clear local data
	c.mtx.Lock()
	c.userName = ""
	c.password = ""
	c.sessionId = ""
	c.mtx.Unlock()

	// Save cleared local data
	err = c.SaveSession()
	if err != nil {
		logger.Println("CloudConnection save session error", err)
	}
}

func (c *Connection) SessionId() string {
	c.mtx.Lock()
	result := c.sessionId
	c.mtx.Unlock()
	return result
}

func (c *Connection) UserName() string {
	return c.userName
}

func (c *Connection) thConn() {
	logger.Println("CloudConnection th begin")
	const inputBufferSize = 100 * 1024
	inputBuffer := make([]byte, inputBufferSize)
	inputBufferOffset := 0
	for !c.stopping {
		if c.conn == nil {
			inputBufferOffset = 0
			if len(c.addr) > 0 {
				var err error
				var conn net.Conn
				conn, err = tls.DialWithDialer(&net.Dialer{Timeout: time.Second * 1}, "tcp", c.addr, &tls.Config{})
				if err != nil {
					c.conn = nil
					logger.Println("CloudConnection th dial error", err)
					time.Sleep(100 * time.Millisecond)
					continue
				}

				logger.Println("Connection connected", c.addr)

				c.mtx.Lock()
				c.conn = conn
				c.mtx.Unlock()

				if c.sessionId == "" {
					c.openSession()
				} else {
					c.regNode()
				}

			} else {
				logger.Println("no addr to connect")
				break
			}
		}

		if inputBufferOffset >= inputBufferSize {
			logger.Println("max buffer size")
			c.applyDisconnected()
			continue
		}

		n, err := c.conn.Read(inputBuffer[inputBufferOffset:])
		if err != nil {
			// connection closed: n = 0; err = EOF
			logger.Println("conn read error:", err)
			c.applyDisconnected()
			continue
		}
		if n == 0 {
			logger.Println("read 0")
			c.applyDisconnected()
			continue
		}

		inputBufferOffset += n

		needExit := false
		processed := 0
		for inputBufferOffset-processed >= 4 {
			frameLen := int(binary.LittleEndian.Uint32(inputBuffer[processed:]))
			if frameLen < 8 || frameLen > inputBufferSize {
				logger.Println("wrong frame len", frameLen)
				needExit = true
				break // critical error
			}
			unprocessedBufferLen := inputBufferOffset - processed
			if unprocessedBufferLen < frameLen {
				break // no enough data
			}

			var frameData BinFrameTask
			frameData.Frame, err = UnmarshalBinFrame(inputBuffer[processed : processed+frameLen])
			if err != nil {
				logger.Println("Error parse frame", err)
			} else {
				frameData.SessionId = c.sessionId
				frameData.Client = c
				c.processData(frameData)
			}

			processed += frameLen
		}

		if needExit {
			c.applyDisconnected()
			continue
		}

		if processed > 0 {
			copy(inputBuffer, inputBuffer[processed:inputBufferOffset])
			inputBufferOffset -= processed
		}
	}

	c.applyDisconnected()

	logger.Println("CloudConnection th end")
	c.started = false
}

func (c *Connection) applyDisconnected() {
	c.mtx.Lock()
	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil
	}
	c.mtx.Unlock()
}

func (c *Connection) SendData(data *BinFrame) {
	c.mtx.Lock()
	conn := c.conn
	if conn != nil {
		data.Header.SessionId = c.sessionId
		frameBytes, _ := data.Marshal()

		fmt.Println("sending cloud ", data.Header.Function)

		sent := 0
		for sent < len(frameBytes) {
			n, err := conn.Write(frameBytes)
			if err != nil {
				break
			}
			sent += n

		}

		fmt.Println("sent cloud ", frameBytes)
	}
	c.mtx.Unlock()
}

func (c *Connection) openSession() {
	var err error
	type OpenSessionRequest struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}

	var frame BinFrame
	var openSessionRequest OpenSessionRequest
	openSessionRequest.UserName = c.userName
	openSessionRequest.Password = c.password

	frame.Header.Src = ""
	frame.Header.Dest = ""
	frame.Header.Function = "#session_open"
	frame.Header.TransactionId = ""
	frame.Header.SessionId = ""
	frame.Data, err = json.Marshal(openSessionRequest)

	if err != nil {
		return
	}
	c.SendData(&frame)
}

func (c *Connection) regNode() {
	var err error
	type RegNodeRequest struct {
		NodeId string `json:"node_id"`
	}

	var frame BinFrame
	var regNodeRequest RegNodeRequest
	regNodeRequest.NodeId = "123"

	frame.Header.Src = ""
	frame.Header.Dest = ""
	frame.Header.Function = "#reg_node"
	frame.Header.TransactionId = ""
	frame.Header.SessionId = ""
	frame.Data, err = json.Marshal(regNodeRequest)

	if err != nil {
		return
	}
	c.SendData(&frame)
}

func (c *Connection) processData(task BinFrameTask) {
	var err error
	logger.Println("processData: ", task.Frame.Data)

	if task.Frame.Header.Function == "#session" {
		logger.Println("CloudConnection #session data received", task.Frame.Data)
		type SessionInfo struct {
			Key string `json:"key"`
		}
		var sessionInfo SessionInfo
		err = json.Unmarshal(task.Frame.Data, &sessionInfo)
		if err == nil {
			logger.Println("CloudConnection #session", sessionInfo.Key)
			c.sessionId = sessionInfo.Key
			c.regNode()
			c.SaveSession()
		}
		return
	}

	if c.requester == nil {
		logger.Println("CloudConnection requester is nil")
		return
	}

	// Frame for the node
	var bs []byte
	bs, err = c.requester.RequestJson(task.Frame.Header.Function, task.Frame.Data, "web")
	if err != nil {
		type ErrorStruct struct {
			Error string `json:"error"`
		}

		var res ErrorStruct
		bs, _ = json.MarshalIndent(res, "", " ")

		var frame BinFrame
		frame.Header.Src = ""
		frame.Header.Dest = ""
		frame.Header.Function = task.Frame.Header.Function
		frame.Header.TransactionId = task.Frame.Header.TransactionId
		frame.Header.SessionId = ""
		frame.Data = bs
		task.Client.SendData(&frame)
		return
	}

	var frame BinFrame
	frame.Header.Src = ""
	frame.Header.Dest = ""
	frame.Header.Function = task.Frame.Header.Function
	frame.Header.TransactionId = task.Frame.Header.TransactionId
	frame.Header.SessionId = ""
	frame.Data = bs

	task.Client.SendData(&frame)
}
