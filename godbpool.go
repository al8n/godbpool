/*
The concurrency fearless of Databases connection pool for Golang.
*/
package godbpool

type SQLType uint8

const (
	MySQL SQLType = iota
	PostgreSQL
	SQLite3
	SQLServer
)

