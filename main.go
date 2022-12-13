package main

import (
	"database/sql"
	"fmt"
	"html/template"
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

var (
	Users []User
	db    *sql.DB
)

func main() {
	log.Println("server start at localhost:8080")
	db, err := sql.Open("sqlite3", "forum.db")
	log.Println(err)

	// createTable(database)
	var u User

	// addUser(db, &User{
	// 	Name:    "Tursynkhan",
	// 	Surname: "Tursunov",
	// 	Gender:  "MALE",
	// 	Email:   "Tursynkhan@mail.ru",
	// 	Pwd:     "hdfuivhishokv",
	// })

	rows, err := db.Query("SELECT * FROM users")
	log.Println(err)
	for rows.Next() {
		rows.Scan(&u.Id, &u.Name, &u.Surname, &u.Gender, &u.Email, &u.Pwd)
		fmt.Println(u.UserInfo())
		Users = append(Users, u)
	}
	fmt.Println(Users)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/addUser", addUserFromSite)
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
	if r.URL.Path != "/" {
		log.Println(r.URL.Path)
		http.NotFound(w, r)
		return
	}

	temp, err := template.ParseFiles("./template/index.html", "./template/header.html", "./template/footer.html")
	if err != nil {
		fmt.Fprintf(w, "Internal Server Error!")
		log.Println(err.Error())
		return
	}

	temp.ExecuteTemplate(w, "index", nil)
}

func addUserFromSite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprintf(w, "Bad Request!")
		log.Println("method not equal post")
		return
	}

	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "Internal Server Error!")
		log.Println(err.Error())
		return
	}

	newUser := User{
		Name:    r.PostFormValue("fName"),
		Surname: r.PostFormValue("surname"),
		Gender:  r.PostFormValue("gender"),
		Email:   r.PostFormValue("email"),
		Pwd:     r.PostFormValue("pwd"),
	}

	addUser(db, &newUser)
}
