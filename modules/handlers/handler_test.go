package handlers_test

// func fixtures() (h *handlers.Handler, m *mocks.Repository) {
// 	m = &mocks.Repository{}
// 	h = handlers.NewHandler(m) // TODO: hand over instance not pointer

// 	return h, m
// }

// func TestFindOne(t *testing.T) {
// 	h, m := fixtures()
// 	e := echo.New()

// 	req := httptest.NewRequest(http.MethodGet, "/", nil)
// 	rec := httptest.NewRecorder()

// 	c := e.NewContext(req, rec)
// 	c.SetParamNames("id")
// 	c.SetParamValues("1")

// 	item := domains.Item{}
// 	m.On("First", &item, uint64(1)).Return(&gorm.DB{Error: nil})
// 	assert.NoError(t, h.FindOne(c))
// }

// // assert.NoError(t, err)
// // if assert.NoError(t, h.FindOne(c)) {
// // Test Blocks
// // item, err := service.Fetch(1)
// // assert.Error(t, err)
// // assert.Nil(t, item)
// // affected, err := service.Insert(&todo.Item{Title: "test"})
// // assert.NoError(t, err)
// // assert.NotZero(t, affected)
// // item, err = service.Fetch(1)
// // assert.NoError(t, err)
// // assert.NotNil(t, item)

// // fmt.Println(rec.Body.String())
// // assert.NotNil(t, rec.Body)
// // assert.Equal(t, mockItemJson, rec.Body.String())
// // }
// // }
