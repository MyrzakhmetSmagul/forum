package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

type CommentReactionQuery interface {
	CommentLike(reaction *model.PostReaction) error
	CommentDislike(reaction *model.PostReaction) error
}

type commentReactionQuery struct {
	db *sql.DB
}

func (cr *commentReactionQuery) CommentLike(reaction *model.CommentReaction) error {
	var (
		sqlStmt string
		err     error
	)

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

func (cr *commentReactionQuery) CommentDislike(reaction *model.CommentReaction) error {
	var (
		sqlStmt string
		err     error
	)

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

	err = query.QueryRow(reaction.Comment.ID, reaction.User.ID).Scan(&reaction.ID, &reaction.Like, &reaction.Dislike)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
