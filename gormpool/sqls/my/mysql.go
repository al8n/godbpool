package my

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type My struct {
	args interface{}
}

func New(args interface{}) My  {
	return My{args: args}
}

func (m My) Open() (*gorm.DB, error)  {
	return gorm.Open("mysql", m.args)
}

func (m My) Args() interface{}  {
	return m.args
}
