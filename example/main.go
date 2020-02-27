package main

import (
	"net/http"

	"github.com/weihongguo/gglmm"
	account "github.com/weihongguo/gglmm-account"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type ExampleService struct {
}

// CustomActions --
func (service *ExampleService) CustomActions() ([]*gglmm.HTTPAction, error) {
	actions := []*gglmm.HTTPAction{
		gglmm.NewHTTPAction("/auth-info", service.AuthInfo, "GET"),
	}
	return actions, nil
}

// RESTAction --
func (service *ExampleService) RESTAction(restAction gglmm.RESTAction) (*gglmm.HTTPAction, error) {
	return nil, nil
}

// AuthInfo --
func (service *ExampleService) AuthInfo(w http.ResponseWriter, r *http.Request) {
	jwtUser := account.JWTUser{}
	err := account.GetJWTClaimsSubjectFromRequest(r, &jwtUser)
	if err != nil {
		gglmm.NewFailResponse("claims subject").WriteJSON(w)
		return
	}
	gglmm.NewSuccessResponse().WriteJSON(w)
}

func main() {
	gglmm.RegisterGormDB("mysql", "example:123456@(127.0.0.1:3306)/example?charset=utf8mb4&parseTime=true&loc=UTC", 10, 5, 600)
	defer gglmm.CloseGormDB()

	gglmm.RegisterRedisCacher("tcp", "127.0.0.1:6379", 10, 5, 3)
	defer gglmm.CloseRedisCacher()

	gglmm.RegisterBasePath("/api/example")

	gglmm.RegisterHTTPHandler(account.NewLoginService(account.Administrator{}, 31536000, "administrator"), "/administrator")

	gglmm.RegisterHTTPHandler(account.NewAuthInfoService(account.Administrator{}), "/administrator").
		Middleware(account.JWTAuthMiddleware([]string{"administrator"}))

	gglmm.RegisterHTTPHandler(account.NewLoginService(account.User{}, 31536000, "user"), "/user")

	gglmm.RegisterHTTPHandler(account.NewAuthInfoService(account.User{}), "/user").
		Middleware(account.JWTAuthMiddleware([]string{"user"}))

	gglmm.RegisterHTTPHandler(account.NewWechatLoginService("appID", "appSecret", 31536000, "user"), "/wechat")

	gglmm.RegisterHTTPHandler(account.NewWechatUserInfoService(31536000, "user"), "/wechat").
		Middleware(account.JWTAuthMiddleware([]string{"user"}))

	gglmm.RegisterHTTPHandler(&ExampleService{}, "/example").
		Middleware(account.JWTAuthMiddleware([]string{"administrator", "user"}))

	gglmm.ListenAndServe(":10000")
}
