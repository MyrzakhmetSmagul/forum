package app

import (
	"errors"
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
	allCategories, err := s.postService.GetAllCategories()
	if err != nil {
		log.Println("ERROR:\ngetNewPost:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	s.render(w, "create-post", http.StatusOK, allCategories)
}

func (s *ServiceServer) PostNewPost(w http.ResponseWriter, r *http.Request, user *model.User) {
	if r.Method != http.MethodPost {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusMethodNotAllowed))
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Println("ERROR:\npostNewPost:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadGateway))
		return
	}

	post := model.Post{User: *user, Title: r.PostFormValue("title"), Content: r.PostFormValue("content")}

	categories := r.Form["categories"]
	allCategories, err := s.postService.GetAllCategories()
	if err != nil {
		log.Println("ERROR:\npostNewPost:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	for i := 0; i < len(categories); i++ {
		status := false
		for j := 0; j < len(allCategories); j++ {
			if categories[i] == allCategories[j].Category {
				post.Categories = append(post.Categories, allCategories[j])
				status = true
				break
			}
		}

		if !status {
			log.Printf("error create post without categories or not exists category: '%s'", categories[i])

			s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadRequest))
			return
		}
	}

	err = s.postService.CreatePost(&post)
	if err != nil {
		log.Println("ERROR:\npostNewPost:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *ServiceServer) Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusMethodNotAllowed))
		return
	}

	_, err := s.getSession(r)
	if err != nil {
		if errors.Is(err, model.ErrNoSession) || errors.Is(err, model.ErrUserNotFound) {
			s.PostUnauth(w, r)
			return
		}
		log.Println("ERROR:\npost:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	post, err := s.getPost(r)
	if err != nil {
		log.Println("ERROR:\npost:", err)

		if errors.Is(err, model.ErrPostNotFound) {
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadRequest))
			return
		}

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	s.render(w, "post", http.StatusOK, post)
}

func (s *ServiceServer) PostUnauth(w http.ResponseWriter, r *http.Request) {
	post, err := s.getPost(r)
	if err != nil {
		log.Println("ERROR:\npostUnauth:", err)

		if errors.Is(err, model.ErrPostNotFound) {
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadRequest))
			return
		}
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	s.render(w, "unauth-view-post", http.StatusOK, post)
}
