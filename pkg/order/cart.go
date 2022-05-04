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

	"github.com/dingdong-grabber/pkg/constants"
	"github.com/dingdong-grabber/pkg/http"
	"k8s.io/klog"
)

// CheckAll 勾选购物车全选按钮
func (o *Order) CheckAll() error {
	var (
		client = http.NewClient(constants.CartCheck)
		params = o.user.QueryParams()
	)
	// 关键参数，必须要带
	client.SetParams(params, map[string]string{
		"ab_config": `{"key_cart_discount_price":"C","key_no_condition_barter":true,"key_show_cart_barter":"0"}`,
		"is_check":  "1",
		"is_load":   "1",
		"is_filter": "0",
	})

	if _, err := client.Get(o.user.Headers(), params); err != nil {
		klog.Infof("勾选购物车全选按钮失败, 错误: %v", err)
		return err
	}

	klog.Info("勾选购物车全选按钮成功")
	return nil
}

// GetCart 获取购物车商品信息
func (o *Order) GetCart() (map[string]interface{}, error) {
	var (
		client = http.NewClient(constants.Cart)
		params = o.user.QueryParams()
	)
	client.SetParams(params, map[string]string{
		"is_filter": "0",                                                                                          // 关键参数，必须要带
		"is_load":   "1",                                                                                          // 关键参数，必须要带
		"ab_config": `{"key_show_cart_barter":"0","key_no_condition_barter":false,"key_cart_discount_price":"C"}`, // 可选参数
	})

	resp, err := client.Get(o.user.Headers(), params)
	if err != nil {
		klog.Errorf("获取购物车商品失败, 错误: %v", err)
		return nil, err
	}

	var cart Cart
	cartBytes, _ := json.Marshal(resp.Data)
	if err := json.Unmarshal(cartBytes, &cart); err != nil {
		klog.Infof("商品解析出错, 错误: %s", err.Error())
		return nil, err
	}

	if len(cart.NewOrderProductList) == 0 {
		klog.Info("购物车无可购买的商品")
		// 人工策略无可购买的商品则停止抢购
		if o.strategy == ManualStrategy {
			o.SetStop(true)
		}
		return nil, nil
	}

	var pl = cart.NewOrderProductList[0]
	for _, p := range pl.Products {
		p["total_money"] = p["total_price"]
		p["total_origin_money"] = p["total_origin_price"]
	}

	klog.Infof("更新购物车数据成功, 订单金额：%s", pl.TotalMoney)
	return map[string]interface{}{
		"products":                  pl.Products,
		"parent_order_sign":         cart.ParentOrderInfo.ParentOrderSign,
		"total_money":               pl.TotalMoney,
		"total_origin_money":        pl.TotalOriginMoney,
		"goods_real_money":          pl.GoodsRealMoney,
		"total_count":               pl.TotalCount,
		"cart_count":                pl.CartCount,
		"is_presale":                pl.IsPresale,
		"instant_rebate_money":      pl.InstantRebateMoney,
		"coupon_rebate_money":       pl.CouponRebateMoney,
		"total_rebate_money":        pl.TotalRebateMoney,
		"used_balance_money":        pl.UsedBalanceMoney,
		"can_used_balance_money":    pl.CanUsedBalanceMoney,
		"used_point_num":            pl.UsedPointNum,
		"used_point_money":          pl.UsedPointMoney,
		"can_used_point_num":        pl.CanUsedPointNum,
		"can_used_point_money":      pl.CanUsedPointMoney,
		"is_share_station":          pl.IsShareStation,
		"only_today_products":       pl.OnlyTodayProducts,
		"only_tomorrow_products":    pl.OnlyTomorrowProducts,
		"package_type":              pl.PackageType,
		"package_id":                pl.PackageId,
		"front_package_text":        pl.FrontPackageText,
		"front_package_type":        pl.FrontPackageType,
		"front_package_stock_color": pl.FrontPackageStockColor,
		"front_package_bg_color":    pl.FrontPackageBgColor,
	}, nil
}
