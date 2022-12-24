package dao

import (
	"database/sql"
	"fmt"
	"forum/models"
	"log"
)

func AddUser(db *sql.DB, u *models.User) error {
	sqlStmt := `INSERT INTO users(uname, email, pwd) 
	VALUES(?, ?, ?)`

	query, err := db.Prepare(sqlStmt)
	if err != nil {
		fmt.Printf("##################################################\n%s\n##################################################\n", err)
		return err
	}

	u.Pwd, err = models.GetPassordHash(u.Pwd)
	if err != nil {
		fmt.Printf("##################################################\n%s\n##################################################\n", err.Error())
		return err
	}

	_, err = query.Exec(u.UName, u.Email, u.Pwd)
	if err != nil {
		fmt.Printf("##################################################\n%s\n##################################################\n", err.Error())
		return err
	}

	log.Println("INSERT INTO OK")

	sqlStmt = `SELECT user_id WHERE uname=?`
	query, err = db.Prepare(sqlStmt)
	err = query.QueryRow(u.UName).Scan(&u.Id)
	return err
}

func UserVerification(db *sql.DB, u *models.User) bool {
	sqlStmt := `SELECT user_id, uname  FROM users WHERE email=? AND pwd=?`
	err := db.QueryRow(sqlStmt, u.Email, u.Pwd).Scan(&u.Id, &u.UName)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func SearchUser(db *sql.DB, u *models.User) bool {
	sqlStmt := `SELECT 1 FROM users WHERE uname=? OR email=? LIMIT 1`
	exist := 0
	err := db.QueryRow(sqlStmt, u.UName, u.Email).Scan(&exist)
	if err != nil {
		log.Println(err)
		return false
	}

	return exist == 1
}
