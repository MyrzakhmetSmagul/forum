package repository

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func NewDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = createTable(db)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println("database was created successfully")
	return db, nil
}
func createTable(db *sql.DB) error {
	query := []string{}

	users := `
	CREATE TABLE IF NOT EXISTS users(
		user_id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		email TEXT NOT NULL,
		password TEXT NOT NULL
	)
	`
	posts := `
	CREATE TABLE IF NOT EXISTS posts(
		post_id INTEGER PRIMARY KEY AUTOINCREMENT,
		category_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		username TEXT NOT NULL,
		like INTEGER NOT NULL,
		dislike INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users (user_id),
		FOREIGN KEY (category_id) REFERENCES categories (category_id)
	)
	`
	categories := `
	CREATE TABLE IF NOT EXISTS categories(
		category_id INTEGER PRIMARY KEY AUTOINCREMENT,
		category TEXT NOT NULL
	)
	`

	session := `
	CREATE TABLE IF NOT EXISTS sessions(
		session_id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		token TEXT NOT NULL,
		expiry DATE NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(user_id)	ON DELETE CASCADE
	)
	`

	comments := `
	CREATE TABLE IF NOT EXISTS comments(
		comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		username TEXT NOT NULL,
		message TEXT NOT NULL,
		like INTEGER NOT NULL,
		dislike INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(user_id),
		FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE
	)
	`

	postLikes := `
	CREATE TABLE IF NOT EXISTS post_likes(
		user_id INTEGER NOT NULL,
		username TEXT NOT NULL,
		post_id INTEGER NOT NULL,
		status INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(user_id),
		FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE
	)
	`

	postDislikes := `
	CREATE TABLE IF NOT EXISTS post_dislikes(
		user_id INTEGER NOT NULL,
		username TEXT NOT NULL,
		post_id INTEGER NOT NULL,
		status INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(user_id),
		FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE
	)
	`

	commentLikes := `CREATE TABLE IF NOT EXISTS comment_likes(
		user_id INTEGER NOT NULL,
		username TEXT NOT NULL, 
		comment_id INTEGER NOT NULL,
		status INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(user_id),
		FOREIGN KEY (comment_id) REFERENCES comments(comment_id) ON DELETE CASCADE
	)
	`
	commentDislikes := `CREATE TABLE IF NOT EXISTS comment_dislikes(
		user_id INTEGER NOT NULL,
		username TEXT NOT NULL, 
		comment_id INTEGER NOT NULL,
		status INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(user_id),
		FOREIGN KEY (comment_id) REFERENCES comments(comment_id) ON DELETE CASCADE
	)
	`

	query = append(query, users, posts, categories, session, comments, postLikes, postDislikes, commentDislikes, commentLikes)
	for _, v := range query {
		_, err := db.Exec(v)
		if err != nil {
			return err
		}
	}
	return nil
}
