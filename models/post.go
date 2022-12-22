package models

type Post struct {
	Id      string
	Title   string
	Content string
	UserId  string
}

func NewPost(id, title, content, userId string) *Post {
	return &Post{id, title, content, userId}
}
