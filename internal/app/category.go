package app

import (
	"errors"
	"html/template"
	"log"
	"net/http"

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

	data, err := s.getPostsOfCategory(r)
	if err != nil {
		log.Println("ERROR:\nUnauthCategory:", err)
		if errors.Is(err, model.ErrValueNotSet) {
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadRequest))
			return
		}

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	t.ExecuteTemplate(w, "index", data)
}

func (s *ServiceServer) Category(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusMethodNotAllowed))
		return
	}

	_, err := s.getSession(r)
	if err != nil {
		if errors.Is(err, model.ErrNoSession) || errors.Is(err, model.ErrUserNotFound) {
			s.UnauthCategory(w, r)
			return
		}
		log.Println("ERROR:\nCategory:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	t, err := template.ParseFiles("./templates/html/index.html")
	if err != nil {
		log.Println("ERROR:\nCategory:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	data, err := s.getPostsOfCategory(r)
	if err != nil {
		log.Println("ERROR:\nCategory:", err)
		if errors.Is(err, model.ErrValueNotSet) {
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadRequest))
			return
		}

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	t.ExecuteTemplate(w, "index", data)
}
