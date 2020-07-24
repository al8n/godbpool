package sqlite

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type SQLite struct {
	args interface{}
}

func New(args interface{}) SQLite  {
	return SQLite{args}
}

func (s SQLite) Open() (*gorm.DB, error)  {
	return gorm.Open("sqlite3", s.args)
}

func (s SQLite) Args() interface{}  {
	return s.args
}
