package data

import (
	"time"
)

type Item struct {
	ID     int64
	Type   string `json:",omitempty"`
	Author string

	ThumbnailURI string `json:",omitempty"`
	ContentURI   string `json:",omitempty"`
	UploadURI    string `json:",omitempty"`

	HasGeo  bool
	Heading float64
	Lat     float64
	Lng     float64

	StartTime time.Time `json:",omitempty"`
	StopTime  time.Time `json:",omitempty"`
}

type ItemComment struct {
	ID     int64
	ItemID int64
	Author string

	Body string `json:",omitempty"`

	UpdateTime time.Time `json:",omitempty"`
}
