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
	"encoding/json"
	"strings"
	"sync"

	"github.com/dingdong-grabber/pkg/constants"
	"github.com/dingdong-grabber/pkg/http"
	"github.com/dingdong-grabber/pkg/user"
	"github.com/google/uuid"
	"k8s.io/klog"
)

type Strategy int

const (
	ManualStrategy   Strategy = 0 // 人工
	TimingStrategy   Strategy = 1 // 定时
	SentinelStrategy Strategy = 2 // 哨兵
)

type Order struct {
	user         *user.User // 订单所属用户
	stop         bool       // 保护线程: 2分钟未下单自动终止,避免对叮咚服务器造成压力,也避免封号
	strategy     Strategy   // 不同策略抢购
	cart         map[string]interface{}
	reservedTime map[string]interface{}
	checkOrder   map[string]interface{}
	mtx          sync.RWMutex
}

func NewOrder(user *user.User, strategy Strategy) *Order {
	return &Order{
		user:     user,
		strategy: strategy,
	}
}

func (o *Order) User() *user.User {
	o.mtx.RLock()
	defer o.mtx.RUnlock()
	return o.user
}

func (o *Order) SetStop(stop bool) {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	o.stop = stop
}

func (o *Order) Stop() bool {
	o.mtx.RLock()
	defer o.mtx.RUnlock()
	return o.stop
}

func (o *Order) SetCart(cart map[string]interface{}) {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	o.cart = cart
}

func (o *Order) Cart() map[string]interface{} {
	o.mtx.RLock()
	defer o.mtx.RUnlock()
	return o.cart
}

func (o *Order) SetReservedTime(reservedTime map[string]interface{}) {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	o.reservedTime = reservedTime
}

func (o *Order) ReservedTime() map[string]interface{} {
	o.mtx.RLock()
	defer o.mtx.RUnlock()
	return o.reservedTime
}

func (o *Order) SetCheckOrder(checkOrder map[string]interface{}) {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	o.checkOrder = checkOrder
}

func (o *Order) CheckOrder() map[string]interface{} {
	o.mtx.RLock()
	defer o.mtx.RUnlock()
	return o.checkOrder
}

// GetMultiReserveTime 获取配送时间
func (o *Order) GetMultiReserveTime() (map[string]interface{}, error) {
	var (
		client      = http.NewClient(constants.ReserveTime)
		body        = o.user.Body()
		products, _ = json.Marshal([]interface{}{o.cart["products"]})
	)

	client.SetBody(body, map[string]string{
		// 关键参数
		"ab_config":  `{"ETA_time_default_selection":"D1.2"}`,
		"address_id": o.user.AddressId(),
		"products":   string(products),
	})

	resp, err := client.Post(o.user.Header(), body)
	if err != nil {
		klog.Errorf("获取预约时间失败, 错误: %v", err)
		return nil, err
	}

	var t []Times
	timesBytes, _ := json.Marshal(resp.Data)
	if err := json.Unmarshal(timesBytes, &t); err != nil {
		klog.Errorf("解析预约时间出错, 错误: %v", err.Error())
		return nil, err
	}

	// 判断叮咚官方是否提供可选的配送时间
	if len(t) == 0 || len(t[0].Times) == 0 {
		klog.Error("叮咚官方未提供任何配送时间")
		return nil, nil
	}

	var reservedTime = make(map[string]interface{})
	for _, d := range t[0].Times[0].Details {
		if d.DisableType == 0 && !strings.Contains(d.SelectMsg, "尽快") {
			reservedTime["reserved_time_start"] = d.StartTimestamp
			reservedTime["reserved_time_end"] = d.EndTimestamp
			klog.Infof("更新配送时间成功, 配送时间段: %s", d.SelectMsg)
			return reservedTime, nil
		}
	}

	var unableInfo = t[0].Times[0].TimeFullTextTip
	if unableInfo == "" && len(t[0].Times[0].Details) > 0 {
		unableInfo = t[0].Times[0].Details[0].DisableMsg
	}
	klog.Errorf("无可选的配送时间, 原因: %s", unableInfo)

	return nil, nil
}

