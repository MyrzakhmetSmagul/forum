package dao

import (
	"database/sql"
	"log"
)

func CreateTable(db *sql.DB) {
	users_table := `CREATE TABLE IF NOT EXISTS users (
		user_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
		name TEXT NOT NULL, 
		surname TEXT NOT NULL, 
		gender TEXT NOT NULL, 
		email TEXT, 
		pwd TEXT NOT NULL)`
	query, err := db.Prepare(users_table)
	if err != nil {
		log.Fatal(err)
	}
	query.Exec()
	log.Println("Table was created")
}
