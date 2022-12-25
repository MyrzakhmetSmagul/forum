package models

type Post struct {
	PostId  string
	UserId  string
	Title   string
	Content string
}

func NewPost(title, content, userId string) *Post {
	return &Post{
		UserId:  userId,
		Title:   title,
		Content: content,
	}
}
