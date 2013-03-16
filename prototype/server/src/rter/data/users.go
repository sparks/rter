package data

import (
	"crypto/md5"
	"fmt"
	"time"
)

type User struct {
	Username string `json:",omitempty"`
	Password string `json:"-"`
	Salt     string `json:"-"`

	Role       string `json:",omitempty"`
	TrustLevel int

	CreateTime time.Time `json:",omitempty"`
}

type UserDirection struct {
	Username     string
	LockUsername string
	Command      string `json:",omitempty"`

	Heading float64
	Lat     float64
	Lng     float64

	UpdateTime time.Time `json:",omitempty"`
}

type Role struct {
	Title       string `json:",omitempty"`
	Permissions int
}

func (user *User) HashAndSalt() {
	t := time.Now()
	hasher := md5.New()

	hasher.Write([]byte(fmt.Sprintf("%v", t.UnixNano())))
	user.Salt = string(hasher.Sum(nil))

	hasher = md5.New()

	hasher.Write([]byte(user.Salt))
	hasher.Write([]byte(user.Password))

	user.Password = string(hasher.Sum(nil))
}

func (user *User) Auth(p string) bool {
	hasher := md5.New()

	hasher.Write([]byte(user.Salt))
	hasher.Write([]byte(p))

	if string(hasher.Sum(nil)) == user.Password {
		return true
	}

	return false
}
