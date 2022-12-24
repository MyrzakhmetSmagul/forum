package main

import (
	"database/sql"
	"fmt"
	"forum/dao"
	"forum/models"
	"log"
	"net/http"
	"text/template"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var (
	Users   []models.User
	db      *sql.DB
	Cookies = make(map[string]string)
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
	userCookie, err := r.Cookie("authToken")
	if err != nil {
		log.Println(err)
	}
	if v, ok := Cookies[userCookie.Value]; ok {
		fmt.Fprintf(w, "<h1>%s</h1>", v)
		log.Println("we have cookie")
		return
	}
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
	user := &models.User{
		Email: r.FormValue("email"),
		Pwd:   r.FormValue("pwd"),
	}
	if !dao.UserVerification(db, user) {
		fmt.Fprintf(w, "<h1>user doesn't exist<h1>")
		return
	}

	token, err := uuid.NewV4()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:  "authToken",
		Value: token.String(),
	}
	Cookies[token.String()] = user.Id
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", 302)
}

func signUpUser(w http.ResponseWriter, r *http.Request) {
	if dao.SearchUser(db, &models.User{UName: r.FormValue("uname"), Email: r.FormValue("email")}) {
		fmt.Fprintf(w, "<h1>user exist</h1>")
		return
	}
}
