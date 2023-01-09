package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

type CommentQuery interface {
	CreateComment(comment *model.Comment) error
	GetPostComments(post *model.Post) error
	CommentSetLike(reaction *model.CommentReaction) error
	CommentSetDislike(reaction *model.CommentReaction) error
	GetCommentInfo(comment *model.Comment) error
}

type commentQuery struct {
	db *sql.DB
}

func (c *commentQuery) CreateComment(comment *model.Comment) error {
	sqlStmt := `INSERT INTO comments (post_id, user_id, username, message) 
	SELECT post_id, ?, ?, ?
	FROM posts
	WHERE EXISTS (SELECT * FROM posts WHERE post_id=?) AND post_id=?`
	query, err := c.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}

	defer query.Close()

	result, err := query.Exec(comment.UserID, comment.Username, comment.Message, comment.PostID, comment.PostID)
	if err != nil {
		log.Println(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("CreateComment result.RowsAffected error", err)
		return err
	}

	if rowsAffected == 0 {
		log.Println("post doesn't exist")
		return errors.New("post doesn't exist")
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
		log.Println("GetPostComments", err)
		return err
	}

	defer query.Close()

	rows, err := query.Query(post.ID)
	if err != nil {
		log.Println("GetPostComments", err)
		return err
	}

	defer rows.Close()

	for rows.Next() {
		comment := model.Comment{}
		err = rows.Scan(&comment.ID, &comment.UserID, &comment.Username, &comment.Message)
		if err != nil {
			log.Println("GetPostComments", err)
			return err
		}

		err = c.getCommentLikesDislikes(&comment)
		if err != nil {
			log.Println("getCOmmentLikesDislikes", err)
			return err
		}

		post.Comments = append(post.Comments, comment)
	}

	return nil
}

func (c *commentQuery) GetCommentInfo(comment *model.Comment) error {
	sqlStmt := `SELECT * FROM comments WHERE comment_id=?`
	err := c.db.QueryRow(sqlStmt, comment.ID).Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Username, &comment.Message)
	if err != nil {
		log.Println("Get Comment Info Error", err)
		return err
	}

	return nil
}
