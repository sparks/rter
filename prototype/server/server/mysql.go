package server

import (
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

var database mysql.Conn

func SetupMySQL() {
	database = mysql.New("tcp", "", "localhost:3306", "root", "", "rter")

	error := database.Connect()
	checkError(error)
}

func queryDatabase() {
	rows, _, error := database.Query("select * from content")
	checkError(error)

	for _, row := range rows {
		for _, col := range row {
			if col == nil {
				fmt.Println("NULL")
			} else {
				// Type assertion required because []interface{} "type" is entirely unknown
				val := col.([]byte)
				fmt.Print(string(val) + "\t|\t")
			}
		}
		fmt.Println()
	}
}
