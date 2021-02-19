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
}

type Call struct {
	function   string
	request    []byte
	response   string
	onResponse func(call *Call)
	err        error
	client     *Client
}

func New(window uiinterfaces.Window) *Client {
	var c Client
	pc := &c

	tr := &http.Transport{}
	c.client = &http.Client{Transport: tr}
	c.client.Timeout = 3 * time.Second

	c.tm = window.NewTimer(100, pc.timer)
	c.tm.StartTimer()

	c.watcher = NewItemsWatcher(&c)

	return pc
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

func (c *Client) thCall(call *Call) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	{
		fw, _ := writer.CreateFormField("func")
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

	addr := "127.0.0.1"
	if false {
		addr = "192.168.24.233"
	}

	response, err := c.client.Post("http://"+addr+":8084/api/request", writer.FormDataContentType(), &body)
	if err != nil {
		call.err = errors.New("no connection to local service")
		logger.Println(err)
	} else {
		content, _ := ioutil.ReadAll(response.Body)
		call.response = strings.TrimSpace(string(content))
		AddStatReceived(len(call.response))
		response.Body.Close()

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
