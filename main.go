package main

import (
	"log"
	"todo-go/modules/domains"
	"todo-go/modules/handlers"
	"todo-go/modules/middlewares"
	"todo-go/modules/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Repository struct {
}

func main() {
	e := echo.New()

	e.Validator = &middlewares.CustomValidator{Validator: validator.New()}
	dialecter := repositories.NewSQLiteDialector()
	repository := repositories.NewRepository(dialecter)
	if err := repository.AutoMigrate(&domains.Item{}); err != nil {
		log.Fatal(err)
	}
	handler := handlers.NewHandler(repository)
	handlers.BindRoutes(e, handler)
	e.Logger.Fatal(e.Start(":8080"))
}
