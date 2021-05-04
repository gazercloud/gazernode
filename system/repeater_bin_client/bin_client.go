package repeater_bin_client

import (
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/protocols/users"
	"net"
	"sync"
	"time"
)

type RepeaterBinClient struct {
	mtx              sync.Mutex
	addr             string
	conn             net.Conn
	chProcessingData chan BinFrameTask
	started          bool
	stopping         bool
	lastError        error
	sessionId        string
	userName         string
	password         string
	auth             *users.Users

	disconnectedSent bool
}

func New(addr string, userName string, password string, chProcessingData chan BinFrameTask) *RepeaterBinClient {
	var c RepeaterBinClient
	c.chProcessingData = chProcessingData
	c.addr = addr
	c.userName = userName
	c.password = password

	return &c
}

func (c *RepeaterBinClient) Connected() bool {
	return c.conn != nil
}

func (c *RepeaterBinClient) Start() {
	if c.started {
		return
	}
	c.started = true
	c.stopping = false
	go c.thConn()
	go c.thBackground()
}

func (c *RepeaterBinClient) Started() bool {
	return c.started
}

func (c *RepeaterBinClient) Stop() {
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

func (c *RepeaterBinClient) SessionId() string {
	c.mtx.Lock()
	result := c.sessionId
	c.mtx.Unlock()
	return result
}

func (c *RepeaterBinClient) SetSession(sessionId string, userName string) {
	c.mtx.Lock()
	c.sessionId = sessionId
	c.userName = userName
	c.mtx.Unlock()
}

func (c *RepeaterBinClient) GetRemoteAddr() string {
	c.mtx.Lock()
	conn := c.conn
	c.mtx.Unlock()
	if conn != nil {
		return conn.RemoteAddr().String()
	}
	return "[no addr]"
}

func (c *RepeaterBinClient) UserName() string {
	return c.userName
}

func (c *RepeaterBinClient) LastError() error {
	return c.lastError
}

func (c *RepeaterBinClient) thBackground() {
	for !c.stopping {
		time.Sleep(200 * time.Millisecond)
	}
}

func (c *RepeaterBinClient) thConn() {
	logger.Println("RepeaterBinClient th started")
	const inputBufferSize = 100 * 1024
	inputBuffer := make([]byte, inputBufferSize)
	inputBufferOffset := 0
	for !c.stopping {
		if c.conn == nil {
			c.sessionId = ""
			inputBufferOffset = 0
			if len(c.addr) > 0 {
				var err error
				var conn net.Conn
				conn, err = tls.DialWithDialer(&net.Dialer{Timeout: time.Second * 1}, "tcp", c.addr, &tls.Config{})
				if err != nil {
					c.lastError = err
					c.conn = nil
					logger.Println("RepeaterBinClient th dial error", err)
					time.Sleep(100 * time.Millisecond)
					continue
				}

				logger.Println("RepeaterBinClient connected", c.addr)

				c.lastError = nil
				c.mtx.Lock()
				c.conn = conn
				c.mtx.Unlock()

				c.applyConnected()
				c.regNode()

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
			frameData.SessionId = c.sessionId
			frameData.Client = c
			frameData.Frame, err = UnmarshalBinFrame(inputBuffer[processed : processed+frameLen])
			if err != nil {
				logger.Println("Error parse frame", err)
			} else {
				if c.chProcessingData != nil {
					c.chProcessingData <- frameData
				} else {
					logger.Println("no processor")
				}
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

	logger.Println("client exit ok")
	c.started = false
}

func (c *RepeaterBinClient) applyConnected() {
	c.mtx.Lock()
	var frame BinFrame
	frame.Header.Function = "#connected#"
	frame.Data = make([]byte, 0)

	var frameData BinFrameTask
	frameData.SessionId = c.sessionId
	frameData.Client = c
	frameData.Frame = &frame
	if c.chProcessingData != nil {
		c.chProcessingData <- frameData
	} else {
	}
	c.mtx.Unlock()
}

func (c *RepeaterBinClient) applyDisconnected() {
	c.mtx.Lock()
	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil

		var frame BinFrame
		frame.Header.Function = "#disconnected#"
		frame.Data = make([]byte, 0)

		var frameData BinFrameTask
		frameData.SessionId = c.sessionId
		frameData.Client = c
		frameData.Frame = &frame
		if c.chProcessingData != nil {
			c.chProcessingData <- frameData
		} else {
		}
	}
	if c.auth != nil {
		c.auth.CloseSession(c.sessionId)
	}
	c.mtx.Unlock()
}

func (c *RepeaterBinClient) SendData(data *BinFrame) {
	c.mtx.Lock()
	conn := c.conn
	if conn != nil {
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

func (c *RepeaterBinClient) SendProxyFrame(function string, transactionId string, sessionId string, data []byte) string {
	var frame BinFrame
	frame.Header.Src = ""
	frame.Header.Dest = ""
	frame.Header.Function = function
	frame.Header.TransactionId = transactionId
	frame.Header.SessionId = sessionId
	frame.Data = data
	c.SendData(&frame)
	return transactionId
}

func (c *RepeaterBinClient) openSession() {
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

	if err != nil {
		return
	}
	c.SendData(&frame)
}

func (c *RepeaterBinClient) regNode() {
	var err error
	type RegNodeRequest struct {
		NodeId string `json:"node_id"`
	}

	var frame BinFrame
	var regNodeRequest RegNodeRequest
	regNodeRequest.NodeId = "123"

	frame.Header.Src = ""
	frame.Header.Dest = ""
	frame.Header.Function = "#reg_node#"
	frame.Header.TransactionId = ""
	frame.Header.SessionId = ""
	frame.Data, err = json.Marshal(regNodeRequest)

	if err != nil {
		return
	}
	c.SendData(&frame)
}
