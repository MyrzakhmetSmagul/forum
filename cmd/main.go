package main

import (
	"log"

	"github.com/MyrzakhmetSmagul/forum/internal/app"
	"github.com/MyrzakhmetSmagul/forum/internal/repository"
	"github.com/MyrzakhmetSmagul/forum/internal/services"
)

func main() {
	db, err := repository.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	dao := repository.NewDao(db)
	authService := services.NewAuthService(dao)
	postService := services.NewPostService(dao)
	sessionService := services.NewSessionService(dao)
	userService := services.NewUserService(dao)
	app := app.NewServiceServer(authService, userService, postService, sessionService)
	err = app.Run()
	if err != nil {
		return
	}
}
