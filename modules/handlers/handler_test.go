package handlers_test

import (
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
	h, m := fixtures()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")

	id := 1
	item := domains.Item{Title: "title1", ID: uint64(id)}
	idMatcher := func(id uint64) bool {
		return id == item.ID
	}

	matchedBy := mock.MatchedBy(idMatcher)
	m.On("FindOne", matchedBy).Return(item, nil)

	// Test
	t.Run("Correct id as param", func(t *testing.T) {
		responseJSON := `{"id":1,"title":"title1","checked":false,"created":"0001-01-01T00:00:00Z","updated":"0001-01-01T00:00:00Z"}`
		c.SetParamValues(strconv.Itoa(id))
		if assert.NoError(t, h.FindOne(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.JSONEq(t, responseJSON, rec.Body.String())
		}
		m.AssertCalled(t, "FindOne", matchedBy)
	})
	t.Run("Incorrect id as param", func(t *testing.T) {
		c.SetParamValues("x")
		assert.Panics(t, func() {
			c.SetParamValues("2")
			assert.Error(t, h.FindOne(c))
			m.AssertCalled(t, "FindOne", matchedBy)
		})
	})
}
