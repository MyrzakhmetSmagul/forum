package app

import (
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

func (s *ServiceServer) NewPost(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("authToken"); err != nil {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusForbidden, StatusText: http.StatusText(http.StatusForbidden)})
		return
	}

	if r.Method == http.MethodGet {
		s.GetNewPost(w, r)
		return
	}
	s.PostNewPost(w, r)
}

func (s *ServiceServer) GetNewPost(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./templates/html/create-post.html")
	if err != nil {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	err = t.ExecuteTemplate(w, "create-post", nil)
	if err != nil {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}
}

func (s *ServiceServer) PostNewPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusMethodNotAllowed, StatusText: http.StatusText(http.StatusMethodNotAllowed)})
		return
	}

	err := r.ParseForm()
	if err != nil {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	cookie, err := r.Cookie("authToken")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	session := model.Session{Token: cookie.Value}
	err = s.sessionService.GetSession(&session)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	user := model.User{ID: session.User.ID}
	err = s.userService.GetUserInfo(&user)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			err = s.sessionService.DeleteSession(&session)
			if err != nil {
				s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
				return
			}

			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	post := model.Post{User: user, Title: r.PostFormValue("title"), Content: r.PostFormValue("content")}
	err = s.postService.CreatePost(&post)
	if err != nil {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *ServiceServer) Post(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("authToken")
	if err != nil {
		s.PostUnauth(w, r)
		return
	}

	if r.Method != http.MethodGet {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusMethodNotAllowed, StatusText: http.StatusText(http.StatusMethodNotAllowed)})
		return
	}

	if r.URL.Query().Get("ID") == "" {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusBadRequest, StatusText: http.StatusText(http.StatusBadRequest)})
		return
	}

	session := model.Session{Token: cookie.Value}
	err = s.sessionService.GetSession(&session)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			s.PostUnauth(w, r)
			return
		}

		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("ID"))
	if err != nil {
		log.Println(err)
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	post := model.Post{ID: int64(id)}
	err = s.postService.GetPost(&post)
	if err != nil {
		if err.Error() != "getPost: sql: no rows in result set" {
			s.ErrorHandler(w, model.Error{StatusCode: http.StatusBadRequest, StatusText: http.StatusText(http.StatusBadRequest)})
			return
		}
		log.Println("error", err)
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	t, err := template.ParseFiles("./templates/html/post.html")
	if err != nil {
		log.Println("template parse error", err)
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	t.ExecuteTemplate(w, "post", post)
}

func (s *ServiceServer) PostUnauth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusMethodNotAllowed, StatusText: http.StatusText(http.StatusMethodNotAllowed)})
		return
	}

	if r.URL.Query().Get("ID") == "" {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusBadRequest, StatusText: http.StatusText(http.StatusBadRequest)})
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("ID"))
	if err != nil {
		log.Println(err)
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	post := model.Post{ID: int64(id)}
	err = s.postService.GetPost(&post)
	if err != nil {
		if err.Error() != "getPost: sql: no rows in result set" {
			s.ErrorHandler(w, model.Error{StatusCode: http.StatusBadRequest, StatusText: http.StatusText(http.StatusBadRequest)})
			return
		}
		log.Println("error", err)
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	t, err := template.ParseFiles("./templates/html/unauth-view-post.html")
	if err != nil {
		log.Println("template parse error", err)
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	t.ExecuteTemplate(w, "unauth-view-post", post)
}
