package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sync"

	"github.com/dingdong-grabber/pkg/constants"

	"github.com/dingdong-grabber/pkg/http"
	"k8s.io/klog"
)

type User struct {
	c         *http.Client
	addressId string
	headers   map[string]string
	body      url.Values
	mtx       sync.RWMutex
}

func NewDefaultUser() *User {
	return &User{
		c: &http.Client{},
	}
}

func (u *User) AddressId() string {
	u.mtx.RLock()
	defer u.mtx.RUnlock()
	return u.addressId
}

func (u *User) SetAddressId(addressId string) {
	u.mtx.RLock()
	defer u.mtx.RUnlock()
	u.addressId = addressId
}

func (u *User) SetClient(url string) {
	u.mtx.Lock()
	defer u.mtx.Unlock()
	u.c.Url = url
}

func (u *User) Client() *http.Client {
	u.mtx.RLock()
	defer u.mtx.RUnlock()
	return u.c
}

func (u *User) SetStationId(stationId string) {
	var (
		headers = u.Headers()
		body    = u.Body()
	)
	u.mtx.Lock()
	defer u.mtx.Unlock()
	headers["ddmc-station-id"] = stationId
	body["station_id"] = []string{stationId}
}

func (u *User) SetCityNumber(cityNumber string) {
	var (
		headers = u.Headers()
		body    = u.Body()
	)
	u.mtx.Lock()
	defer u.mtx.Unlock()
	headers["ddmc-city-number"] = cityNumber
	body["city_number"] = []string{cityNumber}
}

// GetDefaultAddr 获取默认地址 设置配送地址id，必须保证默认收获地址在上海且填写正确作为收获地址，请注意输出信息并确认
func (u *User) GetDefaultAddr() (*Address, error) {
	// body参数为共享，提交购物车时添加了products参数，可能会导致请求参数过长造成invalid character '<' looking for beginning of value，这里重新设置为空字符
	u.SetBody(map[string]string{
		"products":      "",
		"package_order": "",
		"packages":      "",
	})
	u.SetClient(constants.Address)
	resp, err := u.Client().Get(u.HeadersDeepCopy(), u.BodyDeepCopy())
	if err != nil {
		klog.Info(err.Error())
		return nil, err
	}

	bytes, _ := json.Marshal(resp.Data)
	var ads Addresses
	if err := json.Unmarshal(bytes, &ads); err != nil {
		klog.Infof("地址解析出错, 错误: %s", err.Error())
		return nil, err
	}

	if len(ads.ValidAddress) == 0 {
		klog.Info("请添加收货地址，并设置买菜地址为默认地址")
		return nil, errors.New("请添加收货地址，并设置买菜地址为默认地址")
	}

	for _, addr := range ads.ValidAddress {
		if addr.IsDefault {
			klog.Infof("1.默认收货地址：%s%s%s, 手机号: %s", addr.Location.Address, addr.Location.Name, addr.AddrDetail, addr.Mobile)
			klog.Infof("2.该地址对应站点名称为：%s", addr.StationInfo.Name)
			klog.Infof("3.设置买菜地址经度：%v", addr.Location.Location[0])
			klog.Infof("4.设置买菜地址纬度：%v", addr.Location.Location[1])

			u.SetHeaders(map[string]string{
				"longitude": fmt.Sprintf("%v", addr.Location.Location[0]),
				"latitude":  fmt.Sprintf("%v", addr.Location.Location[1]),
			})
			return &addr, nil
		}
	}

	klog.Info("请设置收货地址为默认地址")
	return nil, errors.New("请设置收货地址为默认地址")
}
