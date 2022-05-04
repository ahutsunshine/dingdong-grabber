/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at
  http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package user

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/dingdong-grabber/pkg/config"
	"github.com/dingdong-grabber/pkg/user/ios"

	"github.com/dingdong-grabber/pkg/constants"
	"github.com/dingdong-grabber/pkg/http"
	"k8s.io/klog"
)

type User struct {
	c          *http.Client
	conf       *config.Config
	device     DeviceInterface
	userDetail *UserDetail
	addressId  string
	headers    map[string]string
	params     url.Values
	mtx        sync.RWMutex
}

func NewDefaultUser(conf *config.Config) *User {
	u := &User{
		c:       &http.Client{},
		conf:    conf,
		headers: make(map[string]string),
		params:  url.Values{},
	}
	u.device = getDevice(conf)
	return u
}

func getDevice(conf *config.Config) DeviceInterface {
	var device DeviceInterface
	switch conf.Device {
	case config.IosType:
		device = ios.NewIosDevice()
	case config.AndroidType:
	case config.DefaultType:
		if conf.Cookie == "" {
			klog.Fatal("默认模式下cookie为必填项")
		}
		device = ios.NewIosDevice()
	default:
		klog.Fatalf("暂不支持此设备类型: %s，支持: [ios, android, default]", conf.Device)
	}
	return device
}

func (u *User) AddressId() string {
	u.mtx.RLock()
	defer u.mtx.RUnlock()
	return u.addressId
}

func (u *User) SetAddressId(addressId string) {
	u.mtx.Lock()
	defer u.mtx.Unlock()
	u.addressId = addressId
}

func (u *User) LoadConfig() error {
	var file = "cart.chlsj"
	if u.conf.Device != config.IosType {
		file = "example.chlsj"
	}
	if err := u.device.LoadConfig(file); err != nil {
		klog.Fatal(err)
	}
	// 设置header
	u.SetHeaders(u.device.Headers())
	// 设置 query params
	u.SetQueryParams(u.device.QueryParams())

	if u.conf.Device != config.IosType {
		u.SetHeaders(map[string]string{
			"cookie": u.conf.Cookie,
		})
	}

	ud, err := u.GetUserDetail()
	if err != nil {
		return err
	}

	// 设置header im_secret
	u.SetHeaders(map[string]string{
		"im_secret": ud.UserInfo.ImSecret,
	})

	addr, err := u.GetDefaultAddr()
	if err != nil {
		return err
	}
	// 设置收货地址ID
	u.SetAddressId(addr.Id)

	// 设置Header和Query params收货站ID,城市编码和经纬度
	u.SetHeaders(map[string]string{
		"ddmc-station-id":  addr.StationId,
		"ddmc-city-number": addr.CityNumber,
		"ddmc-longitude":   fmt.Sprintf("%v", addr.Location.Location[0]),
		"ddmc-latitude":    fmt.Sprintf("%v", addr.Location.Location[1]),
	})
	u.SetQueryParams(map[string]string{
		"station_id":  addr.StationId,
		"city_number": addr.CityNumber,
		"longitude":   fmt.Sprintf("%v", addr.Location.Location[0]),
		"latitude":    fmt.Sprintf("%v", addr.Location.Location[1]),
	})

	return nil
}

func (u *User) SetDefaultHeaders(cookie string) {
	u.mtx.Lock()
	defer u.mtx.Unlock()
	if !strings.HasPrefix(cookie, constants.CookiePrefix) {
		cookie = fmt.Sprintf("%s=%s", constants.CookiePrefix, cookie)
	}
	u.headers = map[string]string{
		// Header必填项
		"cookie":           cookie,
		constants.ImSecret: "", // 自动获取

		// 根据cookie动态获取
		"ddmc-uid": "",

		// 设置经纬度, 获取默认地址时会自动添加
		"ddmc-longitude":         "",
		"ddmc-latitude":          "",
		"ddmc-country-code":      "CN",
		"ddmc-locale-identifier": "zh_CN",

		// 下面作为小程序2.85.2版本的默认值
		"ddmc-build-version": "1232",
		"ddmc-sdkversion":    "2.13.2",
		"ddmc-city-number":   "", // 程序会自动获取默认地址的city number填充于此
		"ddmc-station-id":    "", // 程序会自动获取默认地址的station id填充于此
		"ddmc-time":          fmt.Sprintf("%d", time.Now().UnixMilli()/1000),
		"ddmc-channel":       "App Store",
		"ddmc-os-version":    "15.5",
		"ddmc-app-client-id": "1",
		"ddmc-ip":            "",
		"ddmc-language-code": "zh",
		"ddmc-api-version":   "9.50.2",
		"ddmc-device-id":     "",
		"ddmc-device-model":  "",
		"ddmc-device-name":   "",
		"ddmc-device-token":  "",
		"ddmc-idfa":          "",
		"referer":            "https://servicewechat.com/wx1e113254eda17715/435/page-frame.html",
		"content-type":       "application/x-www-form-urlencoded",
		"accept":             "*/*",
		"user-agent":         "user-agent",
		// 不要添加此accept encoding，否则结果会被压缩乱码返回
		//"accept-encoding":    "gzip,compress,br,deflate",
	}
}

