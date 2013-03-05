package data

import (
	"time"
)

type Item struct {
	ID       int64
	Type     string
	AuthorID int64

	ThumbnailURI string
	ContentURI   string
	UploadURI    string

	HasGeo  bool
	Heading float64
	Lat     float64
	Lng     float64

	StartTime time.Time
	StopTime  time.Time
}

type ItemComment struct {
	ID       int64
	ItemID   int64
	AuthorID int64

	Body string

	UpdateTime time.Time
}
