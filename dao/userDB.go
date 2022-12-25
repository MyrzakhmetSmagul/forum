package dao

import (
	"database/sql"
	"forum/models"
	"log"
	"strconv"
)

func AddUser(db *sql.DB, u *models.User) error {
	sqlStmt := `INSERT INTO users(uname, email, pwd) 
	VALUES(?, ?, ?)`

	query, err := db.Prepare(sqlStmt)
	if err != nil {
		return err
	}

	hashedPassword, err := models.PasswordHashing(u.Pwd)
	if err != nil {
		return err
	}

	result, err := query.Exec(u.UName, u.Email, hashedPassword)
	if err != nil {
		return err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.Id = strconv.FormatInt(lastInsertId, 10)

	log.Println("INSERT INTO OK")
	return nil
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

func IsExistUser(db *sql.DB, u *models.User) (bool, error) {
	sqlStmt := `SELECT EXISTS (SELECT 1 FROM users WHERE uname=? OR email=?)`
	var exist bool
	err := db.QueryRow(sqlStmt, u.UName, u.Email).Scan(&exist)
	if err != nil {
		log.Println(err)
		return false, err
	}

	return exist, nil
}
