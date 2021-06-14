package cloud

import (
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazernode/settings"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"
)

type Connection struct {
	mtx  sync.Mutex
	addr string
	conn net.Conn

	httpClient *http.Client

	started   bool
	stopping  bool
	sessionId string
	userName  string
	password  string
	nodeCode  string
	requester common_interfaces.Requester

	connectionStatus string
	loginStatus      string

	disconnectedSent bool

	calls map[string]int64
}

func NewConnection() *Connection {
	var c Connection
	c.addr = ""
	c.userName = ""
	c.password = ""
	c.nodeCode = "1"

	c.calls = make(map[string]int64)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{},
	}

	c.httpClient = &http.Client{Transport: tr, Timeout: 1 * time.Second}

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
			c.userName = config.UserName
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
	config.UserName = c.userName
	bs, err := json.MarshalIndent(config, "", " ")
	if err == nil {
		err = ioutil.WriteFile(settings.ServerDataPath()+"/cloud_session.json", bs, 0600)
	}
	return err
}

func (c *Connection) Login(userName string, password string) {
	logger.Println("CloudConnection login", userName, password)
	c.mtx.Lock()
	c.userName = userName
	c.password = password
	c.sessionId = ""
	c.mtx.Unlock()
	c.openSession()
}

func (c *Connection) Logout() {
	logger.Println("CloudConnection logout")
	c.loginStatus = "logged out"

	var err error
	type CloseSessionRequest struct {
		Key string `json:"key"`
	}

	var frame BinFrame
	var closeSessionRequest CloseSessionRequest
	closeSessionRequest.Key = c.sessionId

	// Send LogOut frame to the cloud
	frame.Header.Src = ""
	frame.Header.Dest = ""
	frame.Header.Function = "session_close"
	frame.Header.TransactionId = ""
	frame.Header.SessionId = ""
	frame.Data, err = json.Marshal(closeSessionRequest)
	c.SendData(&frame)

	// Clear local data
	c.mtx.Lock()
	c.password = ""
	c.sessionId = ""
	c.mtx.Unlock()

	// Save cleared local data
	err = c.SaveSession()
	if err != nil {
		logger.Println("CloudConnection save session error", err)
		return
	}
	logger.Println("CloudConnection logout ok")
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

func (c *Connection) updateCurrentRepeater() {
	logger.Println("updateCurrentRepeater")
	c.connectionStatus = "repeater search"

	IPs, err := net.LookupIP("home.gazer.cloud")
	if err != nil {
		logger.Println("updateCurrentRepeater error (LookupIP)", err)
		return
	}
	logger.Println("updateCurrentRepeater IPs:", IPs)

	//currentIndex := rand.Int() % len(IPs)
	//currentIP := IPs[currentIndex].String()

	req, _ := http.NewRequest("GET", "https://home.gazer.cloud/api/request?fn=s-repeater-for-node", nil)

	response, err := c.httpClient.Transport.RoundTrip(req)

	//link := "https://home.gazer.cloud/api/request?fn=s-repeater-for-node"
	//response, err := c.httpClient.Get(link)

	if err != nil {
		c.connectionStatus = "repeater search error: " + err.Error()
		logger.Println("updateCurrentRepeater error (httpClient)", err)
		return
	}

	type RepeaterForNodeResponse struct {
		NodeId string `json:"node_id"`
		Host   string `json:"host"`
	}

	content, _ := ioutil.ReadAll(response.Body)
	_ = response.Body.Close()

	logger.Println("Content:", string(content))

	var resp RepeaterForNodeResponse
	err = json.Unmarshal(content, &resp)
	if err != nil {
		c.connectionStatus = "repeater search error: " + err.Error()
		logger.Println("updateCurrentRepeater error (Unmarshal)", err)
		return
	}

	c.addr = resp.Host + ":1077"
	c.connectionStatus = "repeater search complete: " + c.addr

	logger.Println("updateCurrentRepeater ok:", resp.Host)
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
				c.connectionStatus = "connecting"
				conn, err = tls.DialWithDialer(&net.Dialer{Timeout: time.Second * 1}, "tcp", c.addr, &tls.Config{})
				if err != nil {
					c.conn = nil
					c.connectionStatus = "connect error:" + err.Error()
					logger.Println("CloudConnection th dial error", err)
					time.Sleep(100 * time.Millisecond)
					c.addr = ""
					continue
				}

				c.connectionStatus = "connected"

				logger.Println("Connection connected", c.addr)

				c.mtx.Lock()
				c.conn = conn
				c.mtx.Unlock()

				if c.sessionId == "" {
					c.openSession()
				} else {
					c.loginStatus = "ok"
					c.regNode()
				}

			} else {
				c.updateCurrentRepeater()
				time.Sleep(100 * time.Millisecond)
				continue
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

		//fmt.Println("sending cloud ", data.Header.Function)

		sent := 0
		for sent < len(frameBytes) {
			n, err := conn.Write(frameBytes)
			if err != nil {
				break
			}
			sent += n

		}

		//fmt.Println("sent cloud ", frameBytes)
	}
	c.mtx.Unlock()
}

