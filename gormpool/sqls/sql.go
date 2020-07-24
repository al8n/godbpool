package sqls

import "github.com/jinzhu/gorm"

type Connector interface {
	Open() (*gorm.DB, error)
	Args() interface{}
}
