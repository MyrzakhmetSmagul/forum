package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

type PostReactionQuery interface {
	PostSetLike(reaction *model.PostReaction) error
	PostSetDislike(reaction *model.PostReaction) error
	GetUserReactionToPost(reaction *model.PostReaction) error
	GetPostReactions(post *model.Post) error
}

type postReactionQuery struct {
	db *sql.DB
}

func (pr *postReactionQuery) createReactionToPost(reaction *model.PostReaction) error {
	sqlStmt := `INSERT INTO posts_likes_dislikes(post_id,user_id, like, dislike)VALUES(?,?,?,?)`
	query, err := pr.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}

	defer query.Close()

	res, err := query.Exec(reaction.Post.ID, reaction.User.ID, 0, 0)
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

func (pr *postReactionQuery) PostSetLike(reaction *model.PostReaction) error {
	var sqlStmt string
	err := pr.GetUserReactionToPost(reaction)
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

func (pr *postReactionQuery) PostSetDislike(reaction *model.PostReaction) error {
	var sqlStmt string
	err := pr.GetUserReactionToPost(reaction)
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

func (pr *postReactionQuery) GetUserReactionToPost(reaction *model.PostReaction) error {
	sqlStmt := `SELECT id, like, dislike FROM posts_likes_dislikes WHERE post_id=? and user_id=?`
	query, err := pr.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}

	defer query.Close()

	err = query.QueryRow(reaction.Post.ID, reaction.User.ID).Scan(&reaction.ID, &reaction.Like, &reaction.Dislike)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return pr.createReactionToPost(reaction)
		}

		log.Println(err)
		return err
	}
	return nil
}

func (p *postReactionQuery) GetPostReactions(post *model.Post) error {
	sqlStmt := `SELECT SUM(like), SUM(dislike) FROM posts_likes_dislikes WHERE comment_id=?`
	query, err := p.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}

	defer query.Close()

	err = query.QueryRow(post.ID).Scan(&post.Like, &post.Dislike)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
