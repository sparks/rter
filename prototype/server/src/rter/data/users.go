package data

import (
	"crypto/md5"
	"fmt"
	"time"
)

type User struct {
	ID       int64
	Username string
	Password string
	Salt     string

	Role       string
	TrustLevel int

	CreateTime time.Time
}

type UserDirection struct {
	UserID     int64
	LockUserID int64
	Command    string

	Heading float64
	Lat     float64
	Lng     float64

	UpdateTime time.Time
}

type Role struct {
	Title       string
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
