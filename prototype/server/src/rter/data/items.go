package data

import (
	"time"
	token "videoserver/auth"
)

type Item struct {
	ID     int64
	Type   string
	Author string

	ThumbnailURI string `json:",omitempty"`
	ContentURI   string `json:",omitempty"`
	UploadURI    string `json:",omitempty"`

	HasHeading bool    `json:",omitempty"`
	Heading    float64 `json:",omitempty"`

	HasGeo bool    `json:",omitempty"`
	Lat    float64 `json:",omitempty"`
	Lng    float64 `json:",omitempty"`

	Live      bool      `json:",omitempty"`
	StartTime time.Time `json:",omitempty"`
	StopTime  time.Time `json:",omitempty"`

	Terms []*Term `json:",omitempty"`

	Token *token.Token `json:",omitempty"`
}

type ItemComment struct {
	ID     int64
	ItemID int64
	Author string

	Body string

	UpdateTime time.Time `json:",omitempty"`
}

func (i *Item) AddTerm(term string, author string) {
	newTerm := new(Term)

	newTerm.Term = term
	newTerm.Author = author

	i.Terms = append(i.Terms, newTerm)
}
