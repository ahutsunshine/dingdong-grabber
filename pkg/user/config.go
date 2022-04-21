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
	u.mtx.RLock()
	defer u.mtx.RUnlock()
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

func (u *User) SetClient(url string) {
	u.mtx.Lock()
	defer u.mtx.Unlock()
	u.c.Url = url
}

func (u *User) Client() *http.Client {
	u.mtx.RLock()
	defer u.mtx.RUnlock()
	return u.c
}

func (u *User) SetStationId(stationId string) {
	var (
		headers = u.Headers()
		body    = u.Body()
	)
	u.mtx.Lock()
	defer u.mtx.Unlock()
	headers["ddmc-station-id"] = stationId
	body["station_id"] = []string{stationId}
}

func (u *User) SetCityNumber(cityNumber string) {
	var (
		headers = u.Headers()
		body    = u.Body()
	)
	u.mtx.Lock()
	defer u.mtx.Unlock()
	headers["ddmc-city-number"] = cityNumber
	body["city_number"] = []string{cityNumber}
}

func (u *User) LoadConfig(cookie string) error {
	if cookie == "" {
		klog.Fatal("请求头cookie为必填项")
	}

	// 设置Header默认请求参数
	u.SetDefaultHeaders(cookie)

	// 设置Body默认请求参数
	u.SetDefaultBody()

	addr, err := u.GetDefaultAddr()
	if err != nil {
		return err
	}
	// 设置收货地址ID
	u.SetAddressId(addr.Id)

	// 设置收货站ID
	u.SetStationId(addr.StationId)

	// 设置城市编码
	u.SetCityNumber(addr.CityNumber)

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
	return nil
}

func (u *User) SetDefaultHeaders(cookie string) {
	u.mtx.RLock()
	defer u.mtx.RUnlock()
	if !strings.HasPrefix(cookie, "DDXQSESSID") {
		cookie = fmt.Sprintf("DDXQSESSID=%s", cookie)
	}
	u.headers = map[string]string{
		// Header必填项
		"cookie": cookie,

		// 根据cookie动态获取
		"ddmc-uid": "",

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
	u.mtx.RLock()
	defer u.mtx.RUnlock()
	for k, v := range headers {
		u.headers[k] = v
	}
}

func (u *User) Headers() map[string]string {
	u.mtx.RLock()
	defer u.mtx.RUnlock()
	return u.headers
}

// HeadersDeepCopy 为了避免多并发造成的并发读写问题: fatal error: concurrent map read and map write
func (u *User) HeadersDeepCopy() map[string]string {
	var headers = u.Headers()
	u.mtx.Lock()
	defer u.mtx.Unlock()
	var cp = make(map[string]string)
	for k, v := range headers {
		cp[k] = v
	}
	return cp
}

// SetDefaultBody 设置默认的用户初始化数据
func (u *User) SetDefaultBody() {
	var headers = u.Headers()
	u.mtx.RLock()
	defer u.mtx.RUnlock()
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

func (u *User) Body() url.Values {
	u.mtx.RLock()
	defer u.mtx.RUnlock()
	return u.body
}

// BodyDeepCopy 为了避免多并发造成的并发读写问题: fatal error: concurrent map read and map write
func (u *User) BodyDeepCopy() url.Values {
	var body = u.Body()
	u.mtx.Lock()
	defer u.mtx.Unlock()
	var cp = make(url.Values)
	for k, v := range body {
		cp[k] = v
	}
	return cp
}

func (u *User) GetUserDetail() (*UserDetail, error) {
	// body参数为共享，提交购物车时添加了products等参数，可能会导致请求参数过长造成invalid character '<' looking for beginning of value，这里重新设置为空字符
	u.SetBody(map[string]string{
		"products":      "",
		"package_order": "",
		"packages":      "",
	})

	u.SetClient(constants.UserDetail)
	resp, err := u.Client().Get(u.HeadersDeepCopy(), u.BodyDeepCopy())
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
