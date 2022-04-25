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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/robertkrimen/otto"
	"k8s.io/klog"
)

type JsSign struct {
	file string
}

func NewDefaultJsSign() SignInterface {
	dir, err := os.Getwd()
	if err != nil {
		klog.Fatal(err)
	}
	return NewSign(fmt.Sprintf("%s%s", dir, "/sign.js"))
}

func NewSign(file string) SignInterface {
	return &JsSign{
		file: file,
	}
}

func (s *JsSign) Sign(data interface{}) (map[string]string, error) {
	bytes, err := ioutil.ReadFile(s.file)
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	vm := otto.New()
	if _, err = vm.Run(string(bytes)); err != nil {
		klog.Error(err)
		return nil, err
	}
	bytes, _ = json.Marshal(data)
	value, err := vm.Call("sign", nil, string(bytes))
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	var signs map[string]string
	if err = json.Unmarshal([]byte(value.String()), &signs); err != nil {
		klog.Errorf("解析签名结果出错，错误: %v", err)
		return nil, err
	}
	return signs, nil
}
