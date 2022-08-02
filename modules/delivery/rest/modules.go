package rest

import (
	"todo-clone/modules/delivery/rest/controller"
	"todo-clone/modules/delivery/rest/service"

	"go.uber.org/fx"
)

var Modules = fx.Options(
	fx.Provide(NewServer),
	fx.Provide(controller.NewController),
	fx.Provide(service.NewService),
	fx.Invoke(registerHook, controller.BindRoutes),
)
