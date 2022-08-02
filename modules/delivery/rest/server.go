package rest

import (
	"context"
	"fmt"
	"log"
	config "todo-clone/modules/confing"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/fx"
)

func NewServer(settings *config.Settings) *echo.Echo {
	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	return e
}

func registerHook(lifecycle fx.Lifecycle, server *echo.Echo, settings *config.Settings) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := server.Start(settings.BindPort())
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("Server running on: %s", settings.BindPort())
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Stop Http Server")
			return server.Shutdown(ctx)
		},
	})
}
