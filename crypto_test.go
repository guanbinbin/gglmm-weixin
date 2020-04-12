package account

import "testing"

func TestBcrypt(t *testing.T) {
	test := "test"

	hashed, err := GeneratePassword(test)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = ComparePassword(hashed, test)
	if err != nil {
		t.Fatalf(err.Error())
	}
}
