package app

import (
	"fmt"
	"log"

	"forum/internal/handler"
	"forum/internal/repository"
	"forum/internal/server"
	"forum/internal/service"

	_ "github.com/mattn/go-sqlite3"
)

const port = "8080"

func Run() {
	db, err := repository.OpenDB("store.db")
	if err != nil {
		log.Fatalf("ERROR ON: opening db: %s", err)
	}

	// repository is the most low level layout of the forum, it work with the db (finding, adding or deletion of information)
	repository := repository.NewRepository(db)

	// service works with repository data (have validation, check and conidered as mid level layout)
	service := service.NewService(repository)

	// initiation of the handlers (using the service data and functions)
	handler := handler.NewHandler(service)
	server := new(server.Server)

	fmt.Printf("Starting server at port %s\nhttp://localhost:%s/\n", port, port)

	if err := server.Run(port, handler.InitRoutes()); err != nil {
		log.Fatalf("error while running the server: %s", err.Error())
	}
}
