package account

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/weihongguo/gglmm"
)

const (
	wechatCode2SessionURL = "https://api.weixin.qq.com/sns/jscode2session?grant_type=authorization_code"
)

// WechatLoginService 微信登录服务
type WechatLoginService struct {
	appID      string
	appSecret  string
	jwtExpires int64
	jwtSecret  string
}

// NewWechatLoginService 新建用户服务
func NewWechatLoginService(appID string, appSecret string, jwtExpires int64, jwtSecret string) *WechatLoginService {
	return &WechatLoginService{
		appID:      appID,
		appSecret:  appSecret,
		jwtExpires: jwtExpires,
		jwtSecret:  jwtSecret,
	}
}

// NewWechatLoginServiceConfig 新建用户服务
func NewWechatLoginServiceConfig(miniProgramConfig ConfigMiniProgram, jwtConfig ConfigJWT) *WechatLoginService {
	if !miniProgramConfig.Check() || !jwtConfig.Check("all") {
		log.Printf("%+v %+v\n", miniProgramConfig, jwtConfig)
		log.Fatal("Config check invalid")
	}
	return NewWechatLoginService(miniProgramConfig.AppID, miniProgramConfig.AppSecret, jwtConfig.Expires, jwtConfig.Secret)
}

// CustomActions --
func (service *WechatLoginService) CustomActions() ([]*gglmm.HTTPAction, error) {
	actions := []*gglmm.HTTPAction{
		gglmm.NewHTTPAction("/login", service.Login, "POST"),
	}
	return actions, nil
}

// RESTAction --
func (service *WechatLoginService) RESTAction(action gglmm.RESTAction) (*gglmm.HTTPAction, error) {
	return nil, nil
}

// Login 登录-微信
// Session已经过期或第一次登录，下发token
func (service *WechatLoginService) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	wechatLoginRequest := WechatLoginRequest{}
	if err := decoder.Decode(&wechatLoginRequest); err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	if !wechatLoginRequest.Check() {
		gglmm.NewFailResponse("请求错误").WriteJSON(w)
		return
	}

	code2SessionRespons, err := service.code2Session(wechatLoginRequest.Code)
	if err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}

	db := gglmm.GormDB()

	wechatUser := WechatUser{}
	if err := db.Where("open_id = ?", code2SessionRespons.OpenID).First(&wechatUser).Error; err != nil && err != gglmm.ErrGormRecordNotFound {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}

	user := User{}
	if gglmm.GormNewRecord(wechatUser) {
		tx := gglmm.GormBegin()

		user.Nickname = "Wechat"
		if err = tx.Create(&user).Error; err != nil {
			tx.Rollback()

			gglmm.NewFailResponse("保存账号信息失败").WriteJSON(w)
			return
		}

		wechatUser.OpenID = code2SessionRespons.OpenID
		wechatUser.SessionKey = code2SessionRespons.SessionKey
		wechatUser.UnionID = code2SessionRespons.UnionID
		wechatUser.Status = gglmm.StatusValid.Value
		wechatUser.UserID = user.ID
		if err = tx.Create(&wechatUser).Error; err != nil {
			tx.Rollback()

			gglmm.NewFailResponse("保存账号信息失败").WriteJSON(w)
			return
		}

		tx.Commit()
	} else {
		if err = db.First(&user, wechatUser.UserID).Error; err != nil {
			gglmm.NewFailResponse(err.Error()).WriteJSON(w)
			return
		}

		if err = db.Model(&wechatUser).Update("session_key", code2SessionRespons.SessionKey).Error; err != nil {
			gglmm.NewFailResponse("更新SessionKey失败").WriteJSON(w)
			return
		}
	}

	if gglmm.GormNewRecord(user) {
		gglmm.NewFailResponse("用户不存在").WriteJSON(w)
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

func (service *WechatLoginService) code2Session(code string) (*WechatCode2SessionResponse, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{
		Timeout:   15 * time.Second,
		Transport: transport,
	}

	url := wechatCode2SessionURL + "&appid=" + service.appID + "&secret=" + service.appSecret + "&js_code=" + code
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	log.Printf("code2session start %s\n", time.Now())
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	log.Printf("code2session end   %s\n", time.Now())
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	response := &WechatCode2SessionResponse{}
	if err := decoder.Decode(response); err != nil {
		return nil, err
	}
	if response.ErrCode != 0 {
		return nil, errors.New(response.ErrMsg)
	}
	return response, nil
}
