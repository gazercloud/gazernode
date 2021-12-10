package bin_client

import (
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"github.com/gazercloud/gazernode/system/protocols/cloud_structures/protocol"
	users2 "github.com/gazercloud/gazernode/system/protocols/users"
	"github.com/gazercloud/gazernode/utilities"
	"github.com/gazercloud/gazernode/utilities/logger"
	"net"
	"sync"
	"time"
)

type BinClient struct {
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
	incomingChannels map[string]bool
	incomingAll      bool
	auth             *users2.Users
	stat             *utilities.Statistics

	disconnectedSent bool
}

func NewByConn(conn net.Conn, chProcessingData chan BinFrameTask, auth *users2.Users, stat *utilities.Statistics) *BinClient {
	var c BinClient
	c.incomingChannels = make(map[string]bool)
	c.chProcessingData = chProcessingData
	c.incomingAll = true
	c.conn = conn
	c.auth = auth
	c.stat = stat
	c.applyConnected()
	return &c
}

func New(addr string, userName string, password string, chProcessingData chan BinFrameTask) *BinClient {
	var c BinClient
	c.incomingChannels = make(map[string]bool)
	c.chProcessingData = chProcessingData
	c.incomingAll = true
	c.addr = addr
	c.userName = userName
	c.password = password
	c.stat = utilities.NewStatistics()
	return &c
}

func (c *BinClient) SetStat(statistics *utilities.Statistics) {
	c.stat = statistics
}

func (c *BinClient) SetIncomingAll(incomingAll bool) {
	c.incomingAll = incomingAll
}

func (c *BinClient) Connected() bool {
	return c.conn != nil
}

func (c *BinClient) SetNeedChannel(channelId string) {
	found := false

	c.mtx.Lock()
	if _, ok := c.incomingChannels[channelId]; ok {
		found = true
	} else {
		c.incomingChannels[channelId] = false
	}
	c.mtx.Unlock()

	if !found {
		c.SendNeedChannel([]string{channelId})
	}
}

func (c *BinClient) SetNeedChannelOK(channels []string) {
	c.mtx.Lock()
	for _, ch := range channels {
		if _, ok := c.incomingChannels[ch]; ok {
			c.incomingChannels[ch] = true
		}
	}
	c.mtx.Unlock()
}

func (c *BinClient) SetNoChannels(channels []string) {
	c.mtx.Lock()
	for _, ch := range channels {
		if _, ok := c.incomingChannels[ch]; ok {
			c.incomingChannels[ch] = false
		}
	}
	c.mtx.Unlock()
}

func (c *BinClient) SetDontNeedChannel(channelId string) {
	c.mtx.Lock()
	delete(c.incomingChannels, channelId)
	c.mtx.Unlock()
}

func (c *BinClient) Start() {
	if c.started {
		return
	}
	c.started = true
	c.stopping = false
	go c.thConn()
	go c.thBackground()
}

func (c *BinClient) Started() bool {
	return c.started
}

