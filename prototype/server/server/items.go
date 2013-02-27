package server

import (
	"time"
)

type Item struct {
	Id       int
	Type     string
	AuthorID int

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
	ID       int
	ItemID   int
	AuthorID int

	Body string

	CreateTime time.Time
}
