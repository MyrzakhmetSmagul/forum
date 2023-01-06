package repository

import (
	"database/sql"
	"log"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

type CategoryQuery interface {
	GetPostCategories(post *model.Post) error
	SetPostCategory(post *model.Post) error
}

type categoryQuery struct {
	db *sql.DB
}

func (c *categoryQuery) GetPostCategories(post *model.Post) error {
	sqlStmt := `SELECT pc.category_id, c.category
	FROM categories c
	INNER JOIN post_categories pc
	ON pc.category_id = c.category_id
	WHERE pc.post_id=?`

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
		var category model.Category
		rows.Scan(&category.ID, &category.Category)
		post.Categories = append(post.Categories, category)
	}

	return nil
}

func (c *categoryQuery) SetPostCategory(post *model.Post) error {
	sqlStmt := `INSERT INTO post_categories(post_id, category_id) VALUES(?,?)`
	query, err := c.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}

	defer query.Close()

	for i := 0; i < len(post.Categories); i++ {
		_, err = query.Exec(post.Categories[i].ID, post.Categories[i].Category)
		if err != nil {
			log.Printf("can't set post categories,\nPost_id = %d, category = %s, number = %d\n ERROR: %s", post.ID, post.Categories[i].Category, i, err.Error())
			return err
		}
	}

	return nil
}

func (c *categoryQuery) CreateCategory(category *model.Category) error {
	sqlStmt := `INSERT INTO categories(category)VALUES(?)`
	query, err := c.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}

	defer query.Close()

	res, err := query.Exec(category.ID)
	if err != nil {
		log.Println(err)
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return err
	}

	category.ID = id
	return nil
}
