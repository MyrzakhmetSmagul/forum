package main

import (
	"log"

	"github.com/MyrzakhmetSmagul/forum/internal/app"
	"github.com/MyrzakhmetSmagul/forum/internal/repository"
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
	// dao := repository.NewDao(db)
	// authService := service.NewAuthService(dao)
	err = app.Run()
	if err != nil {
		return
	}
}
