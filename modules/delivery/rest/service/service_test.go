package service_test

import (
	"context"
	"log"
	"os"
	"testing"
	config "todo-clone/modules/confing"
	"todo-clone/modules/delivery/rest/service"
	"todo-clone/modules/domains"
	"todo-clone/modules/repository"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

// type ServiceTestSuite struct {
// 	suite.Suite
// 	service service.Service
// }

//func NewController(service service.Service) *Controller {
// return &Controller{service}
// }
// func NewServiceTestSuite(service service.Service) *ServiceTestSuite {
// 	return &ServiceTestSuite{
// 		service: service,
// 	}
// }

// func TestServiceSuite(t *testing.T) {
// 	f := func() {
// 		suite.Run(t, new(ServiceTestSuite))
// 	}
// 	app := fxtest.New(t, TestModule, fx.Invoke(f))
// 	app.RequireStart()
// 	defer app.RequireStop()
// }

func registerHook(lifecycle fx.Lifecycle, settings *config.Settings) {
	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return os.Remove("todolist.db")
		},
	})
}

// TODO: mock these dependencies:
var TestModule = fx.Options(
	config.Modules,
	fx.Provide(service.NewService),
	fx.Invoke(registerHook),
	fx.Invoke(func(repo *repository.Repository) {
		if err := repo.AutoMigrate(&domains.Item{}); err != nil {
			log.Fatal(err)
		}
	}),
)

// func (suite *ServiceTestSuite) TestFindOne() {
// 	t := suite.T()
// 	item, err := suite.service.FindOne(1)
// 	assert.Error(t, err)
// 	assert.Nil(t, item)
// }

func TestFindOne(t *testing.T) {
	f := func() {

	}

	app := fxtest.New(t, TestModule, fx.Invoke(f))
	app.RequireStart()
	defer app.RequireStop()
}

// func TestFindAll(t *testing.T) {
// 	f := func() {

// 		result := make([]domains.Item, 0)
// 		result = append(result, domains.Item{
// 			Title: "test",
// 			ID:    1,
// 		})

// 		mockRepository := mocks.Repository{}
// 		mockRepository.On("Find", mock.Anything).Return(result, nil)

// 		s := service.NewService(&mockRepository)

// 		items, err := s.FindAll(0, 100)
// 		assert.NoError(t, err)
// 		assert.Equal(t, 1, len(items))
// 	}

// 	app := fxtest.New(t, TestModule, fx.Invoke(f))
// 	app.RequireStart()
// 	defer app.RequireStop()
// }

// func TestInsert(t *testing.T) {

// 	f := func() {
// 		mockService := &mocks.Service{}
// 		mockService.On("Insert", mock.Anything).Return(int64(1), nil)
// 		affected, err := mockService.Insert(&domains.Item{Title: "test"})
// 		assert.NoError(t, err)
// 		assert.NotZero(t, affected)
// 	}

// 	app := fxtest.New(t, TestModule, fx.Invoke(f))
// 	app.RequireStart()
// 	defer app.RequireStop()
// }
