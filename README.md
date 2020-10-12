<div align="center">
<h1>go-db-pool</h1>
</div>
<div align="center">

[![CI][CI-image]][CI-url]
[![Coverage Status][codecov-image]][codecov-url]
[![Go Report Card][go-report-image]][go-report-url]
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
[CI-image]: https://img.shields.io/github/workflow/status/ALiuGuanyan/godbpool/CI?event=push&style=flat-square
[codecov-image]: https://img.shields.io/codecov/c/gh/ALiuGuanyan/godbpool/master?style=flat-square
[codecov-url]: https://codecov.io/gh/ALiuGuanyan/go-db-pool
[go-report-image]: https://img.shields.io/badge/go%20report-A%2B-brightgreen?style=flat-square
[go-report-url]: https://goreportcard.com/report/github.com/ALiuGuanyan/godbpool
[license-image]: https://img.shields.io/badge/License-MIT-brightgreen.svg?style=flat-square
[license-url]: https://opensource.org/licenses/MIT

