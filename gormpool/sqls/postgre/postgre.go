package postgre

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Postgre struct {
	args interface{}
}

// init
func New(args interface{}) Postgre {
	return Postgre{args}
}

// build a connection with database
func (p Postgre) Open() (*gorm.DB, error) {
	return gorm.Open("postgres", p.args)
}

// return the connection args
func (p Postgre) Args() interface{} {
	return p.args
}
