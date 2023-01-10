package repository

import (
	"database/sql"
	"io"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func NewDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = Migrations(db)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println("database was created successfully")
	return db, nil
}

func Migrations(db *sql.DB) error {
	openTables, err := os.Open("./migrations/init.sql")
	if err != nil {
		return err
	}
	allTables, err := io.ReadAll(openTables)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(allTables))
	if err != nil {
		return err
	}
	return createCategories(db)
}

func createCategories(db *sql.DB) error {
	sqlStmt := `INSERT INTO categories (category)
	SELECT * 
	FROM (SELECT ? as category) AS tmp 
	WHERE NOT EXISTS (SELECT category FROM categories WHERE category=?) LIMIT 1`
	categories := []string{"Adventure stories", "Crime", "Fantasy", "Humore and satire", "Mystery", "Plays", "Romance"}
	for i := 0; i < len(categories); i++ {
		_, err := db.Exec(sqlStmt, categories[i], categories[i])
		if err != nil {
			return err
		}
	}
	return nil
}
