package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

func (c *commentQuery) createReactionToComment(reaction *model.CommentReaction) error {
	sqlStmt := `INSERT INTO comments_likes_dislikes (comment_id, user_id, like, dislike) 
	SELECT comment_id, ?, 0, 0
	FROM comments
	WHERE EXISTS (SELECT * FROM posts WHERE comment_id=?)`
	query, err := c.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}

	defer query.Close()

	res, err := query.Exec(reaction.User.ID, reaction.Comment.ID)
	if err != nil {
		log.Println(err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("CreateReactionToComment result.RowsAffected error", err)
		return err
	}

	if rowsAffected == 0 {
		log.Println("comment doesn't exist")
		return errors.New("comment doesn't exist")
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return err
	}

	reaction.ID = id
	return nil
}

func (c *commentQuery) CommentSetLike(reaction *model.CommentReaction) error {
	var sqlStmt string
	err := c.getUserReactionToComment(reaction)
	if err != nil {
		log.Println(err)
		return err
	}

	if reaction.Like == reaction.Dislike {
		sqlStmt = `UPDATE comments_likes_dislikes SET like=1 WHERE Id=?`
		err = c.updateCommentReaction(sqlStmt, c.db, reaction)
	} else if reaction.Like == 0 {
		sqlStmt = `UPDATE comments_likes_dislikes SET like=1, dislike=0 WHERE Id=?`
		err = c.updateCommentReaction(sqlStmt, c.db, reaction)
	} else {
		sqlStmt = `UPDATE comments_likes_dislikes SET like=0 WHERE Id=?`
		err = c.updateCommentReaction(sqlStmt, c.db, reaction)
	}
	if err != nil {
		log.Println(err)
		return err
	}

	sqlStmt = `SELECT like, dislike FROM comments_likes_dislikes WHERE Id=?`
	err = c.db.QueryRow(sqlStmt, reaction.ID).Scan(&reaction.Like, &reaction.Dislike)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (c *commentQuery) CommentSetDislike(reaction *model.CommentReaction) error {
	var sqlStmt string
	err := c.getUserReactionToComment(reaction)
	if err != nil {
		log.Println(err)
		return err
	}

	if reaction.Like == reaction.Dislike {
		sqlStmt = `UPDATE comments_likes_dislikes SET dislike=1 WHERE Id=?`
		err = c.updateCommentReaction(sqlStmt, c.db, reaction)
	} else if reaction.Dislike == 0 {
		sqlStmt = `UPDATE comments_likes_dislikes SET like=0, dislike=1 WHERE Id=?`
		err = c.updateCommentReaction(sqlStmt, c.db, reaction)
	} else {
		sqlStmt = `UPDATE comments_likes_dislikes SET dislike=0 WHERE Id=?`
		err = c.updateCommentReaction(sqlStmt, c.db, reaction)
	}

	if err != nil {
		log.Println(err)
		return err
	}

	sqlStmt = `SELECT like, dislike FROM comments_likes_dislikes WHERE Id=?`
	err = c.db.QueryRow(sqlStmt, reaction.ID).Scan(&reaction.Like, &reaction.Dislike)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (c *commentQuery) updateCommentReaction(sqlStmt string, db *sql.DB, reaction *model.CommentReaction) error {
	result, err := db.Exec(sqlStmt, reaction.ID)
	if err != nil {
		log.Println(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}

	if rowsAffected == 0 {
		return errors.New("reaction set was failed")
	}
	return nil
}

func (c *commentQuery) getUserReactionToComment(reaction *model.CommentReaction) error {
	sqlStmt := `SELECT id, like, dislike FROM comments_likes_dislikes WHERE comment_id=? and user_id=?`
	query, err := c.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}

	defer query.Close()

	err = query.QueryRow(reaction.Comment.ID, reaction.User.ID).Scan(&reaction.ID, &reaction.Like, &reaction.Dislike)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return c.createReactionToComment(reaction)
		}

		log.Println(err)
		return err
	}
	return nil
}

func (c *commentQuery) getCommentLikesDislikes(comment *model.Comment) error {
	sqlStmt := `SELECT COALESCE(SUM(like), 0), COALESCE(SUM(dislike), 0) FROM comments_likes_dislikes WHERE comment_id=?`
	query, err := c.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}

	defer query.Close()

	err = query.QueryRow(comment.ID).Scan(&comment.Like, &comment.Dislike)
	if err != nil {
		log.Println("getCommentLikesDislikes", err)
		return err
	}

	return nil
}
