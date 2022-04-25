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

package constants

const (
	UserDetail  = "https://sunquan.api.ddxq.mobi/api/v1/user/detail/"
	Address     = "https://sunquan.api.ddxq.mobi/api/v1/user/address/"     // 获取用户默认地址， 此地址会作为买菜地址
	CartCheck   = "https://maicai.api.ddxq.mobi/cart/allCheck"             // 勾选购物车所有商品地址
	Cart        = "https://maicai.api.ddxq.mobi/cart/index"                // 获取购物车商品地址
	ReserveTime = "https://maicai.api.ddxq.mobi/order/getMultiReserveTime" // 预约送达时间地址
	CheckOrder  = "https://maicai.api.ddxq.mobi/order/checkOrder"          // 获取确认订单地址
	SubmitOrder = "https://maicai.api.ddxq.mobi/order/addNewOrder"         // 提交订单地址
	Push        = "http://www.pushplus.plus/send"                          //  推送地址
)
