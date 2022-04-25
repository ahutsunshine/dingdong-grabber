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

package order

import (
	"testing"

	"github.com/dingdong-grabber/pkg/user"
)

const cookie = ""

func TestOrder(t *testing.T) {
	t.Skip("以下为购物流程，便于开发者理解")

	u := user.NewDefaultUser()
	// 1. 初始化用户必须的参数数据
	if err := u.LoadConfig(cookie); err != nil {
		return
	}

	o := NewOrder(u, ManualStrategy)
	// 2. 全选购物车按钮
	if err := o.CheckAll(); err != nil {
		return
	}

	// 3. 获取购物车商品信息
	cart, err := o.GetCart()
	if err != nil {
		return
	}
	o.SetCart(cart)

	// 4. 获取配送时间
	reservedTimes, err := o.GetMultiReserveTime()
	if err != nil {
		return
	}
	o.SetReservedTime(reservedTimes)

	// 5. 确认订单信息
	checkOrder, err := o.GetCheckOrder()
	if err != nil {
		return
	}
	o.SetCheckOrder(checkOrder)

	// 6. 提交订单
	_, err = o.SubmitOrder()
	if err != nil {
		return
	}
}
