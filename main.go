package main

import (
	"log"
	"todo-go/modules/domains"
	"todo-go/modules/handlers"
	"todo-go/modules/repositories"

	"github.com/labstack/echo/v4"
)

type Repository struct {
}

func main() {
	e := echo.New()
	dialecter := repositories.NewSQLiteDialector()
	repository := repositories.NewRepository(dialecter)
	if err := repository.AutoMigrate(&domains.Item{}); err != nil {
		log.Fatal(err)
	}
	handler := handlers.NewHandler(repository)
	handlers.BindRoutes(e, handler)
	e.Logger.Fatal(e.Start(":8080"))
}
