package app

import (
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

func (s *ServiceServer) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.GetSignIn(w, r)
		return
	}
	s.PostSignIn(w, r)
	return
}

func (s *ServiceServer) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.GetSignUp(w, r)
		return
	}
	s.PostSignUp(w, r)
	return
}

func (s *ServiceServer) GetSignIn(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./templates/html/signin.html")
	if err != nil {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	err = t.ExecuteTemplate(w, "signin", nil)
	if err != nil {
		log.Println(err)
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}
}

func (s *ServiceServer) GetSignUp(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./templates/html/signup.html")
	if err != nil {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	err = t.ExecuteTemplate(w, "signup", nil)
	if err != nil {
		log.Println(err)
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}
}

func (s *ServiceServer) PostSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusMethodNotAllowed, StatusText: http.StatusText(http.StatusMethodNotAllowed)})
		return
	}

	err := r.ParseForm()
	if err != nil {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	user := model.User{Email: r.PostFormValue("email"), Password: r.PostFormValue("password")}
	session := model.Session{Expiry: time.Now().Add(time.Minute)}
	err = s.authService.SignIn(&user, &session)
	if err != nil {
		if err.Error() == "the user's email or password is incorrect" || err.Error() == "sql: no rows in result set" {
			http.Redirect(w, r, "/signIn", http.StatusFound)
			return
		}

		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	cookie := http.Cookie{
		Name:     "authToken",
		Value:    session.Token,
		SameSite: http.SameSiteDefaultMode,
		MaxAge:   60,
		Expires:  session.Expiry,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *ServiceServer) PostSignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusMethodNotAllowed, StatusText: http.StatusText(http.StatusMethodNotAllowed)})
		return
	}

	err := r.ParseForm()
	if err != nil {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	user := model.User{Username: r.PostFormValue("username"), Email: r.PostFormValue("email"), Password: r.PostFormValue("password")}
	err = s.authService.SignUp(&user)
	if err != nil {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	http.Redirect(w, r, "/signIn", http.StatusFound)
}

func (s *ServiceServer) SignOut(w http.ResponseWriter, r *http.Request, session *model.Session) {
	if r.Method != http.MethodGet {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusMethodNotAllowed, StatusText: http.StatusText(http.StatusMethodNotAllowed)})
		return
	}

	err := s.authService.SignOut(session)
	if err != nil {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
