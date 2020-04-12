package account

import "github.com/weihongguo/gglmm"

// WechatMiniProgramUser --
type WechatMiniProgramUser struct {
	gglmm.Model
	UserID     int64  `json:"userId"`
	OpenID     string `json:"openId"`
	UnionID    string `json:"unionId"`
	SessionKey string `json:"sessionKey"`
	Nickname   string `json:"nickname"`
	AvatarURL  string `json:"avatarUrl"`
	Gender     int8   `json:"gender"`
	Country    string `json:"country"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Language   string `json:"language"`
}

// ResponseKey --
func (wechatUser WechatMiniProgramUser) ResponseKey() [2]string {
	return [...]string{"wechatMiniProgramUser", "wechatMiniProgramUsers"}
}
