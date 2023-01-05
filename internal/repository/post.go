package repository

import (
	"database/sql"
	"log"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

type PostQuery interface {
	CreatePost(post *model.Post) error
	GetPost(post *model.Post) error
}

type postQuery struct {
	db *sql.DB
}

func (p *postQuery) CreatePost(post *model.Post) error {
	sqlStmt := `INSERT INTO posts (title, content, user_id, username)VALUES(?,?,?,?)`
	query, err := p.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}

	result, err := query.Exec(post.Title, post.Content, post.User.ID, post.User.Username)
	if err != nil {
		log.Println(err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return err
	}

	post.ID = id
	return nil
}

func (p *postQuery) GetPost(post *model.Post) error {
	sqlStmt := `SELECT title,content, user_id, username FROM posts WHERE post_id=?`
	query, err := p.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}

	err = query.QueryRow(post.ID).Scan(&post.Title, &post.Content, post.User.ID, post.User.Username)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
