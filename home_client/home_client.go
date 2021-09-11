package home_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type HomeClient struct {
	userName              string
	password              string
	sessionKey            string
	sessionKeyForActivate string
	httpClient            *http.Client
}

func New() *HomeClient {
	var c HomeClient
	return &c
}

func NewWithSession(sessionKey string) *HomeClient {
	var c HomeClient
	c.sessionKeyForActivate = sessionKey

	tr := &http.Transport{}
	jar, _ := cookiejar.New(nil)
	c.httpClient = &http.Client{Transport: tr, Jar: jar}
	c.httpClient.Timeout = 5 * time.Second

	return &c
}

func (c *HomeClient) SessionActivate() (string, error) {
	var bs []byte
	var err error
	type SessionActivateRequest struct {
		SessionToken string `json:"session_token"`
	}
	var req SessionActivateRequest
	req.SessionToken = c.sessionKeyForActivate
	bs, err = json.Marshal(req)
	if err != nil {
		return "", err
	}
	bs, err = c.Call("session_activate", bs)

	type SessionActivateResponse struct {
		SessionToken string `json:"session_token"`
	}
	var resp SessionActivateResponse
	_ = json.Unmarshal(bs, &resp)

	return resp.SessionToken, err
}

func (c *HomeClient) Call(function string, request []byte) (responseBytes []byte, err error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	{
		fw, _ := writer.CreateFormField("fn")
		_, err = fw.Write([]byte(function))
		if err != nil {
			return
		}
	}
	{
		fw, _ := writer.CreateFormField("rj")
		if request == nil {
			_, err = fw.Write(make([]byte, 0))
			if err != nil {
				return
			}
		} else {
			_, err = fw.Write(request)
			if err != nil {
				return
			}
		}

	}
	err = writer.Close()
	if err != nil {
		return
	}

	var response *http.Response
	response, err = c.Post("https://home.gazer.cloud/api/request", writer.FormDataContentType(), &body, "https://home.gazer.cloud")
	if err == nil {
		responseBytes, err = ioutil.ReadAll(response.Body)
		_ = response.Body.Close()
		if err != nil {
			return
		}
		type ErrorContainer struct {
			Error string `json:"error"`
		}
		var errCont ErrorContainer
		err = json.Unmarshal(responseBytes, &errCont)
		if len(errCont.Error) > 0 {
			err = errors.New(errCont.Error)
		}
	}

	return
}

func (c *HomeClient) Post(url, contentType string, body io.Reader, host string) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Origin", host)
	return c.httpClient.Do(req)
}
