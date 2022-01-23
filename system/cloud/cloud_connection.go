package cloud

import (
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	nodeinterface2 "github.com/gazercloud/gazernode/system/protocols/nodeinterface"
	"github.com/gazercloud/gazernode/utilities/logger"
	"github.com/gazercloud/gazernode/utilities/packer"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Connection struct {
	mtx sync.Mutex
	//ss   *settings.Settings
	addr string
	conn net.Conn

	httpClient *http.Client

	started   bool
	stopping  bool
	sessionId string
	userName  string
	password  string
	nodeId    string
	requester common_interfaces.Requester

	connectionStatus string
	loginStatus      string
	iamStatus        string

	disconnectedSent bool
	targetRepeater   string

	callsSuccess   map[string]int64
	callsUnSuccess map[string]int64

	callPerSecond          float64
	receivedBytesPerSecond float64
	sentBytesPerSecond     float64

	callCount     int64
	receivedBytes int64
	sentBytes     int64

	lastReceivedBytes int64
	lastSentBytes     int64
	lastCallCount     int64

	lastCallCountDT time.Time

	proxyTasks map[string]*ProxyTask

	allowIncomingFunctions map[string]bool

	serverDataPath string
}

func NewConnection(serverDataPath string) *Connection {
	var c Connection
	//c.ss = ss
	c.serverDataPath = serverDataPath
	c.addr = ""
	c.userName = ""
	c.password = ""
	c.nodeId = ""

	c.allowIncomingFunctions = make(map[string]bool)

	for _, f := range nodeinterface2.ApiFunctions() {
		c.allowIncomingFunctions[f] = false
	}

	c.callsSuccess = make(map[string]int64)
	c.callsUnSuccess = make(map[string]int64)
	c.proxyTasks = make(map[string]*ProxyTask)

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

func (c *Connection) CloseConnection() {
	c.mtx.Lock()
	_ = c.conn.Close()
	c.mtx.Unlock()
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

func (c *Connection) LoadSessionData(sessionKey string, userName string, password string) {
	c.sessionId = sessionKey
	c.userName = userName
	c.password = password
}

func (c *Connection) LoadSession() error {
	if len(c.serverDataPath) < 1 {
		return nil
	}

	configString, err := ioutil.ReadFile(c.serverDataPath + "/cloud_session.json")
	if err == nil {
		var config SessionConfig
		err = json.Unmarshal(configString, &config)
		if err == nil {
			logger.Println("CloudConnection LoadSession:", config.Key)
			c.sessionId = config.Key
			c.userName = config.UserName
			c.nodeId = config.NodeId
			for _, f := range config.AllowIncomingFunctions {
				c.allowIncomingFunctions[f] = true
			}
		} else {
			logger.Println("CloudConnection LoadSession unmarshal error:", err)
		}
	} else {
		logger.Println("CloudConnection LoadSession read file error:", err)
	}
	return err
}

func (c *Connection) SaveSession() error {
	if len(c.serverDataPath) < 1 {
		return nil
	}

	var config SessionConfig
	config.Key = c.sessionId
	config.UserName = c.userName
	config.NodeId = c.nodeId
	config.AllowIncomingFunctions = make([]string, 0)
	for key, value := range c.allowIncomingFunctions {
		if value {
			config.AllowIncomingFunctions = append(config.AllowIncomingFunctions, key)
		}
	}
	bs, err := json.MarshalIndent(config, "", " ")
	if err == nil {
		err = ioutil.WriteFile(c.serverDataPath+"/cloud_session.json", bs, 0600)
	}
	return err
}

func (c *Connection) Login(userName string, password string) {
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
	frame.Header.Function = "session_close"
	frame.Header.TransactionId = ""
	frame.Data, err = json.Marshal(closeSessionRequest)
	c.SendData(&frame, true)

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

func (c *Connection) getHomeHost() string {
	IPs, err := net.LookupIP("home.gazer.cloud")
	if err != nil {
		logger.Println("getHomeHost error (LookupIP)", err)
		return ""
	}
	logger.Println("getHomeHost IPs:", IPs)
	return "home.gazer.cloud" // TODO!!!!!!!!!!!!!!!
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

	type RepeaterForNodeResponseItem struct {
		Host  string  `json:"host"`
		Score float64 `json:"score"`
	}

	type RepeaterForNodeResponse struct {
		NodeId string                        `json:"node_id"`
		Items  []RepeaterForNodeResponseItem `json:"items"`
		Host   string                        `json:"host"`
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

	//c.addr = "rep02.gazer.cloud:1077"

	c.connectionStatus = "repeater search complete: " + c.addr
	logger.Println("updateCurrentRepeater ok:", c.addr)

	/*if len(resp.Items) > 0 {
		c.addr = resp.Items[0].Host + ":1077"
		c.connectionStatus = "repeater search complete: " + c.addr
		logger.Println("updateCurrentRepeater ok:", c.addr)
	} else {
		logger.Println("updateCurrentRepeater no items")
	}*/

}

func (c *Connection) updateRepeaterForNode(nodeId string) {
	if nodeId == "" {
		c.addr = c.getHomeHost() + ":1077"
		return
	}

	type SWhereNodeRequest struct {
		NodeId string `json:"node_id"`
	}
	var sWhereNodeRequest SWhereNodeRequest
	sWhereNodeRequest.NodeId = nodeId
	sWhereNodeRequestBytes, _ := json.Marshal(sWhereNodeRequest)

	logger.Println("updateRepeaterForNode")
	c.connectionStatus = "repeater search"

	IPs, err := net.LookupIP("home.gazer.cloud")
	if err != nil {
		logger.Println("updateRepeaterForNode error (LookupIP)", err)
		return
	}
	logger.Println("updateRepeaterForNode IPs:", IPs)

	req, _ := http.NewRequest("GET", "https://home.gazer.cloud/api/request?fn=s-where-node&rj="+string(sWhereNodeRequestBytes), nil)

	response, err := c.httpClient.Transport.RoundTrip(req)

	if err != nil {
		c.connectionStatus = "repeater search error: " + err.Error()
		logger.Println("updateRepeaterForNode error (httpClient)", err)
		return
	}

	type SWhereNodeResponse struct {
		NodeId string `json:"node_id"`
		Host   string `json:"host"`
	}
	var sWhereNodeResponse SWhereNodeResponse

	var content []byte
	content, err = ioutil.ReadAll(response.Body)
	if err == nil {
		json.Unmarshal(content, &sWhereNodeResponse)
		response.Body.Close()
		c.addr = sWhereNodeResponse.Host + ":1077"
		logger.Println("updateRepeaterForNode ok:", c.addr)
	}

	logger.Println("Content:", string(content))

}

func (c *Connection) thConn() {
	logger.Println("CloudConnection th begin")
	const inputBufferSize = 1024 * 1024
	inputBuffer := make([]byte, inputBufferSize)
	inputBufferOffset := 0
	for !c.stopping {
		if c.conn == nil {
			inputBufferOffset = 0
			if len(c.addr) > 0 {
				var err error
				var conn net.Conn
				c.connectionStatus = "connecting"
				conn, err = tls.DialWithDialer(&net.Dialer{Timeout: time.Second * 1, KeepAlive: 5 * time.Second}, "tcp", c.addr, &tls.Config{})
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
				if len(c.nodeId) > 0 {
					c.updateCurrentRepeater()
				} else {
					c.updateRepeaterForNode(c.targetRepeater)
				}
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
				go c.processData(frameData, int64(frameLen))
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

func (c *Connection) SendData(data *BinFrame, isRequest bool) (frameLen int64) {
	c.mtx.Lock()
	conn := c.conn
	if conn != nil {
		data.Header.CloudSessionId = c.sessionId
		data.Header.IsRequest = isRequest
		frameBytes, _ := data.Marshal()
		frameLen += int64(len(frameBytes))

		sent := 0
		for sent < len(frameBytes) {
			n, err := conn.Write(frameBytes)
			if err != nil {
				break
			}
			sent += n

		}
	}
	c.mtx.Unlock()
	return
}

func (c *Connection) openSession() {
	logger.Println("Cloud Connection openSession", c.userName)

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

	frame.Header.Function = "session_open"
	frame.Header.TransactionId = ""
	frame.Data, err = json.Marshal(openSessionRequest)

	c.loginStatus = "processing"

	if err != nil {
		c.loginStatus = "error: " + err.Error()
		return
	}
	c.SendData(&frame, true)
}

func (c *Connection) regNode() {
	var err error
	type RegNodeRequest struct {
		NodeId string `json:"node_id"`
	}

	if c.nodeId == "" {
		c.iamStatus = "no nodeId specified"
		return
	}

	var frame BinFrame
	var regNodeRequest RegNodeRequest
	regNodeRequest.NodeId = c.nodeId

	c.iamStatus = "processing"

	frame.Header.Function = "#iam"
	frame.Header.TransactionId = ""
	frame.Data, err = json.Marshal(regNodeRequest)

	if err != nil {
		return
	}
	c.SendData(&frame, true)
}

func (c *Connection) processData(task BinFrameTask, inputFrameSize int64) {
	var err error

	if time.Now().Sub(c.lastCallCountDT) > 1*time.Second {
		now := time.Now()
		period := now.Sub(c.lastCallCountDT)
		if c.lastCallCount > 0 {
			c.callPerSecond = float64(c.callCount-c.lastCallCount) / period.Seconds()
			c.receivedBytesPerSecond = float64(c.receivedBytes-c.lastReceivedBytes) / period.Seconds()
			c.sentBytesPerSecond = float64(c.sentBytes-c.lastSentBytes) / period.Seconds()
		} else {
			c.callPerSecond = 0
		}

		c.lastCallCount = c.callCount
		c.lastReceivedBytes = c.receivedBytes
		c.lastSentBytes = c.sentBytes
		c.lastCallCountDT = now
	}

	c.callCount++
	c.receivedBytes += inputFrameSize

	responseFromNodeReceived := false
	c.mtx.Lock()
	if tr, ok := c.proxyTasks[task.Frame.Header.TransactionId]; ok {
		tr.ResponseText = task.Frame.Data
		tr.ResponseReceived = true
		responseFromNodeReceived = true
	}
	c.mtx.Unlock()

	if responseFromNodeReceived {
		return
	}

	if task.Frame.Header.Function == "#iam" {
		err = nil
		if task.Frame.Header.Error != "" {
			err = errors.New(task.Frame.Header.Error)
		}

		if err != nil {
			c.iamStatus = err.Error()
			logger.Println("#iam response received. Error received: ", task.Frame.Header.Error)
		} else {
			c.iamStatus = "ok"
			logger.Println("#iam response received. Everything is OK.")
		}
		return
	}

	if task.Frame.Header.Function == "session_open" {
		logger.Println("Cloud Connection session_open data received", task.Frame.Data)
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

	var allowed bool
	{
		var allowedValue bool
		var allowedFound bool
		if allowedValue, allowedFound = c.allowIncomingFunctions[task.Frame.Header.Function]; allowedFound {
			if allowedValue {
				allowed = true
			}
		}
		allowed = true
	}

	if !task.Frame.Header.IsRequest && task.Frame.Header.Function == "s-account-info" {
		logger.Println("!task.Frame.Header.IsRequest && task.Frame.Header.Function == 123123123", task.Frame.Header.TransactionId, task.Frame.Header.Error, string(task.Frame.Data))
		//allowed = true
	}

	var bs []byte
	if allowed {
		if c.requester == nil {
			logger.Println("CloudConnection requester is nil")
			return
		}

		//logger.Println("!!!!!!!!!! CLOUD DATA:", string(task.Frame.Data))

		if len(task.Frame.Data) < 2 {
			logger.Println("CloudConnection len(task.Frame.Data) < 2")
			return
		}

		request := []byte(packer.UnpackString(string(task.Frame.Data)))

		// Frame for the node
		bs, err = c.requester.RequestJson(task.Frame.Header.Function, request, "web", true)
		//logger.Println("CloudConnection REQUEST", task.Frame.Header.Function, "resLen:", len(bs))

		c.addSuccessCallStat(task.Frame.Header.Function)
	} else {
		c.addUnSuccessCallStat(task.Frame.Header.Function)
		logger.Println("NOT ALLOWED FUNCTION:", task.Frame.Header.Function, task.Frame.Header.TransactionId)
		err = errors.New("access denied")
	}

	if err != nil {
		bs = []byte(err.Error())
		var frame BinFrame
		frame.Header.Function = task.Frame.Header.Function
		frame.Header.TransactionId = task.Frame.Header.TransactionId
		frame.Header.Error = err.Error()
		frame.Data = bs
		c.sentBytes += task.Client.SendData(&frame, false)
		return
	}

	var frame BinFrame
	frame.Header.Function = task.Frame.Header.Function
	frame.Header.TransactionId = task.Frame.Header.TransactionId
	if err != nil {
		frame.Header.Error = err.Error()
	}
	//fmt.Println("Cloud Inbound Call", frame.Header.Function)
	bs = packer.PackBytes(bs)
	frame.Data = bs
	c.sentBytes += task.Client.SendData(&frame, false)
}

func (c *Connection) addSuccessCallStat(f string) {
	c.mtx.Lock()
	if value, ok := c.callsSuccess[f]; ok {
		c.callsSuccess[f] = value + 1
	} else {
		c.callsSuccess[f] = 1
	}
	c.mtx.Unlock()
}

func (c *Connection) addUnSuccessCallStat(f string) {
	c.mtx.Lock()
	if value, ok := c.callsUnSuccess[f]; ok {
		c.callsUnSuccess[f] = value + 1
	} else {
		c.callsUnSuccess[f] = 1
	}
	c.mtx.Unlock()
}

type ProxyTask struct {
	TransactionId string
	Function      string
	RequestText   []byte

	ResponseReceived bool
	ErrorReceived    bool
	ResponseText     []byte
}

func (c *Connection) Call(function string, requestText []byte, targetNodeId string) (response []byte, err error) {
	if len(targetNodeId) > 0 {
		c.targetRepeater = targetNodeId
		if c.targetRepeater != "" && strings.HasPrefix(c.addr, "home.gazer.cloud") {
			c.addr = ""
			c.updateRepeaterForNode(c.targetRepeater)
		}
	}

	// Unique Transaction Id
	transactionId := strconv.FormatInt(rand.Int63(), 16) + strconv.FormatInt(time.Now().UnixNano(), 16)

	// ProxyTask
	var task ProxyTask
	task.Function = function
	task.RequestText = requestText
	task.TransactionId = transactionId
	task.ResponseReceived = false
	c.mtx.Lock()
	c.proxyTasks[transactionId] = &task
	c.mtx.Unlock()

	// Send frame to node
	var frame BinFrame
	frame.Header.TargetNodeId = targetNodeId
	frame.Header.Function = function
	frame.Header.TransactionId = transactionId
	frame.Header.IsRequest = true
	frame.Data = requestText
	c.SendData(&frame, true)

	// Waiting for response
	tBegin := time.Now()
	for time.Now().Sub(tBegin) < 3*time.Second && !task.ResponseReceived {
		time.Sleep(10 * time.Millisecond)
	}

	// Remove task
	c.mtx.Lock()
	delete(c.proxyTasks, transactionId)
	c.mtx.Unlock()

	var resultBytes []byte

	if task.ResponseReceived {
		resultBytes = task.ResponseText
		type ErrorStruct struct {
			Error string `json:"error"`
		}
		var errStr ErrorStruct
		err = json.Unmarshal(task.ResponseText, &errStr)
		if errStr.Error != "" {
			err = errors.New(string(resultBytes))
		}
	} else {
		err = errors.New("node timeout")
	}

	return resultBytes, err
}

func (c *Connection) State() (nodeinterface2.CloudStateResponse, error) {
	var resp nodeinterface2.CloudStateResponse
	resp.Connected = c.conn != nil
	resp.LoggedIn = c.sessionId != ""
	resp.UserName = c.userName
	resp.LoginStatus = c.loginStatus
	resp.ConnectionStatus = c.connectionStatus
	resp.IAmStatus = c.iamStatus
	resp.CurrentRepeater = c.addr
	resp.NodeId = c.nodeId
	resp.SessionKey = c.sessionId

	c.mtx.Lock()
	resp.Counters = make([]nodeinterface2.CloudStateResponseItem, 0)
	for key, value := range c.callsSuccess {

		allow, ok := c.allowIncomingFunctions[key]
		if !ok {
			allow = false
		}

		resp.Counters = append(resp.Counters, nodeinterface2.CloudStateResponseItem{
			Name:  key,
			Allow: allow,
			Value: value,
		})
	}
	c.mtx.Unlock()

	return resp, nil
}

func (c *Connection) Nodes() (resp nodeinterface2.CloudNodesResponse, err error) {
	resp.Nodes = make([]nodeinterface2.CloudNodesResponseItem, 0)

	cloudResp, err := c.Call("s-registered-nodes", nil, "")
	if err != nil {
		logger.Println("Connection::Nodes s-registered-nodes error", err)
		return
	}

	type NodeResponseItem struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	type NodesResponse struct {
		Items []NodeResponseItem `json:"items"`
	}

	var nodesResp NodesResponse
	err = json.Unmarshal(cloudResp, &nodesResp)
	if err != nil {
		return
	}

	for _, n := range nodesResp.Items {
		resp.Nodes = append(resp.Nodes, nodeinterface2.CloudNodesResponseItem{
			NodeId: n.Id,
			Name:   n.Name,
		})
	}

	return resp, nil
}

func (c *Connection) AddNode(name string) (resp nodeinterface2.CloudAddNodeResponse, err error) {
	type NodeAddRequest struct {
		Name string `json:"name"`
	}
	var req NodeAddRequest
	req.Name = name
	var bs []byte
	bs, err = json.Marshal(req)
	if err != nil {
		return
	}

	cloudResp, err := c.Call("s-node-add", bs, "")
	if err != nil {
		logger.Println("!!!!!! s-node-add", err)
		return
	}

	logger.Println("!!!!!! s-node-add resp", string(cloudResp))

	type NodeAddResponse struct {
		NodeId string `json:"id"`
	}
	var r NodeAddResponse
	err = json.Unmarshal(cloudResp, &r)
	if err != nil {
		return
	}

	resp.NodeId = r.NodeId
	return
}

func (c *Connection) UpdateNode(nodeId string, name string) (resp nodeinterface2.CloudUpdateNodeResponse, err error) {
	type NodeUpdateRequest struct {
		NodeId string `json:"node_id"`
		Name   string `json:"name"`
	}
	var req NodeUpdateRequest
	req.NodeId = nodeId
	req.Name = name
	var bs []byte
	bs, err = json.Marshal(req)
	if err != nil {
		return
	}

	logger.Println("!!!!!! s-node-update bin:", string(bs))

	cloudResp, err := c.Call("s-node-update", bs, "")
	if err != nil {
		logger.Println("!!!!!! s-node-update error", err)
		return
	}

	type NodeUpdateResponse struct {
	}
	var r NodeUpdateResponse
	err = json.Unmarshal(cloudResp, &r)
	if err != nil {
		return
	}

	return
}

func (c *Connection) RemoveNode(nodeId string) (resp nodeinterface2.CloudRemoveNodeResponse, err error) {
	type NodeRemoveRequest struct {
		NodeId string `json:"node_id"`
	}
	var req NodeRemoveRequest
	req.NodeId = nodeId
	var bs []byte
	bs, err = json.Marshal(req)
	if err != nil {
		return
	}

	logger.Println("!!!!!! s-node-remove bin:", string(bs))

	cloudResp, err := c.Call("s-node-remove", bs, "")
	if err != nil {
		logger.Println("!!!!!! s-node-remove error", err)
		return
	}

	type NodeRemoveResponse struct {
	}
	var r NodeRemoveResponse
	err = json.Unmarshal(cloudResp, &r)
	if err != nil {
		return
	}

	return
}

func (c *Connection) GetSettings(request nodeinterface2.CloudGetSettingsRequest) (nodeinterface2.CloudGetSettingsResponse, error) {
	var resp nodeinterface2.CloudGetSettingsResponse
	resp.Items = make([]*nodeinterface2.CloudGetSettingsResponseItem, 0)
	for _, function := range nodeinterface2.ApiFunctions() {
		v, ok := c.allowIncomingFunctions[function]
		if !ok {
			v = false
		}

		resp.Items = append(resp.Items, &nodeinterface2.CloudGetSettingsResponseItem{
			Function: function,
			Allow:    v,
		})
	}
	return resp, nil
}

func (c *Connection) GetSettingsProfiles(request nodeinterface2.CloudGetSettingsProfilesRequest) (nodeinterface2.CloudGetSettingsProfilesResponse, error) {
	var resp nodeinterface2.CloudGetSettingsProfilesResponse
	resp.Items = make([]*nodeinterface2.CloudGetSettingsProfilesResponseItem, 0)
	for _, role := range nodeinterface2.ApiRoles() {
		resp.Items = append(resp.Items, &nodeinterface2.CloudGetSettingsProfilesResponseItem{
			Code:      role.Code,
			Name:      role.Name,
			Functions: role.Functions,
		})
	}
	return resp, nil
}

func (c *Connection) SetSettings(request nodeinterface2.CloudSetSettingsRequest) (resp nodeinterface2.CloudSetSettingsResponse, err error) {
	for _, item := range request.Items {
		if _, ok := c.allowIncomingFunctions[item.Function]; ok {
			c.allowIncomingFunctions[item.Function] = item.Allow
		}
	}
	err = c.SaveSession()
	return
}

func (c *Connection) AccountInfo(request nodeinterface2.CloudAccountInfoRequest) (resp nodeinterface2.CloudAccountInfoResponse, err error) {
	cloudResp, err := c.Call("s-account-info", []byte("{}"), "")
	if err != nil {
		logger.Println("s-account-info error", err)
		return
	}

	type AccountInfo struct {
		Email         string `json:"email"`
		MaxNodesCount int64  `json:"max_nodes_count"`
	}

	var accountInfoResp AccountInfo
	err = json.Unmarshal(cloudResp, &accountInfoResp)
	if err != nil {
		return
	}

	resp.Email = accountInfoResp.Email
	resp.MaxNodesCount = accountInfoResp.MaxNodesCount

	return resp, nil
}

func (c *Connection) SetCurrentNodeId(request nodeinterface2.CloudSetCurrentNodeIdRequest) (resp nodeinterface2.CloudSetCurrentNodeIdResponse, err error) {
	if request.NodeId != c.nodeId {
		c.nodeId = request.NodeId
		c.CloseConnection()
		err = c.SaveSession()
	}
	return
}

func (c *Connection) Stat() (res common_interfaces.StatGazerCloud) {
	res.CallsPerSecond = c.callPerSecond
	res.ReceiveSpeed = c.receivedBytesPerSecond
	res.SendSpeed = c.sentBytesPerSecond
	return
}
