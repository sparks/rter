package storage

import (
	"database/sql"
	"errors"
	"log"
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
