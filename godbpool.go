/*
Package godbpool : The concurrency fearless of Databases connection pool for Golang.
*/
package godbpool

// SQLType means Support SQL database type
type SQLType uint8

// Support SQL databases
const (
	MySQL SQLType = iota
	PostgreSQL
	SQLite3
	SQLServer
)
