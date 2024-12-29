package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	// Open a database connection
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/database")
	if err != nil {
		log.Fatal(err)
	}
	// Defer the closing of the database connection until the main function returns
	defer db.Close()
	// Prepare a SQL statement
	stmt, err := db.Prepare("INSERT INTO table (name) VALUES (?)")
	if err != nil {
		log.Fatal(err)
	}
	// Defer the closing of the prepared statement until the main function returns
	defer stmt.Close()
	// Execute the SQL statement with a parameter
	_, err = stmt.Exec("example")
	if err != nil {
		log.Fatal(err)
	}
}
