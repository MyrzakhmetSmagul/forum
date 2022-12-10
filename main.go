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

	createTable(database)
	var u User

	addUser(database, &User{
		Id:      1,
		Name:    "Myrzakhmet",
		Surname: "Smagul",
		Gender:  "MALE",
		Email:   "smagul.myrzakhmet@mail.ru",
		Pwd:     "kcnd;lksm",
	})
	rows, err := database.Query("SELECT * FROM users")
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

func addUser(db *sql.DB, u *User) {
	record := `INSERT INTO 	users(name, surname, gender, email, pwd) VALUES(?, ?, ?, ?, ?)`
	query, err := db.Prepare(record)
	if err != nil {
		log.Fatal(err)
	}

	_, err = query.Exec(u.Name, u.Surname, u.Gender, u.Email, u.Pwd)
	if err != nil {
		log.Fatal()
	}

	log.Println("INSERT INTO OK")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
	for _, v := range Users {
		fmt.Fprintf(w, "info")
		fmt.Fprintf(w, v.UserInfo())
	}
	fmt.Fprintf(w, "bye")
}
