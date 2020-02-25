package account

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/weihongguo/gglmm"
)

// WechatUserInfoService 微信登录服务
type WechatUserInfoService struct {
	jwtExpires int64
	jwtSecret  string
}

// NewWechatUserInfoService 新建用户服务
func NewWechatUserInfoService(jwtExpires int64, jwtSecret string) *WechatUserInfoService {
	return &WechatUserInfoService{
		jwtExpires: jwtExpires,
		jwtSecret:  jwtSecret,
	}
}

// NewWechatUserInfoServiceConfig 新建用户服务
func NewWechatUserInfoServiceConfig(config ConfigJWT) *WechatUserInfoService {
	return NewWechatUserInfoService(config.Expires, config.Secret)
}

// CustomActions --
func (service *WechatUserInfoService) CustomActions() ([]*gglmm.HTTPAction, error) {
	actions := []*gglmm.HTTPAction{
		gglmm.NewHTTPAction("/user-info", service.UserInfo, "PUT"),
	}
	return actions, nil
}

// RESTAction --
func (service *WechatUserInfoService) RESTAction(action gglmm.RESTAction) (*gglmm.HTTPAction, error) {
	return nil, nil
}

// UserInfo --
func (service *WechatUserInfoService) UserInfo(w http.ResponseWriter, r *http.Request) {
	jwtUser := JWTUser{}
	err := GetJWTClaimsSubjectFromRequest(r, &jwtUser)
	if err != nil {
		gglmm.NewFailResponse("claims subject").WriteJSON(w)
		return
	}

	decoder := json.NewDecoder(r.Body)
	wechatUserInfoRequest := WechatUserInfoRequest{}
	if err := decoder.Decode(&wechatUserInfoRequest); err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}

	if wechatUserInfoRequest.Check("raw") {
		service.rawUserInfo(w, jwtUser, wechatUserInfoRequest)
		return
	}

	if wechatUserInfoRequest.Check("encrypted") {
		service.encryptedUserInfo(w, jwtUser, wechatUserInfoRequest)
		return
	}

	gglmm.NewFailResponse("check fail").WriteJSON(w)
}

// rawUserInfo --
func (service *WechatUserInfoService) rawUserInfo(w http.ResponseWriter, jwtUser JWTUser, wechatUserInfoRequest WechatUserInfoRequest) {

	db := gglmm.GormDB()

	wechatUser := WechatUser{}
	if err := db.Where("user_id = ?", jwtUser.UserID).First(&wechatUser).Error; err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	if !wechatUserInfoRequest.CheckSignature(wechatUser.SessionKey) {
		gglmm.NewFailResponse("signature chect invalid").WriteJSON(w)
		return
	}

	user := User{}
	if err := db.First(&user, jwtUser.UserID).Error; err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}

	if err := wechatUpdateUser(&wechatUser, &user, &wechatUserInfoRequest.UserInfo); err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}

	authInfo, err := user.GenerateAuthenticationInfo()
	if err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}

	gglmm.NewSuccessResponse().
		AddData("authInfo", authInfo).
		WriteJSON(w)
}

// encryptedUserInfo 解析加密数据
// Session没有过期，下发新token
func (service *WechatUserInfoService) encryptedUserInfo(w http.ResponseWriter, jwtUser JWTUser, wechatUserInfoRequest WechatUserInfoRequest) {

	db := gglmm.GormDB()

	wechatUser := WechatUser{}
	if err := db.First(&wechatUser, jwtUser.UserID).Error; err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}

	wechatUserInfo, err := wechatUserInfoRequest.Decrypt(wechatUser.SessionKey)
	if err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}

	user := User{}
	if err = db.First(&user, jwtUser.UserID).Error; err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}

	if err = wechatUpdateUser(&wechatUser, &user, wechatUserInfo); err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}

	token, jwtClaims, err := user.GenerateJWT(service.jwtExpires, service.jwtSecret)
	if err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}

	authInfo, err := user.GenerateAuthenticationInfo()
	if err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}

	gglmm.NewSuccessResponse().
		AddData("authToken", token).
		AddData("authIssuedAt", jwtClaims.IssuedAt).
		AddData("authExpiresAt", jwtClaims.ExpiresAt).
		AddData("authInfo", authInfo).
		WriteJSON(w)
}

func wechatUpdateUser(wechatUser *WechatUser, user *User, wechatUserInfo *WechatUserInfo) error {
	tx := gglmm.GormBegin()
	if tx == nil {
		return errors.New("do not found")
	}

	wechatUserUpdates := map[string]interface{}{
		"nickname":   wechatUserInfo.Nickname,
		"avatar_url": wechatUserInfo.AvatarURL,
		"gender":     wechatUserInfo.Gender,
		"province":   wechatUserInfo.Province,
		"city":       wechatUserInfo.City,
		"country":    wechatUserInfo.Country,
		"language":   wechatUserInfo.Language,
	}
	if err := tx.Model(&wechatUser).Updates(wechatUserUpdates).Error; err != nil {
		tx.Rollback()

		return err
	}

	userUpdates := map[string]interface{}{
		"nickname":   wechatUserInfo.Nickname,
		"avatar_url": wechatUserInfo.AvatarURL,
	}
	if err := tx.Model(&user).Updates(userUpdates).Error; err != nil {
		tx.Rollback()

		return err
	}

	tx.Commit()

	return nil
}
