package mssql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

type Ms struct {
	args interface{}
}

func New(args interface{}) Ms {
	return Ms{args}
}

func (m Ms) Open() (*gorm.DB, error) {
	return gorm.Open("mysql", m.args)
}

func (m Ms) Args() interface{} {
	return m.args
}
