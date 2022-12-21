package models

import "fmt"

type User struct {
	Id      int
	Name    string
	Surname string
	Gender  string
	Email   string
	Pwd     string
}

func (u *User) UserInfo() string {
	info := fmt.Sprintf("User id: %d", u.Id)
	info += fmt.Sprintf("\nFullname: %s %s", u.Name, u.Surname)
	info += fmt.Sprintf("\nGender: %s", u.Gender)
	return info
}
