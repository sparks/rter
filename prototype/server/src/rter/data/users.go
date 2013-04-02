package data

import (
	"crypto/md5"
	"fmt"
	"time"
)

type User struct {
	Username string
	Password string // Note this shouldn't be sent out, but we cannot set `json:"-"` because then it would never be received. Care must be taken not to send it elsewhere
	Salt     string `json:"-"` //Never send

	Role       string //Tied to Role.Title in DB
	TrustLevel int

	CreateTime time.Time `json:",omitempty"`
}

type UserDirection struct {
	Username     string //Tied to User.Username in DB
	LockUsername string `json:",omitempty"` //User current controlling this User
	Command      string `json:",omitempty"`

	Heading float64 `json:",omitempty"`
	Lat     float64 `json:",omitempty"`
	Lng     float64 `json:",omitempty"`

	UpdateTime time.Time `json:",omitempty"`
}

type Role struct {
	Title       string
	Permissions int
}

// Generates first a Salt for the user, then using that Salt generates a hash of the Password. The Password field should be set on the User when the func is called. The function will then replace the Password field with the hash and populate the Salt field.
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

// Check if the password p matches the stored hashed password for the User. This will use the Salt field of the User to generate the hash of the password p. This hash is then compared to the Password field of the User which should already contain the actual password hash.
func (user *User) Auth(p string) bool {
	hasher := md5.New()

	hasher.Write([]byte(user.Salt))
	hasher.Write([]byte(p))

	if fmt.Sprintf("%x", hasher.Sum(nil)) == user.Password {
		return true
	}

	return false
}
