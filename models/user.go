package models

import "fmt"

type User struct {
	Id    int
	UName string
	Email string
	Pwd   string
}

func (u *User) UserInfo() string {
	info := fmt.Sprintf("User id: %d", u.Id)
	info += fmt.Sprintf("\nUsername: %s", u.UName)
	return info
}
