// Provides an interface to the storage solution with drivers for datatypes from rter/data
//
// Functions are provided to make MySQL storage easier to use within the rtER project. This includes helps for setting up the connection and running queries.
package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

var (
	ErrZeroAffected        = errors.New("Query didn't match anything.")
	ErrUnsupportedDataType = errors.New("Storage doesn't support the given datatype.")
	ErrCannotDelete        = errors.New("Storage doesn't allow deleting that.")
)

// Begin a new transaction against the current db.
func Begin() (*sql.Tx, error) {
	return db.Begin()
}

// Run an exec against the current connected db.
func Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.Exec(query, args...)
}

// Run a query against the current connected db.
func Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(query, args...)
}

// Run an exec against the current connected db. Fatal if the query throws an error.
func MustExec(query string, args ...interface{}) sql.Result {
	res, err := db.Exec(query, args...)
	if err != nil {
		log.Fatalf("Error on Exec %q: %v", query, err)
	}
	return res
}

// Run a query against the current connected db. Fatal if the query throws an error.
func MustQuery(query string, args ...interface{}) *sql.Rows {
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatalf("Error on Query %q: %v", query, err)
	}
	return rows
}

// Print out all the results from an SQL query to STOUT via fmt in a readable format
func DumpRows(rows *sql.Rows) {
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			fmt.Println(columns[i], ": ", value)
		}
		fmt.Println("-----------------------------------")
	}

}
