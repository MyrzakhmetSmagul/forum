package service

import (
	"errors"
	"log"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
	"github.com/MyrzakhmetSmagul/forum/internal/repository"
)

type AuthService interface {
	SignIn(user *model.User, session *model.Session) error
	SignUp(user *model.User) error
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

func (a *authService) SignIn(user *model.User, session *model.Session) error {
	err := a.UserQuery.UserVerification(user)
	if err != nil {
		log.Println(err)
		return err
	}

	err = a.SessionQuery.CreateSession(session)
	if err != nil {
		log.Println("\nuser was confirmed, CREATE SESSION ERROR:", err)
		return err
	}

	return nil
}

func (a *authService) SignUp(user *model.User) error {
	exist, err := a.UserQuery.IsExistUser(user)
	if err != nil {
		log.Println(err)
		return err
	}

	if exist {
		log.Println("user exist")
		return errors.New("user exist")
	}

	err = a.UserQuery.CreateUser(user)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (a *authService) LogOut(session *model.Session) error {
	return a.SessionQuery.DeleteSession(session)
}
