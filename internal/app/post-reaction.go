package app

import (
	"log"
	"net/http"
	"strconv"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

func (s *ServiceServer) PostLike(w http.ResponseWriter, r *http.Request, session *model.Session) {
	if r.Method != http.MethodGet {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusMethodNotAllowed, StatusText: http.StatusText(http.StatusMethodNotAllowed)})
		return
	}

	if r.URL.Query().Get("ID") == "" {
		log.Println("PostID not set")
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusBadRequest, StatusText: http.StatusText(http.StatusBadRequest)})
		return
	}

	postID, err := strconv.Atoi(r.URL.Query().Get("ID"))
	if err != nil {
		log.Println("postID atoi error:", err)
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	post := model.Post{ID: int64(postID)}
	err = s.postService.PostLike(&model.PostReaction{Post: post, User: session.User})
	if err != nil {
		if err.Error() == "post doesn't exist" {
			s.ErrorHandler(w, model.Error{StatusCode: http.StatusBadRequest, StatusText: http.StatusText(http.StatusBadRequest)})
			return
		}
		log.Println("Post like was failed", err)
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	http.Redirect(w, r, "/post?ID="+strconv.Itoa(postID), http.StatusFound)
}

func (s *ServiceServer) PostDislike(w http.ResponseWriter, r *http.Request, session *model.Session) {
	if r.Method != http.MethodGet {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusMethodNotAllowed, StatusText: http.StatusText(http.StatusMethodNotAllowed)})
		return
	}

	if r.URL.Query().Get("ID") == "" {
		log.Println("PostID not set")
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusBadRequest, StatusText: http.StatusText(http.StatusBadRequest)})
		return
	}

	postID, err := strconv.Atoi(r.URL.Query().Get("ID"))
	if err != nil {
		log.Println("postID atoi error:", err)
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	post := model.Post{ID: int64(postID)}
	err = s.postService.PostDislike(&model.PostReaction{Post: post, User: session.User})
	if err != nil {
		if err.Error() == "post doesn't exist" {
			s.ErrorHandler(w, model.Error{StatusCode: http.StatusBadRequest, StatusText: http.StatusText(http.StatusBadRequest)})
			return
		}
		log.Println("Post like was failed", err)
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	http.Redirect(w, r, "/post?ID="+strconv.Itoa(postID), http.StatusFound)
}
