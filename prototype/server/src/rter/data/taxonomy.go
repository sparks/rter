package data

import (
	"time"
)

type Term struct {
	Term string

	Count int

	Automated bool `json:",omitempty"`
	Author    string

	UpdateTime time.Time `json:",omitempty"`
}

type TermRelationship struct {
	Term   string
	ItemID int64
}

type TermRanking struct {
	Term    string
	Ranking string

	UpdateTime time.Time `json:",omitempty"`
}
