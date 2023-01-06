package repository

import (
	"database/sql"
	"log"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

type CommentQuery interface {
}

type commentQuery struct {
	db *sql.DB
}

func (c *commentQuery) CreateComment(comment *model.Comment) error {
	sqlStmt := `INSERT INTO comments(post_id, user_id, username, message) VALUES(?,?,?,?)`
	query, err := c.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}

	result, err := query.Exec(comment.Post.ID, comment.User.ID, comment.User.Username, comment.Message)
	if err != nil {
		log.Println(err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return err
	}

	comment.ID = id
	return nil
}
