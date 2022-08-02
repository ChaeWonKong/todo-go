package controller

import (
	"net/http"
	"strconv"
	"todo-clone/modules/delivery/rest/service"
	"todo-clone/modules/domains"

	echo "github.com/labstack/echo/v4"
)

type Controller struct {
	service.Service
}

type Params struct {
	Page  *int `query:"page"  json:"page,omitempty"  validate:"gte=1" default:"1"  example:"1"`
	Limit *int `query:"limit" json:"limit,omitempty" validate:"gte=1" default:"10" example:"10"`
}

type Body struct {
	Title string `json:"title"`
}

func (p Params) Pagination() (page, limit int) {
	if p.Limit == nil {
		limit = 10
	} else {
		limit = *p.Limit
	}

	if p.Page == nil {
		page = 0
	} else {
		page = limit * (*p.Page - 1)
	}
	return page, limit
}

func NewController(service service.Service) *Controller {
	return &Controller{service}
}

func BindRoutes(server *echo.Echo, controller *Controller) {
	group := server.Group("/todos")
	group.GET("", controller.GetAll)
	group.GET("/:id", controller.GetOne)
	group.POST("", controller.Post)
	group.PATCH("/:id", controller.Patch)
	group.DELETE("/:id", controller.Delete)
}

// GetAll find all todo items
// @ID todo-get-all
// @Tags todo
// @Summary find all todo items
// @Description find all todo items
// @Router /todos [get]
// @Param _ query Params true "params"
// @Produce json
// @Success 200 {object} interface{}
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
func (c Controller) GetAll(ctx echo.Context) error {
	var param Params

	if err := ctx.Bind(&param); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	page, limit := param.Pagination()
	items, err := c.FindAll(page, limit)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, items)
}

// Get todo get
// @ID todo-get
// @Tags todo
// @Summary todo get
// @Description todo get
// @Router /todos/{id} [get]
// @Param id path uint64 true "id"
// @Produce json
// @Success 200 {object} domains.Item
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
func (c Controller) GetOne(ctx echo.Context) error {
	ID := ctx.Param("id")
	id, err := strconv.Atoi(ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	item, err := c.FindOne(uint64(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, item)
}

// Post todo post
// @ID todo-post
// @Tags todo
// @Summary todo post
// @Description todo post
// @Router /todos [post]
// @Param body body Body true "body"
// @Produce json
// @Success 200 {object} domains.Item
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
func (c Controller) Post(ctx echo.Context) error {
	item := &domains.Item{}
	if err := ctx.Bind(item); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	affected, err := c.Insert(item)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if affected == 0 {
		return ctx.JSON(http.StatusConflict, nil)
	}

	return ctx.JSON(http.StatusOK, item)
}

// Patch todo patch
// @ID todo-patch
// @Tags todo
// @Summary todo patch
// @Description todo patch
// @Router /todos/{id} [patch]
// @Param id path uint64 true "id"
// @Param body body Body true "body"
// @Produce json
// @Success 200 {object} interface{}
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
func (c Controller) Patch(ctx echo.Context) error {
	ID := ctx.Param("id")
	id, err := strconv.Atoi(ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	body := &Body{}
	if err := ctx.Bind(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	updated, err := c.Update(uint64(id), body.Title)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if updated == nil {
		return ctx.JSON(http.StatusNotFound, nil)
	}

	return ctx.JSON(http.StatusOK, updated)
}

// Delete todo delete
// @ID todo-delete
// @Tags todo
// @Summary todo delete
// @Description todo delete
// @Router /todos/{id} [delete]
// @Param id path uint64 true "id"
// @Security ApiKeyAuth
// @Produce json
// @Success 204 ""
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
func (c Controller) Delete(ctx echo.Context) error {
	ID := ctx.Param("id")
	id, err := strconv.Atoi(ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	affected, err := c.Service.Delete(uint64(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if affected == 0 {
		return ctx.JSON(http.StatusNotFound, nil)
	}
	return ctx.JSON(http.StatusNoContent, nil)
}
