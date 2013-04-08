package data

import (
	"time"
)

type Term struct {
	Term string

	Count int //Not in DB for convenience oly

	Automated bool   `json:",omitempty"` //Marks if the autotagging is being performed for this tag
	Author    string //Tied to User.Username in DB

	UpdateTime time.Time `json:",omitempty"`
}

func (t *Term) CRUDPrefix() string {
	return "taxonomy"
}

func (t *Term) CRUDPath() string {
	return t.CRUDPrefix() + "/" + t.Term
}

type TermRelationship struct {
	Term   string //Tied to Term.Term in DB
	ItemID int64  //Tied to Item.ID in DB
}

type TermRanking struct {
	Term    string //Tied to Term.Term in DB
	Ranking string //JSON representation of the ranking of Items which has this Term.

	UpdateTime time.Time `json:",omitempty"`
}

func (r *TermRanking) CRUDPrefix() string {
	return "taxonomy/" + r.Term + "/ranking"
}

func (r *TermRanking) CRUDPath() string {
	return r.CRUDPrefix()
}
