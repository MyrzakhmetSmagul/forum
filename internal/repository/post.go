package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

type PostQuery interface {
	CreatePost(post *model.Post) error
	GetPost(post *model.Post) error
	GetAllPosts() ([]model.Post, error)
	PostSetLike(reaction *model.PostReaction) error
	PostSetDislike(reaction *model.PostReaction) error
	SetPostCategory(post *model.Post) error
	CreateCategory(category *model.Category) error
	GetAllCategory() ([]model.Category, error)
	GetPostsOfCategory(category model.Category) ([]model.Post, error)
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

	defer query.Close()

	result, err := query.Exec(post.Title, post.Content, post.User.ID, post.User.Username)
	if err != nil {
		log.Println(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("CREATE POST rows affected error:", err)
		return err
	}

	if rowsAffected == 0 {
		log.Println("create post was failed")
		return errors.New("create post was failed")
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return err
	}

	post.ID = id
	return p.SetPostCategory(post)
}

func (p *postQuery) GetPost(post *model.Post) error {
	sqlStmt := `SELECT title,content, user_id, username FROM posts WHERE post_id=?`
	query, err := p.db.Prepare(sqlStmt)
	if err != nil {
		log.Println("getPost", err)
		return err
	}

	defer query.Close()

	err = query.QueryRow(post.ID).Scan(&post.Title, &post.Content, &post.User.ID, &post.User.Username)
	if err != nil {
		log.Println("getPost", err)
		return errors.New("getPost: " + err.Error())
	}

	log.Println("get post category started")
	err = p.GetPostCategories(post)
	if err != nil {
		log.Println("getPost", err)
		return err
	}

	log.Println("get post category finished")
	log.Println("category", post.Categories[0].Category)

	err = p.GetPostLikesDislikes(post)
	if err != nil {
		log.Println("getPost", err)
		return err
	}

	return nil
}

func (p *postQuery) GetAllPosts() ([]model.Post, error) {
	sqlStmt := `SELECT * FROM posts`
	rows, err := p.db.Query(sqlStmt)
	if err != nil {
		log.Println(err)
		return []model.Post{}, err
	}

	defer rows.Close()

	posts := []model.Post{}
	for rows.Next() {
		post := model.Post{}
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.User.ID, &post.User.Username)
		if err != nil {
			log.Println(err)
			return []model.Post{}, err
		}

		err = p.GetPostCategories(&post)
		if err != nil {
			log.Println("Get post Categories", err)
			return []model.Post{}, err
		}

		err = p.GetPostLikesDislikes(&post)
		if err != nil {
			log.Println(err)
			return []model.Post{}, err
		}

		posts = append(posts, post)
	}
	return posts, nil
}

func (p *postQuery) GetPostsOfCategory(category model.Category) ([]model.Post, error) {
	sqlStmt := `SELECT posts.post_id, posts.title, posts.content, posts.user_id, posts.username FROM posts
	INNER JOIN post_categories ON posts.post_id = post_categories.post_id
	INNER JOIN categories ON post_categories.category_id = categories.category_id
	WHERE categories.category_id = ?`

	query, err := p.db.Prepare(sqlStmt)
	if err != nil {
		log.Println("GetPostsOfCategory ERROR:", err)
		return []model.Post{}, err
	}

	defer query.Close()

	posts := []model.Post{}
	rows, err := query.Query(category.ID)
	if err != nil {
		log.Println("GetPostsOfCategory ERROR:", err)
		return []model.Post{}, err
	}

	for rows.Next() {
		post := model.Post{}
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.User.ID, &post.User.Username)
		if err != nil {
			log.Println("GetPostsOfCategory ERROR:", err)
			return []model.Post{}, err
		}

		err = p.GetPostCategories(&post)
		if err != nil {
			log.Println("GetPostsOfCategory ERROR:", err)
			return []model.Post{}, err
		}

		err = p.GetPostLikesDislikes(&post)
		if err != nil {
			log.Println("GetPostsOfCategory ERROR:", err)
			return []model.Post{}, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}
