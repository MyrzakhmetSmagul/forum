package dao

import (
	"database/sql"
	"forum/models"
	"log"
)

func AddSession(db *sql.DB, session *models.Session) error {
	sqlStmt := `INSERT INTO sessions(token, user_id) VALUES(?,?)`
	query, err := db.Prepare(sqlStmt)
	if err != nil {
		log.Println("ADD SESSION FUNCTION ERROR:", err)
		return err
	}

	defer query.Close()

	_, err = query.Exec(session.Token, session.UserId)
	if err != nil {
		log.Println("ADD SESSION FUNCTION ERROR:", err)
		return err
	}

	sqlStmt = `SELECT token_id FROM sessions WHERE token=? and user_id=?`
	query, err = db.Prepare(sqlStmt)
	if err != nil {
		log.Println("ADD SESSION FUNCTION ERROR:", err)
		return err
	}

	err = query.QueryRow(session.Token, session.UserId).Scan(&session.TokenId)
	if err != nil {
		log.Println("ADD SESSION FUNCTION ERROR:", err)
		return err
	}
	return nil
}

func IsExistSession(db *sql.DB, session *models.Session) (bool, error) {
	sqlStmt := `SELECT 1 FROM sessions WHERE  token = ?`
	var exist bool
	err := db.QueryRow(sqlStmt, session.Token).Scan(&exist)
	if err != nil {
		return false, err
	}

	return exist, err
}

func GetSessionFromDB(db *sql.DB, session *models.Session) error {
	sqlStmt := `SELECT token_id, user_id FROM sessions WHERE token=?`
	query, err := db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}

	err = query.QueryRow(session.Token).Scan(&session.TokenId, &session.UserId)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
