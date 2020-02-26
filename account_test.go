package account

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/weihongguo/gglmm"
)

var jwtSecret1 string = "test1"
var jwtSecret2 string = "test2"

type TestLoginService struct{}

func (service *TestLoginService) CustomActions() ([]*gglmm.HTTPAction, error) {
	actions := []*gglmm.HTTPAction{
		gglmm.NewHTTPAction(
			"/login1",
			func(w http.ResponseWriter, r *http.Request) {
				jwtUser := &JWTUser{}
				jwtUser.UserID = 1
				token, _, err := GenerateJWTToken(jwtUser, JWTExpires, jwtSecret1)
				if err != nil {
					gglmm.NewFailResponse(err.Error()).WriteJSON(w)
					return
				}
				gglmm.NewSuccessResponse().
					AddData("token", token).
					WriteJSON(w)
			},
			"POST",
		),
		gglmm.NewHTTPAction(
			"/login2",
			func(w http.ResponseWriter, r *http.Request) {
				jwtUser := &JWTUser{}
				jwtUser.UserID = 2
				token, _, err := GenerateJWTToken(jwtUser, JWTExpires, jwtSecret2)
				if err != nil {
					gglmm.NewFailResponse(err.Error()).WriteJSON(w)
					return
				}
				gglmm.NewSuccessResponse().
					AddData("token", token).
					WriteJSON(w)
			}, "POST",
		),
	}
	return actions, nil
}

func (service *TestLoginService) RESTAction(action gglmm.RESTAction) (*gglmm.HTTPAction, error) {
	return nil, nil
}

type TestBusinessService struct{}

func (service *TestBusinessService) CustomActions() ([]*gglmm.HTTPAction, error) {
	actions := []*gglmm.HTTPAction{
		gglmm.NewHTTPAction(
			"/test",
			func(w http.ResponseWriter, r *http.Request) {
				jwtUser := &JWTUser{}
				err := GetJWTClaimsSubjectFromRequest(r, jwtUser)
				if err != nil {
					gglmm.NewFailResponse("claims subject").
						WriteJSON(w)
					return
				}
				gglmm.NewSuccessResponse().
					AddData("userID", jwtUser.UserID).
					WriteJSON(w)
			},
			"GET",
		),
	}
	return actions, nil
}

func (service *TestBusinessService) RESTAction(action gglmm.RESTAction) (*gglmm.HTTPAction, error) {
	return nil, nil
}

func TestHTTP(t *testing.T) {

	gglmm.RegisterBasePath("/api")

	gglmm.RegisterHTTPHandler(&TestLoginService{}, "/test")

	gglmm.RegisterHTTPHandler(&TestBusinessService{}, "/business").
		Middleware(
			JWTAuthMiddleware([]string{
				jwtSecret1,
				jwtSecret2,
			}),
		)

	router := gglmm.GenerateHttpRouter()
	check(router, "/api/test/login1", 1, t)
	check(router, "/api/test/login2", 2, t)
}

func check(router *mux.Router, url string, result int64, t *testing.T) {

	wLogin := httptest.NewRecorder()
	rLogin, _ := http.NewRequest("POST", url, nil)

	router.ServeHTTP(wLogin, rLogin)

	responseLogin := gglmm.NewResponse()
	if err := json.Unmarshal(wLogin.Body.Bytes(), responseLogin); err != nil {
		t.Log(err)
		t.Fail()
	}
	if responseLogin.Code != http.StatusOK {
		t.Log(responseLogin.Code)
		t.Fail()
	}

	token := responseLogin.Data["token"].(string)
	t.Log(token)

	wBusiness := httptest.NewRecorder()
	rBusiness, _ := http.NewRequest("GET", "/api/business/test", nil)
	rBusiness.Header.Add("Authorization", "Bearer "+token)

	router.ServeHTTP(wBusiness, rBusiness)

	responseBusiness := gglmm.NewResponse()
	if err := json.Unmarshal(wBusiness.Body.Bytes(), responseBusiness); err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(responseBusiness)
	if responseBusiness.Code != http.StatusOK {
		t.Log(responseBusiness.Code)
		t.Fail()
	}
	userID := responseBusiness.Data["userID"].(float64)
	if userID != float64(result) {
		t.Log(userID)
		t.Fail()
	}
}
