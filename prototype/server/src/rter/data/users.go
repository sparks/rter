package data

import (
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
