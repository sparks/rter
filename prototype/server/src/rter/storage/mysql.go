package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	"log"
)

var db *sql.DB

func SetupMySQL() {
	user := "root"
	pass := ""
	prot := "tcp"
	addr := "localhost:3306"
	dbname := "rter_v1"

	netAddr := fmt.Sprintf("%s(%s)", prot, addr)
	dsn := fmt.Sprintf("%s:%s@%s/%s?charset=utf8", user, pass, netAddr, dbname)
	var err error
	db, err = sql.Open("mysql", dsn)

	if err != nil {
		log.Fatalf("Failed to open connection to database %v", err)
	}
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

func CloseMySQL() {
	db.Close()
}
