package storage

import (
	"database/sql"
	"rter/data"
	"time"
)

func InsertItem(item *data.Item) error {
	ID, err := InsertEntry(
		"INSERT INTO Items (Type, AuthorID, ThumbnailURI, ContentURI, UploadURI, HasGeo, Heading, Lat, Lng, StartTime, StopTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		item.Type,
		item.AuthorID,
		item.ThumbnailURI,
		item.ContentURI,
		item.UploadURI,
		item.HasGeo,
		item.Heading,
		item.Lat,
		item.Lng,
		item.StartTime.UTC(),
		item.StopTime.UTC(),
	)

	if err != nil {
		return err
	}

	item.ID = ID

	return nil
}

func UpdateItem(item *data.Item) error {
	res, err := Exec(
		"UPDATE Items SET Type=?, AuthorID=?, ThumbnailURI=?, ContentURI=?, UploadURI=?, HasGeo=?, Heading=?, Lat=?, Lng=?, StartTime=?, StopTime=? WHERE ID=?",
		item.Type,
		item.AuthorID,
		item.ThumbnailURI,
		item.ContentURI,
		item.UploadURI,
		item.HasGeo,
		item.Heading,
		item.Lat,
		item.Lng,
		item.StartTime.UTC(),
		item.StopTime.UTC(),
		item.ID,
	)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if affected < 1 {
		return ErrZeroMatches
	}

	return nil
}

func SelectAllItems() ([]*data.Item, error) {
	rows, err := Query("SELECT * FROM Items")

	if err != nil {
		return nil, err
	}

	items := make([]*data.Item, 0)

	for rows.Next() {
		item := new(data.Item)
		err = scanItem(item, rows)

		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if len(items) == 0 {
		return nil, ErrZeroMatches
	}

	return items, nil
}

func SelectItem(item *data.Item) error {
	rows, err := Query("SELECT * FROM Items WHERE ID=?", item.ID)

	if err != nil {
		return err
	}

	if !rows.Next() {
		return ErrZeroMatches
	}

	err = scanItem(item, rows)

	return err
}

func DeleteItem(item *data.Item) error {
	return DeleteEntry("DELETE FROM Items WHERE ID=?", item.ID)
}

func scanItem(item *data.Item, rows *sql.Rows) error {
	var startTimeString, stopTimeString string

	err := rows.Scan(
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
		&startTimeString,
		&stopTimeString,
	)

	if err != nil {
		return err
	}

	startTime, err := time.Parse("2006-01-02 15:04:05", startTimeString) // this assumes UTC as timezone

	if err != nil {
		return err
	}

	item.StartTime = startTime

	stopTime, err := time.Parse("2006-01-02 15:04:05", stopTimeString) // this assumes UTC as timezone

	if err != nil {
		return err
	}

	item.StopTime = stopTime

	return nil
}

func InsertItemComment(comment *data.ItemComment) error {
	ID, err := InsertEntry(
		"INSERT INTO ItemComments (ItemID, AuthorID, Body, CreateTime) VALUES (?, ?, ?, ?)",
		comment.ItemID,
		comment.AuthorID,
		comment.Body,
		comment.CreateTime,
	)

	if err != nil {
		return err
	}

	comment.ID = ID

	return nil
}

func SelectItemComment(comment *data.ItemComment) error {
	rows, err := Query("SELECT * FROM ItemComments WHERE ID=?", comment.ID)

	if err != nil {
		return err
	}

	if !rows.Next() {
		return ErrZeroMatches
	}

	var createTimeString string

	err = rows.Scan(
		&comment.ID,
		&comment.ItemID,
		&comment.AuthorID,
		&comment.Body,
		&createTimeString,
	)

	if err != nil {
		return err
	}

	createTime, err := time.Parse("2006-01-02 15:04:05", createTimeString) // this assumes UTC as timezone

	if err != nil {
		return err
	}

	comment.CreateTime = createTime

	return nil
}

func DeleteItemComment(comment *data.ItemComment) error {
	return DeleteEntry("DELETE FROM ItemComments WHERE ID=?", comment.ID)
}
