package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

type PostReactionQuery interface {
	PostLike(reaction *model.PostReaction) error
	PostDislike(reaction *model.PostReaction) error
}

type postReactionQuery struct {
	db *sql.DB
}

func (pr *postReactionQuery) PostLike(reaction *model.PostReaction) error {
	sqlStmt := `SELECT id, like, dislike FROM posts_likes_dislikes WHERE post_id=? and user_id=?`
	query, err := pr.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}

	err = query.QueryRow(reaction.Post.ID, reaction.User.ID).Scan(&reaction.ID, &reaction.Like, &reaction.Dislike)
	if err != nil {
		log.Println(err)
		return err
	}

	if reaction.Like == reaction.Dislike {
		sqlStmt = `UPDATE posts_likes_dislikes SET like=1 WHERE Id=?`
		err = pr.updatePostReaction(sqlStmt, pr.db, reaction)
	} else if reaction.Like == 0 {
		sqlStmt = `UPDATE posts_likes_dislikes SET like=1, dislike=0 WHERE Id=?`
		err = pr.updatePostReaction(sqlStmt, pr.db, reaction)
	} else {
		sqlStmt = `UPDATE posts_likes_dislikes SET like=0 WHERE Id=?`
		err = pr.updatePostReaction(sqlStmt, pr.db, reaction)
	}
	if err != nil {
		log.Println(err)
		return err
	}
	sqlStmt = `SELECT like, dislike FROM posts_likes_dislikes WHERE Id=?`
	err = pr.db.QueryRow(sqlStmt, reaction.ID).Scan(&reaction.Like, &reaction.Dislike)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (pr *postReactionQuery) PostDislike(reaction *model.PostReaction) error {
	sqlStmt := `SELECT id, like, dislike FROM posts_likes_dislikes WHERE post_id=? and user_id=?`
	query, err := pr.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}

	err = query.QueryRow(reaction.Post.ID, reaction.User.ID).Scan(&reaction.ID, &reaction.Like, &reaction.Dislike)
	if err != nil {
		log.Println(err)
		return err
	}

	if reaction.Like == reaction.Dislike {
		sqlStmt = `UPDATE posts_likes_dislikes SET dislike=1 WHERE Id=?`
		err = pr.updatePostReaction(sqlStmt, pr.db, reaction)
	} else if reaction.Dislike == 0 {
		sqlStmt = `UPDATE posts_likes_dislikes SET like=0, dislike=1 WHERE Id=?`
		err = pr.updatePostReaction(sqlStmt, pr.db, reaction)
	} else {
		sqlStmt = `UPDATE posts_likes_dislikes SET dislike=0 WHERE Id=?`
		err = pr.updatePostReaction(sqlStmt, pr.db, reaction)
	}
	if err != nil {
		log.Println(err)
		return err
	}
	sqlStmt = `SELECT like, dislike FROM posts_likes_dislikes WHERE Id=?`
	err = pr.db.QueryRow(sqlStmt, reaction.ID).Scan(&reaction.Like, &reaction.Dislike)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (pr *postReactionQuery) updatePostReaction(sqlStmt string, db *sql.DB, reaction *model.PostReaction) error {
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
