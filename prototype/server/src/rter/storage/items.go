package storage

import (
	"fmt"
	"rter/data"
)

func InsertItem(item *data.Item) error {
	_, err := db.Exec(
		"INSERT INTO Items (ID, Type, AuthorID, ThumbnailURI, ContentURI, UploadURI, HasGeo, Heading, Lat, Lng, StartTime, StopTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?",
		item.ID,
		item.Type,
		item.AuthorID,
		item.ThumbnailURI,
		item.ContentURI,
		item.UploadURI,
		item.HasGeo,
		item.Heading,
		item.Lat,
		item.Lng,
		item.StartTime,
		item.StopTime,
	)

	return err
}

func SelectItem(ID int) (*data.Item, error) {
	rows, err := db.Query("SELECT FROM Items WHERE ID=?", ID)

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, fmt.Errorf("No Item in storage with ID=%v", ID)
	}

	item := new(data.Item)

	err = rows.Scan(
		&item.ID,
		&item.Type,
		&item.AuthorID,
		&item.ThumbnailURI,
		&item.ContentURI,
		&item.UploadURI,
		&item.HasGeo,
		&item.Heading,
		&item.Lat,
		&item.Lng,
		&item.StartTime,
		&item.StopTime,
	)

	if err != nil {
		return nil, err
	}

	return item, nil
}
