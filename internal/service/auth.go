package service

import (
	"github.com/MyrzakhmetSmagul/forum/internal/model"
	"github.com/MyrzakhmetSmagul/forum/internal/repository"
)

type AuthService interface {
	SignIn(user *model.User) error
}

type authService struct {
	repository.UserQuery
	repository.SessionQuery
}

func NewAuthService(dao repository.DAO) AuthService {
	return &authService{
		UserQuery:    dao.NewUserQuery(),
		SessionQuery: dao.NewSessionQuery(),
	}
}

func (a *authService) SignIn(user *model.User) error {

	return nil
}
