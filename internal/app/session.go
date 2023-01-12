package app

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
)

func (s *ServiceServer) getSession(r *http.Request) (model.Session, error) {
	cookie, err := r.Cookie("authToken")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return model.Session{}, fmt.Errorf("getSession: %w", model.ErrNoSession)
		}
		return model.Session{}, fmt.Errorf("getSession: %w", err)
	}

	session := model.Session{Token: cookie.Value}

	err = s.sessionService.GetSession(&session)
	if err != nil {
	}

	if session.Expiry.Before(time.Now()) {
		err = s.sessionService.DeleteSession(&session)
		if err != nil {
			return model.Session{}, fmt.Errorf("getSession: %w", err)
		}
		return model.Session{}, fmt.Errorf("getSession: %w", model.ErrNoSession)
	}

	err = s.userService.GetUserInfo(&session.User)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			err = s.sessionService.DeleteSession(&session)
			if err != nil {
				return model.Session{}, fmt.Errorf("getSession: %w", err)
			}
		}

		return model.Session{}, fmt.Errorf("getSession: %w", err)
	}

	return session, nil
}
