package handlers_test

import (
	"encoding/json"
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

func TestUpdate(t *testing.T) {

	// Setup
	setUpUpdateTest := func(id string, dto map[string]interface{}) (echo.Context, *httptest.ResponseRecorder) {
		e := echo.New()
		e.Validator = &middlewares.CustomValidator{Validator: validator.New()}
		mockCreateJSON, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(string(mockCreateJSON)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(id)

		return c, rec
	}

	t.Run("Success Case", func(t *testing.T) {
		t.Run("Title only", func(t *testing.T) {
			updateDto := map[string]interface{}{
				"Title": "a mock title",
			}
			id := 1
			c, rec := setUpUpdateTest(strconv.Itoa(id), updateDto)
			h, m := fixtures()
			m.On("UpdateOne", uint64(id), updateDto).Return(domains.Item{
				Title: updateDto["Title"].(string),
			}, nil)

			assert.NoError(t, h.Update(c))
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), updateDto["Title"].(string))
		})

		t.Run("Checked only", func(t *testing.T) {
			checked := false
			updateDto := map[string]interface{}{
				"Checked": &checked,
			}
			id := 1
			c, rec := setUpUpdateTest(strconv.Itoa(id), updateDto)
			h, m := fixtures()
			m.On("UpdateOne", uint64(id), updateDto).Return(domains.Item{
				Checked: *updateDto["Checked"].(*bool),
			}, nil)

			assert.NoError(t, h.Update(c))
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), strconv.FormatBool(*updateDto["Checked"].(*bool)))
		})

		t.Run("Both title and checked", func(t *testing.T) {
			checked := false
			updateDto := map[string]interface{}{
				"Title":   "a mock title",
				"Checked": &checked,
			}
			id := 1
			c, rec := setUpUpdateTest(strconv.Itoa(id), updateDto)
			h, m := fixtures()
			m.On("UpdateOne", uint64(id), updateDto).Return(domains.Item{
				Title:   updateDto["Title"].(string),
				Checked: *updateDto["Checked"].(*bool),
			}, nil)

			assert.NoError(t, h.Update(c))
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), updateDto["Title"].(string))
			assert.Contains(t, rec.Body.String(), strconv.FormatBool(*updateDto["Checked"].(*bool)))
		})
	})

	t.Run("Fail Case", func(t *testing.T) {
		t.Run("provide id is not numeric", func(t *testing.T) {
			id := "invalid_id"
			c, _ := setUpUpdateTest(id, map[string]interface{}{})
			h, m := fixtures()
			m.On("UpdateOne", mock.Anything, mock.Anything).Return(domains.Item{}, nil)

			assert.Error(t, h.Update(c))
		})

		t.Run("bind body raises error", func(t *testing.T) {
			e := echo.New()
			e.Validator = &middlewares.CustomValidator{Validator: validator.New()}
			mockCreateJSON, _ := json.Marshal(`{"invalid": "dto"}`)
			id := "1"
			req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(string(mockCreateJSON)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(id)
			h, m := fixtures()
			m.On("UpdateOne", mock.Anything, mock.Anything).Return(domains.Item{}, nil)

			result := h.Update(c).(*echo.HTTPError)
			assert.Error(t, result)
			assert.Equal(t, http.StatusBadRequest, result.Code)
		})

		t.Run("validation failed", func(t *testing.T) {
			id := "1"
			c, _ := setUpUpdateTest(id, map[string]interface{}{})
			h, m := fixtures()
			m.On("UpdateOne", mock.Anything, mock.Anything).Return(domains.Item{}, nil)

			result := h.Update(c).(*echo.HTTPError)
			assert.Error(t, result)
			assert.Equal(t, http.StatusBadRequest, result.Code)
			assert.Equal(t, "At least title or checked should be provided", result.Message)
		})

		// t.Run("both title and checked not provided", func(t *testing.T) {

		// })

		t.Run("repository.UpdateOne returns error", func(t *testing.T) {
			id := "1"
			mockErrMsg := "Raise error"
			checked := false
			updateDto := map[string]interface{}{
				"Title":   "a mock title",
				"Checked": &checked,
			}
			c, _ := setUpUpdateTest(id, updateDto)
			h, m := fixtures()
			m.On("UpdateOne", mock.Anything, mock.Anything).Return(domains.Item{}, errors.New(mockErrMsg))

			result := h.Update(c).(*echo.HTTPError)
			assert.Equal(t, mockErrMsg, result.Message)
		})
	})
}
