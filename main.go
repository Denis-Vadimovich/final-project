package main

import (
	"log"
	"os"

	"github.com/Denis-Vadimovich/final-project/pkg/db"
	"github.com/Denis-Vadimovich/final-project/pkg/server"
)

func main() {
	db.Init("scheduler.db")
	logger := log.New(os.Stdout, `serv `, log.LstdFlags|log.Lshortfile)
	//создание сервера
	serv := server.NewServer(logger)
	//запуск сервера

	err := serv.Server.ListenAndServe()
	if err != nil {
		logger.Fatal(err)
	}
}
