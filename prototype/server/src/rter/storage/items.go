package storage

import (
	"fmt"
	"rter/data"
)

func InsertItem(item *data.Item) error {
	res, err := db.Exec(
		"INSERT INTO Items (ID, Type, AuthorID, ThumbnailURI, ContentURI, UploadURI, HasGeo, Heading, Lat, Lng, StartTime, StopTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
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

	if err != nil {
		return err
	}

	ID, err := res.LastInsertId()

	if err != nil {
		return err
	}

	item.ID = ID

	return nil
}

func SelectItem(ID int64) (*data.Item, error) {
	rows, err := db.Query("SELECT * FROM Items WHERE ID=?", ID)

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

func DeleteItem(item *data.Item) error {
	res, err := db.Exec("DELETE FROM Items WHERE ID=?", item.ID)

	if err == nil {
		affected, _ := res.RowsAffected()
		if affected < 1 {
			return fmt.Errorf("No such Item in storage: %v", item.ID)
		}
	}

	return err
}

func InsertItemComment(comment *data.ItemComment) error {
	_, err := db.Exec(
		"INSERT INTO ItemComments (ID, ItemID, AuthorID, Body, CreateTime) VALUES (?, ?, ?, ?, ?)",
		comment.ID,
		comment.ItemID,
		comment.AuthorID,
		comment.Body,
		comment.CreateTime,
	)

	return err
}

func SelectItemComment(ID int64) (*data.ItemComment, error) {
	rows, err := db.Query("SELECT * FROM ItemComments WHERE ID=?", ID)

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, fmt.Errorf("No ItemComment in storage with ID=%v", ID)
	}

	comment := new(data.ItemComment)

	err = rows.Scan(
		&comment.ID,
		&comment.ItemID,
		&comment.AuthorID,
		&comment.Body,
		&comment.CreateTime,
	)

	if err != nil {
		return nil, err
	}

	return comment, nil
}

func DeleteItemComment(comment *data.ItemComment) error {
	res, err := db.Exec("DELETE FROM ItemComments WHERE ID=?", comment.ID)

	if err == nil {
		affected, _ := res.RowsAffected()
		if affected < 1 {
			return fmt.Errorf("No such ItemComment in storage: %v", comment.ID)
		}
	}

	return err
}
