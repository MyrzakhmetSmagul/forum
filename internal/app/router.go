package app

import (
	"log"
	"net/http"
)

func (s *ServiceServer) Run() error {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./template/static"))))
	mux.HandleFunc("/", s.IndexWihtoutSession)

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