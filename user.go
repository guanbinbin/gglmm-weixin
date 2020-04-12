package account

import (
	"github.com/weihongguo/gglmm"
)

// User 用户模型
type User struct {
	gglmm.Model
	Status   int8    `json:"status"`
	Mobile   *string `json:"mobile"`
	Password *string `json:"-"`

	UserInfo *UserInfo `json:"userInfo"`
}

// ResponseKey --
func (user User) ResponseKey() [2]string {
	return [...]string{"user", "users"}
}

// UserInfo --
type UserInfo struct {
	gglmm.Model
	UserID    int64  `json:"userId"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatarUrl"`
}

// ResponseKey --
func (userInfo UserInfo) ResponseKey() [2]string {
	return [...]string{"userInfo", "userInfos"}
}

// AuthInfo --
func (user User) AuthInfo() *AuthInfo {
	var nickName, avatarURL string
	if user.UserInfo != nil {
		nickName = user.UserInfo.Nickname
		avatarURL = user.UserInfo.AvatarURL
	}
	return &AuthInfo{
		Type:      AuthTypeUser,
		ID:        user.ID,
		Nickname:  nickName,
		AvatarURL: avatarURL,
	}
}
