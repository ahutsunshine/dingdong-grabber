package user

import (
	"fmt"
	"net/url"
	"time"

	"k8s.io/klog"
)

var (
	// --------Header请求必填项-------
	deviceId  = "osP8I0U_MVm-CgBtnDMdBWB16-gQ"                                                                                                                                                // 设置ID
	cookie    = "DDXQSESSID=1eaf61bab8ea4c52079a96618698b1a6"                                                                                                                                 // 用户cookie凭证
	longitude = "121.550668"                                                                                                                                                                  // 定位所在的经度
	latitude  = "31.199737"                                                                                                                                                                   // 定位所在的维度
	uid       = "611477235be794000133f3ab"                                                                                                                                                    // 用户 id
	userAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 15_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.20(0x18001428) NetType/WIFI Language/zh_CN" // 用户使用的设备

	// -------Body请求必填项----------
	sid         = "54ece57125b327fe370a3e66db7b354a"                                                                                                       // sid
	deviceToken = "WHJMrwNw1k/FKPjcOOgRd+ARHDjBJ4zcumZZY909hrTQHxyALdxSlvkVoVwRzcyu1yGtMZyJitHH6ECpmvFnZz8RbSmxIUOBDdCW1tldyDzmauSxIJm5Txg==1487582755342" //  设备token
)

func (u *User) LoadConfig() error {
	if deviceId == "" || cookie == "" || longitude == "" || latitude == "" || uid == "" || userAgent == "" {
		klog.Fatal("Header请求项deviceId, cookie, longitude, latitude, uid, userAgent为必填项")
	}
	u.SetHeaders(deviceId, cookie, longitude, latitude, uid, userAgent)

	if sid == "" || deviceToken == "" {
		klog.Fatal("Body请求项sid, deviceToken为必填项")
	}
	u.SetBody(sid, deviceToken)
	if addr, err := u.GetDefaultAddr(); err != nil {
		return err
	} else {
		// 设置收货站ID
		u.SetStationId(addr.StationId)
		// 设置城市编码
		u.SetCityNumber(addr.CityNumber)
	}

	return nil
}

func (u *User) SetHeaders(deviceId, cookie, longitude, latitude, uid, userAgent string) {
	u.mtx.RLock()
	defer u.mtx.RUnlock()

	u.headers = map[string]string{
		// Header必填项
		"ddmc-device-id": deviceId,
		"cookie":         cookie,
		"ddmc-longitude": longitude,
		"ddmc-latitude":  latitude,
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

func (u *User) Headers() map[string]string {
	u.mtx.RLock()
	defer u.mtx.RUnlock()
	return u.headers
}

func (u *User) SetBody(sid, deviceToken string) {
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

func (u *User) Body() url.Values {
	u.mtx.RLock()
	defer u.mtx.RUnlock()
	return u.body
}
