package app

import (
	"net/http"
	"text/template"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

func (s *ServiceServer) ErrorHandler(w http.ResponseWriter, errorStatus model.Error) {
	t, err := template.ParseFiles("./templates/html/error.html")
	if err != nil {
		http.Error(w, errorStatus.StatusText, errorStatus.StatusCode)
		return
	}
	w.WriteHeader(errorStatus.StatusCode)
	t.ExecuteTemplate(w, "error", errorStatus)
}
