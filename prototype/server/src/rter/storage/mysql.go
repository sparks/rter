package server

import (
	"database/sql"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
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

	checkError(err)
}

func CloseMySQL() {
	db.Close()
}
