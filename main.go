package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id      int
	Name    string
	Surname string
	Gender  string
	Email   string
	Pwd     string
}

func (u *User) UserInfo() string {
	info := fmt.Sprintf("User id: %d", u.Id)
	info += fmt.Sprintf("\nFullname: %s %s", u.Name, u.Surname)
	info += fmt.Sprintf("\nGender: %s", u.Gender)
	return info
}

var Users []User

func main() {
	http.HandleFunc("/", indexHandler)
	log.Println("localhost:8080")
	database, err := sql.Open("sqlite3", "forum.db")
	log.Println(err)
	// statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS users (user_id INTEGER NOT NULL PRIMARY KEY, user_name TEXT NOT NULL, user_surname TEXT NOT NULL, user_gender TEXT NOT NULL, user_email TEXT, user_pwd TEXT NOT NULL)")
	// log.Println(err)
	// statement.Exec()
	rows, err := database.Query("SELECT * FROM users")
	var u User
	log.Println(err)
	for rows.Next() {
		rows.Scan(&u.Id, &u.Name, &u.Surname, &u.Gender, &u.Email, &u.Pwd)
		fmt.Println(u.UserInfo())
		Users = append(Users, u)
	}
	fmt.Println(Users)
	http.ListenAndServe("localhost:8080", nil)
}

func createTable(db *sql.DB) {
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
	for _, v := range Users {
		fmt.Fprintf(w, "info")
		fmt.Fprintf(w, v.UserInfo())
	}
	fmt.Fprintf(w, "bye")
}
