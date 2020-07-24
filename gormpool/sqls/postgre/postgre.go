package postgre


import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Postgre struct {
	args interface{}
}

func New(args interface{}) Postgre  {
	return Postgre{args}
}

func (p Postgre) Open() (*gorm.DB, error)  {
	return gorm.Open("postgres", p.args)
}

func (p Postgre) Args() interface{}  {
	return p.args
}
