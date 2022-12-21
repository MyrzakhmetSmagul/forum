package main

import (
	"database/sql"
	"fmt"
	"forum/dao"
	"forum/models"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var (
	Users []models.User
	db    *sql.DB
)

func main() {
	log.Println("server start at localhost:8080")
	db, _ = sql.Open("sqlite3", "forum.db")

	// createTable(database)
	var u models.User

	rows, err := db.Query("SELECT * FROM users")
	log.Println(err)
	for rows.Next() {
		rows.Scan(&u.Id, &u.Name, &u.Surname, &u.Gender, &u.Email, &u.Pwd)
		fmt.Println(u.UserInfo())
		Users = append(Users, u)
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/write", writeHandler)
	http.HandleFunc("/addUser", addUserHandler)
	http.HandleFunc("/signIn", signInHandler)
	http.ListenAndServe("localhost:8080", nil)
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

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		log.Println("method not equal post")
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	newUser := models.User{
		Name:    r.PostFormValue("fName"),
		Surname: r.PostFormValue("surname"),
		Gender:  r.PostFormValue("gender"),
		Email:   r.PostFormValue("email"),
		Pwd:     r.PostFormValue("pwd"),
	}
	fmt.Println("######################################################\n\nADD USER ACTION WILL BE ACTIVATE\n\n######################################################")
	err = dao.AddUser(db, &newUser)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	temp, err := template.ParseFiles("./template/index.html", "./template/header.html", "./template/footer.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	temp.ExecuteTemplate(w, "index", nil)
}

func writeHandler(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("./template/write.html", "./template/header.html", "./template/footer.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	temp.ExecuteTemplate(w, "write", nil)
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("./template/signin.html", "./template/header.html", "./template/footer.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	temp.ExecuteTemplate(w, "signin", nil)
}
