package dao

import (
	"database/sql"
	"fmt"
	"forum/models"
	"log"
)

func AddPost(db *sql.DB, post *models.Post) error {
	sqlStmt := `INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)`
	query, err := db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}
	defer query.Close()

	_, err = query.Exec(post.UserId, post.Title, post.Content)
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Println("###########################\n\nADD POST SUCCESFULLY\n\n########################")
	return nil
}
