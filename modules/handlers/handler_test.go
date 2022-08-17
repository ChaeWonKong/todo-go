package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	mocks "todo-go/mocks/repositories"
	"todo-go/modules/domains"
	"todo-go/modules/handlers"
	"todo-go/modules/middlewares"

	"github.com/go-playground/validator/v10"
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

func TestCreate(t *testing.T) {
	mockCreateJSON := `{"title": "a mock title"}`
	mockCreateFailJSON := `{}`
	e := echo.New()
	e.Validator = &middlewares.CustomValidator{Validator: validator.New()}

	t.Run("Success case", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(mockCreateJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		h, m := fixtures()
		// item := &domains.Item{}
		m.On("CreateOne", mock.Anything).Return(int64(1), nil)
		if assert.NoError(t, h.Create(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Contains(t, rec.Body.String(), "a mock title")
		}
	})

	t.Run("Fail cases", func(t *testing.T) {
		t.Run("Validation error", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(mockCreateFailJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			h, _ := fixtures()

			assert.Error(t, h.Create(c))
		})

		t.Run("Repository error", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(mockCreateJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			h, m := fixtures()
			m.On("CreateOne", mock.Anything).Return(int64(0), errors.New("An error"))

			assert.Error(t, h.Create(c))
		})

		t.Run("Affected == 0", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(mockCreateJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			h, m := fixtures()
			m.On("CreateOne", mock.Anything).Return(int64(0), nil)

			assert.NoError(t, h.Create(c))
			assert.Equal(t, http.StatusConflict, rec.Code)
		})
	})
}
