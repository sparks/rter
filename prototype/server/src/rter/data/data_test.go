package data

import (
	"testing"
)

func TestCRUDable(test *testing.T) {
	i := new(Item)
	i.ID = 5

	if i.CRUDPrefix() != "items" {
		test.Error("Wrong item CRUDPrefix")
	}

	if i.CRUDPath() != "items/5" {
		test.Error("Wrong item CRUDPath")
	}

	c := new(ItemComment)
	c.ItemID = 9

	if c.CRUDPrefix() != "items/9/comments" {
		test.Error("Wrong comment CRUDPrefix")
	}

	if c.CRUDPath() != "items/9/comments" {
		test.Error("Wrong comment CRUDPath")
	}

	t := new(Term)
	t.Term = "banana"

	if t.CRUDPrefix() != "taxonomy" {
		test.Error("Wrong taxonomy CRUDPrefix")
	}

	if t.CRUDPath() != "taxonomy/banana" {
		test.Error("Wrong taxonomy CRUDPath")
	}

	r := new(TermRanking)
	r.Term = "orapple"

	if r.CRUDPrefix() != "taxonomy/orapple/ranking" {
		test.Error("Wrong ranking CRUDPrefix")
	}

	if r.CRUDPath() != "taxonomy/orapple/ranking" {
		test.Error("Wrong ranking CRUDPath")
	}

	u := new(User)
	u.Username = "janice"

	if u.CRUDPrefix() != "users" {
		test.Error("Wrong user CRUDPrefix")
	}

	if u.CRUDPath() != "users/janice" {
		test.Error("Wrong user CRUDPath")
	}

	d := new(UserDirection)
	d.Username = "yollanda"

	if d.CRUDPrefix() != "users/yollanda/direction" {
		test.Error("Wrong direction CRUDPrefix")
	}

	if d.CRUDPath() != "users/yollanda/direction" {
		test.Error("Wrong direction CRUDPath")
	}

	o := new(Role)
	o.Title = "management"

	if o.CRUDPrefix() != "roles" {
		test.Error("Wrong role CRUDPrefix")
	}

	if o.CRUDPath() != "roles/management" {
		test.Error("Wrong role CRUDPath")
	}

}

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
