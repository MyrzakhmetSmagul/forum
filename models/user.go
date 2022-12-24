package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id    string
	UName string
	Email string
	Pwd   string
}

func NewUser(uname, email, pwd string) *User {
	return &User{
		UName: uname,
		Email: email,
		Pwd:   pwd,
	}
}

func (u *User) UserInfo() string {
	info := fmt.Sprintf("User id: %s", u.Id)
	info += fmt.Sprintf("\nUsername: %s", u.UName)
	return info
}

func GetPassordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(hash), err
}
