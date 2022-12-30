package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id     string
	UName  string
	Email  string
	Passwd string
}

func NewUser(uname, email, passwd string) *User {
	return &User{
		UName:  uname,
		Email:  email,
		Passwd: passwd,
	}
}

func (u *User) UserInfo() string {
	info := fmt.Sprintf("User id: %s", u.Id)
	info += fmt.Sprintf("\nUsername: %s", u.UName)
	return info
}

func PasswordHashing(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(hash), err
}
