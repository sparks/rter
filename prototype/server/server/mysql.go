package server

import (
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	"os"
)

var db mysql.Conn

func SetupMySQL() {
	db = mysql.New("tcp", "", "localhost:3306", "root", "", "rter")

	err := db.Connect()
	checkError(err)

	queryDB()
}

func queryDB() {
	rows, _, err := db.Query("select * from content")
	checkError(err)

	os.Stdout.Write([]byte("Current database contents\n"))

	for _, row := range rows {
		for _, col := range row {
			if col == nil {
				null := []byte("NULL")
				os.Stdout.Write(null)
			} else {
				// Type assertion required because []interface{} "type" is entirely unknown
				val := col.([]byte)
				os.Stdout.Write(append(val, []byte("\t|\t")...))
			}
		}
		os.Stdout.Write([]byte("\n"))
	}
}

// 
