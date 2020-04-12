package main

import (
	"net/http"

	"github.com/weihongguo/gglmm"
	account "github.com/weihongguo/gglmm-account"
	redis "github.com/weihongguo/gglmm-redis"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// ExampleService --
type ExampleService struct {
}

// CustomActions --
func (service *ExampleService) CustomActions() ([]*gglmm.HTTPAction, error) {
	actions := []*gglmm.HTTPAction{
		gglmm.NewHTTPAction("/test", service.Test, "GET"),
	}
	return actions, nil
}

// RESTAction --
func (service *ExampleService) RESTAction(restAction gglmm.RESTAction) (*gglmm.HTTPAction, error) {
	return nil, nil
}

// Test --
func (service *ExampleService) Test(w http.ResponseWriter, r *http.Request) {
	id, err := account.GetAuthID(r, account.AuthTypeUser)
	if err != nil {
		gglmm.NewFailResponse("claims subject").WriteJSON(w)
		return
	}
	gglmm.NewSuccessResponse().
		AddData("userId", id).
		WriteJSON(w)
}

func main() {
	gglmm.RegisterGormRepository("mysql", "example:123456@(127.0.0.1:3306)/example?charset=utf8mb4&parseTime=true&loc=UTC", 10, 5, 600)
	defer gglmm.CloseGormRepository()

	cacher := redis.NewCacher("tcp", "127.0.0.1:6379", 5, 10, 3, 10)
	defer cacher.Close()
	gglmm.RegisterCacher(cacher)

	gglmm.RegisterBasePath("/api/example")

	gglmm.RegisterHTTPHandler(account.NewAdministratorLoginService(31536000, "administrator"), "")
	gglmm.RegisterHTTPHandler(account.NewAdministratorAuthInfoService(), "").
		Middleware(account.JWTAuthMiddleware([]string{"administrator"}))

	gglmm.RegisterHTTPHandler(account.NewUserLoginService(31536000, "user"), "")
	gglmm.RegisterHTTPHandler(account.NewUserAuthInfoService(), "").
		Middleware(account.JWTAuthMiddleware([]string{"user"}))

	gglmm.RegisterHTTPHandler(account.NewWechatMiniProgramLoginService("appID", "appSecret", 31536000, "user"), "")

	gglmm.RegisterHTTPHandler(account.NewWechatMiniProgramUserInfoService(31536000, "user"), "").
		Middleware(account.JWTAuthMiddleware([]string{"user"}))

	gglmm.RegisterHTTPHandler(&ExampleService{}, "/example").
		Middleware(account.JWTAuthMiddleware([]string{"administrator", "user"}))

	gglmm.ListenAndServe(":10000")
}
