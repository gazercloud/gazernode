package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"
	"time"
)

type Client struct {
	window   uiinterfaces.Window
	received []*Call
	mtx      sync.Mutex
	tm       *uievents.FormTimer
	client   *http.Client
	watcher  *ItemsWatcher

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

func New(window uiinterfaces.Window, address string, userName string, password string) *Client {
	var c Client
	c.address = address
	c.userName = userName
	c.password = password
	c.initClient(window)
	return &c
}

func NewWithSessionToken(window uiinterfaces.Window, address string, userName string, sessionToken string) *Client {
	var c Client
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
	c.client = &http.Client{Transport: tr, Jar: jar}
	c.client.Timeout = 5 * time.Second
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
	if !strings.Contains(addr, ":") {
		addr += ":8084"
	}

	response, err := c.client.Post("http://"+addr+"/api/request", writer.FormDataContentType(), &body)
	if err != nil {
		call.err = errors.New("no connection to " + c.address)
		logger.Println(err)
	} else {
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
	}

	//client.CloseIdleConnections()

	call.client.mtx.Lock()
	call.client.received = append(call.client.received, call)
	call.client.mtx.Unlock()
}
