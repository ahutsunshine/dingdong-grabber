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

type Addresses struct {
	ValidAddress    []Address `json:"valid_address"`
	MaxAddressCount int       `json:"max_address_count"`
}

type Address struct {
	Id          string      `json:"id"`
	Gender      int         `json:"gender"`
	Mobile      string      `json:"mobile"`
	Location    location    `json:"location"`
	UserName    string      `json:"user_name"`
	AddrDetail  string      `json:"addr_detail"`
	StationId   string      `json:"station_id"`
	StationName string      `json:"station_name"`
	StationInfo stationInfo `json:"station_info"`
	IsDefault   bool        `json:"is_default"`
	CityNumber  string      `json:"city_number"`
}

type stationInfo struct {
	Id           string `json:"id"`
	Address      string `json:"address"`
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	BusinessTime string `json:"business_time"`
	CityName     string `json:"city_name"`
	CityNumber   string `json:"city_number"`
}

type location struct {
	TypeCode string    `json:"typecode"`
	Address  string    `json:"address"`
	Name     string    `json:"name"`
	Location []float64 `json:"location"`
}

type UserDetail struct {
	// 用户详细信息
	DoingRefundNum      int         `json:"doing_refund_num"`
	NoCommentOrderPoint int         `json:"no_comment_order_point"`
	NameNotice          string      `json:"name_notice"`
	NoPayOrderNum       int         `json:"no_pay_order_num"`
	DoingOrderNum       int         `json:"doing_order_num"`
	UserVip             UserVIP     `json:"user_vip"`
	UserSign            UserSign    `json:"user_sign"`
	NotOnionTip         int         `json:"not_onion_tip"`
	NoDrawCouponMoney   string      `json:"no_draw_coupon_money"`
	PointNum            int         `json:"point_num"`
	Balance             UserBalance `json:"balance"`
	UserInfo            UserInfo    `json:"user_info"`
	CouponNum           int         `json:"coupon_num"`
	NoCommentOrderNum   int         `json:"no_comment_order_num"`
}

type UserVIP struct {
	IsRenew                  int    `json:"is_renew"`
	VipSaveMoneyDescription  string `json:"vip_save_money_description"`
	VipDescription           string `json:"vip_description"`
	VipStatus                int    `json:"vip_status"`
	VipNotice                string `json:"vip_notice"`
	VipExpireTimeDescription string `json:"vip_expire_time_description"`
	VipUrl                   string `json:"vip_url"`
}

type UserSign struct {
	IsTodaySign bool   `json:"is_today_sign"`
	SignSeries  int    `json:"sign_series"`
	SignText    string `json:"sign_text"`
}

type UserBalance struct {
	SetFingerPayPassword int    `json:"set_finger_pay_password"`
	Balance              string `json:"balance"`
	SetPayPassword       int    `json:"set_pay_password"`
}

type UserInfo struct {
	Birthday       string `json:"birthday"`
	ShowInviteCode bool   `json:"show_invite_code"`
	NameInCheck    string `json:"name_in_check"`
	InviteCodeUrl  string `json:"invite_code_url"`
	Sex            int    `json:"sex"`
	Mobile         string `json:"mobile"`
	Avatar         string `json:"avatar"`
	ImUid          int    `json:"im_uid"`
	BindStatus     int    `json:"bind_status"`
	NameStatus     int    `json:"name_status"`
	NewRegister    bool   `json:"new_register"`
	ImSecret       string `json:"im_secret"`
	Name           string `json:"name"`
	Id             string `json:"id"`
	Introduction   string `json:"introduction"`
}
