package sqls

import "github.com/jinzhu/gorm"

// Connector is conn interface and every SQL struct should
// implement this interface
type Connector interface {
	// build a connection with database
	Open() (*gorm.DB, error)
	// return the connection args
	Args() interface{}
}
