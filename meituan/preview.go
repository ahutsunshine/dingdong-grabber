package meituan

import (
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
)

type PreViewBody struct {
	ActionSelect      string `json:"actionSelect"`
	AddressId         string `json:"addressId"`
	AllowZeroPay      string `json:"allowZeroPay"`
	CityId            string `json:"cityId"`
	CouponIds         []int  `json:"couponIds"`
	FromPoiId         string `json:"fromPoiId"`
	FromSource        string `json:"fromSource"`
	IsUseCard         string `json:"isUseCard"`
	Latitude          string `json:"latitude"`
	Longitude         string `json:"longitude"`
	PoiId             string `json:"poiId"`
	SelfLiftingMobile string `json:"selfLiftingMobile"`
	ShippingType      string `json:"shippingType"`
}

type skuProductInfo struct {
	Count               int    `json:"count"`
	FrozenTag           int    `json:"frozenTag"`
	GiftInfoTips        string `json:"giftInfoTips"`
	IsGift              bool   `json:"isGift"`
	ItemTag             int    `json:"itemTag"`
	MemberTag           int    `json:"memberTag"`
	Pic                 string `json:"pic"`
	PoiId               int    `json:"poiId"`
	ProcessingDetail    string `json:"processingDetail"`
	PromotionPrice      int    `json:"promotionPrice"`
	PromotionViewPrice  int    `json:"promotionViewPrice"`
	Scatter             bool   `json:"scatter"`
	SellPrice           int    `json:"sellPrice"`
	SellUnitViewName    string `json:"sellUnitViewName"`
	SellUnitViewPrice   int    `json:"sellUnitViewPrice"`
	SkuId               int    `json:"skuId"`
	SkuName             string `json:"skuName"`
	Spec                string `json:"spec"`
	SpuId               int    `json:"spuId"`
	SubTitle            string `json:"subTitle"`
	TempCount           int    `json:"tempCount"`
	TotalPromotionPrice int    `json:"totalPromotionPrice"`
	TotalSellPrice      int    `json:"totalSellPrice"`
	Unit                string `json:"unit"`
	ViewCount           string `json:"viewCount"`
}

type PreViewPackage struct {
	DateTime                 int              `json:"dateTime"`
	DeliveryCode             int              `json:"deliveryCode"`
	DeliveryEndTime          int              `json:"deliveryEndTime"`
	DeliveryLevel            int              `json:"deliveryLevel"`
	DeliveryStartTime        int              `json:"deliveryStartTime"`
	DeliveryType             int              `json:"deliveryType"`
	DeliveryUuid             string           `json:"deliveryUuid"`
	EarliestUseTime          bool             `json:"earliestUseTime"`
	EstimateTime             int              `json:"estimateTime"`
	EstimateTimeString       string           `json:"estimateTimeString"`
	HalfDayActivityTicketId  int              `json:"halfDayActivityTicketId"`
	PackageId                int              `json:"packageId"`
	PackageLabel             string           `json:"packageLabel"`
	PackageLabelId           int              `json:"packageLabelId"`
	PackageName              string           `json:"packageName"`
	PackageWeight            string           `json:"packageWeight"`
	PredictArrivateTime      int              `json:"predictArrivateTime"`
	SchemeId                 int              `json:"schemeId"`
	SelfLiftingTime          int              `json:"selfLiftingTime"`
	SelfLiftingTimeString    string           `json:"selfLiftingTimeString"`
	ShowDevliveryArrivedIcon bool             `json:"showDevliveryArrivedIcon"`
	SkuProductInfo           []skuProductInfo `json:"skuProductInfo"`
	SpeedyDelivery           struct {
		IsChooseSpeedy  bool `json:"isChooseSpeedy"`
		IsSupportSpeedy bool `json:"isSupportSpeedy"`
	}
	SupportSelfLifting bool `json:"supportSelfLifting"`
	TimeDeliveryPrice  int  `json:"timeDeliveryPrice"`
	TotalCount         int  `json:"totalCount"`
}

type PreViewResult struct {
	Packages []PreViewPackage
	totalPay int `json:"totalPay"`
}

func (cbody *PreViewBody) init(templates *UrlParams) {
	cbody.ActionSelect = "0"
	cbody.AddressId = templates.address_id
	cbody.AllowZeroPay = "true"
	cbody.CityId = "1"
	cbody.CouponIds = []int{-1}
	cbody.FromPoiId = "0"
	cbody.FromSource = "0"
	cbody.IsUseCard = "false"
	cbody.Latitude = "31.148712293836805"
	cbody.Longitude = "121.39732801649305"
	cbody.PoiId = templates.poi
	cbody.SelfLiftingMobile = ""
	cbody.ShippingType = "-1"
}

func (s *MeiTuanSession) PreView() (interface{}, error) {

	// 参数
	templates := s.getParamsTemplates()
	params := s.initUrlParams(&templates)

	//body
	cbody := new(PreViewBody)
	cbody.init(&templates)

	//request
	body, err := httpPost("https://mall.meituan.com/api/c/mallorder/preview", params, *cbody, s.Header)
	if err != nil {
		return nil, err
	}

	// response
	result := gjson.Parse(body)
	if result.Get("code").Num == 0 {
		//print(body)
		var infos []PreViewPackage
		print(result.Get("data.packageInfo").String())
		packageInfoList := []byte(result.Get("data.packageInfo").String())
		jsonRes := json.Unmarshal(packageInfoList, &infos)
		if jsonRes != nil {
			return nil, jsonRes
		}
		body := PreViewResult{
			Packages: infos,
			totalPay: int(result.Get("data.totalPay").Num),
		}
		return body, nil
	}
	return nil, errors.New(body)
}
