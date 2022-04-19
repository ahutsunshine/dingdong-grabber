package order

import (
	"encoding/json"

	"github.com/dingdong-grabber/pkg/constants"
	"k8s.io/klog"
)

// CheckAll 勾选购物车全选按钮
func (o *Order) CheckAll() error {
	// 关键参数，必须要带
	o.user.SetBody(map[string]string{
		"is_check": "1",
		"is_load":  "1",
	})

	o.user.SetClient(constants.CartCheck)
	if _, err := o.user.Client().Get(o.user.HeadersDeepCopy(), o.user.BodyDeepCopy()); err != nil {
		klog.Infof("勾选购物车全选按钮失败, 错误: %v", err)
		return err
	}

	klog.Info("勾选购物车全选按钮成功")
	return nil
}

// GetCart 获取购物车商品信息
func (o *Order) GetCart() (map[string]interface{}, error) {
	o.user.SetBody(map[string]string{
		"is_load":   "1",                                                       // 关键参数，必须要带
		"ab_config": "{\"key_onion\":\"D\",\"key_cart_discount_price\":\"C\"}", // 可选参数
	})

	o.user.SetClient(constants.Cart)
	resp, err := o.user.Client().Get(o.user.HeadersDeepCopy(), o.user.BodyDeepCopy())
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
