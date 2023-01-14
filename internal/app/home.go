package app

import (
	"errors"
	"log"
	"net/http"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

func (s *ServiceServer) IndexWithoutSession(w http.ResponseWriter, r *http.Request) {
	data, err := s.getAllPosts(r)
	if err != nil {
		log.Println("ERROR:\nIndexWithoutSession:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}
	s.render(w, "unauth-index", http.StatusOK, data)
}

func (s *ServiceServer) IndexWithSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusMethodNotAllowed))
		return
	}

	if r.URL.Path != "/" {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusNotFound))
		return
	}

	if _, err := s.getSession(r); err != nil {
		if errors.Is(err, model.ErrUserNotFound) || errors.Is(err, model.ErrNoSession) {
			s.IndexWithoutSession(w, r)
			return
		}
		log.Println("ERROR:\nIndexWithSession:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	data, err := s.getAllPosts(r)
	if err != nil {
		log.Println("ERROR:\nIndexWithSession:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	s.render(w, "index", http.StatusOK, data)
}
