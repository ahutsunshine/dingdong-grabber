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

package user

import (
	"encoding/json"
	"errors"

	"github.com/dingdong-grabber/pkg/constants"
	"github.com/dingdong-grabber/pkg/http"
	"k8s.io/klog"
)

// GetDefaultAddr 获取默认地址 设置配送地址id，必须保证默认收获地址在上海且填写正确作为收获地址，请注意输出信息并确认
func (u *User) GetDefaultAddr() (*Address, error) {
	client := http.NewClient(constants.Address)
	params := u.QueryParams()
	resp, err := client.Get(u.Headers(), params)
	if err != nil {
		return nil, err
	}

	addrs, err := u.DecodeAddress(resp.Data)
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		if addr.IsDefault {
			klog.Infof("1.默认收货地址：%s%s%s, 手机号: %s", addr.Location.Address, addr.Location.Name, addr.AddrDetail, addr.Mobile)
			klog.Infof("2.该地址对应站点名称为：%s", addr.StationInfo.Name)
			klog.Infof("3.设置买菜地址经度：%v", addr.Location.Location[0])
			klog.Infof("4.设置买菜地址纬度：%v", addr.Location.Location[1])

			if station := u.headers["ddmc-station-id"]; station != "" && addr.StationId != "" && station != addr.StationId {
				klog.Errorf("默认地址ddmc-station-id和cart.chlsj设置的值不一致，请先将收获地址设为默认地址，然后重新获取cart.chlsj")
				klog.Infof("默认地址ddmc-station-id: %v, cart.chlsj ddmc-station-id: %v", addr.StationId, station)
				return nil, errors.New("无效ddmc-station-id")
			}
			if cityNumber := u.headers["ddmc-city-number"]; cityNumber != "" && addr.CityNumber != "" && cityNumber != addr.CityNumber {
				klog.Errorf("默认地址ddmc-city-number和cart.chlsj设置的值不一致，请先将收获地址设为默认地址，然后重新获取cart.chlsj")
				klog.Infof("默认地址ddmc-city-number:%v, cart.chlsj ddmc-city-number: %v", addr.CityNumber, cityNumber)
				return nil, errors.New("无效ddmc-city-number")
			}
			return &addr, nil
		}
	}

	klog.Info("请设置收货地址为默认地址")
	return nil, errors.New("请设置收货地址为默认地址")
}

func (u *User) DecodeAddress(data interface{}) ([]Address, error) {
	bytes, _ := json.Marshal(data)
	var addr *Addresses
	if err := json.Unmarshal(bytes, &addr); err != nil {
		klog.Infof("地址解析出错, 错误: %s", err.Error())
		return nil, err
	}
	if len(addr.ValidAddress) == 0 {
		klog.Info("请添加收货地址，并设置买菜地址为默认地址")
		return nil, errors.New("请添加收货地址，并设置买菜地址为默认地址")
	}
	return addr.ValidAddress, nil
}
