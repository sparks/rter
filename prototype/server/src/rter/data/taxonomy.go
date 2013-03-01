package data

import (
	"time"
)

type TaxonomyTerm struct {
	ID   int64
	Term string

	Automated bool
	AuthorID  int64

	CreateTime time.Time
}

type TaxonomyTermRanking struct {
	TermID  int64
	Ranking string

	UpdateTime time.Time
}
