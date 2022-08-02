package config

import "fmt"

func (s Settings) BindPort() string {
	return fmt.Sprintf(":%d", s.App.Port)
}
