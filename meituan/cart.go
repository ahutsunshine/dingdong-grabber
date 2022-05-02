package meituan

import (
	"errors"
	"github.com/tidwall/gjson"
)

type sku struct {
	skuId string `json:"skuId"`
}

type cartBody struct {
	CartOpType   string      `json:"cartOpType"`
	CartOpSource string      `json:"cartOpSource"`
	OpTarget     interface{} `json:"opTarget"`
	PoiId        string      `json:"poiId"`
	ShippingType string      `json:"shippingType"`
}

func (cbody *cartBody) init(templates *UrlParams) {
	cbody.CartOpSource = "CART"
	cbody.CartOpType = "REFRESH"
	cbody.OpTarget = map[string]interface{}{
		"opTargets": []int{},
	}
	cbody.PoiId = templates.poi
	cbody.ShippingType = "0"
}

func (s *MeiTuanSession) CheckCart() error {

	// 参数
	templates := s.getParamsTemplates()
	params := s.initUrlParams(&templates)

	//body
	cbody := new(cartBody)
	cbody.init(&templates)

	//request
	body, err := httpPost("https://mall.meituan.com/api/c/malluser/cart/v2/items", params, *cbody, s.Header)
	if err != nil {
		return err
	}

	// response
	result := gjson.Parse(body)
	if result.Get("code").Num == 0 {
		//print(body)
		return nil
	}

	return errors.New(body)
}
