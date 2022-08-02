package config

import "fmt"

func (s Settings) CreateDSN() string {
	url := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		s.Database.User,
		s.Database.Password,
		s.Database.Host,
		s.Database.Port,
		s.Database.Name,
	)

	return url
}
