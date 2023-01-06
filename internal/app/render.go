package app

import (
	"html/template"
	"net/http"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

func (s *ServiceServer) RenderPage(w http.ResponseWriter, path string, obj interface{}) error {
	t, err := template.ParseFiles(path)
	if err != nil {
		s.ErrorHandler(w, model.Error{StatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)})
	}

	t.Execute(w, obj)
	return nil
}
