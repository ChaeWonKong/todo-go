package handlers

import (
	"net/http"
	"strconv"
	"todo-go/modules/domains"
	"todo-go/modules/repositories"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	repo repositories.Repository
}

func NewHandler(repo repositories.Repository) *Handler {
	return &Handler{repo}
}

func BindRoutes(e *echo.Echo, handler *Handler) {
	e.GET("", handler.FindAll)
	e.GET("/:id", handler.FindOne)
	e.POST("", handler.Create)
	e.PATCH("/:id", handler.Update)
}

func (handler *Handler) FindOne(c echo.Context) error {
	item := *&domains.Item{}
	ID := c.Param("id")

	id, err := strconv.Atoi(ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = handler.repo.First(&item, uint64(id)).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, item)
}

func (handler *Handler) FindAll(c echo.Context) error {
	items := make([]domains.Item, 10)
	err := handler.repo.Find(&items).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
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

type UpdateDto struct {
	Title string `json:"title"`
}

func (handler *Handler) Update(c echo.Context) error {
	ID := c.Param("id")
	id, err := strconv.Atoi(ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	body := &UpdateDto{}
	if err := c.Bind(body); err != nil {

	}

	item := &domains.Item{}
	err = handler.repo.Model(item).Where("id = ?", uint64(id)).Update("title", body.Title).Error

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, item)
}
