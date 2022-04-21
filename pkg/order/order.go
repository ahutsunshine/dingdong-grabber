package order

import (
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/dingdong-grabber/pkg/constants"
	"github.com/dingdong-grabber/pkg/user"
	"github.com/google/uuid"
	"k8s.io/klog"
)

type Strategy int

const (
	ManualStrategy Strategy = 0 // 人工
	TimingStrategy Strategy = 1 // 定时
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
	products, _ := json.Marshal([]interface{}{o.cart["products"]})
	o.user.SetBody(map[string]string{
		// 关键参数
		"address_id":      o.user.AddressId(),
		"products":        string(products),
		"group_config_id": "",
		"isBridge":        "false",
	})

	o.user.SetClient(constants.ReserveTime)
	resp, err := o.user.Client().Post(o.user.HeadersDeepCopy(), o.user.BodyDeepCopy())
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
			klog.Infof("更新配送时间成功, 配送时间段: [%v, %v]", timestamp2Str(d.StartTimestamp), timestamp2Str(d.EndTimestamp))
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

func timestamp2Str(tst int64) string {
	return time.Unix(tst, 0).Format("2006-01-02 15:04:05")
}

// GetCheckOrder 获取订单确认信息
func (o *Order) GetCheckOrder() (map[string]interface{}, error) {
	cart := o.Cart()
	// 构造商品参数信息
	packages := map[string]interface{}{
		"products":                  cart["products"],
		"total_money":               cart["total_money"],
		"total_origin_money":        cart["total_origin_money"],
		"goods_real_money":          cart["goods_real_money"],
		"total_count":               cart["total_count"],
		"cart_count":                cart["cart_count"],
		"is_presale":                cart["is_presale"],
		"instant_rebate_money":      cart["instant_rebate_money"],
		"coupon_rebate_money":       cart["coupon_rebate_money"],
		"total_rebate_money":        cart["total_rebate_money"],
		"used_balance_money":        cart["used_balance_money"],
		"can_used_balance_money":    cart["can_used_balance_money"],
		"used_point_num":            cart["used_point_num"],
		"used_point_money":          cart["used_point_money"],
		"can_used_point_num":        cart["can_used_point_num"],
		"can_used_point_money":      cart["can_used_point_money"],
		"is_share_station":          cart["is_share_station"],
		"only_today_products":       cart["only_today_products"],
		"only_tomorrow_products":    cart["only_tomorrow_products"],
		"package_type":              cart["package_type"],
		"package_id":                cart["package_id"],
		"front_package_text":        cart["front_package_text"],
		"front_package_type":        cart["front_package_type"],
		"front_package_stock_color": cart["front_package_stock_color"],
		"front_package_bg_color":    cart["front_package_bg_color"],
		"reserved_time": map[string]interface{}{
			"reserved_time_start": o.ReservedTime()["reserved_time_start"],
			"reserved_time_end":   o.ReservedTime()["reserved_time_end"],
		},
	}
	packagesBytes, _ := json.Marshal([]interface{}{packages})

	o.user.SetBody(map[string]string{
		// 设置基础参数信息
		"address_id":               o.user.AddressId(),
		"user_ticket_id":           "default",
		"freight_ticket_id":        "default",
		"is_use_point":             "0",
		"is_use_balance":           "0",
		"is_buy_vip":               "0",
		"coupons_id":               "",
		"is_buy_coupons":           "0",
		"check_order_type":         "0",
		"is_support_merge_payment": "1",
		"showData":                 "true",
		"showMsg":                  "false",
		"packages":                 string(packagesBytes),
	})

	o.user.SetClient(constants.CheckOrder)
	resp, err := o.user.Client().Post(o.user.HeadersDeepCopy(), o.user.BodyDeepCopy())
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
	return map[string]interface{}{
		"total_money":            orders.Order.TotalMoney,
		"freight_discount_money": orders.Order.FreightDiscountMoney,
		"freight_money":          orders.Order.FreightMoney,
		"freight_real_money":     orders.Order.FreightRealMoney,
		"user_ticket_id":         orders.Order.DefaultCoupon["default_coupon"].Id,
	}, nil
}

// SubmitOrder 提交订单
func (o *Order) SubmitOrder() (bool, error) {
	reservedTime := o.ReservedTime()
	checkOrder := o.CheckOrder()
	cart := o.Cart()
	paymentOrder := map[string]interface{}{
		"reserved_time_start":    reservedTime["reserved_time_start"],
		"reserved_time_end":      reservedTime["reserved_time_end"],
		"price":                  checkOrder["total_money"],
		"freight_discount_money": checkOrder["freight_discount_money"],
		"freight_money":          checkOrder["freight_money"],
		"order_freight":          checkOrder["freight_real_money"],
		"parent_order_sign":      cart["parent_order_sign"],
		"product_type":           1,
		"address_id":             o.user.AddressId(),
		"form_id":                strings.ReplaceAll(uuid.New().String(), "-", ""),
		"receipt_without_sku":    nil,
		"pay_type":               6, // 2: 支付宝支付, 4: 微信支付，6: 微信小程序支付
		"user_ticket_id":         checkOrder["user_ticket_id"],
		"vip_money":              "",
		"vip_buy_user_ticket_id": "",
		"coupons_money":          "",
		"coupons_id":             "",
	}
	packages := []map[string]interface{}{
		{
			"products":                  cart["products"],
			"total_money":               cart["total_money"],
			"total_origin_money":        cart["total_origin_money"],
			"goods_real_money":          cart["goods_real_money"],
			"total_count":               cart["total_count"],
			"cart_count":                cart["cart_count"],
			"is_presale":                cart["is_presale"],
			"instant_rebate_money":      cart["instant_rebate_money"],
			"coupon_rebate_money":       cart["coupon_rebate_money"],
			"total_rebate_money":        cart["total_rebate_money"],
			"used_balance_money":        cart["used_balance_money"],
			"can_used_balance_money":    cart["can_used_balance_money"],
			"used_point_num":            cart["used_point_num"],
			"used_point_money":          cart["used_point_money"],
			"can_used_point_num":        cart["can_used_point_num"],
			"can_used_point_money":      cart["can_used_point_money"],
			"is_share_station":          cart["is_share_station"],
			"only_today_products":       cart["only_today_products"],
			"only_tomorrow_products":    cart["only_tomorrow_products"],
			"package_type":              cart["package_type"],
			"package_id":                cart["package_id"],
			"front_package_text":        cart["front_package_text"],
			"front_package_type":        cart["front_package_type"],
			"front_package_stock_color": cart["front_package_stock_color"],
			"front_package_bg_color":    cart["front_package_bg_color"],
			"eta_trace_id":              "",
			"reserved_time_start":       reservedTime["reserved_time_start"],
			"reserved_time_end":         reservedTime["reserved_time_end"],
			"soon_arrival":              "",
			"first_selected_big_time":   0,
			"receipt_without_sku":       0,
		},
	}
	payment := map[string]interface{}{
		"payment_order": paymentOrder,
		"packages":      packages,
	}
	paymentBytes, _ := json.Marshal(payment)

	o.user.SetBody(map[string]string{
		"package_order": string(paymentBytes),
		"showMsg":       "false",
		"showData":      "true",
		"ab_config":     `{"key_onion":"C"}`,
	})

	o.user.SetClient(constants.SubmitOrder)
	resp, err := o.user.Client().Post(o.user.HeadersDeepCopy(), o.user.BodyDeepCopy())
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
