package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	mocks "todo-go/mocks/repositories"
	"todo-go/modules/domains"
	"todo-go/modules/handlers"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func fixtures() (h *handlers.Handler, m *mocks.Repository) {
	m = &mocks.Repository{}
	h = handlers.NewHandler(m) // TODO: hand over instance not pointer

	return h, m
}

func TestFindOne(t *testing.T) {
	responseJSON := `{"id":0,"title":"title1","checked":false,"created":"0001-01-01T00:00:00Z","updated":"0001-01-01T00:00:00Z"}`

	h, m := fixtures()
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	item := domains.Item{Title: "title1"}
	m.On("FindOne", uint64(1)).Return(item, nil)
	if assert.NoError(t, h.FindOne(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, responseJSON, rec.Body.String())
	}
}
