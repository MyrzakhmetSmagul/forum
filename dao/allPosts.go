package dao

import (
	"database/sql"
	"forum/models"
	"log"
)

func AllPosts(db *sql.DB) ([]models.Post, error) {
	posts := []models.Post{}
	sqlStmt := `SELECT * from posts`

	rows, err := db.Query(sqlStmt)
	if err != nil {
		log.Println(err.Error())
		return posts, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.PostId, &post.UserId, &post.Title, &post.Content)
		if err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}

	err = rows.Err()
	if err != nil {
		return posts, err
	}
	return posts, nil
}
