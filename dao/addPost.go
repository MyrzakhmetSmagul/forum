package dao

import (
	"database/sql"
	"fmt"
	"forum/models"
	"log"
)

func AddPost(db *sql.DB, post *models.Post) error {
	queryText := `INSERT INTO posts (title, content, user_id)
	VALUES (?, ?, ?) `

	query, err := db.Prepare(queryText)
	if err != nil {
		log.Panicln(err.Error())
		return err
	}
	_, err = query.Exec(post.Title, post.Content)
	if err != nil {
		log.Panic(err.Error())
		return err
	}

	fmt.Println("###########################\n\nADD POST SUCCESFULLY\n\n########################")
	return err
}
