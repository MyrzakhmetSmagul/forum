package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

type CommentReactionQuery interface {
	CommentSetLike(reaction *model.CommentReaction) error
	CommentSetDislike(reaction *model.CommentReaction) error
	GetCommentLikesDislikes(comment *model.Comment) error
}

type commentReactionQuery struct {
	db *sql.DB
}

func (cr *commentReactionQuery) createReactionToComment(reaction *model.CommentReaction) error {
	sqlStmt := `INSERT INTO comments_likes_dislikes(comment_id,user_id, like, dislike)VALUES(?,?,?,?)`
	query, err := cr.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}

	defer query.Close()

	res, err := query.Exec(reaction.Comment.ID, reaction.User.ID, 0, 0)
	if err != nil {
		log.Println(err)
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return err
	}

	reaction.ID = id
	return nil
}

func (cr *commentReactionQuery) CommentSetLike(reaction *model.CommentReaction) error {
	var sqlStmt string
	err := cr.GetUserReactionToComment(reaction)
	if err != nil {
		log.Println(err)
		return err
	}

	if reaction.Like == reaction.Dislike {
		sqlStmt = `UPDATE comments_likes_dislikes SET like=1 WHERE Id=?`
		err = cr.updateCommentReaction(sqlStmt, cr.db, reaction)
	} else if reaction.Like == 0 {
		sqlStmt = `UPDATE comments_likes_dislikes SET like=1, dislike=0 WHERE Id=?`
		err = cr.updateCommentReaction(sqlStmt, cr.db, reaction)
	} else {
		sqlStmt = `UPDATE comments_likes_dislikes SET like=0 WHERE Id=?`
		err = cr.updateCommentReaction(sqlStmt, cr.db, reaction)
	}
	if err != nil {
		log.Println(err)
		return err
	}
	sqlStmt = `SELECT like, dislike FROM comments_likes_dislikes WHERE Id=?`
	err = cr.db.QueryRow(sqlStmt, reaction.ID).Scan(&reaction.Like, &reaction.Dislike)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (cr *commentReactionQuery) CommentSetDislike(reaction *model.CommentReaction) error {
	var sqlStmt string
	err := cr.GetUserReactionToComment(reaction)
	if err != nil {
		log.Println(err)
		return err
	}

	if reaction.Like == reaction.Dislike {
		sqlStmt = `UPDATE comments_likes_dislikes SET dislike=1 WHERE Id=?`
		err = cr.updateCommentReaction(sqlStmt, cr.db, reaction)
	} else if reaction.Dislike == 0 {
		sqlStmt = `UPDATE comments_likes_dislikes SET like=0, dislike=1 WHERE Id=?`
		err = cr.updateCommentReaction(sqlStmt, cr.db, reaction)
	} else {
		sqlStmt = `UPDATE comments_likes_dislikes SET dislike=0 WHERE Id=?`
		err = cr.updateCommentReaction(sqlStmt, cr.db, reaction)
	}

	if err != nil {
		log.Println(err)
		return err
	}

	sqlStmt = `SELECT like, dislike FROM comments_likes_dislikes WHERE Id=?`
	err = cr.db.QueryRow(sqlStmt, reaction.ID).Scan(&reaction.Like, &reaction.Dislike)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (cr *commentReactionQuery) updateCommentReaction(sqlStmt string, db *sql.DB, reaction *model.CommentReaction) error {
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

func (cr *commentReactionQuery) GetUserReactionToComment(reaction *model.CommentReaction) error {
	sqlStmt := `SELECT id, like, dislike FROM comments_likes_dislikes WHERE post_id=? and user_id=?`
	query, err := cr.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}

	defer query.Close()

	err = query.QueryRow(reaction.Comment.ID, reaction.User.ID).Scan(&reaction.ID, &reaction.Like, &reaction.Dislike)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return cr.createReactionToComment(reaction)
		}

		log.Println(err)
		return err
	}
	return nil
}

func (c *commentReactionQuery) GetCommentLikesDislikes(comment *model.Comment) error {
	sqlStmt := `SELECT SUM(like), SUM(dislike) FROM comments_likes_dislikes WHERE comment_id=?`
	query, err := c.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}

	defer query.Close()

	err = query.QueryRow(comment.ID).Scan(&comment.Like, &comment.Dislike)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
