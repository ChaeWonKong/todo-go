package handlers

import (
	"net/http"
	"todo-go/modules/domains"
	"todo-go/modules/repositories"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	repo repositories.Repository // Repository
}

func NewHandler(repo repositories.Repository) *Handler {
	return &Handler{repo}
}

func BindRoutes(e *echo.Echo, handler *Handler) {
	e.GET("", handler.FindAll)
	e.POST("", handler.Create)
}

func (handler *Handler) FindAll(c echo.Context) error {
	items := make([]domains.Item, 10)
	err := handler.repo.Find(&items).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, items)
}

func (handler *Handler) Create(c echo.Context) error {
	item := &domains.Item{}

	if err := c.Bind(item); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tx := handler.repo.Create(item)
	affected, err := tx.RowsAffected, tx.Error

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if affected == 0 {
		return c.JSON(http.StatusConflict, nil)
	}

	return c.JSON(http.StatusCreated, item)
}
