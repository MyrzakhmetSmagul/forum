package app

import "github.com/MyrzakhmetSmagul/forum/internal/service"

type ServiceServer struct {
	authService    service.AuthService
	userService    service.UserService
	postService    service.PostService
	sessionService service.SessionService
}

func NewServiceServer(
	authService service.AuthService,
	userService service.UserService,
	postService service.PostService,
	sessionService service.SessionService,
) ServiceServer {
	return ServiceServer{
		authService:    authService,
		userService:    userService,
		postService:    postService,
		sessionService: sessionService,
	}
}
