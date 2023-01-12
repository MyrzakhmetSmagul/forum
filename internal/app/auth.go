package app

import (
	"database/sql"
	"errors"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/MyrzakhmetSmagul/forum/internal/app/validation"
	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

var (
	SignInMessage string
	SignUpMessage string
)

func (s *ServiceServer) SignIn(w http.ResponseWriter, r *http.Request) {
	_, err := s.getSession(r)
	if err != nil && !errors.Is(err, model.ErrNoSession) && !errors.Is(err, model.ErrUserNotFound) {
		log.Println("ERROR:\nSignIn:", err)
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	} else if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if r.Method == http.MethodGet {
		s.GetSignIn(w, r)
		return
	}
	s.PostSignIn(w, r)
	return
}

func (s *ServiceServer) SignUp(w http.ResponseWriter, r *http.Request) {
	_, err := s.getSession(r)
	if err != nil && !errors.Is(err, model.ErrNoSession) && !errors.Is(err, model.ErrUserNotFound) {
		log.Println("ERROR:\nSignUp:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	} else if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

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
		log.Println("ERROR:\nSignIn with GET Method:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	err = t.ExecuteTemplate(w, "signin", SignInMessage)
	if err != nil {
		log.Println("ERROR:\nSignIn with GET Method:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}
}

func (s *ServiceServer) GetSignUp(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./templates/html/signup.html")
	if err != nil {
		log.Println("ERROR:\nSignUp with GET Method:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	err = t.ExecuteTemplate(w, "signup", SignUpMessage)
	if err != nil {
		log.Println("ERROR:\nSignUp with GET Method:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}
}

func (s *ServiceServer) PostSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusMethodNotAllowed))
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Println("ERROR:\nSignIn With POST Method: ", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadGateway))
		return
	}

	user := model.User{Email: r.PostFormValue("email"), Password: r.PostFormValue("password")}
	if err := validation.ValidationFormSignIn(user.Email, user.Password); err != nil {
		log.Println("ERROR:\nSignIn With POST Method: ", err)

		if errors.Is(err, model.ErrMessageValid) {
			SignInMessage = model.ErrMessageValid.Error()
			http.Redirect(w, r, "/signIn", http.StatusFound)
		} else {
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		}
		return
	}

	session := model.Session{Expiry: time.Now().Add(time.Minute * 15)}

	err := s.authService.SignIn(&user, &session)
	if err != nil {
		log.Println("ERROR:\nSignIn With POST Method: ", err)

		if errors.Is(err, model.ErrUserNotFound) || errors.Is(err, sql.ErrNoRows) {
			SignInMessage = model.ErrMessageValid.Error()
			http.Redirect(w, r, "/signIn", http.StatusFound)
		} else {
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		}
		return
	}

	cookie := http.Cookie{
		Name:     "authToken",
		Value:    session.Token,
		SameSite: http.SameSiteDefaultMode,
		MaxAge:   900,
		Expires:  session.Expiry,
		HttpOnly: true,
	}

	SignInMessage = ""
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *ServiceServer) PostSignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusMethodNotAllowed))
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Println("ERROR:\nSignUp With POST Method: ", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadGateway))
		return
	}

	user := model.User{Username: r.PostFormValue("username"), Email: r.PostFormValue("email"), Password: r.PostFormValue("password"), Password2: r.PostFormValue("password2")}
	if err := validation.ValidationFormSignUp(user.Username, user.Email, user.Password, user.Password2); err != nil {
		log.Println("ERROR:\nSignUp With POST Method: ", err)

		if errors.Is(err, model.ErrMessageValid) {
			SignUpMessage = model.ErrMessageValid.Error()
			http.Redirect(w, r, "/signUp", http.StatusFound)
		} else {
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		}
		return
	}

	err := s.authService.SignUp(&user)
	if err != nil {
		log.Println("ERROR:\nSignUp With POST Method: ", err)

		if errors.Is(err, model.ErrUserExists) {
			SignUpMessage = model.ErrUserExists.Error()
			http.Redirect(w, r, "/signUp", http.StatusFound)
		} else {
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		}
		return
	}

	SignUpMessage = ""
	http.Redirect(w, r, "/signIn", http.StatusFound)
}

func (s *ServiceServer) SignOut(w http.ResponseWriter, r *http.Request, session model.Session) {
	if r.Method != http.MethodGet {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusMethodNotAllowed))
		return
	}

	err := s.authService.SignOut(&session)
	if err != nil {
		log.Println("ERROR:\nSignOut: ", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
