package app

import (
	"log"
	"net/http"
	"strconv"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

func (s *ServiceServer) CommentLike(w http.ResponseWriter, r *http.Request, session model.Session) {
	if r.Method != http.MethodGet {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusMethodNotAllowed))
		return
	}

	commentID, err := s.getID(r)
	if err != nil {
		if err.Error() == "ID not set" {
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadRequest))
			return
		}
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	reaction := model.CommentReaction{Comment: model.Comment{ID: int64(commentID)}, User: session.User}
	err = s.postService.CommentSetLike(&reaction)
	if err != nil {
		if err.Error() == "comment doesn't exist" {
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadRequest))
			return
		}
		log.Println("comment like was failed", err)
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	err = s.postService.GetCommentInfo(&reaction.Comment)
	if err != nil {
		log.Println("Get comment info from comment Like Error", err)
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	http.Redirect(w, r, "/post?ID="+strconv.Itoa(int(reaction.Comment.PostID)), http.StatusFound)
}

func (s *ServiceServer) CommentDislike(w http.ResponseWriter, r *http.Request, session model.Session) {
	if r.Method != http.MethodGet {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusMethodNotAllowed))
		return
	}

	commentID, err := s.getID(r)
	if err != nil {
		if err.Error() == "ID not set" {
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadRequest))
			return
		}
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	reaction := model.CommentReaction{Comment: model.Comment{ID: int64(commentID)}, User: session.User}
	err = s.postService.CommentSetDislike(&reaction)
	if err != nil {
		if err.Error() == "comment doesn't exist" {
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadRequest))
			return
		}
		log.Println("comment dislike was failed", err)
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	err = s.postService.GetCommentInfo(&reaction.Comment)
	if err != nil {
		log.Println("Get comment info from comment Dislike Error", err)
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	http.Redirect(w, r, "/post?ID="+strconv.Itoa(int(reaction.Comment.PostID)), http.StatusFound)
}
