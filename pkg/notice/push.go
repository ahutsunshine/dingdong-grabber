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

package notice

import (
	"encoding/json"
	"fmt"
	nethttp "net/http"

	"github.com/dingdong-grabber/pkg/constants"
	"github.com/dingdong-grabber/pkg/http"
	"k8s.io/klog"
)

type Push struct {
	token   string
	title   string
	content string
}

func NewPush(token, title, content string) NoticeInterface {
	return &Push{
		token:   token,
		title:   title,
		content: content,
	}
}

func (p *Push) Notify() error {
	marshal, err := json.Marshal(p)
	if err != nil {
		klog.Errorf("序列化推送内容失败，错误: %v", err)
		return err
	}

	client := http.NewClient(constants.Push)
	resp, err := client.RawPost(nil, nil, marshal)
	if err != nil {
		klog.Error(err)
		return err
	}

	if resp.StatusCode != nethttp.StatusOK && resp.StatusCode != nethttp.StatusCreated {
		klog.Infof("推送返回不合法的状态值: %d", resp.StatusCode)
		return fmt.Errorf("%v", resp.StatusCode)
	}
	return nil
}
