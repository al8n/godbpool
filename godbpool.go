package godbpool

type SQLType uint8

const (
	MySQL SQLType = iota
	PostgreSQL
	SQLite3
	SQLServer
)

