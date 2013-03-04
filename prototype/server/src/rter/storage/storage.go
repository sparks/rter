package storage

import (
	"database/sql"
	"errors"
	"log"
	"rter/data"
)

var (
	ErrZeroMatches         = errors.New("Query didn't match anything.")
	ErrUnsupportedDataType = errors.New("Storage doesn't support the given datatype.")
	ErrCannotDelete        = errors.New("Storage doesn't allow deleting that.")
)

func Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.Exec(query, args...)
}

func Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(query, args...)
}

func MustExec(query string, args ...interface{}) sql.Result {
	res, err := db.Exec(query, args...)
	if err != nil {
		log.Fatalf("Error on Exec %q: %v", query, err)
	}
	return res
}

func MustQuery(query string, args ...interface{}) *sql.Rows {
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatalf("Error on Query %q: %v", query, err)
	}
	return rows
}

func Insert(val interface{}) error {
	var (
		res sql.Result
		err error
	)

	switch v := val.(type) {
	case *data.Item:
		res, err = Exec(
			"INSERT INTO Items (Type, AuthorID, ThumbnailURI, ContentURI, UploadURI, HasGeo, Heading, Lat, Lng, StartTime, StopTime) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			v.Type,
			v.AuthorID,
			v.ThumbnailURI,
			v.ContentURI,
			v.UploadURI,
			v.HasGeo,
			v.Heading,
			v.Lat,
			v.Lng,
			v.StartTime.UTC(),
			v.StopTime.UTC(),
		)
	case *data.ItemComment:
		res, err = Exec(
			"INSERT INTO ItemComments (ItemID, AuthorID, Body, CreateTime) VALUES (?, ?, ?, ?)",
			v.ItemID,
			v.AuthorID,
			v.Body,
			v.CreateTime,
		)
	case *data.Term:
		res, err = Exec(
			"INSERT INTO Terms (Term, Automated, AuthorID, CreateTime) VALUES (?, ?, ?, ?)",
			v.Term,
			v.Automated,
			v.AuthorID,
			v.CreateTime.UTC(),
		)
	case *data.TermRanking:
		res, err = Exec(
			"INSERT INTO TermRankings (Term, Ranking, UpdateTime) VALUES (?, ?, ?)",
			v.Term,
			v.Ranking,
			v.UpdateTime.UTC(),
		)
	case *data.Role:
		res, err = Exec(
			"INSERT INTO Roles (Title, Permissions) VALUES (?, ?)",
			v.Title,
			v.Permissions,
		)
	case *data.User:
		res, err = Exec(
			"INSERT INTO Users (Username, Password, Salt, Role, TrustLevel, CreateTime) VALUES (?, ?, ?, ?, ?, ?)",
			v.Username,
			v.Password,
			v.Salt,
			v.Role,
			v.TrustLevel,
			v.CreateTime.UTC(),
		)
	case *data.UserDirection:
		res, err = Exec(
			"INSERT INTO UserDirections (UserID, LockUserID, Command, Heading, Lat, Lng, UpdateTime) VALUES (?, ?, ?, ?, ?, ?, ?)",
			v.UserID,
			v.LockUserID,
			v.Command,
			v.Heading,
			v.Lat,
			v.Lng,
			v.UpdateTime.UTC(),
		)
	default:
		return ErrUnsupportedDataType
	}

	if err != nil {
		return err
	}

	ID, err := res.LastInsertId()

	if err != nil {
		return err
	}

	switch v := val.(type) {
	case *data.Item:
		v.ID = ID
	case *data.ItemComment:
		v.ID = ID
	case *data.User:
		v.ID = ID
	}

	return nil
}

func Delete(val interface{}) error {
	var (
		res sql.Result
		err error
	)

	switch v := val.(type) {
	case *data.Item:
		res, err = Exec("DELETE FROM Items WHERE ID=?", v.ID)
	case *data.ItemComment:
		res, err = Exec("DELETE FROM ItemComments WHERE ID=?", v.ID)
	case *data.Term:
		res, err = Exec("DELETE FROM Terms WHERE Term=?", v.Term)
	case *data.TermRanking:
		// res, err = Exec("DELETE FROM TermRankings WHERE Term=?", v.Term)
		return ErrCannotDelete //DB will autodelete
	case *data.Role:
		res, err = Exec("DELETE FROM Roles WHERE Title=?", v.Title)
	case *data.User:
		res, err = Exec("DELETE FROM Users WHERE ID=?", v.ID)
	case *data.UserDirection:
		// res, err = Exec("DELETE FROM UserDirections WHERE UserID=?", v.UserID)
		return ErrCannotDelete //DB will autodelete
	default:
		return ErrUnsupportedDataType
	}

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
