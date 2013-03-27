package data

import (
	"crypto/md5"
	"fmt"
	"time"
)

type User struct {
	Username string
	Password string //TODO: Prevent from being sent
	Salt     string `json:"-"`

	Role       string
	TrustLevel int

	CreateTime time.Time `json:",omitempty"`
}

type UserDirection struct {
	Username     string
	LockUsername string `json:",omitempty"`
	Command      string `json:",omitempty"`

	Heading float64
	Lat     float64
	Lng     float64

	UpdateTime time.Time `json:",omitempty"`
}

type Role struct {
	Title       string
	Permissions int
}

func (user *User) HashAndSalt() {
	t := time.Now()
	hasher := md5.New()

	hasher.Write([]byte(fmt.Sprintf("%v", t.UnixNano())))
	user.Salt = fmt.Sprintf("%x", hasher.Sum(nil))

	hasher = md5.New()

	hasher.Write([]byte(user.Salt))
	hasher.Write([]byte(user.Password))

	user.Password = fmt.Sprintf("%x", hasher.Sum(nil))
}

func (user *User) Auth(p string) bool {
	hasher := md5.New()

	hasher.Write([]byte(user.Salt))
	hasher.Write([]byte(p))

	if fmt.Sprintf("%x", hasher.Sum(nil)) == user.Password {
		return true
	}

	return false
}
