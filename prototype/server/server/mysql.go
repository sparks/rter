package server

import (
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	"os"
)

var database mysql.Conn

func SetupMySQL() {
	database = mysql.New("tcp", "", "localhost:3306", "root", "", "rter")

	error := database.Connect()
	checkError(error)

	queryDatabase()
}

func queryDatabase() {
	rows, _, error := database.Query("select * from content")
	checkError(error)

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
