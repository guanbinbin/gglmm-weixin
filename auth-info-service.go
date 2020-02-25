package account

import (
	"net/http"
	"reflect"

	"github.com/weihongguo/gglmm"
)

// AuthInfoService 登录服务
type AuthInfoService struct {
	modelType  reflect.Type
	repository *gglmm.Repository
}

// NewAuthInfoService 新建用户服务
func NewAuthInfoService(model Authenticationable) *AuthInfoService {
	return &AuthInfoService{
		modelType:  reflect.TypeOf(model),
		repository: gglmm.NewRepository(model),
	}
}

// CustomActions --
func (service *AuthInfoService) CustomActions() ([]*gglmm.HTTPAction, error) {
	actions := []*gglmm.HTTPAction{
		gglmm.NewHTTPAction("/auth-info", service.AuthInfo, "GET"),
	}
	return actions, nil
}

// RESTAction --
func (service *AuthInfoService) RESTAction(action gglmm.RESTAction) (*gglmm.HTTPAction, error) {
	return nil, nil
}

// AuthInfo 信息
func (service *AuthInfoService) AuthInfo(w http.ResponseWriter, r *http.Request) {
	jwtUser := &JWTUser{}
	if err := GetJWTClaimsSubjectFromRequest(r, jwtUser); err != nil {
		gglmm.NewFailResponse("claims subject").WriteJSON(w)
		return
	}

	user := gglmm.ReflectNew(service.modelType)
	if err := service.repository.Get(jwtUser.UserID, user, nil); err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}

	authenticationable := user.(Authenticationable)

	authInfo, err := authenticationable.GenerateAuthenticationInfo()
	if err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}

	gglmm.NewSuccessResponse().
		AddData("authInfo", authInfo).
		WriteJSON(w)
}
