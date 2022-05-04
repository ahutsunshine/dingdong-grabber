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

package strategy

import (
	"context"

	"github.com/dingdong-grabber/pkg/order"
	"k8s.io/klog"
)

// TestScheduler 仅用于测试整个流程是否跑通和可用
type TestScheduler struct {
	o *order.Order
}

func NewTestScheduler(o *order.Order) Interface {
	return &TestScheduler{
		o: o,
	}
}

func (ts *TestScheduler) Schedule(ctx context.Context) error {
	klog.Info("正在使用测试模式测试流程是否跑通")

	// 1. 全选购物车按钮
	err := ts.o.CheckAll()
	if err != nil {
		klog.Error("勾选购物车步骤出错")
		return err
	}

	// 2. 获取购物车订单信息
	cart, err := ts.o.GetCart()
	if err != nil {
		klog.Error("获取购物车订单步骤出错")
		return err
	}
	ts.o.SetCart(cart)

	// 3. 获取预约配送时间
	reservedTimes, err := ts.o.GetMultiReserveTime()
	if err != nil {
		klog.Error("获取预约配送时间步骤出错")
		return err
	}
	ts.o.SetReservedTime(reservedTimes)

	// 4. 确认订单信息
	checkOrder, err := ts.o.GetCheckOrder()
	if err != nil {
		klog.Error("确认订单步骤出错")
		return err
	}
	ts.o.SetCheckOrder(checkOrder)

	// 5. 提交订单
	if ts.o.ReservedTime() != nil {
		_, err = ts.o.SubmitOrder()
		if err != nil {
			klog.Error("提交订单步骤出错")
			return err
		}
	}

	klog.Info("流程测试结束。恭喜你，整个流程畅通无阻，配置正确")
	return nil
}
