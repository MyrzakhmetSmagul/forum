package main

import (
	"database/sql"
	"fmt"
	"forum/models"
	"log"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

var (
	Users []models.User
	db    *sql.DB
)

func main() {
	log.Println("server start at localhost:8080")
	var err error
	db, err = sql.Open("sqlite3", "forum.db")
	if err != nil {
		log.Println(err.Error())
		return
	}

	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/sign-in-form":
		fmt.Println("sign-in-form")
		getSignInUser(w, r)
	case "/sign-up-form":
		getSignUpUser(w, r)
	case "/sign-in":
		signInUser(w, r)
	case "/sign-up":
		signUpUser(w, r)
	default:
		homePage(w, r)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./template/index.html", "./template/header.html", "./template/footer.html")
	t.ExecuteTemplate(w, "index", nil)
}

func getSignInUser(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./template/sign-in.html", "./template/header.html", "./template/footer.html")
	t.ExecuteTemplate(w, "sign-in", nil)
}

func getSignUpUser(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./template/sign-up.html", "./template/header.html", "./template/footer.html")
	t.ExecuteTemplate(w, "sign-up", nil)
}

func signInUser(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func signUpUser(w http.ResponseWriter, r *http.Request) {
	newUser := getUser(r)
	DefaultUserService.createUser(newUser)
}

func getUser(r *http.Request) *models.User {
	email := r.FormValue("email")
	passwd := r.FormValue("pwd")
	return &models.User{
		Email: email,
		Pwd:   passwd,
	}
}
