package dao

import (
	"database/sql"
	"io/ioutil"
	"log"
)

func CreateTable(db *sql.DB) {
	data, err := ioutil.ReadFile("./dao/query/userTable.sql")
	if err != nil {
		log.Fatal(err)
	}
	query, err := db.Prepare(string(data))
	if err != nil {
		log.Fatal(err)
	}
	_, err = query.Exec()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Table was created")
}
