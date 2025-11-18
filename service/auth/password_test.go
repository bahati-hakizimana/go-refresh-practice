package auth

import "testing"

func TestHashPassword(t *testing.T) {

	hash, err := HashPassword("password")

	if err != nil {
		t.Errorf("Error hashing password : %v", err)
	}

	if hash == "" {
		t.Error("Expected hash not to be empty")
	}

	if hash == "password" {
		t.Error("Hash should not be the same as the password")
	}

}

func TestComparePasswords(t *testing.T) {
	hash, err := HashPassword("password")

	if err != nil {
		t.Errorf("Error hashing password: %v", err)
	}

	if !ComparePassword(hash, []byte("password")) {
		t.Errorf("Expected passwords to not mutch hash")
	}

	if ComparePassword(hash, []byte("notpassword")){
		t.Errorf("Expected passwords to not mutch the hash")
	}


}