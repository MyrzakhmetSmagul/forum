package app

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

func (s *ServiceServer) render(w http.ResponseWriter, pageName string, statusCode int, data interface{}) {
	path := fmt.Sprintf("./templates/html/%s.html", pageName)
	t, err := template.ParseFiles(path)
	if err != nil {
		log.Println("ERROR:\nrender:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	w.WriteHeader(statusCode)
	if t.ExecuteTemplate(w, pageName, data) != nil {
		log.Println("ERROR:\nrender:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
	}
}
