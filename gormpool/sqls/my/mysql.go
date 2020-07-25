package my

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// My SQL
type My struct {
	args interface{}
}

// New will init a struct
func New(args interface{}) My {
	return My{args: args}
}

// Open will build a connection with database
func (m My) Open() (*gorm.DB, error) {
	return gorm.Open("mysql", m.args)
}

// Args will return the connection args
func (m My) Args() interface{} {
	return m.args
}
