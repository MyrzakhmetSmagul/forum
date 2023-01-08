package app

import (
	"net/http"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

func (s *ServiceServer) CreateComment(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("authToken"); err != nil {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusForbidden, StatusText: http.StatusText(http.StatusForbidden)})
		return
	}
	if r.Method != http.MethodPost {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusMethodNotAllowed, StatusText: http.StatusText(http.StatusMethodNotAllowed)})
		return
	}
}
