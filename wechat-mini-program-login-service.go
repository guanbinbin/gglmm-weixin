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
	wechatMiniProgramCode2SessionURL = "https://api.weixin.qq.com/sns/jscode2session?grant_type=authorization_code"
)

// WechatMiniProgramLoginService 微信登录服务
type WechatMiniProgramLoginService struct {
	repository *gglmm.GormRepository
	appID      string
	appSecret  string
	jwtExpires int64
	jwtSecret  string
}

// NewWechatMiniProgramLoginService 新建用户服务
func NewWechatMiniProgramLoginService(appID string, appSecret string, jwtExpires int64, jwtSecret string) *WechatMiniProgramLoginService {
	return &WechatMiniProgramLoginService{
		repository: gglmm.DefaultGormRepository(),
		appID:      appID,
		appSecret:  appSecret,
		jwtExpires: jwtExpires,
		jwtSecret:  jwtSecret,
	}
}

// NewWechatMiniProgramLoginServiceConfig 新建用户服务
func NewWechatMiniProgramLoginServiceConfig(miniProgramConfig ConfigMiniProgram, jwtConfig ConfigJWT) *WechatMiniProgramLoginService {
	if !miniProgramConfig.Check() || !jwtConfig.Check("all") {
		log.Printf("%+v %+v\n", miniProgramConfig, jwtConfig)
		log.Fatal("Config check invalid")
	}
	return NewWechatMiniProgramLoginService(miniProgramConfig.AppID, miniProgramConfig.AppSecret, jwtConfig.Expires, jwtConfig.Secret)
}

// CustomActions --
func (service *WechatMiniProgramLoginService) CustomActions() ([]*gglmm.HTTPAction, error) {
	actions := []*gglmm.HTTPAction{
		gglmm.NewHTTPAction("/wechat-mini-program/login", service.MiniProgramLogin, "POST"),
	}
	return actions, nil
}

// RESTAction --
func (service *WechatMiniProgramLoginService) RESTAction(action gglmm.RESTAction) (*gglmm.HTTPAction, error) {
	return nil, nil
}

// MiniProgramLogin 登录-微信
// Session已经过期或第一次登录，下发token
func (service *WechatMiniProgramLoginService) MiniProgramLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	loginRequest := WechatMiniProgramLoginRequest{}
	if err := decoder.Decode(&loginRequest); err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	if !loginRequest.Check() {
		gglmm.NewFailResponse("请求错误").WriteJSON(w)
		return
	}

	code2SessionRespons, err := service.code2Session(loginRequest.Code)
	if err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}

	wechatUser := WechatMiniProgramUser{}
	filterRequest := gglmm.FilterRequest{}
	filterRequest.AddFilter("open_id", gglmm.FilterOperateEqual, code2SessionRespons.OpenID)
	if err := service.repository.Get(&wechatUser, filterRequest); err != nil && err != gglmm.ErrGormRecordNotFound {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}

	user := User{}
	if service.repository.NewRecord(wechatUser) {
		tx := service.repository.Begin()

		user.Nickname = "Wechat"
		if err = tx.Create(&user).Error; err != nil {
			tx.Rollback()

			gglmm.NewFailResponse("保存账号信息失败").WriteJSON(w)
			return
		}

		wechatUser.OpenID = code2SessionRespons.OpenID
		wechatUser.SessionKey = code2SessionRespons.SessionKey
		wechatUser.UnionID = code2SessionRespons.UnionID
		wechatUser.UserID = user.ID
		if err = tx.Create(&wechatUser).Error; err != nil {
			tx.Rollback()

			gglmm.NewFailResponse("保存账号信息失败").WriteJSON(w)
			return
		}

		tx.Commit()
	} else {
		idRequest := gglmm.IDRequest{
			ID: wechatUser.UserID,
		}
		if err = service.repository.Get(&user, idRequest); err != nil {
			gglmm.NewFailResponse(err.Error()).WriteJSON(w)
			return
		}

		wechatUser.SessionKey = code2SessionRespons.SessionKey
		if err = service.repository.Update(wechatUser, wechatUser.ID); err != nil {
			gglmm.NewFailResponse("更新SessionKey失败").WriteJSON(w)
			return
		}
	}

	if service.repository.NewRecord(user) {
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

func (service *WechatMiniProgramLoginService) code2Session(code string) (*WechatCode2SessionResponse, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{
		Timeout:   15 * time.Second,
		Transport: transport,
	}

	url := wechatMiniProgramCode2SessionURL + "&appid=" + service.appID + "&secret=" + service.appSecret + "&js_code=" + code
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
