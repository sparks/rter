package server

import (
	"time"
)

type User struct {
	ID       int
	Username string
	Password string
	Salt     string

	Role       string
	TrustLevel int

	CreateTime time.Time
}

type UserDirection struct {
	UserID     int
	LockUserID int
	Command    string

	Heading float64
	Lat     float64
	Lng     float64

	UpdateTime time.Time
}

type Role struct {
	Role        string
	Permissions int
}
