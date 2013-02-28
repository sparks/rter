package storage

import (
	"database/sql"
	"log"
)

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
