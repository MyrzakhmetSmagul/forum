package service

import (
	"github.com/MyrzakhmetSmagul/forum/internal/model"
	"github.com/MyrzakhmetSmagul/forum/internal/repository"
)

type UserService interface {
	GetUserInfo(user *model.User) error
}

type userService struct {
	repository.UserQuery
}

func (us *userService) GetUserInfo(user *model.User) error {
	return us.UserQuery.GetUserInfo(user)
}
