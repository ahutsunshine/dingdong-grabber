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
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/dingdong-grabber/pkg/constants"
	"github.com/dingdong-grabber/pkg/util"
	"k8s.io/klog"
)

var sign = &Cache{}

type Cache struct {
	file string
	mtx  sync.RWMutex
}

func (c *Cache) CacheFile() string {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.file
}

func (c *Cache) SetCacheFile(file string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.file = file
}

type NodeSign struct {
	file string
	mtx  sync.RWMutex
}

func NewSign(file string) SignInterface {
	return &NodeSign{
		file: file,
	}
}

func NewDefaultJsSign() (SignInterface, error) {
	dir, err := util.GetRootDir()
	if err != nil {
		return nil, err
	}
	template := fmt.Sprintf("%s/%s", dir, constants.SignFile)
	s := NewSign(template).(*NodeSign)

	// 获取cache
	if sign.CacheFile() != "" {
		return s, nil
	}

	// 1. 获取签名算法
	algorithm, err := s.algorithm()
	if err != nil {
		return nil, err
	}

	// 2. 生成实际签名算法
	algorithm, err = s.replace(algorithm)
	if err != nil {
		return nil, err
	}

	// 3.写入临时文件和全局cache中
	if err = s.cache(algorithm); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *NodeSign) Sign(secret string, data interface{}) (map[string]string, error) {
	bytes, _ := json.Marshal(data)
	out, err := Exec("node", []string{
		sign.CacheFile(),
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

func (s *NodeSign) algorithm() (string, error) {
	req, err := http.NewRequest(http.MethodGet, constants.SignAlgorithm, nil)
	if err != nil {
		klog.Error(err)
		return "", err
	}
	var client = &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	req.Header.Add("user-agent", "axios/0.29.0")
	resp, err := client.Do(req)
	if err != nil {
		klog.Error(err)
		return "", err
	}
	algorithm, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		klog.Error(err)
		return "", err
	}
	return string(algorithm), err
}

func (s *NodeSign) replace(algorithm string) (string, error) {
	data, err := ioutil.ReadFile(s.file)
	if err != nil {
		klog.Error(err)
		return "", err
	}
	algorithm = strings.TrimSuffix(algorithm, ";")
	algorithm = strings.ReplaceAll(string(data), "${SIGN}", algorithm)
	return algorithm, nil
}

func (s *NodeSign) cache(algorithm string) error {
	tmp, err := ioutil.TempFile("", "*.js")
	klog.Infof("签名临时文件: %s", tmp.Name())
	if err != nil {
		klog.Error(err)
		return err
	}
	if err = ioutil.WriteFile(tmp.Name(), []byte(algorithm), 0666); err != nil {
		klog.Error(err)
		return err
	}
	sign.SetCacheFile(tmp.Name())
	return nil
}
