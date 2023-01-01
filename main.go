package main

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/dao"
	"forum/models"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var (
	Users   []models.User
	db      *sql.DB
	Cookies = make(map[string]*models.Session)
)

func main() {
	log.Println("server start at localhost:8080")
	var err error

	db, err = sql.Open("sqlite3", "forum.db")
	if err != nil {
		log.Println(err.Error())
		return
	}
	// func() {

	// }()
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
	case "log-out":
		logout(w, r)
	case "/":
		homePage(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./template/index.html", "./template/header.html", "./template/footer.html")
	userCookie, err := r.Cookie("authToken")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Println("cookie not found")
			t.ExecuteTemplate(w, "index", nil)
			return
		}
		log.Println(err)
		t.ExecuteTemplate(w, "index", nil)
		return
	}

	session := models.Session{Token: userCookie.Value}

	if _, ok := Cookies[userCookie.Value]; !ok {
		t.ExecuteTemplate(w, "index", nil)
		return
	}

	err = checkCookie(userCookie, &session)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "<h1>user exist and have token")
}

func getSignInUser(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./template/sign-in.html", "./template/header.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.ExecuteTemplate(w, "sign-in", nil)
}

func getSignUpUser(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./template/sign-up.html", "./template/header.html", "./template/footer.html")
	t.ExecuteTemplate(w, "sign-up", nil)
}

func signInUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{
		Email:  r.FormValue("email"),
		Passwd: r.FormValue("passwd"),
	}

	if dao.UserVerification(db, user) != nil {
		fmt.Fprintf(w, "<h1>user doesn't exist<h1>")
		return
	}

	token, err := uuid.NewV4()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	session := models.Session{Token: token.String(), UserId: user.Id}
	cookie := http.Cookie{
		Name:    "authToken",
		Value:   token.String(),
		Expires: time.Now().Add(time.Second * 15),
	}

	err = dao.AddSession(db, &session)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	Cookies[token.String()] = &session
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", 302)
}

func signUpUser(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./template/index.html", "./template/header.html", "./template/footer.html")
	user := models.NewUser(r.FormValue("uname"), r.FormValue("email"), r.FormValue("passwd"))
	exist, err := dao.IsExistUser(db, user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if exist {
		http.Redirect(w, r, "/sign-in-form", 302)
		log.Println("user exist")
		return
	}

	err = dao.AddUser(db, user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	t.ExecuteTemplate(w, "index", nil)
}

func logout(w http.ResponseWriter, r *http.Request) {
	userCookie, err := r.Cookie("authToken")
	if err != nil {
		http.Error(w, "cant delete cookie", http.StatusInternalServerError)
		return
	}

	err = dao.DeleteSession(db, Cookies[userCookie.Value])
	if err != nil {
		http.Error(w, "cant delete cookie", http.StatusInternalServerError)
		return
	}

	delete(Cookies, userCookie.Value)
	http.Redirect(w, r, "/", 302)
}

func checkCookie(cookie *http.Cookie, session *models.Session) error {
	if time.Now().Before(cookie.Expires) {
		return nil
	}

	err := dao.DeleteSession(db, session)
	if err != nil {
		log.Println("cant delete cookie")
		return err
	}

	delete(Cookies, cookie.Value)
	return errors.New("cookie not valid")
}
