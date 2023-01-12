package app

import (
	"html/template"
	"net/http"
	"time"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

func (s *ServiceServer) UnauthCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusMethodNotAllowed))
		return
	}

	t, err := template.ParseFiles("./templates/html/unauth-index.html")
	if err != nil {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	id, err := s.getID(r)
	if err != nil {
		if err.Error() == "ID not set" {
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadRequest))
			return
		}
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	category := model.Category{ID: int64(id)}
	posts, err := s.postService.GetPostsOfCategory(category)
	if err != nil {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	categories, err := s.postService.GetAllCategory()
	if err != nil {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	data := model.Data{Categories: categories, Posts: posts}
	t.ExecuteTemplate(w, "index", data)
}

func (s *ServiceServer) Category(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusMethodNotAllowed))
		return
	}

	cookie, err := r.Cookie("authToken")
	if err != nil {
		if err == http.ErrNoCookie {
			s.UnauthCategory(w, r)
			return
		}
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadGateway))
		return
	}

	session := model.Session{Token: cookie.Value}
	err = s.sessionService.GetSession(&session)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			s.UnauthCategory(w, r)
			return
		}

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	if session.Expiry.Before(time.Now()) {
		err = s.sessionService.DeleteSession(&session)
		if err != nil {
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
			return
		}

		s.UnauthCategory(w, r)
		return
	}

	t, err := template.ParseFiles("./templates/html/index.html")
	if err != nil {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	id, err := s.getID(r)
	if err != nil {
		if err.Error() == "ID not set" {
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadRequest))
			return
		}
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	category := model.Category{ID: int64(id)}
	posts, err := s.postService.GetPostsOfCategory(category)
	if err != nil {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	categories, err := s.postService.GetAllCategory()
	if err != nil {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	data := model.Data{Categories: categories, Posts: posts}
	t.ExecuteTemplate(w, "index", data)
}
