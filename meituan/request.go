package meituan

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

//func httpGet(urlPath string, params url.Values, reader *bytes.Reader) {
//	Url, _ := url.Parse(urlPath)
//	Url.RawQuery = params.Encode()
//	req, _ := http.NewRequest("GET", Url.String(), nil)
//
//	return req
//}

var client = &http.Client{}

func newReader(cBody *interface{}) (*bytes.Reader, error) {
	bodyBytes, err := json.Marshal(cBody)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(bodyBytes), nil
}

func setHeaders(reqHeader *http.Header, headerParams *Header) {
	t := reflect.TypeOf(*headerParams)
	v := reflect.ValueOf(*headerParams)
	for k := 0; k < t.NumField(); k++ {
		reqHeader.Set(strings.Trim(string(t.Field(k).Tag[5:]), "\""), v.Field(k).String())
	}
}

func httpPost(urlPath string, params url.Values, cBody interface{}, header Header) (string, error) {

	// url
	Url, _ := url.Parse(urlPath)
	Url.RawQuery = params.Encode()
	reader, err2 := newReader(&cBody)
	if err2 != nil {
		return "", err2
	}

	//request
	req, _ := http.NewRequest("POST", Url.String(), reader)

	// header
	setHeaders(&req.Header, &header)

	// response
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == 200 {
		return string(body), nil
	}

	return "", errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body))
}
