package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/dingdong-grabber/pkg/constants"
	"github.com/dingdong-grabber/pkg/sign"
	"k8s.io/klog"
)

type Client struct {
	url string
	mtx sync.RWMutex
}

type Response struct {
	Success *bool       `json:"success"`
	Error   *string     `json:"error"`
	Code    interface{} `json:"code"` // int or string
	Message string      `json:"message"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func NewClient(url string) *Client {
	return &Client{
		url: url,
	}
}
func (c *Client) Url() string {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.url
}

func (c *Client) SetUrl(url string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.url = url
}

func (c *Client) request(method string, headers map[string]string, params url.Values) (*Response, error) {
	req, err := http.NewRequest(method, c.Url(), nil)
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

	if resp.StatusCode == http.StatusRequestURITooLong {
		klog.Info(req.URL.RawQuery)
		return nil, errors.New("请求参数过长")
	}

	if err = json.NewDecoder(resp.Body).Decode(&rsp); err != nil {
		klog.Error(err)
		return nil, err
	}
	return checkSuccess(rsp)
}

func (c *Client) requestForm(method string, headers map[string]string, params url.Values) (*Response, error) {
	req, err := http.NewRequest(method, c.Url(), strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	return decode(req, headers)
}

// checkSuccess 检查返回结果是否出现错误
func checkSuccess(resp *Response) (*Response, error) {
	// 响应结果存在success字段, 且为true，则请求成功
	if resp.Success != nil && *resp.Success {
		return resp, nil
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		klog.Error(err)
	}
	klog.Infof("请求结果有异常, 详情: %s", string(bytes))

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

func (c *Client) Get(header map[string]string, params url.Values) (*Response, error) {
	if err := c.Sign(params); err != nil {
		return nil, err
	}
	return c.request(http.MethodGet, header, params)
}

func (c *Client) Post(header map[string]string, params url.Values) (*Response, error) {
	if err := c.Sign(params); err != nil {
		return nil, err
	}
	return c.requestForm(http.MethodPost, header, params)
}

func (c *Client) RawPost(header map[string]string, params url.Values, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, c.Url(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if len(params) > 0 {
		req.URL.RawQuery = params.Encode()
	}
	var client = &http.Client{}
	for k, v := range header {
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
	return resp, nil
}

func (c *Client) SetBody(body url.Values, params map[string]string) {
	for k, v := range params {
		body[k] = []string{v}
	}
}

func (c *Client) Sign(body url.Values) error {
	signs, err := sign.NewDefaultJsSign().Sign(body)
	if err != nil {
		return err
	}
	body[constants.SignNars] = []string{signs[constants.SignNars]}
	body[constants.SignSesi] = []string{signs[constants.SignSesi]}
	return nil
}
