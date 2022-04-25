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

	"github.com/dingdong-grabber/pkg/constants"
	"github.com/dingdong-grabber/pkg/http"
	"k8s.io/klog"
)

type User struct {
	c          *http.Client
	userDetail *UserDetail
	addressId  string
	headers    map[string]string
	body       url.Values
	mtx        sync.RWMutex
}

func NewDefaultUser() *User {
	return &User{
		c: &http.Client{},
	}
}

func (u *User) SetUserDetail(userDetail *UserDetail) {
	u.mtx.Lock()
	defer u.mtx.Unlock()
	u.userDetail = userDetail
}

func (u *User) UserDetail() *UserDetail {
	u.mtx.RLock()
	defer u.mtx.RUnlock()
	return u.userDetail
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

func (u *User) LoadConfig(cookie string) error {
	if cookie == "" {
		klog.Fatal("请求头cookie为必填项")
	}

	// 设置Header默认请求参数
	u.SetDefaultHeaders(cookie)

	// 设置Body默认请求参数
	u.SetDefaultBody()

	ud, err := u.GetUserDetail()
	if err != nil {
		return err
	}
	// 设置用户详情
	u.SetUserDetail(ud)

	// 设置header ddmc uid
	u.SetHeaders(map[string]string{
		"ddmc-uid": ud.UserInfo.Id,
	})

	// 设置body uid
	u.SetBody(map[string]string{
		"uid": ud.UserInfo.Id,
	})

	addr, err := u.GetDefaultAddr()
	if err != nil {
		return err
	}
	// 设置收货地址ID
	u.SetAddressId(addr.Id)

	// 设置Header和Body收货站ID和城市编码
	u.SetHeaders(map[string]string{
		"ddmc-station-id":  addr.StationId,
		"ddmc-city-number": addr.CityNumber,
	})
	u.SetBody(map[string]string{
		"station_id":  addr.StationId,
		"city_number": addr.CityNumber,
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
		"cookie": cookie,

		// 根据cookie动态获取
		"ddmc-uid": "",

		// 设置经纬度, 获取默认地址时会自动添加
		"ddmc-longitude": "",
		"ddmc-latitude":  "",

		// 下面作为小程序2.83.0版本的默认值
		"ddmc-build-version": "2.83.0",
		"ddmc-city-number":   "", // 程序会自动获取默认地址的city number填充于此
		"ddmc-station-id":    "", // 程序会自动获取默认地址的station id填充于此
		"ddmc-time":          fmt.Sprintf("%d", time.Now().UnixMilli()/1000),
		"ddmc-channel":       "applet",
		"ddmc-os-version":    "[object Undefined]",
		"ddmc-app-client-id": "4",
		"ddmc-ip":            "",
		"ddmc-api-version":   "9.50.0",
		"ddmc-device-id":     "",
		"referer":            "https://servicewechat.com/wx1e113254eda17715/425/page-frame.html",
		"content-type":       "application/x-www-form-urlencoded",
		"accept":             "*/*",
		"user-agent":         "Mozilla/5.0 (iPhone; CPU iPhone OS 11_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E217 MicroMessenger/6.8.0(0x16080000) NetType/WIFI Language/en Branch/Br_trunk MiniProgramEnv/Mac",
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

// Header 返回请求header的复制
func (u *User) Header() map[string]string {
	u.mtx.RLock()
	defer u.mtx.RUnlock()
	var cp = make(map[string]string)
	for k, v := range u.headers {
		cp[k] = v
	}
	return cp
}

// SetDefaultBody 设置默认的用户初始化数据
func (u *User) SetDefaultBody() {
	u.mtx.Lock()
	defer u.mtx.Unlock()
	var headers = u.headers
	u.body = url.Values{
		"s_id":         []string{""},
		"device_token": []string{""},

		// 下面作为小程序2.83.0版本的默认值
		"uid":           []string{headers["ddmc-uid"]},
		"longitude":     []string{headers["ddmc-longitude"]},
		"latitude":      []string{headers["ddmc-latitude"]},
		"station_id":    []string{headers["ddmc-station-id"]},
		"city_number":   []string{headers["ddmc-city-number"]},
		"api_version":   []string{headers["ddmc-api-version"]},
		"app_version":   []string{headers["ddmc-build-version"]},
		"time":          []string{headers["ddmc-time"]},
		"openid":        []string{headers["ddmc-device-id"]},
		"applet_source": []string{""},
		"channel":       []string{"applet"},
		"app_client_id": []string{"4"},
		"sharer_uid":    []string{""},
		"h5_source":     []string{""},
	}
}

// SetBody 设置body参数，避免body因多并发引起的concurrent map writes
func (u *User) SetBody(body map[string]string) {
	u.mtx.Lock()
	defer u.mtx.Unlock()
	for k, v := range body {
		u.body[k] = []string{v}
	}
}

// Body 返回请求body的复制
func (u *User) Body() url.Values {
	u.mtx.RLock()
	defer u.mtx.RUnlock()
	var cp = make(url.Values)
	for k, v := range u.body {
		cp[k] = v
	}
	return cp
}

func (u *User) GetUserDetail() (*UserDetail, error) {
	var (
		client = http.NewClient(constants.UserDetail)
		body   = u.Body()
	)

	resp, err := client.Get(u.Header(), body)
	if err != nil {
		klog.Info(err.Error())
		return nil, err
	}

	var ud UserDetail
	userBytes, _ := json.Marshal(resp.Data)
	if err := json.Unmarshal(userBytes, &ud); err != nil {
		return nil, fmt.Errorf("解析用户数据出错, 错误: %v", err.Error())
	}
	klog.Infof("获取用户信息成功, 用户: %s, id: %s", ud.UserInfo.Name, ud.UserInfo.Id)
	return &ud, nil
}
