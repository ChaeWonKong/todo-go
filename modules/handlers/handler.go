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
	e.DELETE("/:id", handler.Delete)
}

func (handler *Handler) FindOne(c echo.Context) error {
	ID := c.Param("id")

	id, err := strconv.Atoi(ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	item, err := handler.repo.FindOne(uint64(id))

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, item)
}

func (handler *Handler) FindAll(c echo.Context) error {
	items, err := handler.repo.FindAll()
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

	if err := c.Validate(item); err != nil {
		return err
	}

	affected, err := handler.repo.CreateOne(item)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if affected == 0 {
		return c.JSON(http.StatusConflict, nil)
	}

	return c.JSON(http.StatusCreated, item)
}

func (handler *Handler) Update(c echo.Context) error {
	ID := c.Param("id")
	id, err := strconv.Atoi(ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	body := &UpdateDto{}
	if err := c.Bind(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = c.Validate(body); err != nil {
		return err
	}

	if len(body.Title) == 0 && body.Checked == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "At least title or checked should be provided")
	}

	updateDto := make(map[string]interface{}, 2)

	if len(body.Title) != 0 {
		updateDto["Title"] = body.Title
	}
	if body.Checked != nil {
		updateDto["Checked"] = body.Checked
	}

	affected, err := handler.repo.UpdateOne(uint64(id), updateDto)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, affected)
}

func (handler *Handler) Delete(c echo.Context) error {
	ID := c.Param("id")
	id, err := strconv.Atoi(ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	affected, err := handler.repo.DeleteOne(uint64(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if affected == 0 {
		return c.JSON(http.StatusNotFound, nil)
	}

	return c.JSON(http.StatusNoContent, nil)
}
