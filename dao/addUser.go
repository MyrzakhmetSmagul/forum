package dao

import (
	"database/sql"
	"fmt"
	"forum/models"
	"log"
)

func AddUser(db *sql.DB, u *models.User) error {
	queryText := `INSERT INTO users(uname, email, pwd) 
	VALUES(?, ?, ?)`

	query, err := db.Prepare(queryText)
	if err != nil {
		fmt.Printf("##################################################\n%s\n##################################################\n", err)
		return err
	}

	_, err = query.Exec(u.UName, u.Email, u.Pwd)
	if err != nil {
		fmt.Printf("##################################################\n%s\n##################################################\n", err.Error())
		return err
	}
	log.Println("INSERT INTO OK")
	return nil
}
