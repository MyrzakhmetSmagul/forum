package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type user struct {
	id      int
	name    string
	surname string
	gender  string
	email   string
	pwd     string
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe("localhost:8080", nil)

	database, err := sql.Open("sqlite3", "forum.db")
	log.Println(err)
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS users (user_id INTEGER NOT NULL PRIMARY KEY, user_name TEXT NOT NULL, user_surname TEXT NOT NULL,user_gender TEXT NOT NULL, user_email TEXT, user_pwd TEXT NOT NULL)")
	log.Println(err)
	statement.Exec()
	// statement, err = database.Prepare("INSERT INTO users (user_id, user_name, user_surname, user_gender, user_email, user_pwd) VALUES (?, ?, ?, ?, ?, ? )")
	//log.Println(err)
	// statement.Exec(2, "John", "Doe", "MALE", "jahndoe@email.com", "SHA256")
	rows, err := database.Query("SELECT * FROM users")
	var u user
	var users []user
	log.Println(err)
	for rows.Next() {
		rows.Scan(&u.id, &u.name, &u.surname, &u.gender, &u.email, &u.pwd)
		users = append(users, u)
	}
	fmt.Println(users)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

}
