package services

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

func NewUserService(dao repository.DAO) UserService {
	return &userService{
		dao.NewUserQuery(),
	}
}

func (u *userService) GetUserInfo(user *model.User) error {
	return u.UserQuery.GetUserInfo(user)
}
