package data

import (
	"time"
)

type Taxonomy struct {
	ID   int64
	Term string

	Automated bool
	AuthorID  int64

	CreateTime time.Time
}

type TaxonomyRanking struct {
	ID         int64
	RawRanking string

	TaxonomyID int64

	UpdateTime time.Time
}
