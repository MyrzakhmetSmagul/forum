package app

import (
	"log"
	"net/http"
	"strconv"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

func (s *ServiceServer) PostLike(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("authToken")
	if err != nil {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusForbidden, StatusText: http.StatusText(http.StatusForbidden)})
		return
	}

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

	session := model.Session{Token: cookie.Value}
	err = s.sessionService.GetSession(&session)
	if err != nil {
		log.Println("PostLike Get Session Info error:", err)
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	user := model.User{ID: session.User.ID}
	err = s.userService.GetUserInfo(&user)
	if err != nil {
		log.Println("postLike cann't get user info by ID:", err)
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	post := model.Post{ID: int64(postID)}
	err = s.postService.PostLike(&model.PostReaction{Post: post, User: user})
	if err != nil {
		log.Println("Post like was failed", err)
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	http.Redirect(w, r, "/post?ID="+strconv.Itoa(postID), http.StatusFound)
}

func (s *ServiceServer) PostDislike(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("authToken")
	if err != nil {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusForbidden, StatusText: http.StatusText(http.StatusForbidden)})
		return
	}

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

	session := model.Session{Token: cookie.Value}
	err = s.sessionService.GetSession(&session)
	if err != nil {
		log.Println("PostLike Get Session Info error:", err)
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	user := model.User{ID: session.User.ID}
	err = s.userService.GetUserInfo(&user)
	if err != nil {
		log.Println("postDislike cann't get user info by ID:", err)
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	post := model.Post{ID: int64(postID)}
	err = s.postService.PostDislike(&model.PostReaction{Post: post, User: user})
	if err != nil {
		log.Println("Post like was failed", err)
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	http.Redirect(w, r, "/post?ID="+strconv.Itoa(postID), http.StatusFound)
}
