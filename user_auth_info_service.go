package account

import (
	"net/http"

	"github.com/weihongguo/gglmm"
)

// UserAuthInfoService 登录服务
type UserAuthInfoService struct {
	repository *gglmm.GormRepository
}

// NewUserAuthInfoService 新建用户服务
func NewUserAuthInfoService() *UserAuthInfoService {
	return &UserAuthInfoService{
		repository: gglmm.DefaultGormRepository(),
	}
}

// CustomActions --
func (service *UserAuthInfoService) CustomActions() ([]*gglmm.HTTPAction, error) {
	actions := []*gglmm.HTTPAction{
		gglmm.NewHTTPAction("/user/auth-info", service.AuthInfo, "GET"),
	}
	return actions, nil
}

// RESTAction --
func (service *UserAuthInfoService) RESTAction(action gglmm.RESTAction) (*gglmm.HTTPAction, error) {
	return nil, nil
}

// AuthInfo 信息
func (service *UserAuthInfoService) AuthInfo(w http.ResponseWriter, r *http.Request) {
	userID, err := GetAuthID(r, AuthTypeUser)
	if err != nil {
		gglmm.NewFailResponse("claims subject").WriteJSON(w)
		return
	}
	user := User{}
	idRequest := gglmm.IDRequest{
		ID:       userID,
		Preloads: []string{"UserInfo"},
	}
	if err := service.repository.Get(&user, idRequest); err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	gglmm.NewSuccessResponse().
		AddData("authInfo", user.AuthInfo()).
		WriteJSON(w)
}
