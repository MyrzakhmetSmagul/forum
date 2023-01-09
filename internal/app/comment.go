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

	postID, err := s.getID(r)
	if err != nil {
		if err.Error() == "ID not set" {
			s.ErrorHandler(w, model.Error{StatusCode: http.StatusBadRequest, StatusText: http.StatusText(http.StatusBadRequest)})
			return
		}
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	comment := model.Comment{PostID: int64(postID), UserID: session.User.ID, Username: session.User.Username, Message: r.PostFormValue("comment")}
	s.postService.CreateComment(&comment)
}
