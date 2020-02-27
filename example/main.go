package main

import (
	"github.com/weihongguo/gglmm"
	account "github.com/weihongguo/gglmm-account"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

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

	gglmm.ListenAndServe(":10000")
}
