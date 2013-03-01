package data

import (
	"time"
)

type Term struct {
	ID   int64
	Term string

	Automated bool
	AuthorID  int64

	CreateTime time.Time
}

type TermRanking struct {
	TermID  int64
	Ranking string

	UpdateTime time.Time
}
