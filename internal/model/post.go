package model

type Post struct {
	ID         int64
	CategoryID int64
	Title      string
	Content    string
	UserID     int64
	Username   string
	Like       int
	Dislike    int
	Comments   []Comment
}
