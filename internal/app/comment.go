package app

import (
	"net/http"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

func (s *ServiceServer) CreateComment(w http.ResponseWriter, r *http.Request, session *model.Session) {
	if r.Method != http.MethodPost {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusMethodNotAllowed, StatusText: http.StatusText(http.StatusMethodNotAllowed)})
		return
	}
}
