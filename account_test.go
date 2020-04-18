package account

/*
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

type Test struct {
	ID int64
}

func (test Test) AuthInfo() *AuthInfo {
	return &AuthInfo{
		Type:      "test",
		ID:        test.ID,
		Nickname:  "test",
		AvatarURL: "test",
	}
}

type TestLoginService struct{}

func (service *TestLoginService) CustomActions() ([]*gglmm.HTTPAction, error) {
	actions := []*gglmm.HTTPAction{
		gglmm.NewHTTPAction(
			"/login1",
			func(w http.ResponseWriter, r *http.Request) {
				test := Test{
					ID: 1,
				}
				authToken, _, err := GenerateAuthToken(test, JWTExpires, jwtSecret1)
				if err != nil {
					gglmm.FailResponse(err.Error()).JSON(w)
					return
				}
				gglmm.OkResponse().
					AddData("authToken", authToken).
					AddData("authInfo", test.AuthInfo()).
					JSON(w)
			},
			"POST",
		),
		gglmm.NewHTTPAction(
			"/login2",
			func(w http.ResponseWriter, r *http.Request) {
				test := Test{
					ID: 2,
				}
				authToken, _, err := GenerateAuthToken(test, JWTExpires, jwtSecret1)
				if err != nil {
					gglmm.FailResponse(err.Error()).JSON(w)
					return
				}
				gglmm.OkResponse().
					AddData("authToken", authToken).
					AddData("authInfo", test.AuthInfo()).
					JSON(w)
			}, "POST",
		),
	}
	return actions, nil
}

func (service *TestLoginService) Action(action string) (*gglmm.HTTPAction, error) {
	return nil, nil
}

type TestBusinessService struct{}

func (service *TestBusinessService) CustomActions() ([]*gglmm.HTTPAction, error) {
	actions := []*gglmm.HTTPAction{
		gglmm.NewHTTPAction(
			"/test",
			func(w http.ResponseWriter, r *http.Request) {
				id, err := GetAuthID(r, "test")
				if err != nil {
					gglmm.FailResponse("claims subject").JSON(w)
					return
				}
				gglmm.OkResponse().
					AddData("userID", id).
					JSON(w)
			},
			"GET",
		),
	}
	return actions, nil
}

func (service *TestBusinessService) Action(action string) (*gglmm.HTTPAction, error) {
	return nil, nil
}

func TestHTTP(t *testing.T) {

	gglmm.BasePath("/api")

	gglmm.HandleHTTP(&TestLoginService{}, "/test")

	gglmm.HandleHTTP(&TestBusinessService{}, "/business").
		Middleware(
			JWTAuthMiddleware([]string{
				jwtSecret1,
				jwtSecret2,
			}),
		)

	router := gglmm.handleHTTPRouter()
	check(router, "/api/test/login1", 1, t)
	check(router, "/api/test/login2", 2, t)
}

func check(router *mux.Router, url string, result int64, t *testing.T) {

	wLogin := httptest.NewRecorder()
	rLogin, _ := http.NewRequest("POST", url, nil)

	router.ServeHTTP(wLogin, rLogin)

	responseLogin := gglmm.OkResponse()
	if err := json.Unmarshal(wLogin.Body.Bytes(), responseLogin); err != nil {
		t.Log(err)
		t.Fail()
	}
	if responseLogin.Code != http.StatusOK {
		t.Log(responseLogin.Code)
		t.Fail()
	}

	authToken := responseLogin.Data["authToken"].(string)
	t.Log(authToken)

	wBusiness := httptest.NewRecorder()
	rBusiness, _ := http.NewRequest("GET", "/api/business/test", nil)
	rBusiness.Header.Add("Authorization", "Bearer "+authToken)

	router.ServeHTTP(wBusiness, rBusiness)

	responseBusiness := gglmm.OkResponse()
	if err := json.Unmarshal(wBusiness.Body.Bytes(), responseBusiness); err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Logf("%+v", responseBusiness)
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
*/
