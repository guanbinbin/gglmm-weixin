package wechat

// MiniProgramResponse --
type MiniProgramResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// MiniProgramCode2SessionResponse --
type MiniProgramCode2SessionResponse struct {
	MiniProgramResponse
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
}
