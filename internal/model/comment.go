package model

type Comment struct {
	ID      int64
	Post    Post
	User    User
	Message string
}
