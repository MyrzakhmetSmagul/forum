package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"forum/dao"
// 	"forum/models"
// 	"log"
// 	"net/http"
// 	"text/template"

// 	"github.com/gofrs/uuid"
// 	_ "github.com/mattn/go-sqlite3"
// )

// var (
// 	Users   []models.User
// 	db      *sql.DB
// 	Cookies = make(map[string]string)
// )

// func main() {
// 	log.Println("server start at localhost:8080")
// 	var err error

// 	db, err = sql.Open("sqlite3", "forum.db")
// 	if err != nil {
// 		log.Println(err.Error())
// 		return
// 	}

// 	http.HandleFunc("/", indexHandler)
// 	http.ListenAndServe(":8080", nil)
// }

// func indexHandler(w http.ResponseWriter, r *http.Request) {
// 	switch r.URL.Path {
// 	case "/sign-in-form":
// 		getSignInUser(w, r)
// 	case "/sign-up-form":
// 		getSignUpUser(w, r)
// 	case "/sign-in":
// 		signInUser(w, r)
// 	case "/sign-up":
// 		signUpUser(w, r)
// 	default:
// 		homePage(w, r)
// 	}
// }

// // HomePage displays the homepage of the forum
// func homePage(w http.ResponseWriter, r *http.Request) {
// 	t, _ := template.ParseFiles("./template/index.html", "./template/header.html", "./template/footer.html")
// 	userCookie, err := r.Cookie("authToken")
// 	if err != nil {
// 		if err == http.ErrNoCookie {
// 			log.Println("cookie not found")
// 			t.ExecuteTemplate(w, "index", nil)
// 			return
// 		}
// 		log.Println(err)
// 		t.ExecuteTemplate(w, "index", nil)
// 		return
// 	}

// 	session := models.Session{Token: userCookie.Value}
// 	exist, err := dao.IsExistSession(db, &session)
// 	if err != nil {
// 		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 		log.Println(err)
// 		return
// 	}

// 	if !exist {
// 		t.ExecuteTemplate(w, "index", nil)
// 		return
// 	}

// 	err = dao.GetSessionFromDB(db, &session)
// 	if err != nil {
// 		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 		log.Println(err)
// 		return
// 	}

// 	t.ExecuteTemplate(w, "index", nil)
// }

// // GetSignInUser displays the sign in page
// func getSignInUser(w http.ResponseWriter, r *http.Request) {
// 	t, _ := template.ParseFiles("./template/sign-in.html", "./template/header.html", "./template/footer.html")
// 	t.ExecuteTemplate(w, "sign-in", nil)
// }

// // GetSignUpUser displays the sign up page
// func getSignUpUser(w http.ResponseWriter, r *http.Request) {
// 	t, _ := template.ParseFiles("./template/sign-up.html", "./template/header.html", "./template/footer.html")
// 	t.ExecuteTemplate(w, "sign-up", nil)
// }

// // SignInUser handles the user sign in request
// func signInUser(w http.ResponseWriter, r *http.Request) {
// 	user := &models.User{
// 		Email:  r.FormValue("email"),
// 		Passwd: r.FormValue("passwd"),
// 	}
// 	if !dao.UserVerification(db, user) {
// 		fmt.Fprintf(w, "<h1>user doesn't exist</h1>")
// 		return
// 	}

// 	token, err := uuid.NewV4()
// 	if err != nil {
// 		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 		return
// 	}
// 	session := models.Session{Token: token.String(), UserId: user.Id}
// 	cookie := http.Cookie{
// 		Name:  "authToken",
// 		Value: token.String(),
// 	}

// 	err = dao.AddSession(db, &session)
// 	if err != nil {
// 		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 		return
// 	}

// 	Cookies[token.String()] = user.Id
// 	http.SetCookie(w, &cookie)
// 	http.Redirect(w, r, "/", 302)
// }

// // SignUpUser handles the user sign up request
// func signUpUser(w http.ResponseWriter, r *http.Request) {
// 	t, _ := template.ParseFiles("./template/index.html", "./template/header.html", "./template/footer.html")
// 	user := models.NewUser(r.FormValue("uname"), r.FormValue("email"), r.FormValue("passwd"))
// 	exist, err := dao.IsExistUser(db, user)
// 	if err != nil {
// 		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 		log.Println(err)
// 		return
// 	}

// 	if exist {
// 		fmt.Fprintln(w, "<h1>User exist</h1>")
// 		return
// 	}

// 	err = dao.AddUser(db, user)
// 	if err != nil {
// 		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 		log.Println(err)
// 		return
// 	}

// 	t.ExecuteTemplate(w, "index", nil)
// }
