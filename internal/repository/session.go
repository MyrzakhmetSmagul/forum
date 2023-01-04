package repository

import "database/sql"

type SessionQuery interface{}

type sessionQuery struct {
	db *sql.DB
}
