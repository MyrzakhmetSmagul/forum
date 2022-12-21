package dao

import (
	"database/sql"
	"forum/models"
	"io/ioutil"
	"log"
)

func AllPosts(db *sql.DB) ([]models.Post, error) {
	var posts []models.Post
	var post models.Post
	data, err := ioutil.ReadFile("./dao/sqlQuery/allPosts.sql")
	if err != nil {
		log.Println(err.Error())
		return posts, err
	}

	rows, err := db.Query(string(data))
	if err != nil {
		log.Println(err.Error())
		return posts, err
	}

	for rows.Next() {
		rows.Scan(&post.Id, &post.Title, &post.Content)
		posts = append(posts, post)
	}
	return posts, nil
}
