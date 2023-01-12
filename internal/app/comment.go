package app

import (
	"net/http"
	"strconv"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

func (s *ServiceServer) CreateComment(w http.ResponseWriter, r *http.Request, session model.Session) {
	if r.Method != http.MethodPost {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusMethodNotAllowed))
		return
	}

	postID, err := s.getID(r)
	if err != nil {
		if err.Error() == "ID not set" {
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadRequest))
			return
		}
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	if r.ParseForm() != nil {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadGateway))
		return
	}

	comment := model.Comment{PostID: int64(postID), UserID: session.User.ID, Username: session.User.Username, Message: r.PostFormValue("comment")}
	err = s.postService.CreateComment(&comment)
	if err != nil {
		if err.Error() == "post doesn't exist" {
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadRequest))
			return
		}
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	http.Redirect(w, r, "/post?ID="+strconv.Itoa(postID), http.StatusFound)
}
