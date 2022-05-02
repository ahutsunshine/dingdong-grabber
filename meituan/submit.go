package meituan

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"strconv"
)

type PackageInfo struct {
	DateTime          int    `json:"dateTime"`
	DeliverType       int    `json:"deliverType"`
	DeliveryEndTime   int    `json:"deliveryEndTime"`
	DeliveryLevel     int    `json:"deliveryLevel"`
	DeliveryStartTime int    `json:"deliveryStartTime"`
	DeliveryUuid      string `json:"deliveryUuid"`
	EstimateTime      int    `json:"estimateTime"`
	IsSpeedy          bool   `json:"isSpeedy"`
	PackageId         int    `json:"packageId"`
	SchemeId          int    `json:"schemeId"`
}

type SubmitBody struct {
	ActionSelect       int           `json:"actionSelect"`
	AddressId          int           `json:"addressId"`
	AllowZeroPay       bool          `json:"allowZeroPay"`
	AppId              string        `json:"appId"`
	CityId             int           `json:"cityId"`
	CouponAssignId     []int         `json:"couponAssignId"`
	CouponIds          []int         `json:"couponIds"`
	PackageInfoList    []PackageInfo `json:"packageInfo"`
	OpenId             string        `json:"openId"`
	PoiId              int           `json:"poiId"`
	Remark             string        `json:"remark"`
	SelfLiftingAddress string        `json:"selfLiftingAddress"`
	SelfLiftingMobile  string        `json:"selfLiftingMobile"`
	ShippingType       int           `json:"shippingType"`
	StockPois          []int         `json:"stockPois"`
	TotalPay           int           `json:"totalPay"`
}

func (cbody *SubmitBody) init(templates *UrlParams, s *MeiTuanSession, pView PreViewResult) {

	// 默认参数
	cbody.ActionSelect = 0
	cbody.AddressId, _ = strconv.Atoi(templates.address_id)
	cbody.AllowZeroPay = true
	cbody.AppId = s.UserInfo.AppId
	cbody.CityId = 1
	cbody.CouponAssignId = []int{}
	cbody.CouponIds = []int{}
	cbody.OpenId = templates.openId
	cbody.PoiId, _ = strconv.Atoi(templates.poi)
	cbody.Remark = ""
	cbody.SelfLiftingAddress = ""
	cbody.SelfLiftingMobile = ""
	cbody.ShippingType = 0
	stockPois, _ := strconv.Atoi(templates.stockPois)
	cbody.StockPois = []int{stockPois}
	// 及其总金额
	cbody.TotalPay = pView.totalPay

	// 包裹信息
	for _, Package := range pView.Packages {
		cbody.PackageInfoList = append(cbody.PackageInfoList, PackageInfo{
			DateTime:          Package.DateTime,
			DeliverType:       Package.DeliveryType,
			DeliveryEndTime:   (Package.DeliveryEndTime),
			DeliveryLevel:     Package.DeliveryLevel,
			DeliveryStartTime: (Package.DeliveryStartTime),
			DeliveryUuid:      Package.DeliveryUuid,
			EstimateTime:      (Package.EstimateTime),
			IsSpeedy:          false,
			PackageId:         Package.PackageId,
			SchemeId:          Package.SchemeId,
		})
	}

	fmt.Printf("订单总金额：%v\n", cbody.TotalPay)
}

func (s *MeiTuanSession) Submit(pViewResult interface{}) error {

	// 参数
	templates := s.getParamsTemplates()
	params := s.initUrlParams(&templates)

	//body
	cbody := new(SubmitBody)
	cbody.init(&templates, s, pViewResult.(PreViewResult))

	//request
	body, err := httpPost("https://mall.meituan.com/api/c/mallorder/submit", params, *cbody, s.Header)
	if err != nil {
		return err
	}

	// response
	result := gjson.Parse(body)
	if result.Get("code").Num == 0 {
		return nil
	}

	return errors.New(body)
}
