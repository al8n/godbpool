package sqlite

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type SQLite struct {
	args interface{}
}

// init
func New(args interface{}) SQLite {
	return SQLite{args}
}

// build a connection with database
func (s SQLite) Open() (*gorm.DB, error) {
	return gorm.Open("sqlite3", s.args)
}

// return the connection args
func (s SQLite) Args() interface{} {
	return s.args
}
