package account

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/weihongguo/gglmm"
)

// AdministratorLoginService 登录服务
type AdministratorLoginService struct {
	repository *gglmm.GormRepository
	jwtExpires int64
	jwtSecret  string
}

// NewAdministratorLoginService 新建用户服务
func NewAdministratorLoginService(jwtExpires int64, jwtSecret string) *AdministratorLoginService {
	return &AdministratorLoginService{
		repository: gglmm.DefaultGormRepository(),
		jwtExpires: jwtExpires,
		jwtSecret:  jwtSecret,
	}
}

// NewAdministratorLoginServiceConfig 新建用户服务
func NewAdministratorLoginServiceConfig(jwtConfig ConfigJWT) *AdministratorLoginService {
	if !jwtConfig.Check("all") {
		log.Printf("%+v\n", jwtConfig)
		log.Fatal("Config check invalid")
	}
	return NewAdministratorLoginService(jwtConfig.Expires, jwtConfig.Secret)
}

// CustomActions --
func (service *AdministratorLoginService) CustomActions() ([]*gglmm.HTTPAction, error) {
	actions := []*gglmm.HTTPAction{
		gglmm.NewHTTPAction("/administrator/login", service.Login, "POST"),
	}
	return actions, nil
}

// RESTAction --
func (service *AdministratorLoginService) RESTAction(action gglmm.RESTAction) (*gglmm.HTTPAction, error) {
	return nil, nil
}

// Login 登录
func (service *AdministratorLoginService) Login(w http.ResponseWriter, r *http.Request) {
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
	administrator := Administrator{}
	filterRequest := gglmm.FilterRequest{
		Preloads: []string{"AdministratorInfo"},
	}
	filterRequest.AddFilter("mobile", gglmm.FilterOperateEqual, loginRequest.Mobile)
	if err := service.repository.Get(&administrator, filterRequest); err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	if err := ComparePassword(administrator.Password, loginRequest.Password); err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	authToken, jwtClaims, err := GenerateAuthToken(administrator, service.jwtExpires, service.jwtSecret)
	if err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	gglmm.NewSuccessResponse().
		AddData("authToken", authToken).
		AddData("authIssuedAt", jwtClaims.IssuedAt).
		AddData("authExpiresAt", jwtClaims.ExpiresAt).
		AddData("authInfo", administrator.AuthInfo()).
		WriteJSON(w)
}
