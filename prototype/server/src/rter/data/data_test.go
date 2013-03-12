package data

import (
	"testing"
)

func TestUserPass(t *testing.T) {
	user := new(User)
	password := "RightPassword"
	user.Password = password

	user.HashAndSalt()

	t.Log("Salt:", user.Salt)
	t.Log("Pass Hash:", user.Password)

	if !user.Auth(password) {
		t.Error("Failed to Auth user")
	}

	if user.Auth("WrongPassword") {
		t.Error("Shouldn't have Auth user")
	}
}
