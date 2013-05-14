package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net"
)

var db *sql.DB

// Open a connection the storage solution: a MySQL db
func OpenStorage(user, pass, protocol, addr, dbname string) error {
	netAddr := fmt.Sprintf("%s(%s)", protocol, addr)
	dsn := fmt.Sprintf("%s:%s@%s/%s?charset=utf8", user, pass, netAddr, dbname)
	var err error
	db, err = sql.Open("mysql", dsn)

	if err != nil {
		return err
	}

	con, err := net.Dial(protocol, addr)

	if err == nil {
		con.Close()
	}

	return err
}

// Close connection the storage solution.
func CloseStorage() {
	db.Close()
}
