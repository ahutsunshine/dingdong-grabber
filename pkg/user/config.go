package user

import (
	"fmt"
	"net/url"
	"time"

	"k8s.io/klog"
)

func (u *User) LoadConfig(deviceId, cookie, uid, userAgent, sid, deviceToken string) error {
	if deviceId == "" || cookie == "" || uid == "" || userAgent == "" {
		klog.Fatal("Header请求项deviceId, cookie, uid, userAgent为必填项")
	}
	if sid == "" || deviceToken == "" {
		klog.Fatal("Body请求项sid, deviceToken为必填项")
	}

	// 设置Header默认请求参数
	u.SetDefaultHeaders(deviceId, cookie, uid, userAgent)

	// 设置Body默认请求参数
	u.SetDefaultBody(sid, deviceToken)

	if addr, err := u.GetDefaultAddr(); err != nil {
		return err
	} else {
		// 设置收货地址ID
		u.SetAddressId(addr.Id)
		// 设置收货站ID
		u.SetStationId(addr.StationId)
		// 设置城市编码
		u.SetCityNumber(addr.CityNumber)
	}

	return nil
}

func (u *User) SetDefaultHeaders(deviceId, cookie, uid, userAgent string) {
	u.mtx.RLock()
	defer u.mtx.RUnlock()

	u.headers = map[string]string{
		// Header必填项
		"ddmc-device-id": deviceId,
		"cookie":         cookie,
		"ddmc-uid":       uid,
		"user-agent":     userAgent,

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
		"referer":            "https://servicewechat.com/wx1e113254eda17715/425/page-frame.html",
		"content-type":       "application/x-www-form-urlencoded",
		"accept":             "*/*",
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

// SetDefaultBody 设置默认的用户初始化数据
func (u *User) SetDefaultBody(sid, deviceToken string) {
	var headers = u.Headers()
	u.mtx.RLock()
	defer u.mtx.RUnlock()
	u.body = url.Values{
		// Body必填项
		"s_id":         []string{sid},
		"device_token": []string{deviceToken},

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
	u.mtx.RLock()
	defer u.mtx.RUnlock()
	for k, v := range body {
		u.body[k] = []string{v}
	}
}

func (u *User) Body() url.Values {
	u.mtx.RLock()
	defer u.mtx.RUnlock()
	return u.body
}
