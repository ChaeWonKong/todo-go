package config

import (
	"encoding/json"
	"log"

	"github.com/Netflix/go-env"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/fx"
)

const (
	// AppVersion version
	AppVersion = "21.2.0"
	// AppName name
	AppName = "todolist-api"
)

type Settings struct {
	App struct {
		Version     string `env:"-,default=1.0.0" json:"version"`
		Name        string `env:"-,default=todolist-api" json:"name"`
		Port        int    `env:"PORT,default=5000" json:"port"`
		Environment string `env:"ENVIRONMENT,default=development" json:"environment"` // development, stage, production
	}
	Database struct {
		User     string `env:"DATABASE_USER,required=true" json:"user"`
		Password string `env:"DATABASE_PASSWORD,required=true" json:"password"`
		Name     string `env:"DATABASE_NAME,default=todolist" json:"name"`
		Host     string `env:"DATABASE_HOST,default=localhost" json:"host"`
		Port     int    `env:"DATABASE_PORT,default=3306" json:"port"`
	}
	Extras env.EnvSet `json:"-"`
}

// JSON 설정 값 출력
func JSON() string {
	settings := NewSettings()
	jsonBytes, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(jsonBytes)
}

func NewSettings() *Settings {
	var settings Settings
	extras, err := env.UnmarshalFromEnviron(&settings)

	if err != nil {
		// log.Fatal(err)
	}
	settings.Extras = extras
	settings.App.Version = AppVersion
	settings.App.Name = AppName
	return &settings
}

// Debug mode
func (s Settings) Debug() bool {
	return s.App.Environment == "development"
}

var Modules = fx.Options(fx.Provide(NewSettings))
