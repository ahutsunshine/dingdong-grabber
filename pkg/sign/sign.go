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

package sign

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/dingdong-grabber/pkg/constants"
	"github.com/dingdong-grabber/pkg/util"
	"k8s.io/klog"
)

type JsSign struct {
	file string
}

func NewDefaultJsSign() (SignInterface, error) {
	data, err := util.SignFile()
	if err != nil {
		klog.Fatal(err)
	}

	resp, err := rawRequest(constants.SignSwitch, http.MethodGet, map[string]string{
		"user-agent": "axios/0.29.0",
	}, nil, nil)

	if err != nil {
		return nil, err
	}
	sign, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	signStr := strings.TrimSuffix(string(sign), ";")

	data = strings.ReplaceAll(data, "${SIGN}", signStr)
	tmp, _ := util.SignConfigFilePath()
	s := &JsSign{
		file: tmp,
	}
	if err = s.Write([]byte(data)); err != nil {
		return nil, err
	}
	return s, nil
}

func rawRequest(url, method string, header map[string]string, params url.Values, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		klog.Error(err)
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

func (s *JsSign) Write(data []byte) error {
	return ioutil.WriteFile(s.file, data, 0666)
}

func (s *JsSign) Sign(secret string, data interface{}) (map[string]string, error) {
	bytes, _ := json.Marshal(data)
	out, err := Exec("node", []string{
		s.file,
		secret,
		string(bytes),
	})
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	var sign map[string]string
	if err = json.Unmarshal([]byte(out), &sign); err != nil {
		klog.Error(err)
		return nil, err
	}
	return sign, nil
}
