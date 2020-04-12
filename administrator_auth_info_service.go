package account

import (
	"net/http"

	"github.com/weihongguo/gglmm"
)

// AdministratorAuthInfoService 登录服务
type AdministratorAuthInfoService struct {
	repository *gglmm.GormRepository
}

// NewAdministratorAuthInfoService 新建用户服务
func NewAdministratorAuthInfoService() *AdministratorAuthInfoService {
	return &AdministratorAuthInfoService{
		repository: gglmm.DefaultGormRepository(),
	}
}

// CustomActions --
func (service *AdministratorAuthInfoService) CustomActions() ([]*gglmm.HTTPAction, error) {
	actions := []*gglmm.HTTPAction{
		gglmm.NewHTTPAction("/administrator/auth-info", service.AuthInfo, "GET"),
	}
	return actions, nil
}

// RESTAction --
func (service *AdministratorAuthInfoService) RESTAction(action gglmm.RESTAction) (*gglmm.HTTPAction, error) {
	return nil, nil
}

// AuthInfo 信息
func (service *AdministratorAuthInfoService) AuthInfo(w http.ResponseWriter, r *http.Request) {
	administratorID, err := GetAuthID(r, AuthTypeAdministrator)
	if err != nil {
		gglmm.NewFailResponse("claims subject").WriteJSON(w)
		return
	}
	administrator := Administrator{}
	idRequest := gglmm.IDRequest{
		ID:       administratorID,
		Preloads: []string{"AdministratorInfo"},
	}
	if err := service.repository.Get(&administrator, idRequest); err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	gglmm.NewSuccessResponse().
		AddData("authInfo", administrator.AuthInfo()).
		WriteJSON(w)
}
