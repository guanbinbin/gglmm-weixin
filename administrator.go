package account

import (
	"github.com/weihongguo/gglmm"
)

// Administrator 用户模型
type Administrator struct {
	gglmm.Model
	Status   int8   `json:"status"`
	Mobile   string `json:"mobile"`
	Password string `json:"-"`

	AdministratorInfo *AdministratorInfo `json:"administratorInfo"`
}

// ResponseKey --
func (administrator Administrator) ResponseKey() [2]string {
	return [...]string{"administrator", "administrators"}
}

// AdministratorInfo --
type AdministratorInfo struct {
	gglmm.Model
	AdministratorID int64  `json:"administratorId"`
	Nickname        string `json:"nickname"`
	AvatarURL       string `json:"avatarUrl"`
}

// ResponseKey --
func (administratorInfo AdministratorInfo) ResponseKey() [2]string {
	return [...]string{"administratorInfo", "administratorInfos"}
}

// AuthInfo --
func (administrator Administrator) AuthInfo() *AuthInfo {
	var nickName, avatarURL string
	if administrator.AdministratorInfo != nil {
		nickName = administrator.AdministratorInfo.Nickname
		avatarURL = administrator.AdministratorInfo.AvatarURL
	}
	return &AuthInfo{
		Type:      AuthTypeAdministrator,
		ID:        administrator.ID,
		Nickname:  nickName,
		AvatarURL: avatarURL,
	}
}
