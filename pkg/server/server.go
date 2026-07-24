package server

import (
	"log"
	"net/http"
	"time"

	"github.com/Denis-Vadimovich/final-project/pkg/api"
)

type Server struct {
	Logger *log.Logger
	Server http.Server
}

func NewServer(logger *log.Logger) *Server {

	mux := http.NewServeMux()
	api.Init(mux)

	srv := &http.Server{
		Addr:         "0.0.0.0:7540",
		Handler:      mux,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		Logger: logger,
		Server: *srv,
	}
}
