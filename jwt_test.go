package account

import (
	"net/http"
	"strconv"
	"testing"
)

type UserID int64

func (id UserID) ID() int64 {
	return int64(id)
}

func (id UserID) ToString() (string, error) {
	return strconv.FormatInt(int64(id), 10), nil
}

func (id UserID) FromString(value string) error {
	result, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		return err
	}
	id = UserID(result)
	return nil
}

func TestRequestContext(t *testing.T) {
	r1, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fail()
	}
	userID := UserID(1)
	jwtClaims, err := generateJWTClaims(&userID, JWTExpires)
	if err != nil {
		t.Fail()
	}
	r2 := setJWTClaimsToRequest(r1, jwtClaims)

	var userIDOut UserID
	err = GetJWTClaimsSubjectFromRequest(r2, &userIDOut)
	if err != nil {
		t.Fail()
	}
	if userIDOut != UserID(1) {
		t.Fail()
	}
}
