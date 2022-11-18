package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database, err := sql.Open("sqlite3", "forum.db")
	log.Println(err)
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS users (user_id INTEGER NOT NULL PRIMARY KEY, user_name TEXT NOT NULL, user_surname TEXT NOT NULL,user_gender TEXT NOT NULL, user_email TEXT, user_pwd TEXT NOT NULL)")
	log.Println(err)
	statement.Exec()
	statement, err = database.Prepare("INSERT INTO users (user_id, user_name, user_surname, user_gender, user_email, user_pwd) VALUES (1, 'Myrzakhmet', 'Smagul', 'MALE', 'smagul.myrzakhmet@mail.ru', 'SHA256' )")
	log.Println(err)
	statement.Exec()

}
