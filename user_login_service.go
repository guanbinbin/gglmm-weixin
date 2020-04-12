package account

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/weihongguo/gglmm"
)

// UserLoginService 登录服务
type UserLoginService struct {
	repository *gglmm.GormRepository
	jwtExpires int64
	jwtSecret  string
}

// NewUserLoginService 新建用户服务
func NewUserLoginService(jwtExpires int64, jwtSecret string) *UserLoginService {
	return &UserLoginService{
		repository: gglmm.DefaultGormRepository(),
		jwtExpires: jwtExpires,
		jwtSecret:  jwtSecret,
	}
}

// NewUserLoginServiceConfig 新建用户服务
func NewUserLoginServiceConfig(jwtConfig ConfigJWT) *UserLoginService {
	if !jwtConfig.Check("all") {
		log.Printf("%+v\n", jwtConfig)
		log.Fatal("Config check invalid")
	}
	return NewUserLoginService(jwtConfig.Expires, jwtConfig.Secret)
}

// CustomActions --
func (service *UserLoginService) CustomActions() ([]*gglmm.HTTPAction, error) {
	actions := []*gglmm.HTTPAction{
		gglmm.NewHTTPAction("/user/login", service.Login, "POST"),
	}
	return actions, nil
}

// RESTAction --
func (service *UserLoginService) RESTAction(action gglmm.RESTAction) (*gglmm.HTTPAction, error) {
	return nil, nil
}

// Login 登录
func (service *UserLoginService) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	loginRequest := LoginRequest{}
	if err := decoder.Decode(&loginRequest); err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	if !loginRequest.Check() {
		gglmm.NewFailResponse("请求错误").WriteJSON(w)
		return
	}
	user := User{}
	filterRequest := gglmm.FilterRequest{
		Preloads: []string{"UserInfo"},
	}
	filterRequest.AddFilter("mobile", gglmm.FilterOperateEqual, loginRequest.Mobile)
	if err := service.repository.Get(&user, filterRequest); err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	if err := ComparePassword(*user.Password, loginRequest.Password); err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	authToken, jwtClaims, err := GenerateAuthToken(user, service.jwtExpires, service.jwtSecret)
	if err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	gglmm.NewSuccessResponse().
		AddData("authToken", authToken).
		AddData("authIssuedAt", jwtClaims.IssuedAt).
		AddData("authExpiresAt", jwtClaims.ExpiresAt).
		AddData("authInfo", user.AuthInfo()).
		WriteJSON(w)
}