// GetCheckOrder 获取订单确认信息
func (o *Order) GetCheckOrder() (map[string]interface{}, error) {
	var (
		client = http.NewClient(constants.CheckOrder)
		body   = o.user.Body()
		cart   = o.Cart()
	)

	// 构造商品参数信息
	packages := map[string]interface{}{
		"products":                cart["products"],
		"real_match_supply_order": false,
		"is_supply_order":         false,
		"package_type":            1,
		"package_id":              1,
		"reserved_time": map[string]interface{}{
			"time_biz_type":       0,
			"reserved_time_start": o.ReservedTime()["reserved_time_start"],
			"reserved_time_end":   o.ReservedTime()["reserved_time_end"],
		},
	}
	packagesBytes, _ := json.Marshal([]interface{}{packages})

	client.SetBody(body, map[string]string{
		// 设置基础参数信息
		"address_id":        o.user.AddressId(),
		"user_ticket_id":    "default",
		"freight_ticket_id": "default",
		"is_use_point":      "0",
		"is_use_balance":    "1",
		"is_buy_vip":        "0",
		"coupons_id":        "",
		"is_buy_coupons":    "0",
		"check_order_type":  "0",
		"packages":          string(packagesBytes),
	})

	resp, err := client.Post(o.user.Header(), body)
	if err != nil {
		klog.Errorf("获取订单确认信息失败, 错误: %s", err.Error())
		return nil, err
	}

	var orders Orders
	ordersBytes, _ := json.Marshal(resp.Data)
	if err := json.Unmarshal(ordersBytes, &orders); err != nil {
		klog.Infof("商品价格配送信息解析出错, 错误: %s", err.Error())
		return nil, err
	}

	klog.Info("确认订单信息成功")
	var ticketId = orders.Order.DefaultCoupon["default_coupon"]
	if ticketId != nil {
		switch ticketId.(type) {
		case map[string]interface{}:
			ticketId = ticketId.(map[string]interface{})["_id"]
		case map[string]string:
			ticketId = ticketId.(map[string]string)["_id"]
		case map[interface{}]interface{}:
			ticketId = ticketId.(map[interface{}]interface{})["_id"]
		default:
			ticketId = nil
		}
	}
	return map[string]interface{}{
		"total_money":            orders.Order.TotalMoney,
		"freight_discount_money": orders.Order.FreightDiscountMoney,
		"freight_money":          orders.Order.FreightMoney,
		"freight_real_money":     orders.Order.FreightRealMoney,
		"user_ticket_id":         ticketId,
	}, nil
}

// SubmitOrder 提交订单
func (o *Order) SubmitOrder() (bool, error) {
	var (
		client       = http.NewClient(constants.SubmitOrder)
		body         = o.user.Body()
		reservedTime = o.ReservedTime()
		checkOrder   = o.CheckOrder()
		cart         = o.Cart()
	)

	paymentOrder := map[string]interface{}{
		"price":                  checkOrder["total_money"],
		"freight_discount_money": checkOrder["freight_discount_money"],
		"freight_money":          checkOrder["freight_money"],
		"order_freight":          checkOrder["freight_real_money"],
		"parent_order_sign":      cart["parent_order_sign"],
		"address_id":             o.user.AddressId(),
		"form_id":                strings.ReplaceAll(uuid.New().String(), "-", ""),
		"receipt_without_sku":    "0",
		"pay_type":               2,
		"user_ticket_id":         checkOrder["user_ticket_id"],
		"current_position":       []string{body["latitude"][0], body["longitude"][0]},
	}
	packages := []map[string]interface{}{
		{
			"products":                cart["products"],
			"total_money":             cart["total_money"],
			"total_origin_money":      cart["total_origin_money"],
			"goods_real_money":        cart["goods_real_money"],
			"total_count":             cart["total_count"],
			"cart_count":              cart["cart_count"],
			"is_presale":              cart["is_presale"],
			"instant_rebate_money":    cart["instant_rebate_money"],
			"coupon_rebate_money":     cart["coupon_rebate_money"],
			"total_rebate_money":      cart["total_rebate_money"],
			"used_balance_money":      cart["used_balance_money"],
			"can_used_balance_money":  cart["can_used_balance_money"],
			"used_point_num":          cart["used_point_num"],
			"used_point_money":        cart["used_point_money"],
			"can_used_point_num":      cart["can_used_point_num"],
			"can_used_point_money":    cart["can_used_point_money"],
			"is_share_station":        cart["is_share_station"],
			"only_today_products":     cart["only_today_products"],
			"only_tomorrow_products":  cart["only_tomorrow_products"],
			"package_type":            cart["package_type"],
			"package_id":              cart["package_id"],
			"eta_trace_id":            "",
			"reserved_time_start":     reservedTime["reserved_time_start"],
			"reserved_time_end":       reservedTime["reserved_time_end"],
			"soon_arrival":            "",
			"first_selected_big_time": 0,
			"receipt_without_sku":     0,
		},
	}
	payment := map[string]interface{}{
		"payment_order": paymentOrder,
		"packages":      packages,
	}
	paymentBytes, _ := json.Marshal(payment)

	client.SetBody(body, map[string]string{
		"package_order": string(paymentBytes),
		"ab_config":     `{"key_no_condition_barter":false}`,
	})

	resp, err := client.Post(o.user.Header(), body)
	if err != nil {
		klog.Errorf("提交订单失败, 错误: %s 当前下单总金额：%v", err, cart["total_money"])
		return false, err
	}
	bytes, err := json.Marshal(resp)
	if err != nil {
		klog.Errorf("打印下单成功结果出错, 错误:  %v", err)
	}
	klog.Infof("下单成功返回的响应信息: %s", string(bytes))
	klog.Infof("恭喜你，已成功下单，当前下单总金额：%s", cart["total_money"])
	return true, nil
}
