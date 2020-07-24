package sqls

import "github.com/jinzhu/gorm"

// conn interface
type Connector interface {
	// build a connection with database
	Open() (*gorm.DB, error)
	// return the connection args
	Args() interface{}
}
