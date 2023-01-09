package app

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

func (s *ServiceServer) authMiddleware(next func(http.ResponseWriter, *http.Request, model.Session)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("authToken")
		if err != nil {
			if err == http.ErrNoCookie {
				if r.URL.Path == "/post" {
					s.PostUnauth(w, r)
					return
				}
				s.ErrorHandler(w, model.Error{StatusCode: http.StatusUnauthorized, StatusText: http.StatusText(http.StatusUnauthorized)})
				return
			}
			s.ErrorHandler(w, model.Error{StatusCode: http.StatusBadGateway, StatusText: http.StatusText(http.StatusBadGateway)})
			return
		}

		session := model.Session{Token: cookie.Value}
		err = s.sessionService.GetSession(&session)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				if r.URL.Path == "/post" {
					s.PostUnauth(w, r)
					return
				}

				s.ErrorHandler(w, model.Error{StatusCode: http.StatusUnauthorized, StatusText: http.StatusText(http.StatusUnauthorized)})
				return
			}

			s.ErrorHandler(w, model.Error{StatusCode: http.StatusBadGateway, StatusText: http.StatusText(http.StatusBadGateway)})
			return
		}

		if session.Expiry.Before(time.Now()) {
			err = s.sessionService.DeleteSession(&session)
			if err != nil {
				s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
				return
			}

			if r.URL.Path == "/post" {
				s.PostUnauth(w, r)
				return
			}

			s.ErrorHandler(w, model.Error{StatusCode: http.StatusUnauthorized, StatusText: http.StatusText(http.StatusUnauthorized)})
			return
		}

		err = s.userService.GetUserInfo(&session.User)
		if err != nil {
			if err.Error() == "user doesn't exist" {
				err = s.sessionService.DeleteSession(&session)
				if err != nil {
					s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
					return
				}
			} else {
				s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
				return
			}

			if r.URL.Path == "/post" {
				s.PostUnauth(w, r)
				return
			}

			s.ErrorHandler(w, model.Error{StatusCode: http.StatusUnauthorized, StatusText: http.StatusText(http.StatusUnauthorized)})
			return
		}

		next(w, r, session)
	}
}

func (s *ServiceServer) getID(r *http.Request) (int, error) {
	if r.URL.Query().Get("ID") == "" {
		log.Println("ID not set")
		return 0, errors.New("ID not set")
	}

	id, err := strconv.Atoi(r.URL.Query().Get("ID"))
	if err != nil {
		log.Println("ID atoi error:", err)
		return 0, err
	}
	return id, nil
}
