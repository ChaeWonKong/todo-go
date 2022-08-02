package repository

import (
	config "todo-clone/modules/confing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLDialector gorm.Dialector

func NewMySQLDialector(settings *config.Settings) MySQLDialector {
	return mysql.Open(settings.CreateDSN())
}
