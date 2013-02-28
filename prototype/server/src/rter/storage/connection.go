package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	"log"
)

var db *sql.DB

func OpenStorage() {
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

func CloseStorage() {
	db.Close()
}
