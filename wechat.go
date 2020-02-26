package account

import "github.com/weihongguo/gglmm"

// WechatAccount --
type WechatAccount struct {
	gglmm.Model
	OpenID     string `json:"openID"`
	UnionID    string `json:"unionID"`
	SessionKey string `json:"sessionKey"`
	Nickname   string `json:"nickname"`
	AvatarURL  string `json:"avatarURL"`
	Gender     int8   `json:"gender"`
	Country    string `json:"country"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Language   string `json:"language"`
}

// WechatUser 用户微信表
type WechatUser struct {
	WechatAccount
	Status int8  `json:"status" gorm:"not null;"`
	UserID int64 `json:"userID"`
}

// ResponseKey --
func (wechatUser WechatUser) ResponseKey() [2]string {
	return [...]string{"wechatUser", "wechatUsers"}
}
