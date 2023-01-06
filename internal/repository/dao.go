package repository

import "database/sql"

type DAO interface {
	NewUserQuery() UserQuery
	NewSessionQuery() SessionQuery
	NewCategoryQuery() CategoryQuery
	NewCommentQuery() CommentQuery
	NewCommentReactionQuery() CommentReactionQuery
	NewPostQuery() PostQuery
	NewPostReactionQuery() PostReactionQuery
}

type dao struct {
	db *sql.DB
}

func NewDao(db *sql.DB) DAO {
	return &dao{
		db: db,
	}
}

func (dao *dao) NewUserQuery() UserQuery {
	return &userQuery{
		db: dao.db,
	}
}

func (dao *dao) NewSessionQuery() SessionQuery {
	return &sessionQuery{
		db: dao.db,
	}
}

func (dao *dao) NewCategoryQuery() CategoryQuery {
	return &categoryQuery{
		db: dao.db,
	}
}

func (dao *dao) NewCommentQuery() CommentQuery {
	return &commentQuery{
		db: dao.db,
	}
}

func (dao *dao) NewCommentReactionQuery() CommentReactionQuery {
	return &commentReactionQuery{
		db: dao.db,
	}
}

func (dao *dao) NewPostQuery() PostQuery {
	return &postQuery{
		db: dao.db,
	}
}
func (dao *dao) NewPostReactionQuery() PostReactionQuery {
	return &postReactionQuery{
		db: dao.db,
	}
}
