package wechat

// WechatResponse --
type WechatResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// WechatCode2SessionResponse --
type WechatCode2SessionResponse struct {
	WechatResponse
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
}
