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
	GetUserInfo(user *model.User) error
}

type userQuery struct {
	db *sql.DB
}

func (u *userQuery) CreateUser(user *model.User) error {
	sqlStmt := `INSERT INTO users(username, email, password) 
	VALUES(?, ?, ?)`

	query, err := u.db.Prepare(sqlStmt)
	if err != nil {
		log.Println("CreateUser ERROR:", err)
		return err
	}

	defer query.Close()

	result, err := query.Exec(user.Username, user.Email, user.Password)
	if err != nil {
		log.Println("CreateUser ERROR:", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("CreateUser ERROR:", err)
		return err
	}

	user.ID = id
	return nil
}

func (u *userQuery) UserVerification(user *model.User) error {
	sqlStmt := `SELECT user_id, username, password FROM users WHERE email=?`
	query, err := u.db.Prepare(sqlStmt)
	if err != nil {
		log.Println("UserVerification ERROR:", err)
		return err
	}

	defer query.Close()

	tempPasswd := user.Password
	err = query.QueryRow(user.Email).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		log.Println("UserVerification", err)
		if err.Error() == "sql: no rows in result ser" {
			return errors.New("the user's email or password is incorrect")
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(tempPasswd))
	if err != nil {
		log.Println("the user's email or password is incorrect", err)
		return errors.New("the user's email or password is incorrect")
	}

	return nil
}

func (u *userQuery) DeleteUser(userID int64) error {
	sqlStmt := `DELETE FROM users WHERE user_id=?`
	result, err := u.db.Exec(sqlStmt, userID)
	if err != nil {
		log.Println("DeleteUser ERROR:", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("DeleteUser ERROR:", err)
		return err
	}

	if rowsAffected == 0 {
		log.Println("DeleteUser ERROR:", err)
		return errors.New("delete user was failed")
	}

	return nil
}

func (u *userQuery) IsExistUser(user *model.User) (bool, error) {
	sqlStmt := `SELECT EXISTS(SELECT 1 FROM users WHERE username=? OR email=? LIMIT 1)`
	var exist bool

	err := u.db.QueryRow(sqlStmt, user.Username, user.Email).Scan(&exist)
	if err != nil {
		log.Println("is exist user query ERROR:", err)
		return false, err
	}

	return exist, nil
}

func (u *userQuery) GetUserInfo(user *model.User) error {
	sqlStmt := `SELECT username, email, password FROM users WHERE user_id=?`
	query, err := u.db.Prepare(sqlStmt)
	if err != nil {
		log.Println("GetUserInfo ERROR", err)
		return err
	}

	defer query.Close()

	err = query.QueryRow(user.ID).Scan(&user.Username, &user.Email, &user.Password)
	if err != nil {
		log.Println("GetUserInfo ERROR", err)
		if err.Error() == "sql: no rows in result set" {
			log.Println("user doesn't exist")
			return errors.New("user doesn't exist")
		}
		return err
	}

	return nil
}
