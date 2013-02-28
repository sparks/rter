package server

import (
	"time"
)

type Taxonomy struct {
	ID   int
	Term string

	Automated bool
	AuthorID  int

	CreateTime time.Time
}

type TaxonomyRanking struct {
	ID         int
	RawRanking string

	TaxonomyID int

	UpdateTime time.Time
}
