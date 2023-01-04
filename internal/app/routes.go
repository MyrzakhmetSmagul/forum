package app

import (
	"log"
	"net/http"
)

func Run() error {
	mux := http.NewServeMux()
	// mux.HandleFunc("/", index)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Printf("server started at http://localhost%s", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Println("error when starting the server", err)
		return err
	}
	return nil
}
