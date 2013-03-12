package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	"net"
)

var db *sql.DB

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

func CloseStorage() {
	db.Close()
}
