package account

import (
	"fmt"
	"net/http"

	"github.com/weihongguo/gglmm"
)

// JWTAuthMiddleware JWT通用认证中间件
func JWTAuthMiddleware(secrets []string) gglmm.Middleware {
	return gglmm.Middleware{
		Name: fmt.Sprintf("%s%+v", "JWTAuth", secrets),
		Func: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				for _, secret := range secrets {
					authInfo, _, err := ParseAuthToken(GetAuthToken(r), secret)
					if err == nil {
						r = RequestWithAuthInfo(r, authInfo)
						next.ServeHTTP(w, r)
						return
					}
				}
				gglmm.WriteUnauthorized(w)
			})
		},
	}
}