func (c *BinClient) Stop() {
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

func (c *BinClient) SessionId() string {
	c.mtx.Lock()
	result := c.sessionId
	c.mtx.Unlock()
	return result
}

func (c *BinClient) SetSession(sessionId string, userName string) {
	c.mtx.Lock()
	c.sessionId = sessionId
	c.userName = userName
	c.mtx.Unlock()
}

func (c *BinClient) GetRemoteAddr() string {
	c.mtx.Lock()
	conn := c.conn
	c.mtx.Unlock()
	if conn != nil {
		return conn.RemoteAddr().String()
	}
	return "[no addr]"
}

func (c *BinClient) UserName() string {
	return c.userName
}

func (c *BinClient) LastError() error {
	return c.lastError
}

func (c *BinClient) thBackground() {
	for !c.stopping {
		time.Sleep(200 * time.Millisecond)
		if c.sessionId != "" {
			channelsToNeedSend := make([]string, 0)
			c.mtx.Lock()
			for ch, sentNeed := range c.incomingChannels {
				if !sentNeed {
					channelsToNeedSend = append(channelsToNeedSend, ch)
				}
			}
			c.mtx.Unlock()

			if len(channelsToNeedSend) > 0 {
				c.SendNeedChannel(channelsToNeedSend)
			}
		}
	}
}

func (c *BinClient) thConn() {
	logger.Println("binClient th started")
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
				//conn, err = tls.DialWithDialer(&net.Dialer{Timeout: time.Second * 1}, "tcp", c.addr, &tls.Config{})
				conn, err = tls.DialWithDialer(&net.Dialer{Timeout: time.Second * 1}, "tcp", c.addr, &tls.Config{InsecureSkipVerify: true})

				if err != nil {
					c.lastError = err
					c.conn = nil
					logger.Println("binClient th dial error", err)
					time.Sleep(100 * time.Millisecond)
					continue
				}
				c.lastError = nil
				c.mtx.Lock()
				c.conn = conn

				channels := make([]string, 0)
				for ch, _ := range c.incomingChannels {
					channels = append(channels, ch)
				}
				for _, ch := range channels {
					c.incomingChannels[ch] = false
				}
				c.mtx.Unlock()

				c.applyConnected()
				c.openSession()

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

		c.stat.Add("rcv", n)

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
				c.mtx.Lock()
				needToPass := false
				if c.incomingAll {
					needToPass = true
				} else {
				}
				c.mtx.Unlock()

				// Special
				if len(frameData.Frame.Channel) > 0 {
					if frameData.Frame.Channel[0] == '#' {
						c.processSpecial(frameData)
						needToPass = true
					} else {
						if _, ok := c.incomingChannels[frameData.Frame.Channel]; ok {
							needToPass = true
						}
					}
				}

				if needToPass {
					if c.chProcessingData != nil {
						c.chProcessingData <- frameData
					} else {
						logger.Println("no processor")
					}
				} else {
					c.SendDontNeedChannel([]string{frameData.Frame.Channel})
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

func (c *BinClient) processSpecial(frameData BinFrameTask) {
	if frameData.Frame.Channel == "#needChannelsOK#" {
		channels := make([]string, 0)
		err := json.Unmarshal(frameData.Frame.Data, &channels)
		if err != nil {
			logger.Println("System #needChannelsOK# error", err)
			return
		}
		logger.Println("System #needChannelsOK#", channels)
		c.SetNeedChannelOK(channels)
		return
	}

	if frameData.Frame.Channel == "#openSession#" {
		if c.auth != nil {
			userAndPassword := make([]string, 0)
			err := json.Unmarshal(frameData.Frame.Data, &userAndPassword)
			if err != nil {
				logger.Println("System #openSession# error", err)
				return
			}
			logger.Println("System #openSession#", userAndPassword)
			if len(userAndPassword) != 2 {
				logger.Println("System #openSession# error != 2")
				return
			}

			var session *users2.Session
			session, err = c.auth.OpenSession(userAndPassword[0], userAndPassword[1])
			if err == nil {
				logger.Println("session opened: ", userAndPassword[0], ":", session.Id())
				c.SetSession(session.Id(), userAndPassword[0])
				c.SendOpenedSession(session.Id())
			} else {
				logger.Println("can not create session:", err)
			}
		}

		return
	}

	if frameData.Frame.Channel == "#openedSession#" {
		sessionId := make([]string, 0)
		err := json.Unmarshal(frameData.Frame.Data, &sessionId)
		if err != nil {
			logger.Println("System #openedSession# error", err)
			return
		}
		if len(sessionId) != 1 {
			logger.Println("System #openedSession# error != 1")
			return
		}
		c.SetSession(sessionId[0], frameData.Client.UserName())
		logger.Println("System #openedSession#")
		return
	}

}

func (c *BinClient) applyConnected() {
	c.mtx.Lock()
	var frame BinFrame
	frame.Channel = "#connected#"
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

func (c *BinClient) applyDisconnected() {
	c.mtx.Lock()
	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil

		var frame BinFrame
		frame.Channel = "#disconnected#"
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

func (c *BinClient) SendData(data *BinFrame) {
	c.mtx.Lock()
	conn := c.conn
	if conn != nil {
		frameBytes, _ := data.Marshal()

		sent := 0
		for sent < len(frameBytes) {
			n, err := conn.Write(frameBytes)
			if err != nil {
				break
			}
			sent += n
			c.stat.Add("snd", n)
		}
	}
	c.mtx.Unlock()
}

func (c *BinClient) SendNeedChannel(channels []string) {
	var err error
	var data BinFrame
	data.Channel = "#needChannels#"
	data.Password = ""
	data.Data, err = json.Marshal(channels)
	if err != nil {
		return
	}
	c.SendData(&data)
}

func (c *BinClient) SendCountOfSubscribers(channel string, count int) {
	var err error

	var resp protocol.ResponseCountOfSubscribers
	resp.Channel = channel
	resp.Count = count

	var data BinFrame
	data.Channel = "#count_of_subscribers#"
	data.Password = ""
	data.Data, err = json.Marshal(&resp)
	if err != nil {
		return
	}
	c.SendData(&data)
}

func (c *BinClient) Stat() *utilities.Statistics {
	return c.stat
}

func (c *BinClient) SendNeedChannels() {
	channels := make([]string, 0)
	c.mtx.Lock()
	for a := range c.incomingChannels {
		channels = append(channels, a)
	}
	c.mtx.Unlock()

	var err error
	var data BinFrame
	data.Channel = "#needChannels#"
	data.Password = ""
	data.Data, err = json.Marshal(channels)
	if err != nil {
		return
	}
	c.SendData(&data)
}

func (c *BinClient) SendDontNeedChannel(channels []string) {
	var err error
	var data BinFrame
	data.Channel = "#dontNeedChannels#"
	data.Password = ""
	data.Data, err = json.Marshal(channels)
	if err != nil {
		return
	}
	c.SendData(&data)
}

func (c *BinClient) openSession() {
	var err error
	var data BinFrame
	userAndPassword := make([]string, 2)
	userAndPassword[0] = c.userName
	userAndPassword[1] = c.password
	data.Channel = "#openSession#"
	data.Password = ""
	data.Data, err = json.Marshal(userAndPassword)
	if err != nil {
		return
	}
	c.SendData(&data)
}

func (c *BinClient) SendOpenedSession(sessionId string) {
	var err error
	var data BinFrame
	sessionIDs := make([]string, 1)
	sessionIDs[0] = sessionId
	data.Channel = "#openedSession#"
	data.Password = ""
	data.Data, err = json.Marshal(sessionIDs)
	if err != nil {
		return
	}
	c.SendData(&data)
}

func (c *BinClient) SendNeedChannelOK(channels []string) {
	var err error
	var data BinFrame
	data.Channel = "#needChannelsOK#"
	data.Password = ""
	data.Data, err = json.Marshal(channels)
	if err != nil {
		return
	}
	c.SendData(&data)
}

func (c *BinClient) SendNoChannels(channels []string) {
	var err error
	var data BinFrame
	data.Channel = "#noChannels#"
	data.Password = ""
	data.Data, err = json.Marshal(channels)
	if err != nil {
		return
	}
	c.SendData(&data)
}
