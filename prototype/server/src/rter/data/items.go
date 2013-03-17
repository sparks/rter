package data

import (
	"time"
)

type Item struct {
	ID     int64
	Type   string
	Author string

	ThumbnailURI string `json:",omitempty"`
	ContentURI   string `json:",omitempty"`
	UploadURI    string `json:",omitempty"`

	HasGeo  bool    `json:",omitempty"`
	Heading float64 `json:",omitempty"`
	Lat     float64 `json:",omitempty"`
	Lng     float64 `json:",omitempty"`

	StartTime time.Time `json:",omitempty"`
	StopTime  time.Time `json:",omitempty"`

	Terms []*Term `json:",omitempty"`
}

type ItemComment struct {
	ID     int64
	ItemID int64
	Author string

	Body string

	UpdateTime time.Time `json:",omitempty"`
}
