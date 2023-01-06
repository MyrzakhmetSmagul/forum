package repository

import (
	"database/sql"
	"log"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

type CommentQuery interface {
	CreateComment(comment *model.Comment) error
	GetPostComments(post *model.Post) error
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

	defer query.Close()

	result, err := query.Exec(comment.PostID, comment.UserID, comment.Username, comment.Message)
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

func (c *commentQuery) GetPostComments(post *model.Post) error {
	sqlStmt := `SELECT comment_id, user_id, username, message FROM comments WHERE post_id=?`
	query, err := c.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}

	defer query.Close()

	rows, err := query.Query(post.ID)
	if err != nil {
		log.Println(err)
		return err
	}

	defer rows.Close()

	for rows.Next() {
		comment := model.Comment{}
		err = rows.Scan(&comment.ID, &comment.UserID, &comment.Username, &comment.Message)
		if err != nil {
			log.Println(err)
			return err
		}

		post.Comments = append(post.Comments, comment)
	}

	return nil
}
