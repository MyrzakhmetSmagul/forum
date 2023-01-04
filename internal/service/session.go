package service

import (
	"github.com/MyrzakhmetSmagul/forum/internal/model"
	"github.com/MyrzakhmetSmagul/forum/internal/repository"
)

type SessionService interface {
	NewSession(user *model.User) error
}

type sessionService struct {
	repository.SessionQuery
}
