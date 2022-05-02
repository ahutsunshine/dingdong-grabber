package meituan

type UserInfo struct {
	OpenId       string `json:"openId"`
	T            string `json:"t"`
	OpenIdCipher string `json:"openIdCipher"`
	UUID         string `json:"uuid"`
	AddressId    string `json:"address_id"`
	UserId       string `json:"userid"`
	AppId        string `json:"appId"`
}

var user = new(UserInfo)

func GetUserInfo() *UserInfo {
	return user
}
