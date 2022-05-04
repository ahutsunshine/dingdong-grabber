package ios

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sync"

	"github.com/dingdong-grabber/charles"
	"github.com/dingdong-grabber/pkg/util"
	"k8s.io/klog"
)

type session struct {
	headers map[string]string
	params  map[string]string
}

func NewDefaultSession() *session {
	return &session{
		headers: make(map[string]string),
		params:  make(map[string]string),
	}
}

type ios struct {
	s   *session
	mtx sync.RWMutex
}

func NewIosDevice() *ios {
	return &ios{
		s: NewDefaultSession(),
	}
}

func (i *ios) LoadConfig(file string) error {
	// 1. 获取根目录
	dir, err := util.GetRootDir()
	if err != nil {
		return err
	}

	// 2. 读取配置文件
	data, err := util.ReadFile(fmt.Sprintf("%s/%s/%s", dir, "charles/ios", file))
	if err != nil {
		return err
	}

	// 3. 解析配置文件
	return i.Decode(data)
}

func (i *ios) Decode(data []byte) error {
	var s []charles.Session
	if err := json.Unmarshal(data, &s); err != nil {
		klog.Error(err)
		return err
	}
	if len(s) == 0 {
		klog.Error("无效的cart.chlsj文件，无headers参数，请参考charles/ios/example.chlsj文件")
		return errors.New("无效文件")
	}
	if err := i.decodeHeader(s[0].Request.Header.Headers); err != nil {
		return err
	}

	return i.decodeParams(s[0].Query)
}

func (i *ios) decodeParams(queryStr string) error {
	values, err := url.ParseQuery(queryStr)
	if err != nil {
		klog.Errorf("解析请求参数出错, 详情: %v", err)
		return err
	}
	var params = map[string]string{
		"api_version":      values["api_version"][0],
		"app_client_id":    values["app_client_id"][0],
		"app_type":         values["app_type"][0],
		"buildVersion":     values["buildVersion"][0],
		"channel":          values["channel"][0],
		"city_number":      values["city_number"][0],
		"countryCode":      values["countryCode"][0],
		"device_id":        values["device_id"][0],
		"device_model":     values["device_model"][0],
		"device_name":      values["device_name"][0],
		"device_token":     values["device_token"][0],
		"idfa":             values["idfa"][0],
		"ip":               values["ip"][0],
		"languageCode":     values["languageCode"][0],
		"latitude":         values["latitude"][0],
		"localeIdentifier": values["localeIdentifier"][0],
		"longitude":        values["longitude"][0],
		"os_version":       values["os_version"][0],
		"seqid":            values["seqid"][0],
		"station_id":       values["station_id"][0],
		"time":             values["time"][0],
		"uid":              values["uid"][0],
	}
	for k, v := range values {
		params[k] = v[0]
	}
	i.SetParams(params)
	return nil
}

func (i *ios) decodeHeader(headers []charles.HeaderEntry) error {
	if len(headers) == 0 {
		klog.Error("无效的cart.chlsj文件，无headers参数，请参考charles/ios/example.chlsj文件")
		return errors.New("无效文件")
	}
	var header = make(map[string]string)
	for _, h := range headers {
		header[h.Name] = h.Value
	}
	header = map[string]string{
		"accept": header["accept"],
		// 不开始压缩
		"accept-encoding":        header["accept-encoding"],
		"accept-language":        header["accept-language"],
		"content-type":           "application/x-www-form-urlencoded",
		"cookie":                 header["cookie"],
		"x-tingyun-id":           header["x-tingyun-id"],
		"x-tingyun":              header["x-tingyun"],
		"ddmc-api-version":       header["ddmc-api-version"],
		"ddmc-app-client-id":     header["ddmc-app-client-id"],
		"ddmc-build-version":     header["ddmc-build-version"],
		"ddmc-channel":           header["ddmc-channel"],
		"ddmc-city-number":       header["ddmc-city-number"],
		"ddmc-country-code":      header["ddmc-country-code"],
		"ddmc-device-id":         header["ddmc-device-id"],
		"ddmc-device-model":      header["ddmc-device-model"],
		"ddmc-device-name":       header["ddmc-device-name"],
		"ddmc-device-token":      header["ddmc-device-token"],
		"ddmc-idfa":              header["ddmc-idfa"],
		"ddmc-ip":                header["ddmc-ip"],
		"ddmc-language-code":     header["ddmc-language-code"],
		"ddmc-latitude":          header["ddmc-latitude"],
		"ddmc-locale-identifier": header["ddmc-locale-identifier"],
		"ddmc-longitude":         header["ddmc-longitude"],
		"ddmc-os-version":        header["ddmc-os-version"],
		"ddmc-station-id":        header["ddmc-station-id"],
		"ddmc-uid":               header["ddmc-uid"],
		"time":                   header["time"],
		"user-agent":             header["user-agent"],
	}
	i.SetHeaders(header)
	return nil
}

// Headers 返回请求header的复制
func (i *ios) Headers() map[string]string {
	i.mtx.RLock()
	defer i.mtx.RUnlock()
	var cp = make(map[string]string)
	for k, v := range i.s.headers {
		cp[k] = v
	}
	return cp
}

func (i *ios) SetHeaders(header map[string]string) {
	i.mtx.Lock()
	defer i.mtx.Unlock()
	i.s.headers = header
}

func (i *ios) QueryParams() map[string]string {
	i.mtx.RLock()
	defer i.mtx.RUnlock()
	var cp = make(map[string]string)
	for k, v := range i.s.params {
		cp[k] = v
	}
	return cp
}

func (i *ios) SetParams(params map[string]string) {
	i.mtx.Lock()
	defer i.mtx.Unlock()
	for k, v := range params {
		i.s.params[k] = v
	}
}