// SetHeaders 设置header参数，避免header因多并发引起的concurrent map writes
func (u *User) SetHeaders(headers map[string]string) {
	u.mtx.Lock()
	defer u.mtx.Unlock()
	for k, v := range headers {
		u.headers[k] = v
	}
}

// Headers 返回请求headers的复制
func (u *User) Headers() map[string]string {
	u.mtx.RLock()
	defer u.mtx.RUnlock()
	var cp = make(map[string]string)
	for k, v := range u.headers {
		cp[k] = v
	}
	return cp
}

// SetDefaultParams 设置默认的用户初始化数据
func (u *User) SetDefaultParams() {
	u.mtx.Lock()
	defer u.mtx.Unlock()
	var headers = u.headers
	u.params = url.Values{
		"s_id":         []string{""},
		"device_token": []string{""},

		// 下面作为小程序2.85.2版本的默认值
		"uid":              []string{headers["ddmc-uid"]},
		"longitude":        []string{headers["ddmc-longitude"]},
		"latitude":         []string{headers["ddmc-latitude"]},
		"station_id":       []string{headers["ddmc-station-id"]},
		"city_number":      []string{headers["ddmc-city-number"]},
		"api_version":      []string{headers["ddmc-api-version"]},
		"app_version":      []string{headers["ddmc-build-version"]},
		"time":             []string{headers["ddmc-time"]},
		"openid":           []string{headers["ddmc-device-id"]},
		"buildVersion":     []string{headers["ddmc-device-id"]},
		"countryCode":      []string{headers["ddmc-country-code"]},
		"idfa":             []string{headers["ddmc-idfa"]},
		"ip":               []string{headers["ddmc-ip"]},
		"localeidentifier": []string{headers["ddmc-locale-identifier"]},
		"app_client_id":    []string{headers["ddmc-app-client-id"]},
		"applet_source":    []string{""},
		"channel":          []string{"App Store"},
		"sharer_uid":       []string{""},
		"h5_source":        []string{""},
	}
}

// SetQueryParams 设置query params参数，避免因多并发引起的concurrent map writes
func (u *User) SetQueryParams(params map[string]string) {
	u.mtx.Lock()
	defer u.mtx.Unlock()
	for k, v := range params {
		u.params[k] = []string{v}
	}
}

// QueryParams 返回query params的复制
func (u *User) QueryParams() url.Values {
	u.mtx.RLock()
	defer u.mtx.RUnlock()
	var cp = make(url.Values)
	for k, v := range u.params {
		cp[k] = v
	}
	return cp
}

func (u *User) GetUserDetail() (*UserDetail, error) {
	client := http.NewClient(constants.UserDetail)
	params := u.QueryParams()

	resp, err := client.Get(u.Headers(), params)
	if err != nil {
		return nil, err
	}

	return u.DecodeUser(resp.Data)
}

func (u *User) DecodeUser(data interface{}) (*UserDetail, error) {
	bytes, _ := json.Marshal(data)
	if err := json.Unmarshal(bytes, &u.userDetail); err != nil {
		return nil, fmt.Errorf("解析用户数据出错, 错误: %v", err.Error())
	}
	klog.Infof("获取用户信息成功, 用户: %s, id: %s", u.userDetail.UserInfo.Name, u.userDetail.UserInfo.Id)
	return u.userDetail, nil
}
