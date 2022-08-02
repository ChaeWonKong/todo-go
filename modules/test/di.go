package test

import (
	"testing"
	config "todo-clone/modules/confing"
	"todo-clone/modules/delivery/rest"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func NewForTest(tb testing.TB, opts ...fx.Option) *fxtest.App {
	defaultOptions := []fx.Option{
		config.Modules,
		rest.Modules,
	}
	opts = append(defaultOptions, opts...)
	return fxtest.New(tb, opts...)
}
