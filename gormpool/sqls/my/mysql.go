package my

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type My struct {
	args interface{}
}

// init
func New(args interface{}) My {
	return My{args: args}
}

// build a connection with database
func (m My) Open() (*gorm.DB, error) {
	return gorm.Open("mysql", m.args)
}

// return the connection args
func (m My) Args() interface{} {
	return m.args
}
