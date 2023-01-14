package app

import (
	"net/http"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

func (s *ServiceServer) CreatedPosts(w http.ResponseWriter, r *http.Request, session model.Session) {
	if r.Method != http.MethodGet {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusMethodNotAllowed))
		return
	}

	posts, err := s.postService.GetUserPosts(session.User)
	if err != nil {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	categories, err := s.postService.GetAllCategories()
	if err != nil {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	data := model.Data{Categories: categories, Posts: posts}

	s.render(w, "index", http.StatusOK, data)
}

func (s *ServiceServer) RatedPosts(w http.ResponseWriter, r *http.Request, session model.Session) {
	if r.Method != http.MethodGet {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusMethodNotAllowed))
		return
	}

	posts, err := s.postService.GetRatedPosts(session.User)
	if err != nil {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	categories, err := s.postService.GetAllCategories()
	if err != nil {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	data := model.Data{Categories: categories, Posts: posts}

	s.render(w, "index", http.StatusOK, data)
}
