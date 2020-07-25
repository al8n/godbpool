package mssql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

// Ms SQL Server
type Ms struct {
	args interface{}
}

// New: init a struct
func New(args interface{}) Ms {
	return Ms{args}
}

// Open: build a connection with database
func (m Ms) Open() (*gorm.DB, error) {
	return gorm.Open("mysql", m.args)
}

// Args: return the connection args
func (m Ms) Args() interface{} {
	return m.args
}
