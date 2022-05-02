package meituan

type UrlParams struct {
	uuid       string `json:"uuid"`
	xuuid      string `json:"xuuid"`
	reqTraceID string `json:"__reqTraceID"`
	platform   string `json:"platform"`
	utm_medium string `json:"utm_medium"`
	brand      string `json:"brand"`
	tenantId   string `json:"tenantId"`
	utm_term   string `json:"utm_term"`
	cacheId    string `json:"cacheId"`
	abGroup    string `json:"abGroup"`
	poi        string `json:"poi"`
	stockPois  string `json:"stockPois"`
	ci         string `json:"ci"`
	bizId      string `json:"bizId"`
	openId     string `json:"openId"`
	address_id string `json:"address_id"`
	sysName    string `json:"sysName"`
	sysVerion  string `json:"sysVerion"`
	app_tag    string `json:"app_tag"`
	uci        string `json:"uci"`
	userid     string `json:"userid"`
}

func (s *MeiTuanSession) getParamsTemplates() UrlParams {
	u := new(UrlParams)
	u.reqTraceID = "be364734-82b5-ec7a-2f1c-700a8362932f"
	u.platform = "ios"
	u.utm_medium = "wxapp"
	u.brand = "xiaoxiangmaicai"
	u.tenantId = "1"
	u.utm_term = "5.32.6"
	u.cacheId = "1514760679758925874"
	u.abGroup = "3"
	u.poi = "217"
	u.stockPois = "217"
	u.ci = "1"
	u.bizId = "2"
	u.sysName = "iOS"
	u.sysVerion = "15.4"
	u.app_tag = "union"
	u.uci = "1"
	return *u
}
