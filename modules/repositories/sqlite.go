package repositories

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLiteDialector gorm.Dialector

func NewSQLiteDialector() SQLiteDialector {
	return sqlite.Open("todolist.db")
}
