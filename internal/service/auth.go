package service

import (
	"errors"
	"fmt"
	"log"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
	"github.com/MyrzakhmetSmagul/forum/internal/repository"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignIn(user *model.User, session *model.Session) error
	SignUp(user *model.User) error
	SignOut(session *model.Session) error
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
	log.Println(user.Email)
	err := a.UserQuery.UserVerification(user)
	if err != nil {
		log.Println(err)
		return err
	}

	token, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
		return err
	}

	session.Token = token.String()
	session.User = *user
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

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		log.Println("Generate from password ERROR:", err)
		return err
	}

	user.Password = string(hash)
	err = a.UserQuery.CreateUser(user)
	if err != nil {
		log.Println("CREATE USER ERROR:", err)
		return err
	}
	fmt.Println("user create, username:", user.Username)
	return nil
}

func (a *authService) SignOut(session *model.Session) error {
	return a.SessionQuery.DeleteSession(session)
}
