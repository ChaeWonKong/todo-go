package main

import (
	"log"
	"net/http"
	"todo-go/modules/domains"
	"todo-go/modules/handlers"
	"todo-go/modules/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Repository struct {
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func main() {
	e := echo.New()

	e.Validator = &CustomValidator{validator: validator.New()}
	dialecter := repositories.NewSQLiteDialector()
	repository := repositories.NewRepository(dialecter)
	if err := repository.AutoMigrate(&domains.Item{}); err != nil {
		log.Fatal(err)
	}
	handler := handlers.NewHandler(repository)
	handlers.BindRoutes(e, handler)
	e.Logger.Fatal(e.Start(":8080"))
}
