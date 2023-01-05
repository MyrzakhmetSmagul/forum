package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UserQuery interface {
	CreateUser(user *model.User) error
	DeleteUser(userID int64) error
	UserVerification(user *model.User) error
	IsExistUser(user *model.User) (bool, error)
}

type userQuery struct {
	db *sql.DB
}

func (u *userQuery) CreateUser(user *model.User) error {
	sqlStmt := `INSERT INTO users(uname, email, passwd) 
	VALUES(?, ?, ?)`

	query, err := u.db.Prepare(sqlStmt)
	if err != nil {
		return err
	}

	defer query.Close()

	result, err := query.Exec(user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}

func (u *userQuery) UserVerification(user *model.User) error {
	sqlStmt := `SELECT user_id, username, password  FROM users WHERE email=?`
	query, err := u.db.Prepare(sqlStmt)
	if err != nil {
		return err
	}

	defer query.Close()

	tempPasswd := user.Password
	err = query.QueryRow(user.Email).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		log.Println(err)
		return err
	}

	bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(tempPasswd))
	if err != nil {
		log.Println("the user's password is incorrect", err)
		return errors.New("the user's password is incorrect")
	}

	return nil
}

func (u *userQuery) DeleteUser(userID int64) error {
	sqlStmt := `DELETE FROM users WHERE user_id=?`
	result, err := u.db.Exec(sqlStmt, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("delete user was failed")
	}

	return nil
}

func (u *userQuery) IsExistUser(user *model.User) (bool, error) {
	sqlStmt := `SELECT EXISTS(SELECT 1 FROM users WHERE uname=? OR email=? LIMIT 1)`
	var exist bool
	err := u.db.QueryRow(sqlStmt, user.Username, user.Email).Scan(&exist)
	if err != nil {
		log.Println(err)
		return false, err
	}

	return exist, nil
}
