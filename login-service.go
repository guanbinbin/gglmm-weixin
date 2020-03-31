package account

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"

	"github.com/weihongguo/gglmm"
)

// LoginService 登录服务
type LoginService struct {
	repository *gglmm.GormRepository
	modelType  reflect.Type

	jwtExpires int64
	jwtSecret  string
}

// NewLoginService 新建用户服务
func NewLoginService(model Authenticationable, jwtExpires int64, jwtSecret string) *LoginService {
	return &LoginService{
		repository: gglmm.DefaultGormRepository(),
		modelType:  reflect.TypeOf(model),
		jwtExpires: jwtExpires,
		jwtSecret:  jwtSecret,
	}
}

// NewLoginServiceConfig 新建用户服务
func NewLoginServiceConfig(model Authenticationable, jwtConfig ConfigJWT) *LoginService {
	if !jwtConfig.Check("all") {
		log.Printf("%+v\n", jwtConfig)
		log.Fatal("Config check invalid")
	}
	return NewLoginService(model, jwtConfig.Expires, jwtConfig.Secret)
}

// CustomActions --
func (service *LoginService) CustomActions() ([]*gglmm.HTTPAction, error) {
	actions := []*gglmm.HTTPAction{
		gglmm.NewHTTPAction("/login", service.Login, "POST"),
	}
	return actions, nil
}

// RESTAction --
func (service *LoginService) RESTAction(action gglmm.RESTAction) (*gglmm.HTTPAction, error) {
	return nil, nil
}

// Login 登录
func (service *LoginService) Login(w http.ResponseWriter, r *http.Request) {
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

	filterRequest := gglmm.FilterRequest{}
	filterRequest.AddFilter("mobile", gglmm.FilterOperateEqual, loginRequest.Mobile)
	user := reflect.New(service.modelType).Interface()
	if err := service.repository.Get(user, filterRequest); err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}

	authenticationable := user.(Authenticationable)

	if err := authenticationable.CheckPassword(loginRequest.Password); err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}

	token, jwtClaims, err := authenticationable.GenerateJWT(service.jwtExpires, service.jwtSecret)
	if err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}

	authInfo, err := authenticationable.GenerateAuthenticationInfo()
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
