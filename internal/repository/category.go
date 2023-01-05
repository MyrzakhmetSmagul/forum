package repository

import (
	"database/sql"
	"log"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

type CategoryQuery interface {
}

type categoryQuery struct {
	db *sql.DB
}

func (c *categoryQuery) GetCategory(post *model.Post) error {
	sqlStmt := `SELECT category_id FROM post_categories WHERE post_id=?`
	query, err := c.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}

	defer query.Close()
	return nil
}
