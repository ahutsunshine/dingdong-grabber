package meituan

type Header struct {
	Host         string `json:"Host"`
	OpenId       string `json:"openId"`
	ContentType  string `json:"content-type"`
	Traceids     string `json:"traceids"`
	ReqOfMaicai  string `json:"req_of_maicai"`
	T            string `json:"t"`
	OpenIdCipher string `json:"openIdCipher"`
	UserAgent    string `json:"User-Agent"`
	Referer      string `json:"Referer"`
}

func newDefaultHeader() Header {
	h := new(Header)
	h.Host = "mall.meituan.com"
	h.OpenId = ""
	h.ContentType = "application/json"
	h.Traceids = "5fc24523"
	h.ReqOfMaicai = "1"
	h.T = ""
	h.UserAgent = "User-Agent: Mozilla/5.0 (iPhone; CPU iPhone OS 15_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.20(0x18001428) NetType/4G Language/zh_CN"
	h.OpenIdCipher = ""
	h.Referer = "https://servicewechat.com/wx92916b3adca84096/227/page-frame.html"
	return *h
}
