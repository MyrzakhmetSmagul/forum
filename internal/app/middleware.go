package app

import (
	"log"
	"net/http"
	"time"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

func (s *ServiceServer) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("authToken")
		if err == http.ErrNoCookie {
			log.Println("middleware no cookie", r.URL.Path)
			if r.URL.Path == "/" || r.URL.Path == "/signOut" {
				s.IndexWithoutSession(w, r)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		session := model.Session{Token: cookie.Value}
		err = s.sessionService.GetSession(&session)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				if r.URL.Path == "/signIn" || r.URL.Path == "/signUp" {
					next.ServeHTTP(w, r)
					return
				}

				if r.URL.Path == "/" {
					s.IndexWithoutSession(w, r)
					return
				}

				http.Redirect(w, r, "/", http.StatusFound)
				return
			}

			s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
			return
		}

		if session.Expiry.Before(time.Now()) {
			err = s.sessionService.DeleteSession(&session)
			if err != nil {
				s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
				return
			}

			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		if r.URL.Path == "/signIn" || r.URL.Path == "/signUp" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	}
}
