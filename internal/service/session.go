package service

import (
	"github.com/MyrzakhmetSmagul/forum/internal/model"
	"github.com/MyrzakhmetSmagul/forum/internal/repository"
)

type SessionService interface {
	CreateSession(session *model.Session) error
	DeleteSession(session *model.Session) error
	GetSession(session *model.Session) error
}

type sessionService struct {
	repository.SessionQuery
}

func NewSessionService(dao repository.DAO) SessionService {
	return &sessionService{
		SessionQuery: dao.NewSessionQuery(),
	}
}

func (s *sessionService) CreateSession(session *model.Session) error {
	return s.CreateSession(session)
}

func (s *sessionService) DeleteSession(session *model.Session) error {
	return s.DeleteSession(session)
}

func (s *sessionService) GetSession(session *model.Session) error {
	return s.GetSession(session)
}
