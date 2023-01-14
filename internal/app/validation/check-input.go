package validation

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

func CheckInput(r *http.Request, post *model.Post, allCategories []model.Category) error {
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("checkInput: %w", err)
	}

	post.Title = r.PostFormValue("title")

	fmt.Println("##################\ntitle:", post.Title)

	if temp := strings.Trim(post.Title, " "); temp == "" || len(post.Title) > 255 {
		log.Println("temp := strings.Trim(post.Title, \" \")", len(post.Title) > 255, temp == "")
		return fmt.Errorf("checkInput: %w", model.ErrMessageInvalid)
	}

	post.Content = r.PostFormValue("content")
	log.Println("##################\ncontent:", post.Content)
	if temp := strings.Trim(post.Content, " "); temp == "" {
		log.Println("temp := strings.Trim(post.Content, \" \");", temp == "")
		return fmt.Errorf("checkInput: %w", model.ErrMessageInvalid)
	}

	categories := r.Form["categories"]
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
			return fmt.Errorf("checkInput: %w", model.ErrMessageInvalid)
		}
	}

	return nil
}
