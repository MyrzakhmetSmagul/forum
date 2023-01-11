package app

import (
	"html/template"
	"log"
	"net/http"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

func (s *ServiceServer) NewPost(w http.ResponseWriter, r *http.Request, session model.Session) {
	if r.Method == http.MethodGet {
		s.GetNewPost(w, r)
		return
	}
	s.PostNewPost(w, r, &session.User)
}

func (s *ServiceServer) GetNewPost(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./templates/html/create-post.html")
	if err != nil {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	allCategories, err := s.postService.GetAllCategory()
	if err != nil {
		log.Println("get all categories error", err)
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	err = t.ExecuteTemplate(w, "create-post", allCategories)
	if err != nil {
		log.Println("template error")
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}
}

func (s *ServiceServer) PostNewPost(w http.ResponseWriter, r *http.Request, user *model.User) {
	if r.Method != http.MethodPost {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusMethodNotAllowed, StatusText: http.StatusText(http.StatusMethodNotAllowed)})
		return
	}

	if r.ParseForm() != nil {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusBadGateway, StatusText: http.StatusText(http.StatusBadGateway)})
		return
	}

	post := model.Post{User: *user, Title: r.PostFormValue("title"), Content: r.PostFormValue("content")}

	categories := r.Form["categories"]
	allCategories, err := s.postService.GetAllCategory()
	if err != nil {
		log.Println("get all categories error", err)
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	status := false
	for i := 0; i < len(categories); i++ {
		for j := 0; j < len(allCategories); j++ {
			if categories[i] == allCategories[j].Category {
				post.Categories = append(post.Categories, allCategories[j])
				status = true
				break
			}
		}
	}

	if !status {
		log.Println("error create post without categories")
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusBadRequest, StatusText: http.StatusText(http.StatusBadRequest)})
		return
	}

	err = s.postService.CreatePost(&post)
	if err != nil {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *ServiceServer) Post(w http.ResponseWriter, r *http.Request, session model.Session) {
	if r.Method != http.MethodGet {
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

	post := model.Post{ID: int64(postID)}
	err = s.postService.GetPost(&post)
	if err != nil {
		if err.Error() != "getPost sql: no rows in result set" {
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

	postID, err := s.getID(r)
	if err != nil {
		if err.Error() == "ID not set" {
			s.ErrorHandler(w, model.Error{StatusCode: http.StatusBadRequest, StatusText: http.StatusText(http.StatusBadRequest)})
			return
		}
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}

	post := model.Post{ID: int64(postID)}
	err = s.postService.GetPost(&post)
	if err != nil {
		if err.Error() != "getPost sql: no rows in result set" {
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
