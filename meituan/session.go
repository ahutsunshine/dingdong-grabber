package meituan

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type MeiTuanSession struct {
	Client   *http.Client `json:"client"`
	Header   Header       `json:"header"`
	UserInfo *UserInfo    `json:"UserInfo"`
}

func (s *MeiTuanSession) InitSession(userInfo *UserInfo) error {
	fmt.Println("########## 美团初始化 ##########")
	s.Client = &http.Client{}
	s.Header = newDefaultHeader()
	s.UserInfo = userInfo
	s.Header.OpenId = userInfo.OpenId
	s.Header.T = userInfo.T
	s.Header.OpenIdCipher = userInfo.OpenIdCipher

	return nil
}

func (s *MeiTuanSession) initUrlParams(templates *UrlParams) url.Values {
	// 参数
	templates.uuid = s.UserInfo.UUID
	templates.xuuid = s.UserInfo.UUID
	templates.userid = s.UserInfo.UserId
	templates.address_id = s.UserInfo.AddressId
	templates.openId = s.UserInfo.OpenId
	params := url.Values{}
	t := reflect.TypeOf(*templates)
	v := reflect.ValueOf(*templates)
	for k := 0; k < t.NumField(); k++ {
		params.Set(strings.Trim(string(t.Field(k).Tag[5:]), "\""), v.Field(k).String())
	}
	return params
}

func (s *MeiTuanSession) initClientHeader(req *http.Request) *http.Request {
	// 参数
	//header
	t := reflect.TypeOf(s.Header)
	v := reflect.ValueOf(s.Header)
	for k := 0; k < t.NumField(); k++ {
		req.Header.Set(string(t.Field(k).Tag[5:]), v.Field(k).String())
	}
	return req
}
