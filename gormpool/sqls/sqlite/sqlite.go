package sqlite

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// SQLite
type SQLite struct {
	args interface{}
}

// New will init a struct
func New(args interface{}) SQLite {
	return SQLite{args}
}

// Open will build a connection with database
func (s SQLite) Open() (*gorm.DB, error) {
	return gorm.Open("sqlite3", s.args)
}

// Args will return the connection args
func (s SQLite) Args() interface{} {
	return s.args
}
