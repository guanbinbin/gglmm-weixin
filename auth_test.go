package account

import (
	"net/http"
	"testing"
)

type UserID int64

func (id *UserID) AuthInfo() *AuthInfo {
	return &AuthInfo{
		Type: "test",
		ID:   int64(*id),
	}
}

func TestRequestContext(t *testing.T) {
	userID := UserID(1)
	authToken, jwtClaims, err := GenerateAuthToken(&userID, JWTExpires, "test")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(authToken, jwtClaims)

	r1, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	r1.Header.Add("Authorization", "Bearer "+authToken)

	authToken = GetAuthToken(r1)
	authInfo, jwtClaims, err := ParseAuthToken(authToken, "test")
	t.Log(authInfo, jwtClaims)

	r2 := RequestWithAuthInfo(r1, userID.AuthInfo())
	id, err := GetAuthID(r2, "test")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if id != 1 {
		t.Log(id)
		t.Fail()
	}
}
