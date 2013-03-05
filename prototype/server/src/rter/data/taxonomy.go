package data

import (
	"time"
)

type Term struct {
	Term string

	Automated bool
	AuthorID  int64

	UpdateTime time.Time
}

type TermRanking struct {
	Term    string
	Ranking string

	UpdateTime time.Time
}
