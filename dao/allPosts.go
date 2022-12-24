package dao

import (
	"database/sql"
	"forum/models"
	"log"
)

func AllPosts(db *sql.DB) ([]models.Post, error) {
	var posts []models.Post
	var post models.Post

	sqlStmt := `SELECT * from posts`

	rows, err := db.Query(sqlStmt)
	if err != nil {
		log.Println(err.Error())
		return posts, err
	}

	for rows.Next() {
		rows.Scan(&post.Id, &post.Title, &post.Content, &post.UserId)
		posts = append(posts, post)
	}
	return posts, nil
}
