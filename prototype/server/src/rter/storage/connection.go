package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	"log"
)

var db *sql.DB

func OpenStorage(user, pass, protocol, addr, dbname string) {
	netAddr := fmt.Sprintf("%s(%s)", protocol, addr)
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
