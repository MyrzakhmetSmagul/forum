package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

type SessionQuery interface {
	CreateSession(session *model.Session) error
	DeleteSession(session *model.Session) error
	GetSession(session *model.Session) error
}

type sessionQuery struct {
	db *sql.DB
}

func (s *sessionQuery) CreateSession(session *model.Session) error {
	sqlStmt := `INSERT INTO sessions(user_id, token, expiry) VALUES(?,?,?)`
	query, err := s.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}

	defer query.Close()

	result, err := query.Exec(session.User.ID, session.Token, session.Expiry)
	if err != nil {
		log.Println(err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return err
	}

	session.ID = id
	log.Println("CREATE SESSION WAS SUCCESFULE")
	return nil
}

func (s *sessionQuery) DeleteSession(session *model.Session) error {
	sqlStmt := `DELETE FROM sessions WHERE token=?`
	query, err := s.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}

	defer query.Close()

	result, err := query.Exec(session.Token)
	if err != nil {
		log.Println(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}

	if rowsAffected == 0 {
		log.Println("delete session was failed")
		return errors.New("delete session was failed")
	}

	return nil
}

func (s *sessionQuery) GetSession(session *model.Session) error {
	sqlStmt := `SELECT session_id, user_id, expiry FROM sessions WHERE token=?`
	query, err := s.db.Prepare(sqlStmt)
	if err != nil {
		log.Println("sessionQuery.GetSession", err)
		return err
	}

	defer query.Close()

	err = query.QueryRow(session.Token).Scan(&session.ID, &session.User.ID, &session.Expiry)
	if err != nil {
		log.Println("sessionQuery.GetSession", err)
		return err
	}
	log.Println("GET SESSION WAS SUCCESFULE")
	return nil
}
