package http

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"k8s.io/klog"
)

type Client struct {
	Url string
	mtx sync.RWMutex
}

type Response struct {
	Success bool        `json:"success"`
	Error   *string     `json:"error"`
	Code    interface{} `json:"code"` // int or string
	Message string      `json:"message"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func (c *Client) request(method string, headers map[string]string, params url.Values) (*Response, error) {
	req, err := http.NewRequest(method, c.Url, nil)
	if err != nil {
		return nil, err
	}
	if len(params) > 0 {
		req.URL.RawQuery = params.Encode()
	}
	return decode(req, headers)
}

func decode(req *http.Request, headers map[string]string) (rsp *Response, err error) {
	var client = &http.Client{}
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	if req.URL.Scheme == "https" {
		client = &http.Client{
			Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	if err = json.NewDecoder(resp.Body).Decode(&rsp); err != nil {
		klog.Error(err)
		return nil, err
	}
	return checkSuccess(rsp)
}

func (c *Client) requestForm(method string, headers map[string]string, params url.Values) (*Response, error) {
	encode := c.ParamsEncode(params)
	req, err := http.NewRequest(method, c.Url, strings.NewReader(encode))
	if err != nil {
		return nil, err
	}
	return decode(req, headers)
}

// ParamsEncode 为了避免多并发造成的并发读写问题: fatal error: concurrent map read and map write
func (c *Client) ParamsEncode(params url.Values) string {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return params.Encode()
}

// checkSuccess 检查返回结果是否出现错误
func checkSuccess(resp *Response) (*Response, error) {
	// 如果存在success并且为true，则请求成功
	if resp.Success {
		return resp, nil
	}
	if resp.Error != nil {
		return nil, errors.New(*resp.Error)
	}
	if resp.Message != "" {
		return nil, errors.New(resp.Message)
	}
	if resp.Msg != "" {
		return nil, errors.New(resp.Msg)
	}
	return nil, fmt.Errorf("%v", resp.Code)
}

func (c *Client) Get(headers map[string]string, params url.Values) (*Response, error) {
	return c.request(http.MethodGet, headers, params)
}

func (c *Client) Post(header map[string]string, params url.Values) (*Response, error) {
	return c.requestForm(http.MethodPost, header, params)
}
