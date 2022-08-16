package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	mocks "todo-go/mocks/repositories"
	"todo-go/modules/domains"
	"todo-go/modules/handlers"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func fixtures() (h *handlers.Handler, m *mocks.Repository) {
	m = &mocks.Repository{}
	h = handlers.NewHandler(m) // TODO: hand over instance not pointer

	return h, m
}

func TestFindOne(t *testing.T) {
	// Setup
	h, m := fixtures()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")

	id := 1
	item := domains.Item{Title: "title1", ID: uint64(id)}

	// Mock FindOne
	m.On("FindOne", mock.Anything).Return(func(id uint64) interface{} {
		if id == item.ID {
			return item
		}
		return nil
	}, func(id uint64) error {
		if id == item.ID {
			return nil
		}
		return errors.New("ID not found")
	})

	// Test cases
	t.Run("Correct id provided", func(t *testing.T) {
		responseJSON := `{"id":1,"title":"title1","checked":false,"created":"0001-01-01T00:00:00Z","updated":"0001-01-01T00:00:00Z"}`
		c.SetParamValues(strconv.Itoa(id))
		if assert.NoError(t, h.FindOne(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.JSONEq(t, responseJSON, rec.Body.String())
		}
	})

	t.Run("Not numeric id provided", func(t *testing.T) {
		c.SetParamValues("x")
		assert.Error(t, h.FindOne(c))
	})

	t.Run("Non-existing id provided", func(t *testing.T) {
		c.SetParamValues("-1")
		assert.Error(t, h.FindOne(c))
	})
}

func TestFindAll(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	items := make([]domains.Item, 0)
	t.Run("Success Case", func(t *testing.T) {
		h, m := fixtures()
		m.On("FindAll").Return(items, nil)
		if assert.NoError(t, h.FindAll(c)) {
			assert.JSONEq(t, "[]", rec.Body.String())
		}
	})

	t.Run("Fail Case", func(t *testing.T) {
		h, m := fixtures()
		m.On("FindAll").Return(nil, errors.New("An Error"))
		assert.Error(t, h.FindAll(c))
	})
}
