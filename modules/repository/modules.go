package repository

import (
	"log"
	config "todo-clone/modules/confing"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

// type Repository interface {
// 	Create(value interface{}) (tx *gorm.DB)
// 	Limit(limit int) (tx *gorm.DB)
// 	Offset(offset int) (tx *gorm.DB)
// 	First(dest interface{}, conds ...interface{}) (tx *gorm.DB)
// 	Model(value interface{}) (tx *gorm.DB)
// 	Delete(value interface{}, conds ...interface{}) (tx *gorm.DB)
// }

type Repository struct {
	*gorm.DB
}

func NewRepository(
	settings *config.Settings,
	mysqlDialector MySQLDialector,
	sqliteDialector SQLiteDialector,
) *Repository {
	var dialector gorm.Dialector = mysqlDialector
	if settings.Debug() {
		dialector = sqliteDialector
	}

	db, err := gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	return &Repository{db}
}

var Modules = fx.Options(
	fx.Provide(NewRepository),
	fx.Provide(NewMySQLDialector),
	fx.Provide(NewSQLiteDialector),
	// fx.Provide(NewCreator),
	// fx.Provide(NewFirster),
	// fx.Provide(NewFinder),
)
