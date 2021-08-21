package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"
	"time"
)

type TransportType string

const TransportTypeLocal TransportType = "local"
const TransportTypeCloudHttps TransportType = "cloud_https"
const TransportTypeCloudBin TransportType = "cloud_bin"

type Client struct {
	window        uiinterfaces.Window
	received      []*Call
	mtx           sync.Mutex
	tm            *uievents.FormTimer
	httpClient    *http.Client
	watcher       *ItemsWatcher
	transportType TransportType

	address      string
	userName     string
	password     string
	sessionToken string

	OnSessionOpen  func()
	OnSessionClose func()
}

type Call struct {
	function   string
	request    []byte
	response   string
	onResponse func(call *Call)
	err        error
	client     *Client
}

func New(window uiinterfaces.Window, address string, userName string, password string, transportType string) *Client {
	var c Client
	c.address = address
	c.userName = userName
	c.password = password

	c.transportType = TransportTypeLocal
	if transportType == "cloud_https" {
		c.transportType = TransportTypeCloudHttps
	}
	if transportType == "cloud_bin" {
		c.transportType = TransportTypeCloudBin
	}

	c.initClient(window)
	return &c
}

func NewWithSessionToken(window uiinterfaces.Window, address string, userName string, sessionToken string, transportType string) *Client {
	var c Client

	c.transportType = TransportTypeLocal
	if transportType == "cloud_https" {
		c.transportType = TransportTypeCloudHttps
	}
	if transportType == "cloud_bin" {
		c.transportType = TransportTypeCloudBin
	}

	c.address = address
	c.userName = userName
	c.sessionToken = sessionToken
	c.initClient(window)
	c.SessionActivate(sessionToken, nil)
	return &c
}

func (c *Client) initClient(window uiinterfaces.Window) {
	tr := &http.Transport{}
	jar, _ := cookiejar.New(nil)
	c.httpClient = &http.Client{Transport: tr, Jar: jar}
	c.httpClient.Timeout = 5 * time.Second
	c.tm = window.NewTimer(100, c.timer)
	c.tm.StartTimer()
	c.watcher = NewItemsWatcher(c)
}

func (c *Client) timer() {
	c.mtx.Lock()
	for _, call := range c.received {
		call.onResponse(call)
	}
	c.received = make([]*Call, 0)
	c.mtx.Unlock()
}

func (c *Client) GetItemValue(name string) common_interfaces.ItemValue {
	return c.watcher.Get(name)
}

func (c *Client) Address() string {
	return c.address
}

func (c *Client) Transport() string {
	return string(c.transportType)
}

func (c *Client) UserName() string {
	return c.userName
}

func (c *Client) Password() string {
	return c.password
}

func (c *Client) SessionToken() string {
	return c.sessionToken
}

func (c *Client) thCall(call *Call) {
	//logger.Println("CLIENT CALL:", call.function)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	{
		fw, _ := writer.CreateFormField("fn")
		fw.Write([]byte(call.function))
	}
	{
		fw, _ := writer.CreateFormField("rj")
		if call.request == nil {
			fw.Write(make([]byte, 0))
		} else {
			fw.Write(call.request)
		}

	}
	writer.Close()

	AddStatSent(body.Len())

	addr := c.address

	if call.function == "session_open" {
		logger.Println("!!!!!!!!!!!!!!!!REMOTE CLIENT:", "Session_open")
	}

	if c.transportType == TransportTypeLocal {
		if !strings.Contains(addr, ":") {
			addr += ":8084"
		}

		response, err := c.Post("http://"+addr+"/api/request", writer.FormDataContentType(), &body, "https://"+addr)

		if err != nil {
			call.err = errors.New("no connection to " + c.address)
			logger.Println(err)
		} else {
			if call.function == "session_open" {
				logger.Println("!!!!!!!!!!!!!!!!REMOTE CLIENT:", "response ok")
			}
			var content []byte
			content, err = ioutil.ReadAll(response.Body)
			if err == nil {
				call.response = strings.TrimSpace(string(content))
				AddStatReceived(len(call.response))
				response.Body.Close()
			} else {
				call.err = err
			}

			type ErrorContainer struct {
				Error string `json:"error"`
			}
			var errCont ErrorContainer
			json.Unmarshal([]byte(call.response), &errCont)
			if len(errCont.Error) > 0 {
				call.err = errors.New(errCont.Error)
			}

			if call.function == "session_open" {
				logger.Println("!!!!!!!!!!!!!!!!REMOTE CLIENT:", "response ok", call.response)
			}
		}
	}

	if c.transportType == TransportTypeCloudHttps {
		response, err := c.Post("https://rep02.gazer.cloud/api/request", writer.FormDataContentType(), &body, "https://"+addr+"-n.gazer.cloud")

		if err != nil {
			call.err = errors.New("no connection to " + c.address)
			logger.Println(err)
		} else {
			if call.function == "session_open" {
				logger.Println("!!!!!!!!!!!!!!!!REMOTE CLIENT:", "response ok")
			}
			var content []byte
			content, err = ioutil.ReadAll(response.Body)
			if err == nil {
				call.response = strings.TrimSpace(string(content))
				AddStatReceived(len(call.response))
				response.Body.Close()
			} else {
				call.err = err
			}

			type ErrorContainer struct {
				Error string `json:"error"`
			}
			var errCont ErrorContainer
			json.Unmarshal([]byte(call.response), &errCont)
			if len(errCont.Error) > 0 {
				call.err = errors.New(errCont.Error)
			}

			if call.function == "session_open" {
				logger.Println("!!!!!!!!!!!!!!!!REMOTE CLIENT:", "response ok", call.response)
			}
		}
	}
	//logger.Println("!!!!!!!!!!!!!!!!REMOTE CLIENT:", "post", c.sessionToken)
	//logger.Println("!!!!!!!!!!!!!!!!REMOTE CLIENT:", "response")

	//client.CloseIdleConnections()

	call.client.mtx.Lock()
	call.client.received = append(call.client.received, call)
	call.client.mtx.Unlock()
}

func (c *Client) Post(url, contentType string, body io.Reader, host string) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Origin", host)
	return c.httpClient.Do(req)
}
