package gazer_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/system/cloud"
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

type GazerNodeClient struct {
	received   []*Call
	mtx        sync.Mutex
	httpClient *http.Client

	binCloudConnection *cloud.Connection

	repeater            string
	needSessionActivate bool

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
	client     *GazerNodeClient
}

func New(addr string) *GazerNodeClient {
	var c GazerNodeClient
	c.address = addr
	c.initClient()
	return &c
}

func NewWithSession(addr string, sessionKey string) *GazerNodeClient {
	var c GazerNodeClient
	c.address = addr
	c.sessionToken = sessionKey
	c.initClient()
	return &c
}

func (c *GazerNodeClient) initClient() {
	tr := &http.Transport{}
	jar, _ := cookiejar.New(nil)
	c.httpClient = &http.Client{Transport: tr, Jar: jar}
	c.httpClient.Timeout = 5 * time.Second
}

func (c *GazerNodeClient) Address() string {
	return c.address
}

func (c *GazerNodeClient) UserName() string {
	return c.userName
}

func (c *GazerNodeClient) Password() string {
	return c.password
}

func (c *GazerNodeClient) SessionToken() string {
	return c.sessionToken
}

func (c *GazerNodeClient) thCall(call *Call) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	{
		fw, _ := writer.CreateFormField("fn")
		fw.Write([]byte(call.function))
	}
	{
		fw, _ := writer.CreateFormField("s")
		fw.Write([]byte(c.sessionToken))
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

	addr := c.address

	response, err := c.Post(addr+"/api/request", writer.FormDataContentType(), &body, "https://"+addr)

	if err != nil {
		call.err = errors.New("no connection to " + c.address)
	} else {
		var content []byte
		content, err = ioutil.ReadAll(response.Body)
		if err == nil {
			call.response = strings.TrimSpace(string(content))
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

	call.client.mtx.Lock()
	call.client.received = append(call.client.received, call)
	call.client.mtx.Unlock()
}

func (c *GazerNodeClient) updateRepeater() {

	type SWhereNodeRequest struct {
		NodeId string `json:"node_id"`
	}
	var sWhereNodeRequest SWhereNodeRequest
	sWhereNodeRequest.NodeId = c.address
	request, _ := json.Marshal(sWhereNodeRequest)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	{
		fw, _ := writer.CreateFormField("fn")
		fw.Write([]byte("s-where-node"))
	}
	{
		fw, _ := writer.CreateFormField("rj")
		fw.Write(request)

	}
	writer.Close()

	response, err := c.Post("https://home.gazer.cloud/api/request", writer.FormDataContentType(), &body, "https://home.gazer.cloud")

	if err == nil {
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
		}
		c.repeater = sWhereNodeResponse.Host
	}

}

func (c *GazerNodeClient) Post(url, contentType string, body io.Reader, host string) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Origin", host)
	return c.httpClient.Do(req)
}
