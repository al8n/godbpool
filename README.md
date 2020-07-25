<div align="center">
<h1>go-db-pool</h1>
</div>
<div align="center">

[![CI][CI-image]][CI-url]
[![Coverage Status][codecov-image]][codecov-url]
[![Go Report Card](https://goreportcard.com/badge/github.com/ALiuGuanyan/godbpool)](https://goreportcard.com/report/github.com/ALiuGuanyan/godbpool)
[![License: MIT][license-image]][license-url]

English | [简体中文](README-zh_CN.md)

The concurrency fearless of Databases connection pool for Golang.

</div>

## Support & TODO
| Database/ORM | Status | Official Site | Client Repo |
| :---: | :---: | :---: | :---: |
| [gorm](https://github.com/ALiuGuanyan/godbpool/gormpool/) | <div align="center"><img src="images/correct.svg" width="24px" height="24px" /></div> | https://gorm.io | [gorm repo](https://github.com/go-gorm/gorm) |
| [xorm]() |  | https://xorm.io | [xorm repo](https://gitea.com/xorm/xorm) |
| [MongoDB]() |  | https://www.mongodb.com | [mongo-go-driver](https://github.com/mongodb/mongo-go-driver) |
| [etcd]() |  | https://etcd.io | [etcd v3 client](https://go.etcd.io/etcd/v3/client) |
| [TiKV]() |  | https://tikv.org | [TiKV client](https://github.com/tikv/client-go) |
| [Apache Cassandra]() |  | https://cassandra.apache.org | [gocql](https://github.com/gocql/gocql) |


## Examples
- [gorm pool](https://github.com/ALiuGuanyan/go-db-pool/blob/master/examples/gorm/main.go)
  ```go
  package main
  
  import (
  	"context"  
  	"github.com/ALiuGuanyan/godbpool"
  	"github.com/ALiuGuanyan/godbpool/gormpool"
  	"log"  
  	"time"
  )
  
  func main() {
  	// config options
  	opts := gormpool.Options{
  		Type:            godbpool.MySQL,
  		Args:            "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True",
  		KeepConn:        2,
  		Capacity:        5,
  		MaxWaitDuration: 2000 * time.Millisecond,
  	}
  
  	// create a pool  
  	p, err := gormpool.NewPool(context.Background(), opts)
  	if err != nil {
  		log.Println(err)
  		return
  	}
    
    // Get a conn
    conn, err := p.Get()
    if err != nil {
        log.Println(err)
    }
  
    // ... do some CURD
    
    // put conn back to pool 
    p.Put(conn)
  
    // check the pool status
    p.Status()
  
    // close pool
    p.Close()
  }
  ```
 
 
[CI-url]: https://github.com/ALiuGuanyan/go-db-pool/actions?query=workflow%3ACI
[CI-image]: https://github.com/ALiuGuanyan/go-db-pool/workflows/CI/badge.svg?branch=master
[codecov-image]: https://codecov.io/gh/ALiuGuanyan/go-db-pool/branch/master/graph/badge.svg
[codecov-url]: https://codecov.io/gh/ALiuGuanyan/go-db-pool
[go-report-image]: https://goreportcard.com/badge/github.com/ALiuGuanyan/godbpool
[go-report-url]: https://goreportcard.com/report/github.com/ALiuGuanyan/godbpool
[license-image]: https://img.shields.io/badge/License-MIT-blue.svg
[license-url]: https://opensource.org/licenses/MIT