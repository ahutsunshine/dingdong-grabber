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

type Cart struct {
	NewOrderProductList []newOrderProductList `json:"new_order_product_list"`
	ParentOrderInfo     ParentOrderInfo       `json:"parent_order_info"`
}

type newOrderProductList struct {
	Products               []map[string]interface{} `json:"products"`
	ParentOrderInfo        map[string]interface{}   `json:"parent_order_info"`
	TotalMoney             string                   `json:"total_money"`
	TotalOriginMoney       string                   `json:"total_origin_money"`
	GoodsRealMoney         string                   `json:"goods_real_money"`
	TotalCount             int                      `json:"total_count"`
	CartCount              int                      `json:"cart_count"`
	IsPresale              int                      `json:"is_presale"`
	InstantRebateMoney     string                   `json:"instant_rebate_money"`
	CouponRebateMoney      string                   `json:"coupon_rebate_money"`
	TotalRebateMoney       string                   `json:"total_rebate_money"`
	UsedBalanceMoney       string                   `json:"used_balance_money"`
	CanUsedBalanceMoney    string                   `json:"can_used_balance_money"`
	UsedPointNum           int                      `json:"used_point_num"`
	UsedPointMoney         string                   `json:"used_point_money"`
	CanUsedPointNum        int                      `json:"can_used_point_num"`
	CanUsedPointMoney      string                   `json:"can_used_point_money"`
	IsShareStation         int                      `json:"is_share_station"`
	OnlyTodayProducts      interface{}              `json:"only_today_products"`
	OnlyTomorrowProducts   interface{}              `json:"only_tomorrow_products"`
	PackageType            int                      `json:"package_type"`
	PackageId              int                      `json:"package_id"`
	FrontPackageText       string                   `json:"front_package_text"`
	FrontPackageType       int                      `json:"front_package_type"`
	FrontPackageStockColor string                   `json:"front_package_stock_color"`
	FrontPackageBgColor    string                   `json:"front_package_bg_color"`
}

type ParentOrderInfo struct {
	ParentOrderSign  string `json:"parent_order_sign"`
	TotalRebateMoney string `json:"total_rebate_money"`
	TotalOriginMoney string `json:"total_origin_money"`
}

type Times struct {
	StationId string `json:"station_id"`
	Times     []Time `json:"time"`
}
type Time struct {
	TimeFullTextTip string       `json:"time_full_text_tip"`
	IsInvalid       bool         `json:"is_invalid"`
	Details         []TimeDetail `json:"times"`
}

type TimeDetail struct {
	DisableType    int    `json:"disableType"`
	DisableMsg     string `json:"disableMsg"`
	StartTimestamp int64  `json:"start_timestamp"`
	EndTimestamp   int64  `json:"end_timestamp"`
	SelectMsg      string `json:"select_msg"`
}

type Orders struct {
	Order OrderDetail `json:"order"`
}

type OrderDetail struct {
	FreightDiscountMoney string            `json:"freight_discount_money"` // 配送费折扣
	FreightMoney         string            `json:"freight_money"`          // 配送费
	TotalMoney           string            `json:"total_money"`            // 订单总价格
	FreightRealMoney     string            `json:"freight_real_money"`     // 最终的配送费
	DefaultCoupon        map[string]Coupon `json:"default_coupon"`         // 购物券
}

type Coupon struct {
	Id string `json:"_id"`
}
