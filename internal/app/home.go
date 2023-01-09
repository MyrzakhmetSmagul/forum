package app

import (
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

func (s *ServiceServer) IndexWithoutSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusMethodNotAllowed, StatusText: http.StatusText(http.StatusMethodNotAllowed)})
		return
	}

	if r.URL.Path != "/" {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusNotFound, StatusText: http.StatusText(http.StatusNotFound)})
		return
	}

	t, err := template.ParseFiles("./templates/html/unauth-index.html")
	if err != nil {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	categories, err := s.postService.GetAllCategory()
	if err != nil {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	posts, err := s.postService.GetAllPosts()
	if err != nil {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	data := model.Data{Categories: categories, Posts: posts}

	err = t.ExecuteTemplate(w, "index", data)
	if err != nil {
		log.Println(err)

		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}
}

func (s *ServiceServer) IndexWithSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusMethodNotAllowed, StatusText: http.StatusText(http.StatusMethodNotAllowed)})
		return
	}

	if r.URL.Path != "/" {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusNotFound, StatusText: http.StatusText(http.StatusNotFound)})
		return
	}

	cookie, err := r.Cookie("authToken")
	if err != nil {
		if err == http.ErrNoCookie {
			s.IndexWithoutSession(w, r)
			return
		}
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusBadGateway, StatusText: http.StatusText(http.StatusBadGateway)})
		return
	}

	session := model.Session{Token: cookie.Value}
	err = s.sessionService.GetSession(&session)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			s.IndexWithoutSession(w, r)
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

		s.IndexWithoutSession(w, r)
		return
	}

	t, err := template.ParseFiles("./templates/html/index.html")
	if err != nil {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	categories, err := s.postService.GetAllCategory()
	if err != nil {
		log.Println("get all categories error", err)
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	posts, err := s.postService.GetAllPosts()
	if err != nil {
		log.Println("get all posts error", err)
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	data := model.Data{Categories: categories, Posts: posts}

	err = t.ExecuteTemplate(w, "index", data)
	if err != nil {
		log.Println(err)

		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}
}