func (c *Connection) openSession() {
	if c.password == "" {
		c.loginStatus = "error: no password provided"
		logger.Println("CloudConnection openSession - no auth data")
		return
	}

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
	frame.Header.Function = "session_open"
	frame.Header.TransactionId = ""
	frame.Header.SessionId = ""
	frame.Data, err = json.Marshal(openSessionRequest)

	c.loginStatus = "processing"

	if err != nil {
		c.loginStatus = "error: " + err.Error()
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
	regNodeRequest.NodeId = c.nodeCode

	frame.Header.Src = ""
	frame.Header.Dest = ""
	frame.Header.Function = "#iam"
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

	if task.Frame != nil {
		c.mtx.Lock()
		if value, ok := c.calls[task.Frame.Header.Function]; ok {
			c.calls[task.Frame.Header.Function] = value + 1
		} else {
			c.calls[task.Frame.Header.Function] = 1
		}
		c.mtx.Unlock()
	}

	if task.Frame.Header.Function == "#iam" {
		logger.Println("#iam response received")
		return
	}

	if task.Frame.Header.Function == "session_open" {
		logger.Println("CloudConnection session_open data received", task.Frame.Data)
		type SessionInfo struct {
			SessionToken string `json:"session_token"`
			Error        string `json:"error"`
		}
		var sessionInfo SessionInfo
		err = json.Unmarshal(task.Frame.Data, &sessionInfo)
		if err == nil {
			if sessionInfo.Error == "" {
				logger.Println("CloudConnection session_open", sessionInfo.SessionToken)
				c.sessionId = sessionInfo.SessionToken
				c.regNode()
				c.SaveSession()
				c.loginStatus = "ok"
			} else {
				c.loginStatus = "error: " + sessionInfo.Error
			}
		} else {
			c.loginStatus = "error: " + err.Error()
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

func (c *Connection) State() (nodeinterface.CloudStateResponse, error) {
	var resp nodeinterface.CloudStateResponse
	resp.Connected = c.conn != nil
	resp.LoggedIn = c.sessionId != ""
	resp.UserName = c.userName
	resp.LoginStatus = c.loginStatus
	resp.Status = c.connectionStatus
	resp.CurrentRepeater = c.addr
	resp.NodeId = c.nodeCode

	c.mtx.Lock()
	resp.Counters = make([]nodeinterface.CloudStateResponseItem, 0)
	for key, value := range c.calls {
		resp.Counters = append(resp.Counters, nodeinterface.CloudStateResponseItem{
			Name:  key,
			Value: value,
		})
	}
	c.mtx.Unlock()

	return resp, nil
}

func (c *Connection) Nodes() (nodeinterface.CloudNodesResponse, error) {
	var resp nodeinterface.CloudNodesResponse
	return resp, nil
}

func (c *Connection) AddNode() (nodeinterface.CloudAddNodeResponse, error) {
	var resp nodeinterface.CloudAddNodeResponse
	return resp, nil
}

func (c *Connection) UpdateNode() (nodeinterface.CloudUpdateNodeResponse, error) {
	var resp nodeinterface.CloudUpdateNodeResponse
	return resp, nil
}

func (c *Connection) RemoveNode() (nodeinterface.CloudRemoveNodeResponse, error) {
	var resp nodeinterface.CloudRemoveNodeResponse
	return resp, nil
}

func (c *Connection) GetSettings() (nodeinterface.CloudGetSettingsResponse, error) {
	var resp nodeinterface.CloudGetSettingsResponse
	return resp, nil
}

func (c *Connection) SetSettings() (nodeinterface.CloudSetSettingsResponse, error) {
	var resp nodeinterface.CloudSetSettingsResponse
	return resp, nil
}
