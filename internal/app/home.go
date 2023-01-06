package app

import (
	"log"
	"net/http"
	"text/template"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

func (s *ServiceServer) IndexWihtoutSession(w http.ResponseWriter, r *http.Request) {
	log.Println("indexWithoutSession1")
	if r.Method != http.MethodGet {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusMethodNotAllowed, StatusText: http.StatusText(http.StatusMethodNotAllowed)})
		return
	}
	log.Println("indexWithoutSession2")

	if r.URL.Path != "/" {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusNotFound, StatusText: http.StatusText(http.StatusNotFound)})
		return
	}
	log.Println("indexWithoutSession3")

	t, err := template.ParseFiles("./template/html/withoutSession/index.html")
	if err != nil {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}
	log.Println("indexWithoutSession4")

	categories, err := s.postService.GetAllCategory()

	if err != nil {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}
	log.Println("indexWithoutSession5")

	posts, err := s.postService.GetAllPosts()
	if err != nil {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}
	log.Println("indexWithoutSession6")

	data := model.Data{Categories: categories, Posts: posts}

	err = t.Execute(w, data)
	if err != nil {
		log.Println(err)

		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
		return
	}
	log.Println("indexWithoutSession7")

}
