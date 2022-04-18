package order

import (
	"testing"

	"github.com/dingdong-grabber/pkg/user"
)

const (
	deviceId    = ""
	cookie      = ""
	uid         = ""
	userAgent   = ""
	sid         = ""
	deviceToken = ""
)

func TestOrder(t *testing.T) {
	//t.Skip("以下为购物流程，便于开发者理解")

	u := user.NewDefaultUser()
	// 1. 初始化用户必须的参数数据
	if err := u.LoadConfig(deviceId, cookie, uid, userAgent, sid, deviceToken); err != nil {
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
