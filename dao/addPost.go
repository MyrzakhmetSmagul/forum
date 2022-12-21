package dao

import (
	"database/sql"
	"fmt"
	"forum/models"
	"io/ioutil"
	"log"
)

func AddPost(db *sql.DB, post *models.Post) error {
	data, err := ioutil.ReadFile("./dao/sqlQuery/postInsert.sql")
	if err != nil {
		log.Panicln(err.Error())
		return err
	}
	query, err := db.Prepare(string(data))
	if err != nil {
		log.Panicln(err.Error())
		return err
	}
	_, err = query.Exec(post.Title, post.Content)
	if err != nil {
		log.Panic(err.Error())
		return err
	}

	fmt.Println("###########################\n\nADD POST SUCCESFULLY\n\n########################")
	return err
}
