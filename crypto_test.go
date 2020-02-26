package account

import "testing"

func TestBcrypt(t *testing.T) {
	test := "test"

	hashed, err := BcryptGenerateFromPassword(test)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = BcryptCompareHashAndPassword(hashed, test)
	if err != nil {
		t.Fatalf(err.Error())
	}
}
