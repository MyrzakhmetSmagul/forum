package model

type Dislike struct {
	UserID   int64
	Username string
	PostID   int64
	Status   bool
}
